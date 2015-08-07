package main

import (
	"encoding/json"
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
	data, err := json.MarshalIndent(monitors, "", "    ")
	logger.DieIf(err)

	_, err = file.WriteString(string(data))
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
			diff = appendDiff(diff, name, fA.Interface(), fB.Interface())
			if fA.Interface() != fB.Interface() {
				diffNum++
			}
		} else {
			name := strings.Replace(sAType.Field(i).Tag.Get("json"), ",omitempty", "", 1)
			//diff = append(diff, fmt.Sprintf("    \"%s\": [", sAType.Field(i).Name))
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
		//fmt.Println(strings.Join(diff, "\n"))
		return strings.Join(diff, "\n")
	}
	return ""
}

func isSameMonitor(a *mkr.Monitor, b *mkr.Monitor) (string, bool) {
	if a == nil || b == nil {
		return "", false
	}
	if reflect.DeepEqual(*a, *b) {
		return "", true
	}
	if a.ID == b.ID {
		diff := diffMonitor(a, b)
		if diff != "" {
			return diff, false
		}
		return "", true
	}
	return "", false
}

func doMonitorsDiff(c *cli.Context) {
	conffile := c.GlobalString("conf")

	monitorsRemote, err := newMackerel(conffile).FindMonitors()
	logger.DieIf(err)

	monitorsLocal, err := monitorLoadRules()
	logger.DieIf(err)

	var onlyRemote []*(mkr.Monitor)
	var onlyLocal []*(mkr.Monitor)
	var diffs []string
	counter := map[string]uint64{"diff": 0, "remote": 0, "local": 0}

	for _, remote := range monitorsRemote {
		found := false
		for i, local := range monitorsLocal {
			diff, isSame := isSameMonitor(remote, local)
			if isSame {
				monitorsLocal[i] = nil
				found = true
				break
			} else if diff != "" {
				monitorsLocal[i] = nil
				diffs = append(diffs, diff)
				counter["diff"]++
				found = true
				break
			}
		}
		if found == false {
			onlyRemote = append(onlyRemote, remote)
			counter["remote"]++
		}
	}
	for _, local := range monitorsLocal {
		if local != nil {
			onlyLocal = append(onlyLocal, local)
			counter["local"]++
		}
	}
	fmt.Printf("Summary: %d modify, %d append, %d remove\n\n", counter["diff"], counter["local"], counter["remote"])
	for _, diff := range diffs {
		fmt.Println(diff)
	}
	for _, m := range onlyRemote {
		printMonitor(m, "-")
	}
	for _, m := range onlyLocal {
		printMonitor(m, "+")
	}
}
