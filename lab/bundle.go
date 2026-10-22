package lab

import (
	"encoding/json"
	"fmt"
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
