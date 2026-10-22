package lab

import "github.com/bubskee/muster/engine"

// GoldenRun is one replay reduced to what the browser lab must reproduce. The
// browser build replays every choice subset of every trace and must match
// these exactly.
type GoldenRun struct {
	Lab      string        `json:"lab"`
	Trace    string        `json:"trace"`
	Choices  []string      `json:"choices"`
	Harm     []string      `json:"harm"`
	Terminal []string      `json:"terminal"`
	Events   []GoldenEvent `json:"events"`
}

type GoldenEvent struct {
	ID        string   `json:"id"`
	Status    string   `json:"status"`
	Acted     []string `json:"acted"`
	Observed  []string `json:"observed"`
	Blocked   []string `json:"blocked"`
	Responded []string `json:"responded"`
}

type GoldenSearch struct {
	Lab    string       `json:"lab"`
	Trace  string       `json:"trace"`
	Result SearchResult `json:"result"`
}

type GoldenRobust struct {
	Lab    string       `json:"lab"`
	Result RobustResult `json:"result"`
}

type GoldenFile struct {
	Runs     []GoldenRun    `json:"runs"`
	Searches []GoldenSearch `json:"searches"`
	Robust   []GoldenRobust `json:"robust"`
}

// Golden replays every choice subset of every trace of every lab, and runs
// search and robust on each, for checking the browser engine.
func Golden(labs []*Lab) (GoldenFile, error) {
	var file GoldenFile
	runs, err := goldenRuns(labs)
	if err != nil {
		return file, err
	}
	file.Runs = runs
	for _, l := range labs {
		for _, trace := range l.Traces() {
			result, err := l.Search(trace)
			if err != nil {
				return file, err
			}
			file.Searches = append(file.Searches, GoldenSearch{Lab: l.Manifest.ID, Trace: trace.ID, Result: result})
		}
		robust, err := l.Robust()
		if err != nil {
			return file, err
		}
		file.Robust = append(file.Robust, GoldenRobust{Lab: l.Manifest.ID, Result: robust})
	}
	return file, nil
}

func goldenRuns(labs []*Lab) ([]GoldenRun, error) {
	var out []GoldenRun
	for _, l := range labs {
		keys := l.ChoiceKeys()
		for _, trace := range l.Traces() {
			for mask := 0; mask < 1<<len(keys); mask++ {
				choices := maskKeys(keys, mask)
				run, err := l.Run(trace, choices)
				if err != nil {
					return nil, err
				}
				g := GoldenRun{
					Lab:      l.Manifest.ID,
					Trace:    trace.ID,
					Choices:  choices,
					Harm:     orEmpty(l.Harm(run)),
					Terminal: orEmpty(run.TerminalState),
				}
				for _, e := range run.Trace {
					g.Events = append(g.Events, GoldenEvent{
						ID:        e.EventID,
						Status:    string(e.Status),
						Acted:     orEmpty(actedInOrder(e)),
						Observed:  orEmpty(e.ObservedBy),
						Blocked:   orEmpty(e.BlockedBy),
						Responded: orEmpty(e.RespondedBy),
					})
				}
				out = append(out, g)
			}
		}
	}
	return out, nil
}

func actedInOrder(entry engine.TraceEntry) []string {
	var ids []string
	for _, c := range entry.Controls {
		if c.Acted {
			ids = append(ids, c.ControlID)
		}
	}
	return ids
}

func orEmpty(s []string) []string {
	if s == nil {
		return []string{}
	}
	return s
}
