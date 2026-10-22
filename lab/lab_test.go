package lab_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/bubskee/muster/lab"
)

func load(t *testing.T, path string) *lab.Lab {
	t.Helper()
	l, err := lab.Load(path)
	if err != nil {
		t.Fatal(err)
	}
	return l
}

func TestToySearchFindsEarlyIsolationReversal(t *testing.T) {
	l := load(t, "../examples/toy-jackpot/lab.yaml")

	result, err := l.Search(l.Baseline)
	if err != nil {
		t.Fatal(err)
	}

	var found *lab.Finding
	for i, f := range result.Findings {
		if f.Kind == lab.FindingWorse && f.Added == "R2" && reflect.DeepEqual(f.Context, []string{"V1", "V2", "R1"}) {
			found = &result.Findings[i]
		}
	}
	if found == nil {
		t.Fatalf("expected V1+V2+R1 + R2 reversal, got %+v", result.Findings)
	}
	if len(found.Before) != 0 {
		t.Fatalf("containing Watchline should reach no adverse outcome, got %v", found.Before)
	}
	want := []string{"attacker:node-access", "attacker:cluster-credential"}
	if !reflect.DeepEqual(found.NewHarm, want) {
		t.Fatalf("new harm = %v, want %v", found.NewHarm, want)
	}
	if found.Divergence != "cluster-credential-theft" {
		t.Fatalf("divergence = %q, want cluster-credential-theft", found.Divergence)
	}
}

func TestToyReversalSurvivesSideBranchOnly(t *testing.T) {
	l := load(t, "../examples/toy-jackpot/lab.yaml")
	f := lab.Finding{
		Context: []string{"V1", "V2", "R1"},
		Added:   "R2",
		NewHarm: []string{"attacker:node-access"},
	}

	side, _ := l.Trace("side-branch")
	if ok, err := l.Holds(side, f); err != nil || !ok {
		t.Fatalf("reversal should survive the unrelated side branch (ok=%v err=%v)", ok, err)
	}

	// When the credential is replayed before discovery, revocation is already
	// too late, so the containing Watchline no longer contains and there is
	// no node-access reversal left to find.
	early, _ := l.Trace("replay-first")
	if ok, err := l.Holds(early, f); err != nil || ok {
		t.Fatalf("node-access reversal should vanish when replay precedes discovery (ok=%v err=%v)", ok, err)
	}
}

func TestHFHasNoNonMonotoneInteraction(t *testing.T) {
	l := load(t, "../examples/hf-july-2026/lab.yaml")
	for _, trace := range l.Traces() {
		result, err := l.Search(trace)
		if err != nil {
			t.Fatal(err)
		}
		if len(result.Findings) != 0 {
			t.Fatalf("%s: expected no findings, got %+v", trace.ID, result.Findings)
		}
	}
}

func TestHFLabMatchesExperimentTable(t *testing.T) {
	l := load(t, "../examples/hf-july-2026/lab.yaml")
	cases := map[string][]string{
		"V1":    {"attacker:node-root", "attacker:internal-network-access", "attacker:cluster-admin", "attacker:worker-secrets"},
		"V1,P1": {"attacker:worker-secrets"},
		"V1,P2": {"attacker:node-root", "attacker:internal-network-access", "attacker:worker-secrets"},
		"V1,R1": {"attacker:worker-secrets"},
		"V1,P3": {"attacker:node-root", "attacker:internal-network-access", "attacker:cluster-admin"},
	}
	for keys, want := range cases {
		run, err := l.Run(l.Baseline, splitKeys(keys))
		if err != nil {
			t.Fatal(err)
		}
		if got := l.Harm(run); !reflect.DeepEqual(got, want) {
			t.Errorf("%s: harm = %v, want %v", keys, got, want)
		}
	}
}

func TestEveryExampleLabLoadsAndBundles(t *testing.T) {
	labs, err := lab.LoadAll("../examples")
	if err != nil {
		t.Fatal(err)
	}
	if len(labs) < 3 {
		t.Fatalf("expected at least 3 labs, got %d", len(labs))
	}

	data, err := json.Marshal(lab.NewBundle(labs))
	if err != nil {
		t.Fatal(err)
	}
	rebuilt, err := lab.LoadBundle(data)
	if err != nil {
		t.Fatal(err)
	}

	// The browser replays from the bundle; it must agree with the CLI.
	for i, l := range labs {
		r := rebuilt[i]
		for _, preset := range l.Manifest.Presets {
			for j, trace := range l.Traces() {
				a, err := l.Run(trace, preset.Choices)
				if err != nil {
					t.Fatal(err)
				}
				b, err := r.Run(r.Traces()[j], preset.Choices)
				if err != nil {
					t.Fatal(err)
				}
				if !reflect.DeepEqual(l.Harm(a), r.Harm(b)) {
					t.Fatalf("%s/%s/%s: bundle disagrees with CLI", l.Manifest.ID, trace.ID, preset.Label)
				}
			}
		}
		if _, err := l.Robust(); err != nil {
			t.Fatalf("%s: robust: %v", l.Manifest.ID, err)
		}
	}
}

func splitKeys(s string) []string {
	var out []string
	start := 0
	for i := 0; i <= len(s); i++ {
		if i == len(s) || s[i] == ',' {
			out = append(out, s[start:i])
			start = i + 1
		}
	}
	return out
}
