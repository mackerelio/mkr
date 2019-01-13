package wrap

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"

	"github.com/Songmu/wrapcommander"
	"github.com/mackerelio/mackerel-agent/config"
	"golang.org/x/sync/errgroup"
	cli "gopkg.in/urfave/cli.v1"
)

// CommandPlugin is definition of mkr plugin
var Command = cli.Command{
	Name:      "wrap",
	Usage:     "wrap command status",
	ArgsUsage: "[--verbose | -v] [--name | -n <name>] [--memo | -m <memo>] -- /path/to/batch",
	Description: `
    wrap command line
`,
	Action: doWrap,
	Flags: []cli.Flag{
		cli.StringFlag{Name: "name, n", Value: "", Usage: "monitor <name>"},
		cli.BoolFlag{Name: "verbose, v", Usage: "verbose output mode"},
		cli.StringFlag{Name: "memo, m", Value: "", Usage: "monitor <memo>"},
		cli.StringFlag{Name: "H, host", Value: "", Usage: "<hostId>"},
		cli.BoolFlag{Name: "warning, w", Usage: "alert as warning"},
	},
}

func doWrap(c *cli.Context) error {
	confFile := c.GlobalString("conf")
	conf, err := config.LoadConfig(confFile)
	if err != nil {
		return err
	}
	apibase := c.GlobalString("apibase")
	apikey := conf.Apikey
	if apikey == "" {
		apikey = os.Getenv("MACKEREL_APIKEY")
	}
	if apikey == "" {
		return fmt.Errorf(`MACKEREL_APIKEY environment variable is not set. (Try "export MACKEREL_APIKEY='<Your apikey>'`)
	}
	hostID, _ := conf.LoadHostID()
	if c.String("host") != "" {
		hostID = c.String("host")
	}
	if hostID == "" {
		return fmt.Errorf("failed to load hostID. (Try to specify -host option explicitly)")
	}
	cmd := c.Args()
	if len(cmd) > 0 && cmd[0] == "--" {
		cmd = cmd[1:]
	}
	if len(cmd) < 1 {
		return fmt.Errorf("no command specified")
	}

	return (&app{
		apibase: apibase,
		name:    c.String("name"),
		verbose: c.Bool("verbose"),
		memo:    c.String("memo"),
		warning: c.Bool("warning"),
		hostID:  hostID,
		apikey:  apikey,
		cmd:     cmd,
	}).run()
}

type app struct {
	apibase string
	name    string
	verbose bool
	memo    string
	warning bool
	hostID  string
	apikey  string
	cmd     []string
}

type result struct {
	cmd        []string
	name, memo string

	output, stdout, stderr string
	pid, exitCode          *int
	signaled               bool
	startAt, endAt         time.Time

	msg     string
	success bool
}

func (ap *app) run() error {
	re := ap.runCmd()
	_ = re
	return nil
}

func (ap *app) runCmd() *result {
	cmd := exec.Command(ap.cmd[0], ap.cmd[1:]...)
	re := &result{
		cmd:  ap.cmd,
		name: ap.name,
		memo: ap.memo,
	}

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		re.msg = fmt.Sprintf("command invocation failed with follwing error: %s", err)
		return re
	}
	defer stdoutPipe.Close()

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		re.msg = fmt.Sprintf("command invocation failed with follwing error: %s", err)
		return re
	}
	defer stderrPipe.Close()

	var (
		bufStdout = &bytes.Buffer{}
		bufStderr = &bytes.Buffer{}
		bufMerged = &bytes.Buffer{}
	)
	stdoutPipe2 := io.TeeReader(stdoutPipe, io.MultiWriter(bufStdout, bufMerged))
	stderrPipe2 := io.TeeReader(stderrPipe, io.MultiWriter(bufStderr, bufMerged))

	re.startAt = time.Now()
	err = cmd.Start()
	if err != nil {
		re.msg = fmt.Sprintf("command invocation failed with follwing error: %s", err)
		return re
	}
	re.pid = &cmd.Process.Pid
	eg := &errgroup.Group{}

	eg.Go(func() error {
		defer stdoutPipe.Close()
		_, err := io.Copy(os.Stdout, stdoutPipe2)
		return err
	})
	eg.Go(func() error {
		defer stderrPipe.Close()
		_, err := io.Copy(os.Stderr, stderrPipe2)
		return err
	})
	eg.Wait()

	cmdErr := cmd.Wait()
	re.endAt = time.Now()
	ex := wrapcommander.ResolveExitCode(cmdErr)
	re.exitCode = &ex
	if *re.exitCode > 128 {
		w, ok := wrapcommander.ErrorToWaitStatus(cmdErr)
		if ok {
			re.signaled = w.Signaled()
		}
	}
	if !re.signaled {
		re.msg = fmt.Sprintf("command exited with code: %d", *re.exitCode)
	} else {
		re.msg = fmt.Sprintf("command died with signal: %d", *re.exitCode&127)
	}
	re.stdout = bufStdout.String()
	re.stderr = bufStderr.String()
	re.output = bufMerged.String()

	re.success = *re.exitCode == 0
	return re
}
