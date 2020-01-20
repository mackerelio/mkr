package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/mackerelio/mackerel-client-go"
	"github.com/mackerelio/mkr/format"
	"github.com/mackerelio/mkr/logger"
	"github.com/mackerelio/mkr/mackerelclient"
	"github.com/urfave/cli"
	"github.com/yudai/gojsondiff"
	"github.com/yudai/gojsondiff/formatter"
)

var commandMonitors = cli.Command{
	Name:  "monitors",
	Usage: "Manipulate monitors",
	Description: `
    Manipulate monitor rules. With no subcommand specified, this will show all monitor rules.
    Requests APIs under "/api/v0/monitors". See https://mackerel.io/api-docs/entry/monitors .
`,
	Action: doMonitorsList,
	Subcommands: []cli.Command{
		{
			Name:      "pull",
			Usage:     "pull rules",
			ArgsUsage: "[--file-path | -F <file>] [--verbose | -v]",
			Description: `
    Pull monitor rules from Mackerel server and save them to a file. The file can be specified by filepath argument <file>. The default is 'monitors.json'.
`,
			Action: doMonitorsPull,
			Flags: []cli.Flag{
				cli.StringFlag{Name: "file-path, F", Value: "", Usage: "Filename to store monitor rule definitions. default: monitors.json"},
				cli.BoolFlag{Name: "verbose, v", Usage: "Verbose output mode"},
			},
		},
		{
			Name:  "diff",
			Usage: "diff rules",
			Description: `
    Show difference of monitor rules between Mackerel and a file. The file can be specified by filepath argument <file>. The default is 'monitors.json'.
`,
			ArgsUsage: "[--file-path | -F <file>]",
			Action:    doMonitorsDiff,
			Flags: []cli.Flag{
				cli.BoolFlag{Name: "exit-code, e", Usage: "Make mkr exit with code 1 if there are differences and 0 if there aren't. This is similar to diff(1)"},
				cli.StringFlag{Name: "file-path, F", Value: "", Usage: "Filename to store monitor rule definitions. default: monitors.json"},
				cli.BoolFlag{Name: "reverse", Usage: "The difference on the remote server is represented by plus and the difference on the local file is represented by minus"},
			},
		},
		{
			Name:      "push",
			Usage:     "push rules",
			ArgsUsage: "[--dry-run | -d] [--file-path | -F <file>] [--verbose | -v]",
			Description: `
    Push monitor rules stored in a file to Mackerel. The file can be specified by filepath argument <file>. The default is 'monitors.json'.
`,
			Action: doMonitorsPush,
			Flags: []cli.Flag{
				cli.StringFlag{Name: "file-path, F", Value: "", Usage: "Filename to store monitor rule definitions. default: monitors.json"},
				cli.BoolFlag{Name: "dry-run, d", Usage: "Show which apis are called, but not execute."},
				cli.BoolFlag{Name: "verbose, v", Usage: "Verbose output mode"},
			},
		},
	},
}

func monitorSaveRules(rules []mackerel.Monitor, optFilePath string) error {
	filePath := "monitors.json"
	if optFilePath != "" {
		filePath = optFilePath
	}
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	monitors := map[string]interface{}{"monitors": rules}
	data := format.JSONMarshalIndent(monitors, "", "    ") + "\n"

	_, err = file.WriteString(data)
	if err != nil {
		return err
	}
	return nil
}

func monitorLoadRules(optFilePath string) ([]mackerel.Monitor, error) {
	filePath := "monitors.json"
	if optFilePath != "" {
		filePath = optFilePath
	}

	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	return decodeMonitors(f)
}

// decodeMonitors decodes monitors JSON.
//
// There are almost same code in mackerel-client-go.
func decodeMonitors(r io.Reader) ([]mackerel.Monitor, error) {
	var data struct {
		Monitors []json.RawMessage `json:"monitors"`
	}
	if err := json.NewDecoder(r).Decode(&data); err != nil {
		return nil, err
	}
	ms := make([]mackerel.Monitor, 0, len(data.Monitors))
	for _, rawmes := range data.Monitors {
		m, err := decodeMonitor(rawmes)
		if err != nil {
			return nil, err
		}
		ms = append(ms, m)
	}
	return ms, nil
}

// decodeMonitor decodes json.RawMessage and returns monitor.
//
// There are almost same code in mackerel-client-go.
func decodeMonitor(mes json.RawMessage) (mackerel.Monitor, error) {
	var typeData struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(mes, &typeData); err != nil {
		return nil, err
	}
	var m mackerel.Monitor
	switch typeData.Type {
	case "connectivity":
		m = &mackerel.MonitorConnectivity{}
	case "host":
		m = &mackerel.MonitorHostMetric{}
	case "service":
		m = &mackerel.MonitorServiceMetric{}
	case "external":
		m = &mackerel.MonitorExternalHTTP{}
	case "expression":
		m = &mackerel.MonitorExpression{}
	case "anomalyDetection":
		m = &mackerel.MonitorAnomalyDetection{}
	}
	if err := json.Unmarshal(mes, m); err != nil {
		return nil, err
	}
	return m, nil
}

func doMonitorsList(c *cli.Context) error {
	monitors, err := mackerelclient.NewFromContext(c).FindMonitors()
	logger.DieIf(err)

	format.PrettyPrintJSON(os.Stdout, monitors)
	return nil
}

func doMonitorsPull(c *cli.Context) error {
	isVerbose := c.Bool("verbose")
	filePath := c.String("file-path")

	monitors, err := mackerelclient.NewFromContext(c).FindMonitors()
	logger.DieIf(err)

	monitorSaveRules(monitors, filePath)

	if isVerbose {
		format.PrettyPrintJSON(os.Stdout, monitors)
	}

	if filePath == "" {
		filePath = "monitors.json"
	}
	logger.Log("info", fmt.Sprintf("Monitor rules are saved to '%s' (%d rules).", filePath, len(monitors)))
	return nil
}

func stringifyMonitor(a mackerel.Monitor, prefix string) string {
	return prefix + format.JSONMarshalIndent(a, prefix, "  ") + ","
}

// diffMonitor returns JSON diff between monitors.
// In order to use `mkr monitors` without pull and to manage monitors by name
// only, it skips top level "id" field
func diffMonitor(a mackerel.Monitor, b mackerel.Monitor) string {
	as := filterIDLine(format.JSONMarshalIndent(a, " ", "  "))
	bs := filterIDLine(format.JSONMarshalIndent(b, " ", "  "))
	diff, err := gojsondiff.New().Compare([]byte(as), []byte(bs))
	if err != nil || !diff.Modified() {
		return ""
	}
	var left map[string]interface{}
	json.Unmarshal([]byte(as), &left)
	result, err := formatter.NewAsciiFormatter(left, formatter.AsciiFormatterDefaultConfig).Format(diff)
	if err != nil {
		return ""
	}
	return strings.TrimRight(result, "\n") + ","
}

func filterIDLine(s string) string {
	lines := strings.Split(s, "\n")
	filtered := make([]string, 0, len(lines))
	for _, l := range lines {
		if strings.HasPrefix(l, `   "id":`) {
			continue
		}
		filtered = append(filtered, l)
	}
	return strings.Join(filtered, "\n")
}

func isSameMonitor(a mackerel.Monitor, b mackerel.Monitor, flagNameUniqueness bool) (string, bool) {
	if a == nil || b == nil {
		return "", false
	}
	if reflect.DeepEqual(a, b) {
		return "", true
	}
	aID := a.MonitorID()
	bID := b.MonitorID()
	if aID == bID || (flagNameUniqueness == true && bID == "" && a.MonitorName() == b.MonitorName()) {
		diff := diffMonitor(a, b)
		if diff != "" {
			return diff, false
		}
		return "", true
	}
	return "", false
}

func validateRuleAnomalyDetectionScopes(v reflect.Value, label string) (bool, error) {
	f := "Scopes"
	vf := v.FieldByName("Scopes")
	if !vf.IsValid() {
		return false, fmt.Errorf("Monitor '%s' should have '%s': %s", label, f, v.FieldByName(f).Interface())
	}
	scopes, ok := vf.Interface().([]string)
	if !ok {
		return false, fmt.Errorf("Monitor '%s' has invalid '%s': %s", label, f, v.FieldByName(f).Interface())
	} else if len(scopes) == 0 {
		return false, fmt.Errorf("Monitor '%s' has empty '%s'", label, f)
	}
	return true, nil
}

func validateRules(monitors []mackerel.Monitor, label string) (bool, error) {

	flagNameUniqueness := true
	// check each monitor
	for _, monitor := range monitors {
		v := reflect.ValueOf(monitor).Elem()
		for _, f := range []string{"Type"} {
			vf := v.FieldByName(f)
			if !vf.IsValid() || (vf.Type().String() == "string" && vf.Interface() == "") {
				return false, fmt.Errorf("Monitor '%s' should have '%s': %s", label, f, v.FieldByName(f).Interface())
			}
		}
		switch m := monitor.(type) {
		case *mackerel.MonitorHostMetric, *mackerel.MonitorServiceMetric:
			for _, f := range []string{"Name", "Metric"} {
				vf := v.FieldByName(f)
				if !vf.IsValid() || (vf.Type().String() == "string" && vf.Interface() == "") {
					return false, fmt.Errorf("Monitor '%s' should have '%s': %s", label, f, v.FieldByName(f).Interface())
				}
			}
		case *mackerel.MonitorExternalHTTP:
			for _, f := range []string{"Name", "URL"} {
				vf := v.FieldByName(f)
				if !vf.IsValid() || (vf.Type().String() == "string" && vf.Interface() == "") {
					return false, fmt.Errorf("Monitor '%s' should have '%s': %s", label, f, v.FieldByName(f).Interface())
				}
			}
		case *mackerel.MonitorExpression:
			for _, f := range []string{"Name", "Expression"} {
				vf := v.FieldByName(f)
				if !vf.IsValid() || (vf.Type().String() == "string" && vf.Interface() == "") {
					return false, fmt.Errorf("Monitor '%s' should have '%s': %s", label, f, v.FieldByName(f).Interface())
				}
			}
		case *mackerel.MonitorAnomalyDetection:
			for _, f := range []string{"Name"} {
				vf := v.FieldByName(f)
				if !vf.IsValid() || (vf.Type().String() == "string" && vf.Interface() == "") {
					return false, fmt.Errorf("Monitor '%s' should have '%s': %s", label, f, v.FieldByName(f).Interface())
				}
			}
			return validateRuleAnomalyDetectionScopes(v, label)
		case *mackerel.MonitorConnectivity:
		default:
			return false, fmt.Errorf("Unknown type is found: %s", m.MonitorType())
		}
	}

	// check name uniqueness
	names := map[string]bool{}
	for _, m := range monitors {
		name := m.MonitorName()
		if names[name] {
			logger.Log("Warning: ", fmt.Sprintf("Names of %s are not unique.", label))
			flagNameUniqueness = false
		}
		names[name] = true
	}
	return flagNameUniqueness, nil
}

type monitorDiffPair struct {
	remote mackerel.Monitor
	local  mackerel.Monitor
}

type monitorDiff struct {
	onlyRemote []mackerel.Monitor
	onlyLocal  []mackerel.Monitor
	diff       []*monitorDiffPair
}

func checkMonitorsDiff(c *cli.Context) monitorDiff {
	filePath := c.String("file-path")

	var monitorDiff monitorDiff

	monitorsRemote, err := mackerelclient.NewFromContext(c).FindMonitors()
	logger.DieIf(err)
	flagNameUniquenessRemote, err := validateRules(monitorsRemote, "remote rules")
	logger.DieIf(err)

	monitorsLocal, err := monitorLoadRules(filePath)
	logger.DieIf(err)
	flagNameUniquenessLocal, err := validateRules(monitorsLocal, "local rules")
	logger.DieIf(err)

	flagNameUniqueness := flagNameUniquenessLocal && flagNameUniquenessRemote

	for _, remote := range monitorsRemote {
		found := false
		for i, local := range monitorsLocal {
			diff, isSame := isSameMonitor(remote, local, flagNameUniqueness)
			if isSame || diff != "" {
				monitorsLocal[i] = nil
				found = true
				if diff != "" {
					monitorDiff.diff = append(monitorDiff.diff, &monitorDiffPair{remote, local})
				}
				break
			}
		}
		if found == false {
			monitorDiff.onlyRemote = append(monitorDiff.onlyRemote, remote)
		}
	}
	for _, local := range monitorsLocal {
		if local != nil {
			monitorDiff.onlyLocal = append(monitorDiff.onlyLocal, local)
		}
	}

	return monitorDiff
}

func doMonitorsDiff(c *cli.Context) error {
	monitorDiff := checkMonitorsDiff(c)
	isExitCode := c.Bool("exit-code")
	isReverse := c.Bool("reverse")

	var diffs []string
	for _, d := range monitorDiff.diff {
		var diff string
		if isReverse {
			diff = diffMonitor(d.local, d.remote)
		} else {
			diff = diffMonitor(d.remote, d.local)
		}
		diffs = append(diffs, diff)
	}

	var monitorOnlyFrom []mackerel.Monitor
	var monitorOnlyTo []mackerel.Monitor
	if isReverse {
		monitorOnlyFrom = monitorDiff.onlyLocal
		monitorOnlyTo = monitorDiff.onlyRemote
	} else {
		monitorOnlyFrom = monitorDiff.onlyRemote
		monitorOnlyTo = monitorDiff.onlyLocal
	}

	fmt.Printf("Summary: %d modify, %d append, %d remove\n\n", len(monitorDiff.diff), len(monitorOnlyTo), len(monitorOnlyFrom))
	noDiff := true
	for _, diff := range diffs {
		fmt.Println(diff)
		noDiff = false
	}
	for _, m := range monitorOnlyFrom {
		fmt.Println(stringifyMonitor(m, "-"))
		noDiff = false
	}
	for _, m := range monitorOnlyTo {
		fmt.Println(stringifyMonitor(m, "+"))
		noDiff = false
	}
	if isExitCode == true && noDiff == false {
		os.Exit(1)
	}
	return nil
}

func doMonitorsPush(c *cli.Context) error {
	monitorDiff := checkMonitorsDiff(c)
	isDryRun := c.Bool("dry-run")
	isVerbose := c.Bool("verbose")

	client := mackerelclient.NewFromContext(c)
	if isVerbose {
		client.Verbose = true
	}

	for _, m := range monitorDiff.onlyLocal {
		logger.Log("info", "Create a new rule.")
		fmt.Println(stringifyMonitor(m, ""))
		if !isDryRun {
			_, err := client.CreateMonitor(m)
			logger.DieIf(err)
		}
	}
	for _, m := range monitorDiff.onlyRemote {
		logger.Log("info", "Delete a rule.")
		fmt.Println(stringifyMonitor(m, ""))
		if !isDryRun {
			_, err := client.DeleteMonitor(m.MonitorID())
			logger.DieIf(err)
		}
	}
	for _, d := range monitorDiff.diff {
		logger.Log("info", "Update a rule.")
		fmt.Println(stringifyMonitor(d.local, ""))
		if !isDryRun {
			_, err := client.UpdateMonitor(d.remote.MonitorID(), d.local)
			logger.DieIf(err)
		}
	}
	return nil
}
