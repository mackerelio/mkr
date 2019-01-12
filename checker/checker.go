package checker

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/mackerelio/mackerel-agent/config"
	cli "gopkg.in/urfave/cli.v1"
)

var Command = cli.Command{
	Name:  "run-checks",
	Usage: "run check commands in mackerel-agent.conf",
	Description: `
    Execute command of check plugins in mackerel-agent.conf all at once.
    It is used for checking setting and operation of the check plugins.
`,
	Action: doRunChecks,
}

func doRunChecks(c *cli.Context) error {
	confFile := c.GlobalString("conf")
	conf, err := config.LoadConfig(confFile)
	if err != nil {
		return err
	}
	return runChecks(conf.CheckPlugins)
}

type result struct {
	name, memo, cmd string
	stdout, stderr  string
	exitCode        int
	err             error
}

func runChecks(plugins map[string]*config.CheckPlugin) error {
	ch := make(chan *result)
	go func() {
		wg := &sync.WaitGroup{}
		wg.Add(len(plugins))
		for name, p := range plugins {
			go func(name string, p *config.CheckPlugin) {
				defer wg.Done()
				stdout, stderr, exitCode, err := p.Command.Run()
				cmdStr := p.Command.Cmd
				if cmdStr == "" {
					b, _ := json.Marshal(p.Command.Args)
					cmdStr = string(b)
				}
				ch <- &result{
					name:     name,
					memo:     p.Memo,
					cmd:      cmdStr,
					stdout:   stdout,
					stderr:   stderr,
					exitCode: exitCode,
					err:      err,
				}
			}(name, p)
		}
		wg.Wait()
		close(ch)
	}()

	for re := range ch {
		fmt.Println(re)
	}
	return nil
}
