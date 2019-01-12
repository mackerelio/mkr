package checker

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
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
	return runChecks(checkers, os.Stdout)
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
		exitCode: exitCode,
		stdout:   stdout,
		stderr:   stderr,
		err:      err,
	}
}

type checker interface {
	check() *result
}

func runChecks(checkers []checker, w io.Writer) error {
	ch := make(chan *result)
	total := len(checkers)
	go func() {
		wg := &sync.WaitGroup{}
		wg.Add(total)
		for _, c := range checkers {
			go func(c checker) {
				defer wg.Done()
				ch <- c.check()
			}(c)
		}
		wg.Wait()
		close(ch)
	}()
	fmt.Fprintln(w, "TAP version 13")
	fmt.Fprintf(w, "1..%d", total)
	testNum := 1
	for re := range ch {
		fmt.Fprintln(w, re)
	}
	return nil
}
