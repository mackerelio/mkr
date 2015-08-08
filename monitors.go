package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"sort"
	"strings"

	"github.com/codegangsta/cli"
	mkr "github.com/mackerelio/mackerel-client-go"
	"github.com/mackerelio/mkr/logger"
)

func monitorSaveRules(rules []*(mkr.Monitor)) error {
	file, err := os.Create("monitors.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	monitors := map[string]interface{}{"monitors": rules}
	dataRaw, err := json.MarshalIndent(monitors, "", "    ")
	logger.DieIf(err)
	data := strings.Replace(string(dataRaw), "\\u003c", "<", -1)
	data = strings.Replace(data, "\\u003e", ">", -1)

	_, err = file.WriteString(data)
	if err != nil {
		return err
	}
	return nil
}

func monitorLoadRules() ([]*(mkr.Monitor), error) {
	buff, err := ioutil.ReadFile("monitors.json")
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

func doMonitorsList(c *cli.Context) {
	conffile := c.GlobalString("conf")

	monitors, err := newMackerel(conffile).FindMonitors()
	logger.DieIf(err)

	PrettyPrintJSON(monitors)
}

func doMonitorsPull(c *cli.Context) {
	conffile := c.GlobalString("conf")
	isVerbose := c.Bool("verbose")

	monitors, err := newMackerel(conffile).FindMonitors()
	logger.DieIf(err)

	monitorSaveRules(monitors)

	if isVerbose {
		PrettyPrintJSON(monitors)
	}
}

func appendDiff(src []string, name string, a interface{}, b interface{}) []string {
	diff := []string{}
	aType := reflect.TypeOf(a).String()
	format := "\"%s\""
	switch aType {
	case "uint64":
		format = "%d"
	case "float64":
		format = "%f"
	}
	if b != nil && a != b {
		diff = append(src, fmt.Sprintf("-  \"%s\": "+format+",", name, a))
		diff = append(diff, fmt.Sprintf("+  \"%s\": "+format+",", name, b))
	} else {
		diff = append(src, fmt.Sprintf("   \"%s\": "+format+",", name, a))
	}
	return diff
}

func printMonitor(a *mkr.Monitor, prefix string) {
	sA := reflect.ValueOf(a).Elem()
	diff := []string{" {"}
	for i := 0; i < sA.NumField(); i++ {
		fA := sA.Field(i)
		sAType := sA.Type()
		name := strings.Replace(sAType.Field(i).Tag.Get("json"), ",omitempty", "", 1)
		if sAType.Field(i).Type.String() != "[]string" {
			if name == "id" && fA.Interface() == "" {
				continue
			}
			diff = appendDiff(diff, name, fA.Interface(), nil)
		} else {
			diff = append(diff, fmt.Sprintf("   \"%s\": [", name))
			sortA := fA.Interface().([]string)
			sort.Strings(sortA)
			i := 0
			for i < len(sortA) {
				diff = append(diff, fmt.Sprintf("     \"%s\",", sortA[i]))
				i++
			}
			diff = append(diff, "   ],")
		}
	}
	diff = append(diff, " },")
	for _, d := range diff {
		fmt.Println(prefix + d)
	}
}

func diffMonitor(a *mkr.Monitor, b *mkr.Monitor) string {
	diff := []string{"  {"}
	diffNum := 0
	sA := reflect.ValueOf(a).Elem()
	sB := reflect.ValueOf(b).Elem()
	for i := 0; i < sA.NumField(); i++ {
		fA := sA.Field(i)
		fB := sB.Field(i)
		sAType := sA.Type()
		if sAType.Field(i).Type.String() != "[]string" {
			name := strings.Replace(sAType.Field(i).Tag.Get("json"), ",omitempty", "", 1)
			if name == "id" {
				continue
			}
			diff = appendDiff(diff, name, fA.Interface(), fB.Interface())
			if fA.Interface() != fB.Interface() {
				diffNum++
			}
		} else {
			name := strings.Replace(sAType.Field(i).Tag.Get("json"), ",omitempty", "", 1)
			diff = append(diff, fmt.Sprintf("    \"%s\": [", name))
			sortA := fA.Interface().([]string)
			sortB := fB.Interface().([]string)
			sort.Strings(sortA)
			sort.Strings(sortB)
			i := 0
			j := 0
			for i < len(sortA) && j < len(sortB) {
				if sortA[i] == sortB[j] {
					diff = append(diff, fmt.Sprintf("      \"%s\",", sortA[i]))
					i++
					j++
				} else if sortA[i] < sortB[j] {
					diff = append(diff, fmt.Sprintf("-     \"%s\",", sortA[i]))
					i++
					diffNum++
				} else if sortB[j] < sortA[i] {
					diff = append(diff, fmt.Sprintf("+     \"%s\",", sortB[j]))
					j++
					diffNum++
				}
			}
			diff = append(diff, "    ],")
		}
	}

	if diffNum > 0 {
		diff = append(diff, "  },")
		return strings.Join(diff, "\n")
	}
	return ""
}

func isSameMonitor(a *mkr.Monitor, b *mkr.Monitor, flagNameUniqueness bool) (string, bool) {
	if a == nil || b == nil {
		return "", false
	}
	if reflect.DeepEqual(*a, *b) {
		return "", true
	}
	if a.ID == b.ID || (flagNameUniqueness == true && a.Name == b.Name) {
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
		if m.Name == "" {
			return false, errors.New("Monitor should have 'Name'")
		}
		if m.Type != "host" && m.Type != "service" && m.Type != "external" && m.Type != "passive" {
			return false, fmt.Errorf("Unknown type is found: %s", m.Type)
		}
	}

	// check name uniqueness
	names := map[string]bool{}
	for _, m := range monitors {
		if names[m.Name] {
			logger.Log("Warning: ", fmt.Sprintf("Names of %s is not unique.", label))
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
	conffile := c.GlobalString("conf")

	var monitorDiff monitorDiff

	monitorsRemote, err := newMackerel(conffile).FindMonitors()
	logger.DieIf(err)
	flagNameUniquenessRemote, err := validateRules(monitorsRemote, "remote rules")
	logger.DieIf(err)

	monitorsLocal, err := monitorLoadRules()
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

func doMonitorsDiff(c *cli.Context) {
	monitorDiff := checkMonitorsDiff(c)

	var diffs []string
	for _, d := range monitorDiff.diff {
		diffs = append(diffs, diffMonitor(d.remote, d.local))
	}

	fmt.Printf("Summary: %d modify, %d append, %d remove\n\n", len(monitorDiff.diff), len(monitorDiff.onlyLocal), len(monitorDiff.onlyRemote))
	for _, diff := range diffs {
		fmt.Println(diff)
	}
	for _, m := range monitorDiff.onlyRemote {
		printMonitor(m, "-")
	}
	for _, m := range monitorDiff.onlyLocal {
		printMonitor(m, "+")
	}
}

func doMonitorsPush(c *cli.Context) {
	monitorDiff := checkMonitorsDiff(c)
	isDryRun := c.Bool("dryRun")
	isVerbose := c.Bool("verbose")

	conffile := c.GlobalString("conf")
	client := newMackerel(conffile)
	if isVerbose {
		client.Verbose = true
	}

	for _, m := range monitorDiff.onlyLocal {
		logger.Log("info", "Create a new rule.")
		printMonitor(m, "")
		if !isDryRun {
			_, err := client.CreateMonitor(m)
			logger.DieIf(err)
		}
	}
	for _, m := range monitorDiff.onlyRemote {
		logger.Log("info", "Delete a rule.")
		printMonitor(m, "")
		if !isDryRun {
			_, err := client.DeleteMonitor(m.ID)
			logger.DieIf(err)
		}
	}
	for _, d := range monitorDiff.diff {
		logger.Log("info", "Update a rule.")
		printMonitor(d.local, "")
		if !isDryRun {
			_, err := client.UpdateMonitor(d.remote.ID, d.local)
			logger.DieIf(err)
		}
	}

}
