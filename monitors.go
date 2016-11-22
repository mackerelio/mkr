package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strings"

	mkr "github.com/mackerelio/mackerel-client-go"
	"github.com/mackerelio/mkr/logger"
	"github.com/yudai/gojsondiff"
	"github.com/yudai/gojsondiff/formatter"
	"gopkg.in/urfave/cli.v1"
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
			Name:        "pull",
			Usage:       "pull rules",
			Description: "Pull monitor rules from Mackerel server and save them to a file. The file can be specified by filepath argument <file>. The default is 'monitors.json'.",
			Action:      doMonitorsPull,
			Flags: []cli.Flag{
				cli.StringFlag{Name: "file-path, F", Value: "", Usage: "Filename to store monitor rule definitions. default: monitors.json"},
				cli.BoolFlag{Name: "verbose, v", Usage: "Verbose output mode"},
			},
		},
		{
			Name:        "diff",
			Usage:       "diff rules",
			Description: "Show difference of monitor rules between Mackerel and a file. The file can be specified by filepath argument <file>. The default is 'monitors.json'.",
			Action:      doMonitorsDiff,
			Flags: []cli.Flag{
				cli.BoolFlag{Name: "exit-code, e", Usage: "Make mkr exit with code 1 if there are differences and 0 if there aren't. This is similar to diff(1)"},
				cli.StringFlag{Name: "file-path, F", Value: "", Usage: "Filename to store monitor rule definitions. default: monitors.json"},
			},
		},
		{
			Name:        "push",
			Usage:       "push rules",
			Description: "Push monitor rules stored in a file to Mackerel. The file can be specified by filepath argument <file>. The default is 'monitors.json'.",
			Action:      doMonitorsPush,
			Flags: []cli.Flag{
				cli.StringFlag{Name: "file-path, F", Value: "", Usage: "Filename to store monitor rule definitions. default: monitors.json"},
				cli.BoolFlag{Name: "dry-run, d", Usage: "Show which apis are called, but not execute."},
				cli.BoolFlag{Name: "verbose, v", Usage: "Verbose output mode"},
			},
		},
	},
}

func monitorSaveRules(rules []*(mkr.Monitor), optFilePath string) error {
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
	data := JSONMarshalIndent(monitors, "", "    ") + "\n"

	_, err = file.WriteString(data)
	if err != nil {
		return err
	}
	return nil
}

func monitorLoadRules(optFilePath string) ([]*(mkr.Monitor), error) {
	filePath := "monitors.json"
	if optFilePath != "" {
		filePath = optFilePath
	}

	buff, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var data struct {
		Monitors []*(mkr.Monitor) `json:"monitors"`
	}

	err = json.Unmarshal(buff, &data)
	if err != nil {
		return nil, err
	}
	return data.Monitors, nil
}

func doMonitorsList(c *cli.Context) error {
	monitors, err := newMackerelFromContext(c).FindMonitors()
	logger.DieIf(err)

	PrettyPrintJSON(monitors)
	return nil
}

func doMonitorsPull(c *cli.Context) error {
	isVerbose := c.Bool("verbose")
	filePath := c.String("file-path")

	monitors, err := newMackerelFromContext(c).FindMonitors()
	logger.DieIf(err)

	monitorSaveRules(monitors, filePath)

	if isVerbose {
		PrettyPrintJSON(monitors)
	}

	if filePath == "" {
		filePath = "monitors.json"
	}
	logger.Log("info", fmt.Sprintf("Monitor rules are saved to '%s' (%d rules).", filePath, len(monitors)))
	return nil
}

func isEmpty(a interface{}) bool {
	switch a.(type) {
	case bool:
		if reflect.ValueOf(a).Bool() == false {
			return true
		}
	case int, int8, int16, int32, int64:
		if reflect.ValueOf(a).Int() == 0 {
			return true
		}
	case uint, uint8, uint16, uint32, uint64:
		if reflect.ValueOf(a).Uint() == 0 {
			return true
		}
	case float32, float64:
		if reflect.ValueOf(a).Float() == 0.0 {
			return true
		}
	case string:
		if reflect.ValueOf(a).String() == "" {
			return true
		}
	}
	return false
}

func appendDiff(src []string, name string, a interface{}, b interface{}) []string {
	diff := src
	aType := reflect.TypeOf(a).String()
	format := "\"%s\""
	isAEmpty := isEmpty(a)
	isBEmpty := isEmpty(b)
	switch aType {
	case "bool":
		format = "%t"
	case "uint64":
		format = "%d"
	case "float64":
		format = "%f"
	}
	if isAEmpty == false || isBEmpty == false {
		if a != b {
			diff = append(diff, fmt.Sprintf("-   \"%s\": "+format+",", name, a))
			diff = append(diff, fmt.Sprintf("+   \"%s\": "+format+",", name, b))
		} else {
			diff = append(diff, fmt.Sprintf("    \"%s\": "+format+",", name, a))
		}
	}
	return diff
}

func stringifyMonitor(a *mkr.Monitor, prefix string) string {
	return prefix + JSONMarshalIndent(a, prefix, "  ") + ","
}

// diffMonitor returns JSON diff between monitors.
// In order to use `mkr monitors` without pull and to manage monitors by name
// only, it skips top level "id" field
func diffMonitor(a *mkr.Monitor, b *mkr.Monitor) string {
	as := filterIDLine(JSONMarshalIndent(a, " ", "  "))
	bs := filterIDLine(JSONMarshalIndent(b, " ", "  "))
	diff, err := gojsondiff.New().Compare([]byte(as), []byte(bs))
	if err != nil || !diff.Modified() {
		return ""
	}
	var left map[string]interface{}
	json.Unmarshal([]byte(as), &left)
	result, err := formatter.NewAsciiFormatter(left).Format(diff)
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

func isSameMonitor(a *mkr.Monitor, b *mkr.Monitor, flagNameUniqueness bool) (string, bool) {
	if a == nil || b == nil {
		return "", false
	}
	if reflect.DeepEqual(*a, *b) {
		return "", true
	}
	if a.ID == b.ID || (flagNameUniqueness == true && b.ID == "" && a.Name == b.Name) {
		diff := diffMonitor(a, b)
		if diff != "" {
			return diff, false
		}
		return "", true
	}
	return "", false
}

func validateRules(monitors []*(mkr.Monitor), label string) (bool, error) {

	flagNameUniqueness := true
	// check each monitor
	for _, m := range monitors {
		v := reflect.ValueOf(m).Elem()
		for _, f := range []string{"Type"} {
			vf := v.FieldByName(f)
			if !vf.IsValid() || (vf.Type().String() == "string" && vf.Interface() == "") {
				return false, fmt.Errorf("Monitor '%s' should have '%s': %s", label, f, v.FieldByName(f).Interface())
			}
		}
		switch m.Type {
		case "host", "service":
			for _, f := range []string{"Name", "Metric"} {
				vf := v.FieldByName(f)
				if !vf.IsValid() || (vf.Type().String() == "string" && vf.Interface() == "") {
					return false, fmt.Errorf("Monitor '%s' should have '%s': %s", label, f, v.FieldByName(f).Interface())
				}
			}
		case "external":
			for _, f := range []string{"Name", "URL"} {
				vf := v.FieldByName(f)
				if !vf.IsValid() || (vf.Type().String() == "string" && vf.Interface() == "") {
					return false, fmt.Errorf("Monitor '%s' should have '%s': %s", label, f, v.FieldByName(f).Interface())
				}
			}
		case "expression":
			for _, f := range []string{"Name", "Expression"} {
				vf := v.FieldByName(f)
				if !vf.IsValid() || (vf.Type().String() == "string" && vf.Interface() == "") {
					return false, fmt.Errorf("Monitor '%s' should have '%s': %s", label, f, v.FieldByName(f).Interface())
				}
			}
		case "connectivity":
		default:
			return false, fmt.Errorf("Unknown type is found: %s", m.Type)
		}
	}

	// check name uniqueness
	names := map[string]bool{}
	for _, m := range monitors {
		if names[m.Name] {
			logger.Log("Warning: ", fmt.Sprintf("Names of %s are not unique.", label))
			flagNameUniqueness = false
		}
		names[m.Name] = true
	}
	return flagNameUniqueness, nil
}

type monitorDiffPair struct {
	remote *mkr.Monitor
	local  *mkr.Monitor
}

type monitorDiff struct {
	onlyRemote []*(mkr.Monitor)
	onlyLocal  []*(mkr.Monitor)
	diff       []*monitorDiffPair
}

func checkMonitorsDiff(c *cli.Context) monitorDiff {
	filePath := c.String("file-path")

	var monitorDiff monitorDiff

	monitorsRemote, err := newMackerelFromContext(c).FindMonitors()
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

	var diffs []string
	for _, d := range monitorDiff.diff {
		diffs = append(diffs, diffMonitor(d.remote, d.local))
	}

	fmt.Printf("Summary: %d modify, %d append, %d remove\n\n", len(monitorDiff.diff), len(monitorDiff.onlyLocal), len(monitorDiff.onlyRemote))
	noDiff := true
	for _, diff := range diffs {
		fmt.Println(diff)
		noDiff = false
	}
	for _, m := range monitorDiff.onlyRemote {
		fmt.Println(stringifyMonitor(m, "-"))
		noDiff = false
	}
	for _, m := range monitorDiff.onlyLocal {
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

	client := newMackerelFromContext(c)
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
			_, err := client.DeleteMonitor(m.ID)
			logger.DieIf(err)
		}
	}
	for _, d := range monitorDiff.diff {
		logger.Log("info", "Update a rule.")
		fmt.Println(stringifyMonitor(d.local, ""))
		if !isDryRun {
			_, err := client.UpdateMonitor(d.remote.ID, d.local)
			logger.DieIf(err)
		}
	}
	return nil
}
