package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/bubskee/muster/lab"
)

func labRun(args []string) error {
	flags := flag.NewFlagSet("lab", flag.ContinueOnError)
	flags.SetOutput(os.Stderr)
	path := flags.String("file", "", "lab manifest YAML")
	choices := flags.String("choices", "", "comma-separated choice keys, e.g. V1,R1")
	preset := flags.String("preset", "", "preset label to run instead of --choices")
	variant := flags.String("variant", "", "variant trace id (default: baseline)")
	asJSON := flags.Bool("json", false, "emit replay result as JSON")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if *path == "" {
		return fmt.Errorf("--file is required")
	}

	l, err := lab.Load(*path)
	if err != nil {
		return err
	}
	trace, err := l.Trace(*variant)
	if err != nil {
		return err
	}

	keys := splitList(*choices)
	if *preset != "" {
		found := false
		for _, p := range l.Manifest.Presets {
			if p.Label == *preset {
				keys, found = p.Choices, true
			}
		}
		if !found {
			return fmt.Errorf("unknown preset %q", *preset)
		}
	}

	result, err := l.Run(trace, keys)
	if err != nil {
		return err
	}
	harm := l.Harm(result)

	if *asJSON {
		return writeJSON(os.Stdout, map[string]any{
			"lab_id":  l.Manifest.ID,
			"trace":   trace.ID,
			"choices": keys,
			"harm":    harm,
			"result":  result,
		})
	}

	fmt.Printf("lab: %s  trace: %s  choices: [%s]\n", l.Manifest.ID, trace.ID, strings.Join(keys, ","))
	for _, entry := range result.Trace {
		label := l.Manifest.Events[entry.EventID]
		if label == "" {
			label = entry.EventID
		}
		marks := ""
		if len(entry.ObservedBy) > 0 {
			marks += " 🔭" + strings.Join(entry.ObservedBy, ",")
		}
		if len(entry.BlockedBy) > 0 {
			marks += " 🛡️" + strings.Join(entry.BlockedBy, ",")
		}
		if len(entry.RespondedBy) > 0 {
			marks += " 🔧" + strings.Join(entry.RespondedBy, ",")
		}
		fmt.Printf("  %-8s %s%s\n", entry.Status, label, marks)
	}
	fmt.Printf("adverse: %s\n", factLabels(l, harm))
	return nil
}

func searchRun(args []string) error {
	flags := flag.NewFlagSet("search", flag.ContinueOnError)
	flags.SetOutput(os.Stderr)
	path := flags.String("file", "", "lab manifest YAML")
	variant := flags.String("variant", "", "variant trace id (default: baseline)")
	asJSON := flags.Bool("json", false, "emit findings as JSON")
	tradeoffs := flags.Bool("tradeoffs", false, "also list tradeoffs (one outcome traded for another)")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if *path == "" {
		return fmt.Errorf("--file is required")
	}

	l, err := lab.Load(*path)
	if err != nil {
		return err
	}
	trace, err := l.Trace(*variant)
	if err != nil {
		return err
	}
	result, err := l.Search(trace)
	if err != nil {
		return err
	}
	if *asJSON {
		return writeJSON(os.Stdout, result)
	}

	fmt.Printf("# Search: %s (%s)\n\n", l.Manifest.ID, trace.ID)
	fmt.Printf("%d choices, %d Watchlines replayed.\n\n", result.Choices, result.Watchline)

	worse, traded := 0, 0
	for _, f := range result.Findings {
		if f.Kind == lab.FindingWorse {
			worse++
		} else {
			traded++
		}
	}
	fmt.Printf("Strictly worse after adding one defense: %d minimal case(s)\n", worse)
	fmt.Printf("Tradeoffs (fixes one outcome, opens another): %d\n\n", traded)

	for _, f := range result.Findings {
		if f.Kind == lab.FindingTradeoff && !*tradeoffs {
			continue
		}
		fmt.Printf("- [%s] %s + %s\n", f.Kind, bracket(choiceLabels(l, f.Context)), choiceLabel(l, f.Added))
		fmt.Printf("    before: %s\n", factLabels(l, f.Before))
		fmt.Printf("    after:  %s\n", factLabels(l, f.After))
		if f.Divergence != "" {
			fmt.Printf("    first divergence: %s\n", eventLabel(l, f.Divergence))
		}
	}
	return nil
}

func robustRun(args []string) error {
	flags := flag.NewFlagSet("robust", flag.ContinueOnError)
	flags.SetOutput(os.Stderr)
	path := flags.String("file", "", "lab manifest YAML")
	asJSON := flags.Bool("json", false, "emit robustness result as JSON")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if *path == "" {
		return fmt.Errorf("--file is required")
	}
	l, err := lab.Load(*path)
	if err != nil {
		return err
	}
	result, err := l.Robust()
	if err != nil {
		return err
	}
	if *asJSON {
		return writeJSON(os.Stdout, result)
	}
	return writeRobustMarkdown(os.Stdout, l, result)
}

func writeRobustMarkdown(w io.Writer, l *lab.Lab, r lab.RobustResult) error {
	fmt.Fprintf(w, "# Robustness: %s\n\n", l.Manifest.Title)
	fmt.Fprintln(w, "Adverse outcomes reached by each Watchline, as recorded and under each alternative trace.")
	fmt.Fprintln(w, "↑ = worse than as recorded, ↓ = better, ≠ = different.")
	fmt.Fprintln(w)

	header := []string{"Watchline"}
	for _, t := range r.Traces {
		header = append(header, t.Label)
	}
	fmt.Fprintf(w, "| %s |\n", strings.Join(header, " | "))
	fmt.Fprintf(w, "|%s\n", strings.Repeat(" --- |", len(header)))

	short := shortOutcomes(l)
	for _, row := range r.Rows {
		cells := []string{row.Label}
		for _, c := range row.Cells {
			text := "none"
			if len(c.Harm) > 0 {
				parts := make([]string, len(c.Harm))
				for i, h := range c.Harm {
					parts[i] = short[h]
				}
				text = strings.Join(parts, ", ")
			}
			switch c.Change {
			case "worse":
				text += " ↑"
			case "better":
				text += " ↓"
			case "different":
				text += " ≠"
			}
			cells = append(cells, text)
		}
		fmt.Fprintf(w, "| %s |\n", strings.Join(cells, " | "))
	}

	fmt.Fprintln(w, "\nOutcomes:")
	for _, o := range l.Manifest.Adverse {
		fmt.Fprintf(w, "- %s = %s\n", short[o.Fact], o.Label)
	}

	hasWorse := false
	for _, fc := range r.Findings {
		if fc.Finding.Kind == lab.FindingWorse {
			hasWorse = true
			break
		}
	}
	if hasWorse {
		fmt.Fprintln(w, "\n## Do the baseline non-monotone findings survive?")
		fmt.Fprintln(w)
		for _, fc := range r.Findings {
			f := fc.Finding
			if f.Kind != lab.FindingWorse {
				continue
			}
			fmt.Fprintf(w, "- %s + **%s**:", bracket(choiceLabels(l, f.Context)), choiceLabel(l, f.Added))
			for _, t := range l.Variants {
				mark := "✗"
				if fc.Holds[t.ID] {
					mark = "✓"
				}
				fmt.Fprintf(w, " %s %s;", mark, t.Label)
			}
			fmt.Fprintln(w)
		}
	}

	for _, t := range l.Variants {
		emergent := r.Emergent[t.ID]
		count := 0
		for _, f := range emergent {
			if f.Kind == lab.FindingWorse {
				count++
			}
		}
		if count == 0 {
			continue
		}
		fmt.Fprintf(w, "\n## New non-monotone cases under “%s”\n\n", t.Label)
		for _, f := range emergent {
			if f.Kind != lab.FindingWorse {
				continue
			}
			fmt.Fprintf(w, "- %s + **%s** → also %s\n", bracket(choiceLabels(l, f.Context)), choiceLabel(l, f.Added), factLabels(l, f.NewHarm))
		}
	}
	return nil
}

func bundleRun(args []string) error {
	flags := flag.NewFlagSet("bundle", flag.ContinueOnError)
	flags.SetOutput(os.Stderr)
	root := flags.String("root", "examples", "directory searched for lab.yaml files")
	out := flags.String("out", "", "output file (default stdout)")
	golden := flags.String("golden", "", "also write every replay of every lab here, for checking the browser engine")
	if err := flags.Parse(args); err != nil {
		return err
	}
	labs, err := lab.LoadAll(*root)
	if err != nil {
		return err
	}
	if *golden != "" {
		runs, err := lab.Golden(labs)
		if err != nil {
			return err
		}
		if err := encodeFile(*golden, runs); err != nil {
			return err
		}
	}
	bundle := lab.NewBundle(labs)
	if *out != "" {
		if err := encodeFile(*out, bundle); err != nil {
			return err
		}
	} else if err := json.NewEncoder(os.Stdout).Encode(bundle); err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "bundled %d labs\n", len(labs))
	return nil
}

// encodeFile writes v as JSON to path, reporting a failed Close so a short
// write is not mistaken for a complete file.
func encodeFile(path string, v any) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	if err := json.NewEncoder(f).Encode(v); err != nil {
		f.Close()
		return err
	}
	return f.Close()
}

func writeJSON(w io.Writer, v any) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(v)
}

func splitList(raw string) []string {
	var out []string
	for _, part := range strings.Split(raw, ",") {
		if s := strings.TrimSpace(part); s != "" {
			out = append(out, s)
		}
	}
	return out
}

func choiceLabel(l *lab.Lab, key string) string {
	for _, c := range l.Manifest.Choices {
		if c.Key == key {
			return fmt.Sprintf("%s %s", key, c.Label)
		}
	}
	return key
}

func choiceLabels(l *lab.Lab, keys []string) []string {
	out := make([]string, len(keys))
	for i, k := range keys {
		out[i] = choiceLabel(l, k)
	}
	return out
}

func bracket(items []string) string {
	if len(items) == 0 {
		return "[no defense]"
	}
	return "[" + strings.Join(items, "; ") + "]"
}

func eventLabel(l *lab.Lab, id string) string {
	if label := l.Manifest.Events[id]; label != "" {
		return label
	}
	return id
}

func factLabels(l *lab.Lab, facts []string) string {
	if len(facts) == 0 {
		return "none"
	}
	labels := map[string]string{}
	for _, o := range l.Manifest.Adverse {
		labels[o.Fact] = o.Label
	}
	out := make([]string, len(facts))
	for i, f := range facts {
		out[i] = labels[f]
		if out[i] == "" {
			out[i] = f
		}
	}
	return strings.Join(out, "; ")
}

// shortOutcomes assigns O1..On to adverse outcomes for compact tables.
func shortOutcomes(l *lab.Lab) map[string]string {
	short := map[string]string{}
	for i, o := range l.Manifest.Adverse {
		short[o.Fact] = fmt.Sprintf("O%d", i+1)
	}
	return short
}

func surveyRun(args []string) error {
	flags := flag.NewFlagSet("survey", flag.ContinueOnError)
	flags.SetOutput(os.Stderr)
	root := flags.String("root", "examples", "directory searched for lab.yaml files")
	if err := flags.Parse(args); err != nil {
		return err
	}
	labs, err := lab.LoadAll(*root)
	if err != nil {
		return err
	}

	fmt.Println("# Survey of every lab")
	fmt.Println()
	fmt.Println("Backfires: minimal cases where adding one defense strictly grows the set of")
	fmt.Println("adverse outcomes, on the recorded trace / on any alternative trace.")
	fmt.Println("Fragile presets: presets whose outcome changes in at least one alternative trace.")
	fmt.Println()
	fmt.Println("| Lab | Domain | Kind | Watchlines | Backfires (recorded / any variant) | Tradeoffs (recorded) | Fragile presets |")
	fmt.Println("| --- | --- | --- | ---: | ---: | ---: | --- |")

	for _, l := range labs {
		searches, err := l.SearchAll()
		if err != nil {
			return err
		}
		base := searches[0]
		worse, traded := 0, 0
		for _, f := range base.Findings {
			if f.Kind == lab.FindingWorse {
				worse++
			} else {
				traded++
			}
		}

		robust, err := l.RobustWith(searches)
		if err != nil {
			return err
		}
		emergent := 0
		for _, fs := range robust.Emergent {
			for _, f := range fs {
				if f.Kind == lab.FindingWorse {
					emergent++
				}
			}
		}
		fragile, presets := 0, 0
		for _, row := range robust.Rows[1:] {
			presets++
			for _, c := range row.Cells {
				if c.Change != "same" {
					fragile++
					break
				}
			}
		}

		fragileText := "—"
		if len(l.Variants) > 0 {
			fragileText = fmt.Sprintf("%d of %d", fragile, presets)
		}
		fmt.Printf("| %s | %s | %s | %d | %d / %d | %d | %s |\n",
			l.Manifest.Title, l.Manifest.Domain, l.Manifest.Kind, base.Watchline,
			worse, worse+emergent, traded, fragileText)
	}
	return nil
}
