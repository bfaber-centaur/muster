package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/bubskee/muster/config"
	"github.com/bubskee/muster/engine"
	"github.com/bubskee/muster/experiment"
)

type replayOutput struct {
	ControlSetID string           `json:"control_set_id"`
	Result       engine.RunResult `json:"result"`
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}

	switch os.Args[1] {
	case "replay":
		if err := replay(os.Args[2:]); err != nil {
			fmt.Fprintln(os.Stderr, "muster:", err)
			os.Exit(1)
		}
	case "experiment":
		if err := runExperiment(os.Args[2:]); err != nil {
			fmt.Fprintln(os.Stderr, "muster:", err)
			os.Exit(1)
		}
	case "lab", "search", "robust", "bundle", "survey":
		commands := map[string]func([]string) error{
			"lab":    labRun,
			"search": searchRun,
			"robust": robustRun,
			"bundle": bundleRun,
			"survey": surveyRun,
		}
		if err := commands[os.Args[1]](os.Args[2:]); err != nil {
			fmt.Fprintln(os.Stderr, "muster:", err)
			os.Exit(1)
		}
	case "help", "-h", "--help":
		usage()
	default:
		fmt.Fprintf(os.Stderr, "muster: unknown command %q\n\n", os.Args[1])
		usage()
		os.Exit(2)
	}
}

func replay(args []string) error {
	flags := flag.NewFlagSet("replay", flag.ContinueOnError)
	flags.SetOutput(os.Stderr)

	scenarioPath := flags.String("scenario", "", "scenario YAML file")
	controlsPath := flags.String("controls", "", "control-set YAML file")
	onlyControls := flags.String(
		"only-controls",
		"",
		"comma-separated control IDs to enable; empty means all controls",
	)
	asJSON := flags.Bool("json", false, "emit replay result as JSON")

	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 0 {
		return fmt.Errorf("unexpected positional arguments: %v", flags.Args())
	}
	if *scenarioPath == "" {
		return fmt.Errorf("--scenario is required")
	}
	if *controlsPath == "" {
		return fmt.Errorf("--controls is required")
	}

	scenario, err := config.LoadScenario(*scenarioPath)
	if err != nil {
		return err
	}
	controlSet, err := config.LoadControlSet(*controlsPath)
	if err != nil {
		return err
	}

	selected, err := selectControls(controlSet.Controls, *onlyControls)
	if err != nil {
		return err
	}

	result := engine.Replay(scenario, selected...)

	if *asJSON {
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(replayOutput{
			ControlSetID: controlSet.ID,
			Result:       result,
		})
	}

	fmt.Printf("control_set: %s\n", controlSet.ID)
	if strings.TrimSpace(*onlyControls) != "" {
		fmt.Printf("enabled_controls: %s\n", strings.Join(controlIDs(selected), ","))
	}
	fmt.Print(engine.Summarize(result).String())
	return nil
}

func selectControls(controls []engine.Control, raw string) ([]engine.Control, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return append([]engine.Control(nil), controls...), nil
	}

	wanted := make(map[string]struct{})
	for _, part := range strings.Split(raw, ",") {
		id := strings.TrimSpace(part)
		if id == "" {
			return nil, fmt.Errorf("--only-controls contains an empty control id")
		}
		if _, exists := wanted[id]; exists {
			return nil, fmt.Errorf("duplicate control id %q in --only-controls", id)
		}
		wanted[id] = struct{}{}
	}

	selected := make([]engine.Control, 0, len(wanted))
	found := make(map[string]struct{}, len(wanted))

	for _, control := range controls {
		id, ok := controlID(control)
		if !ok {
			return nil, fmt.Errorf("control %T does not expose an EventControl id", control)
		}
		if _, enabled := wanted[id]; enabled {
			selected = append(selected, control)
			found[id] = struct{}{}
		}
	}

	var missing []string
	for id := range wanted {
		if _, ok := found[id]; !ok {
			missing = append(missing, id)
		}
	}
	if len(missing) > 0 {
		sort.Strings(missing)
		return nil, fmt.Errorf("unknown control id(s): %s", strings.Join(missing, ", "))
	}

	return selected, nil
}

func controlID(control engine.Control) (string, bool) {
	switch c := control.(type) {
	case engine.EventControl:
		return c.ID, c.ID != ""
	case *engine.EventControl:
		return c.ID, c.ID != ""
	default:
		return "", false
	}
}

func controlIDs(controls []engine.Control) []string {
	ids := make([]string, 0, len(controls))
	for _, control := range controls {
		if id, ok := controlID(control); ok {
			ids = append(ids, id)
		}
	}
	return ids
}

func runExperiment(args []string) error {
	flags := flag.NewFlagSet("experiment", flag.ContinueOnError)
	flags.SetOutput(os.Stderr)

	path := flags.String("file", "", "experiment YAML file")
	format := flags.String("format", "markdown", "output format: markdown, json, or csv")

	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 0 {
		return fmt.Errorf("unexpected positional arguments: %v", flags.Args())
	}
	if *path == "" {
		return fmt.Errorf("--file is required")
	}

	result, err := experiment.ExecuteFile(*path)
	if err != nil {
		return err
	}

	switch *format {
	case "markdown", "md":
		return experiment.WriteMarkdown(os.Stdout, result)
	case "csv":
		return experiment.WriteCSV(os.Stdout, result)
	case "json":
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(result)
	default:
		return fmt.Errorf("unknown --format %q; want markdown, json, or csv", *format)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, `usage:
  muster replay --scenario SCENARIO.yaml --controls CONTROLS.yaml [--only-controls ID,ID,...] [--json]
  muster experiment --file EXPERIMENT.yaml [--format markdown|json|csv]
  muster lab --file LAB.yaml [--choices V1,R1 | --preset LABEL] [--variant ID] [--json]
  muster search --file LAB.yaml [--variant ID] [--tradeoffs] [--json]
  muster robust --file LAB.yaml [--json]
  muster survey [--root examples]
  muster bundle [--root examples] [--out labs.json] [--golden golden.json]

examples:
  muster replay --scenario examples/gate2-between.yaml --controls examples/gate2-broad-correlated.yaml
  muster replay --scenario examples/notebook/toy-jackpot.yaml --controls examples/notebook/toy-controls.yaml --only-controls worker-credential-theft,worker-k8s-discovery,revoke-cluster-credential --json
  muster experiment --file experiments/hf-july-2026.yaml
  muster search --file examples/toy-jackpot/lab.yaml
  muster robust --file examples/hf-july-2026/lab.yaml`)
}
