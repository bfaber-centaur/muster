package lab

import (
	"fmt"
	"math/bits"
	"sort"

	"github.com/bubskee/muster/engine"
)

// MaxSearchChoices bounds exhaustive search; 2^n replays are run.
const MaxSearchChoices = 18

// Finding is a pair of Watchlines that differ by one added choice.
//
// Kind "worse": the added choice strictly grows the set of adverse outcomes.
// Kind "tradeoff": it removes some adverse outcomes and adds others.
type Finding struct {
	Kind       string   `json:"kind"`
	Context    []string `json:"context"`
	Added      string   `json:"added"`
	Before     []string `json:"before"`
	After      []string `json:"after"`
	NewHarm    []string `json:"new_harm"`
	Divergence string   `json:"divergence,omitempty"`
}

const (
	FindingWorse    = "worse"
	FindingTradeoff = "tradeoff"
)

type SearchResult struct {
	LabID     string    `json:"lab_id"`
	TraceID   string    `json:"trace_id"`
	Choices   int       `json:"choices"`
	Watchline int       `json:"watchlines"`
	Findings  []Finding `json:"findings"`
}

type harmTable struct {
	keys []string
	harm []uint64
	runs []engine.RunResult
}

func (l *Lab) table(trace Trace) (harmTable, error) {
	keys := l.ChoiceKeys()
	if len(keys) > MaxSearchChoices {
		return harmTable{}, fmt.Errorf("%d choices exceeds exhaustive search limit of %d", len(keys), MaxSearchChoices)
	}

	n := len(keys)
	t := harmTable{
		keys: keys,
		harm: make([]uint64, 1<<n),
		runs: make([]engine.RunResult, 1<<n),
	}
	for mask := 0; mask < 1<<n; mask++ {
		run, err := l.Run(trace, maskKeys(keys, mask))
		if err != nil {
			return harmTable{}, err
		}
		t.runs[mask] = run
		t.harm[mask] = l.harmBits(run)
	}
	return t, nil
}

func (l *Lab) harmBits(run engine.RunResult) uint64 {
	terminal := make(map[string]struct{}, len(run.TerminalState))
	for _, fact := range run.TerminalState {
		terminal[fact] = struct{}{}
	}
	var b uint64
	for i, outcome := range l.Manifest.Adverse {
		if _, ok := terminal[outcome.Fact]; ok {
			b |= 1 << i
		}
	}
	return b
}

func (l *Lab) harmFacts(b uint64) []string {
	facts := []string{}
	for i, outcome := range l.Manifest.Adverse {
		if b&(1<<i) != 0 {
			facts = append(facts, outcome.Fact)
		}
	}
	return facts
}

func maskKeys(keys []string, mask int) []string {
	out := []string{}
	for i, key := range keys {
		if mask&(1<<i) != 0 {
			out = append(out, key)
		}
	}
	return out
}

// Search replays every subset of the lab's choices against trace and reports
// each place where adding one choice makes the modeled outcome worse.
//
// Only minimal contexts are reported: a finding is dropped when a smaller
// context already shows the same added choice producing the same new harm.
func (l *Lab) Search(trace Trace) (SearchResult, error) {
	t, err := l.table(trace)
	if err != nil {
		return SearchResult{}, err
	}
	n := len(t.keys)

	type raw struct {
		kind  string
		ctx   int
		added int
		new   uint64
	}
	var found []raw

	for mask := 0; mask < 1<<n; mask++ {
		for c := 0; c < n; c++ {
			bit := 1 << c
			if mask&bit != 0 {
				continue
			}
			before, after := t.harm[mask], t.harm[mask|bit]
			gained := after &^ before
			if gained == 0 {
				continue
			}
			kind := FindingWorse
			if before&^after != 0 {
				kind = FindingTradeoff
			}
			found = append(found, raw{kind: kind, ctx: mask, added: c, new: gained})
		}
	}

	sort.SliceStable(found, func(i, j int) bool {
		return bits.OnesCount(uint(found[i].ctx)) < bits.OnesCount(uint(found[j].ctx))
	})

	var kept []raw
	for _, f := range found {
		redundant := false
		for _, k := range kept {
			if k.kind == f.kind && k.added == f.added && k.new == f.new && k.ctx&f.ctx == k.ctx {
				redundant = true
				break
			}
		}
		if !redundant {
			kept = append(kept, f)
		}
	}

	result := SearchResult{
		LabID:     l.Manifest.ID,
		TraceID:   trace.ID,
		Choices:   n,
		Watchline: 1 << n,
		Findings:  []Finding{},
	}
	for _, k := range kept {
		withAdded := k.ctx | 1<<k.added
		result.Findings = append(result.Findings, Finding{
			Kind:       k.kind,
			Context:    maskKeys(t.keys, k.ctx),
			Added:      t.keys[k.added],
			Before:     l.harmFacts(t.harm[k.ctx]),
			After:      l.harmFacts(t.harm[withAdded]),
			NewHarm:    l.harmFacts(k.new),
			Divergence: FirstDivergence(t.runs[k.ctx], t.runs[withAdded]),
		})
	}

	// Worse before tradeoff, then smaller contexts first.
	sort.SliceStable(result.Findings, func(i, j int) bool {
		a, b := result.Findings[i], result.Findings[j]
		if a.Kind != b.Kind {
			return a.Kind == FindingWorse
		}
		return len(a.Context) < len(b.Context)
	})
	return result, nil
}

// Holds reports whether a finding reproduces in trace: the same context plus
// the same added choice still produces at least the same new harm.
func (l *Lab) Holds(trace Trace, f Finding) (bool, error) {
	before, err := l.Run(trace, f.Context)
	if err != nil {
		return false, err
	}
	after, err := l.Run(trace, append(append([]string(nil), f.Context...), f.Added))
	if err != nil {
		return false, err
	}
	gained := l.harmBits(after) &^ l.harmBits(before)
	var want uint64
	for i, outcome := range l.Manifest.Adverse {
		for _, fact := range f.NewHarm {
			if outcome.Fact == fact {
				want |= 1 << i
			}
		}
	}
	return gained&want == want, nil
}

// FirstDivergence returns the first event whose status, acting controls, or
// resulting state differs between two runs of the same scenario.
func FirstDivergence(a, b engine.RunResult) string {
	for i := range a.Trace {
		if i >= len(b.Trace) {
			return a.Trace[i].EventID
		}
		ea, eb := a.Trace[i], b.Trace[i]
		if ea.Status != eb.Status || !equal(ea.After, eb.After) || !equal(acted(ea), acted(eb)) {
			return ea.EventID
		}
	}
	return ""
}

func acted(entry engine.TraceEntry) []string {
	var ids []string
	for _, c := range entry.Controls {
		if c.Acted {
			ids = append(ids, c.ControlID)
		}
	}
	sort.Strings(ids)
	return ids
}

func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
