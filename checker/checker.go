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
	checkers := make([]checker, len(conf.CheckPlugins))
	i := 0
	for name, p := range conf.CheckPlugins {
		checkers[i] = &checkPluginChecker{
			name: name,
			cp:   p,
		}
		i++
	}
	return runChecks(checkers)
}

type result struct {
	name, memo, cmd string
	stdout, stderr  string
	exitCode        int
	err             error
}

type checkPluginChecker struct {
	name string
	cp   *config.CheckPlugin
}

func (cpc *checkPluginChecker) check() *result {
	p := cpc.cp
	stdout, stderr, exitCode, err := p.Command.Run()
	cmdStr := p.Command.Cmd
	if cmdStr == "" {
		b, _ := json.Marshal(p.Command.Args)
		cmdStr = string(b)
	}
	return &result{
		name:     cpc.name,
		memo:     p.Memo,
		cmd:      cmdStr,
		stdout:   stdout,
		stderr:   stderr,
		exitCode: exitCode,
		err:      err,
	}
}

type checker interface {
	check() *result
}

func runChecks(checkers []checker) error {
	ch := make(chan *result)
	go func() {
		wg := &sync.WaitGroup{}
		wg.Add(len(checkers))
		for _, c := range checkers {
			go func(c checker) {
				defer wg.Done()
				ch <- c.check()
			}(c)
		}
		wg.Wait()
		close(ch)
	}()

	for re := range ch {
		fmt.Println(re)
	}
	return nil
}
