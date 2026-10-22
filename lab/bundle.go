package lab

import (
	"encoding/json"
	"fmt"

	"github.com/bubskee/muster/engine"
)

// Bundle is every lab serialized with its raw YAML, for environments without
// a filesystem (the browser build). The YAML is re-parsed by the same config
// loader the CLI uses, so the browser replays exactly what the CLI replays.
type Bundle struct {
	Labs []BundleLab `json:"labs"`
}

type BundleLab struct {
	Manifest     Manifest        `json:"manifest"`
	ScenarioYAML string          `json:"scenario_yaml"`
	ControlsYAML string          `json:"controls_yaml"`
	Variants     []BundleVariant `json:"variants"`

	// Compiled holds every trace already parsed, for the browser engine,
	// which has no YAML parser. LoadBundle ignores it.
	Compiled []CompiledTrace `json:"compiled,omitempty"`
}

type CompiledTrace struct {
	ID       string                `json:"id"`
	Label    string                `json:"label"`
	Kind     string                `json:"kind,omitempty"`
	Note     string                `json:"note,omitempty"`
	Scenario engine.Scenario       `json:"scenario"`
	Controls []engine.EventControl `json:"controls"`
}

type BundleVariant struct {
	ID           string `json:"id"`
	ScenarioYAML string `json:"scenario_yaml"`
	ControlsYAML string `json:"controls_yaml"`
}

// Manifest JSON hides file paths, so carry the variant list explicitly.
func NewBundle(labs []*Lab) Bundle {
	b := Bundle{Labs: make([]BundleLab, 0, len(labs))}
	for _, l := range labs {
		bl := BundleLab{
			Manifest:     l.Manifest,
			ScenarioYAML: l.Baseline.ScenarioYAML,
			ControlsYAML: l.Baseline.ControlsYAML,
		}
		for _, v := range l.Variants {
			bl.Variants = append(bl.Variants, BundleVariant{
				ID:           v.ID,
				ScenarioYAML: v.ScenarioYAML,
				ControlsYAML: v.ControlsYAML,
			})
		}
		for _, t := range l.Traces() {
			ct := CompiledTrace{ID: t.ID, Label: t.Label, Scenario: t.Scenario}
			for _, mv := range l.Manifest.Variants {
				if mv.ID == t.ID {
					ct.Kind, ct.Note = mv.Kind, mv.Note
				}
			}
			for _, c := range t.Controls.Controls {
				if ec, ok := c.(engine.EventControl); ok {
					ct.Controls = append(ct.Controls, ec)
				}
			}
			bl.Compiled = append(bl.Compiled, ct)
		}
		b.Labs = append(b.Labs, bl)
	}
	return b
}

// LoadBundle rebuilds labs from bundle JSON.
func LoadBundle(data []byte) ([]*Lab, error) {
	var b Bundle
	if err := json.Unmarshal(data, &b); err != nil {
		return nil, fmt.Errorf("decode bundle: %w", err)
	}
	labs := make([]*Lab, 0, len(b.Labs))
	for _, bl := range b.Labs {
		docs := make(map[string][2]string, len(bl.Variants))
		for _, v := range bl.Variants {
			docs[v.ID] = [2]string{v.ScenarioYAML, v.ControlsYAML}
		}
		l, err := Build(bl.Manifest, bl.ScenarioYAML, bl.ControlsYAML, docs)
		if err != nil {
			return nil, fmt.Errorf("lab %q: %w", bl.Manifest.ID, err)
		}
		labs = append(labs, l)
	}
	return labs, nil
}
