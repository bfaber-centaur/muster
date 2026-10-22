package lab

// Cell is the adverse outcome of one Watchline replayed against one trace.
type Cell struct {
	TraceID string   `json:"trace_id"`
	Harm    []string `json:"harm"`
	// Change relative to the baseline trace: "same", "worse", "better", or
	// "different" (some outcomes gained, others lost).
	Change string `json:"change"`
}

type RobustRow struct {
	Label   string   `json:"label"`
	Choices []string `json:"choices"`
	Cells   []Cell   `json:"cells"`
}

// FindingCheck records whether a baseline search finding reproduces in each
// variant trace.
type FindingCheck struct {
	Finding Finding         `json:"finding"`
	Holds   map[string]bool `json:"holds"`
}

type RobustResult struct {
	LabID    string         `json:"lab_id"`
	Traces   []TraceInfo    `json:"traces"`
	Rows     []RobustRow    `json:"rows"`
	Findings []FindingCheck `json:"findings"`
	// Findings that exist in a variant but not in the baseline.
	Emergent map[string][]Finding `json:"emergent"`
}

type TraceInfo struct {
	ID    string `json:"id"`
	Label string `json:"label"`
}

// Robust replays every preset against the baseline and every variant, and
// checks whether each baseline search finding survives in each variant.
func (l *Lab) Robust() (RobustResult, error) {
	traces := l.Traces()
	out := RobustResult{LabID: l.Manifest.ID, Emergent: map[string][]Finding{}}
	for _, t := range traces {
		out.Traces = append(out.Traces, TraceInfo{ID: t.ID, Label: t.Label})
	}

	rows := []Preset{{Label: "No added defense", Choices: []string{}}}
	rows = append(rows, l.Manifest.Presets...)

	for _, preset := range rows {
		row := RobustRow{Label: preset.Label, Choices: preset.Choices}
		var base uint64
		for i, trace := range traces {
			run, err := l.Run(trace, preset.Choices)
			if err != nil {
				return RobustResult{}, err
			}
			harm := l.harmBits(run)
			if i == 0 {
				base = harm
			}
			row.Cells = append(row.Cells, Cell{
				TraceID: trace.ID,
				Harm:    l.harmFacts(harm),
				Change:  change(base, harm),
			})
		}
		out.Rows = append(out.Rows, row)
	}

	baseline, err := l.Search(l.Baseline)
	if err != nil {
		return RobustResult{}, err
	}
	for _, f := range baseline.Findings {
		check := FindingCheck{Finding: f, Holds: map[string]bool{}}
		for _, trace := range l.Variants {
			ok, err := l.Holds(trace, f)
			if err != nil {
				return RobustResult{}, err
			}
			check.Holds[trace.ID] = ok
		}
		out.Findings = append(out.Findings, check)
	}

	for _, trace := range l.Variants {
		search, err := l.Search(trace)
		if err != nil {
			return RobustResult{}, err
		}
		for _, f := range search.Findings {
			if !containsFinding(baseline.Findings, f) {
				out.Emergent[trace.ID] = append(out.Emergent[trace.ID], f)
			}
		}
	}
	return out, nil
}

func change(base, harm uint64) string {
	switch {
	case harm == base:
		return "same"
	case harm&^base == 0:
		return "better"
	case base&^harm == 0:
		return "worse"
	default:
		return "different"
	}
}

func containsFinding(list []Finding, f Finding) bool {
	for _, g := range list {
		if g.Kind == f.Kind && g.Added == f.Added && equal(g.Context, f.Context) && equal(g.NewHarm, f.NewHarm) {
			return true
		}
	}
	return false
}
