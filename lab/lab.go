// Package lab loads Muster Lab manifests.
//
// A manifest is presentation and experiment-selection metadata layered over
// an ordinary scenario and control pool: which facts count as adverse, which
// human-facing choices expand to which engine controls, presets, labels, and
// alternative traces for robustness checks. Replay semantics stay in the
// scenario/control YAML and the engine; a manifest only selects.
package lab

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/bubskee/muster/config"
	"github.com/bubskee/muster/engine"
	"gopkg.in/yaml.v3"
)

type Outcome struct {
	Fact  string `yaml:"fact" json:"fact"`
	Label string `yaml:"label" json:"label"`
}

type Choice struct {
	Key      string   `yaml:"key" json:"key"`
	Role     string   `yaml:"role" json:"role"`
	Label    string   `yaml:"label" json:"label"`
	Help     string   `yaml:"help" json:"help,omitempty"`
	Controls []string `yaml:"controls" json:"controls"`
}

type Preset struct {
	Label   string   `yaml:"label" json:"label"`
	Choices []string `yaml:"choices" json:"choices"`
	Note    string   `yaml:"note" json:"note,omitempty"`
}

// Variant is an alternative trace of the same incident: a plausible branch
// the record does not show, or a branch that opens if a defense that held in
// the real incident had failed.
type Variant struct {
	ID       string `yaml:"id" json:"id"`
	Label    string `yaml:"label" json:"label"`
	Kind     string `yaml:"kind" json:"kind,omitempty"`
	Scenario string `yaml:"scenario" json:"-"`
	Controls string `yaml:"controls" json:"-"`
	Note     string `yaml:"note" json:"note,omitempty"`
}

type Source struct {
	Title string `yaml:"title" json:"title"`
	URL   string `yaml:"url" json:"url,omitempty"`
}

type Manifest struct {
	ID       string            `yaml:"id" json:"id"`
	Title    string            `yaml:"title" json:"title"`
	Domain   string            `yaml:"domain" json:"domain"`
	Kind     string            `yaml:"kind" json:"kind"`
	Summary  string            `yaml:"summary" json:"summary"`
	Note     string            `yaml:"note" json:"note,omitempty"`
	Order    int               `yaml:"order" json:"order"`
	Scenario string            `yaml:"scenario" json:"-"`
	Controls string            `yaml:"controls" json:"-"`
	Budget   int               `yaml:"budget" json:"budget,omitempty"`
	Adverse  []Outcome         `yaml:"adverse" json:"adverse"`
	Events   map[string]string `yaml:"events" json:"events"`
	Facts    map[string]string `yaml:"facts" json:"facts"`
	Choices  []Choice          `yaml:"choices" json:"choices"`
	Presets  []Preset          `yaml:"presets" json:"presets"`
	Variants []Variant         `yaml:"variants" json:"variants,omitempty"`
	Sources  []Source          `yaml:"sources" json:"sources,omitempty"`
}

// Trace is one loaded scenario/control pair: the baseline or a variant.
type Trace struct {
	ID           string
	Label        string
	Scenario     engine.Scenario
	Controls     config.ControlSet
	ScenarioYAML string
	ControlsYAML string
}

type Lab struct {
	Manifest Manifest
	Path     string
	Baseline Trace
	Variants []Trace
}

// Kinds of lab.
const (
	KindPublicRecord = "public-record"
	KindSynthetic    = "synthetic"
	KindIllustrative = "illustrative"
)

func Load(path string) (*Lab, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read lab %q: %w", path, err)
	}
	manifest, err := ParseManifest(data)
	if err != nil {
		return nil, fmt.Errorf("parse lab %q: %w", path, err)
	}

	dir := filepath.Dir(path)
	read := func(rel string) (string, error) {
		b, err := os.ReadFile(filepath.Join(dir, rel))
		if err != nil {
			return "", err
		}
		return string(b), nil
	}

	scenarioYAML, err := read(manifest.Scenario)
	if err != nil {
		return nil, fmt.Errorf("lab %q: %w", manifest.ID, err)
	}
	controlsYAML, err := read(manifest.Controls)
	if err != nil {
		return nil, fmt.Errorf("lab %q: %w", manifest.ID, err)
	}

	variants := make(map[string][2]string, len(manifest.Variants))
	for _, v := range manifest.Variants {
		s, err := read(v.Scenario)
		if err != nil {
			return nil, fmt.Errorf("lab %q variant %q: %w", manifest.ID, v.ID, err)
		}
		c := controlsYAML
		if v.Controls != "" {
			if c, err = read(v.Controls); err != nil {
				return nil, fmt.Errorf("lab %q variant %q: %w", manifest.ID, v.ID, err)
			}
		}
		variants[v.ID] = [2]string{s, c}
	}

	l, err := Build(manifest, scenarioYAML, controlsYAML, variants)
	if err != nil {
		return nil, fmt.Errorf("lab %q: %w", path, err)
	}
	l.Path = path
	return l, nil
}

// Build assembles a lab from a manifest and raw YAML documents. variants maps
// variant ID to its {scenario, controls} YAML. It is shared by the CLI loader
// and the browser build, which has no filesystem.
func Build(manifest Manifest, scenarioYAML, controlsYAML string, variants map[string][2]string) (*Lab, error) {
	baseline, err := loadTrace("baseline", "As recorded", scenarioYAML, controlsYAML)
	if err != nil {
		return nil, err
	}

	l := &Lab{Manifest: manifest, Baseline: baseline}

	for _, v := range manifest.Variants {
		docs, ok := variants[v.ID]
		if !ok {
			return nil, fmt.Errorf("variant %q: missing YAML", v.ID)
		}
		trace, err := loadTrace(v.ID, v.Label, docs[0], docs[1])
		if err != nil {
			return nil, fmt.Errorf("variant %q: %w", v.ID, err)
		}
		l.Variants = append(l.Variants, trace)
	}

	if err := l.validate(); err != nil {
		return nil, err
	}
	return l, nil
}

func loadTrace(id, label, scenarioYAML, controlsYAML string) (Trace, error) {
	scenario, err := config.ParseScenario([]byte(scenarioYAML))
	if err != nil {
		return Trace{}, fmt.Errorf("scenario: %w", err)
	}
	controls, err := config.ParseControlSet([]byte(controlsYAML))
	if err != nil {
		return Trace{}, fmt.Errorf("controls: %w", err)
	}
	return Trace{
		ID:           id,
		Label:        label,
		Scenario:     scenario,
		Controls:     controls,
		ScenarioYAML: scenarioYAML,
		ControlsYAML: controlsYAML,
	}, nil
}

func ParseManifest(data []byte) (Manifest, error) {
	var m Manifest
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	decoder.KnownFields(true)
	if err := decoder.Decode(&m); err != nil {
		return Manifest{}, fmt.Errorf("decode YAML: %w", err)
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		return Manifest{}, fmt.Errorf("multiple YAML documents are not supported")
	}

	if strings.TrimSpace(m.ID) == "" {
		return Manifest{}, fmt.Errorf("lab id is required")
	}
	if m.Scenario == "" || m.Controls == "" {
		return Manifest{}, fmt.Errorf("lab scenario and controls are required")
	}
	switch m.Kind {
	case KindPublicRecord, KindSynthetic, KindIllustrative:
	default:
		return Manifest{}, fmt.Errorf("lab kind %q: want %s, %s, or %s", m.Kind, KindPublicRecord, KindSynthetic, KindIllustrative)
	}
	if len(m.Adverse) == 0 {
		return Manifest{}, fmt.Errorf("lab must define at least one adverse outcome")
	}
	if len(m.Choices) == 0 {
		return Manifest{}, fmt.Errorf("lab must define at least one choice")
	}
	return m, nil
}

func (l *Lab) validate() error {
	m := l.Manifest

	keys := make(map[string]struct{}, len(m.Choices))
	for _, choice := range m.Choices {
		if choice.Key == "" {
			return fmt.Errorf("choice %q: key is required", choice.Label)
		}
		if _, dup := keys[choice.Key]; dup {
			return fmt.Errorf("duplicate choice key %q", choice.Key)
		}
		keys[choice.Key] = struct{}{}

		switch engine.Role(choice.Role) {
		case engine.Vedette, engine.Picket, engine.Reserve:
		default:
			return fmt.Errorf("choice %q: unknown role %q", choice.Key, choice.Role)
		}
		if len(choice.Controls) == 0 {
			return fmt.Errorf("choice %q: controls must not be empty", choice.Key)
		}
	}

	for _, preset := range m.Presets {
		for _, key := range preset.Choices {
			if _, ok := keys[key]; !ok {
				return fmt.Errorf("preset %q: unknown choice %q", preset.Label, key)
			}
		}
	}

	// Every choice must resolve in every trace, so a choice means the same
	// defense whichever variant is replayed.
	for _, trace := range l.Traces() {
		ids := controlIDs(trace.Controls)
		for _, choice := range m.Choices {
			for _, id := range choice.Controls {
				if _, ok := ids[id]; !ok {
					return fmt.Errorf("%s: choice %q references unknown control %q", trace.ID, choice.Key, id)
				}
			}
		}
	}

	events := make(map[string]struct{})
	for _, trace := range l.Traces() {
		for _, e := range trace.Scenario.Events {
			events[e.ID] = struct{}{}
		}
	}
	for id := range m.Events {
		if _, ok := events[id]; !ok {
			return fmt.Errorf("event label for unknown event %q", id)
		}
	}
	return nil
}

func controlIDs(set config.ControlSet) map[string]struct{} {
	ids := make(map[string]struct{}, len(set.Controls))
	for _, control := range set.Controls {
		if c, ok := control.(engine.EventControl); ok {
			ids[c.ID] = struct{}{}
		}
	}
	return ids
}

// Traces returns the baseline followed by every variant.
func (l *Lab) Traces() []Trace {
	return append([]Trace{l.Baseline}, l.Variants...)
}

// Trace returns the baseline for "" or "baseline", otherwise the named variant.
func (l *Lab) Trace(id string) (Trace, error) {
	if id == "" || id == "baseline" {
		return l.Baseline, nil
	}
	for _, v := range l.Variants {
		if v.ID == id {
			return v, nil
		}
	}
	return Trace{}, fmt.Errorf("unknown variant %q", id)
}

// Expand maps choice keys to the engine controls they enable in trace.
// Controls keep their declaration order in the control pool.
func (l *Lab) Expand(trace Trace, keys []string) ([]engine.Control, error) {
	wanted := make(map[string]struct{})
	for _, key := range keys {
		choice, ok := l.choice(key)
		if !ok {
			return nil, fmt.Errorf("unknown choice %q", key)
		}
		for _, id := range choice.Controls {
			wanted[id] = struct{}{}
		}
	}

	var selected []engine.Control
	for _, control := range trace.Controls.Controls {
		if c, ok := control.(engine.EventControl); ok {
			if _, enabled := wanted[c.ID]; enabled {
				selected = append(selected, control)
			}
		}
	}
	return selected, nil
}

func (l *Lab) choice(key string) (Choice, bool) {
	for _, c := range l.Manifest.Choices {
		if c.Key == key {
			return c, true
		}
	}
	return Choice{}, false
}

// Run replays trace with the given choices enabled.
func (l *Lab) Run(trace Trace, keys []string) (engine.RunResult, error) {
	controls, err := l.Expand(trace, keys)
	if err != nil {
		return engine.RunResult{}, err
	}
	return engine.Replay(trace.Scenario, controls...), nil
}

// Harm returns the adverse facts present in a run's terminal state, in
// manifest order.
func (l *Lab) Harm(result engine.RunResult) []string {
	terminal := make(map[string]struct{}, len(result.TerminalState))
	for _, fact := range result.TerminalState {
		terminal[fact] = struct{}{}
	}
	var harm []string
	for _, outcome := range l.Manifest.Adverse {
		if _, ok := terminal[outcome.Fact]; ok {
			harm = append(harm, outcome.Fact)
		}
	}
	return harm
}

// ChoiceKeys returns every choice key in manifest order.
func (l *Lab) ChoiceKeys() []string {
	keys := make([]string, len(l.Manifest.Choices))
	for i, c := range l.Manifest.Choices {
		keys[i] = c.Key
	}
	return keys
}

// LoadAll loads every lab.yaml under root, ordered by manifest order then ID.
func LoadAll(root string) ([]*Lab, error) {
	var paths []string
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && d.Name() == "lab.yaml" {
			paths = append(paths, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	labs := make([]*Lab, 0, len(paths))
	for _, path := range paths {
		l, err := Load(path)
		if err != nil {
			return nil, err
		}
		labs = append(labs, l)
	}
	sort.SliceStable(labs, func(i, j int) bool {
		a, b := labs[i].Manifest, labs[j].Manifest
		if a.Order != b.Order {
			return a.Order < b.Order
		}
		return a.ID < b.ID
	})
	return labs, nil
}
