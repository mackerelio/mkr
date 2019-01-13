package wrap

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/Songmu/retry"
	"github.com/Songmu/wrapcommander"
	"github.com/mackerelio/mackerel-agent/config"
	mackerel "github.com/mackerelio/mackerel-client-go"
	"github.com/mackerelio/mkr/logger"
	"golang.org/x/sync/errgroup"
	cli "gopkg.in/urfave/cli.v1"
)

// Command is definition of mkr wrap
var Command = cli.Command{
	Name:      "wrap",
	Usage:     "wrap command status",
	ArgsUsage: "[--verbose | -v] [--name | -n <name>] [--memo | -m <memo>] -- /path/to/batch",
	Description: `
    wrap command line
`,
	Action: doWrap,
	Flags: []cli.Flag{
		cli.StringFlag{Name: "name, n", Value: "", Usage: "monitored `check-name` which must be unique on a host"},
		cli.BoolFlag{Name: "verbose, v", Usage: "verbose output"},
		cli.StringFlag{Name: "memo, m", Value: "", Usage: "`memo` of the job"},
		cli.StringFlag{Name: "H, host", Value: "", Usage: "`hostID`"},
		cli.BoolFlag{Name: "warning, w", Usage: "alerts as warning"},
	},
}

func doWrap(c *cli.Context) error {
	confFile := c.GlobalString("conf")
	conf, err := config.LoadConfig(confFile)
	if err != nil {
		return err
	}
	apibase := c.GlobalString("apibase")
	if apibase == "" {
		apibase = conf.Apibase
	}

	apikey := conf.Apikey
	if apikey == "" {
		apikey = os.Getenv("MACKEREL_APIKEY")
	}
	if apikey == "" {
		logger.Log("error", "[mkr wrap] failed to detect Mackerel APIKey. Try to specify in mackerel-agent.conf or export MACKEREL_APIKEY='<Your apikey>'")
	}
	hostID, _ := conf.LoadHostID()
	if c.String("host") != "" {
		hostID = c.String("host")
	}
	if hostID == "" {
		logger.Log("error", "[mkr wrap] failed to load hostID. Try to specify -host option explicitly")
	}
	// Since command execution has the highest priority, even when apikey or
	// hostID is empty, we don't return errors and only output the log here.

	cmd := c.Args()
	if len(cmd) > 0 && cmd[0] == "--" {
		cmd = cmd[1:]
	}
	if len(cmd) < 1 {
		return fmt.Errorf("no commands specified")
	}

	return (&wrap{
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

type wrap struct {
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
	Cmd        []string
	Name, Memo string

	Output, Stdout, Stderr string `json:"-"`
	Pid                    int
	ExitCode               int
	Signaled               bool
	StartAt, EndAt         time.Time

	Msg     string
	Success bool
}

var reg = regexp.MustCompile(`[^-a-zA-Z0-9_]`)

func normalizeName(devName string) string {
	return reg.ReplaceAllString(strings.TrimSpace(devName), "_")
}

func (re *result) checkName() string {
	if re.Name != "" {
		return re.Name
	}
	sum := md5.Sum([]byte(strings.Join(re.Cmd, " ")))
	return fmt.Sprintf("mkrwrap-%s-%x",
		normalizeName(filepath.Base(re.Cmd[0])),
		sum[0:3])
}

func (re *result) resultFile() string {
	return filepath.Join(os.TempDir(), fmt.Sprintf("mkrwrap-%s.json", re.checkName()))
}

func (re *result) loadLastResult() (*result, error) {
	prevRe := &result{}
	fname := re.resultFile()

	f, err := os.Open(fname)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer f.Close()

	err = json.NewDecoder(f).Decode(prevRe)
	return prevRe, err
}

func (re *result) saveResult() error {
	fname := re.resultFile()
	tmpf, err := ioutil.TempFile(filepath.Dir(fname), "tmp-mkrwrap")
	if err != nil {
		return err
	}
	defer func(tmpfname string) {
		tmpf.Close()
		os.Remove(tmpfname)
	}(tmpf.Name())

	if err := json.NewEncoder(tmpf).Encode(re); err != nil {
		return err
	}
	if err := tmpf.Close(); err != nil {
		return err
	}
	return os.Rename(tmpf.Name(), fname)
}

func (re *result) errorEnd(format string, err error) *result {
	re.Msg = fmt.Sprintf(format, err)
	re.ExitCode = wrapcommander.ResolveExitCode(err)
	return re
}

func (wr *wrap) run() error {
	re := wr.runCmd()
	if err := wr.report(re); err != nil {
		logger.Logf("error", "failed to post following report to Mackerel: %s\n%s",
			err, wr.buildMsg(re))
	}
	if !re.Success {
		return cli.NewExitError(re.Msg, re.ExitCode)
	}
	return nil
}

func (wr *wrap) runCmd() *result {
	cmd := exec.Command(wr.cmd[0], wr.cmd[1:]...)
	re := &result{
		Cmd:  wr.cmd,
		Name: wr.name,
		Memo: wr.memo,
	}

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return re.errorEnd("command invocation failed with follwing error: %s", err)
	}
	defer stdoutPipe.Close()

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return re.errorEnd("command invocation failed with follwing error: %s", err)
	}
	defer stderrPipe.Close()

	var (
		bufStdout = &bytes.Buffer{}
		bufStderr = &bytes.Buffer{}
		bufMerged = &bytes.Buffer{}
	)
	stdoutPipe2 := io.TeeReader(stdoutPipe, io.MultiWriter(bufStdout, bufMerged))
	stderrPipe2 := io.TeeReader(stderrPipe, io.MultiWriter(bufStderr, bufMerged))

	re.StartAt = time.Now()
	err = cmd.Start()
	if err != nil {
		return re.errorEnd("command invocation failed with follwing error: %s", err)
	}
	re.Pid = cmd.Process.Pid
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
	re.EndAt = time.Now()
	re.ExitCode = wrapcommander.ResolveExitCode(cmdErr)
	if re.ExitCode > 128 {
		w, ok := wrapcommander.ErrorToWaitStatus(cmdErr)
		if ok {
			re.Signaled = w.Signaled()
		}
	}
	if !re.Signaled {
		re.Msg = fmt.Sprintf("command exited with code: %d", re.ExitCode)
	} else {
		re.Msg = fmt.Sprintf("command died with signal: %d", re.ExitCode&127)
	}
	re.Stdout = bufStdout.String()
	re.Stderr = bufStderr.String()
	re.Output = bufMerged.String()

	re.Success = re.ExitCode == 0
	return re
}

func (wr *wrap) report(re *result) error {
	defer func() {
		err := re.saveResult()
		if err != nil {
			logger.Logf("error", "failed to save result: %s", err)
		}
	}()

	if wr.apikey == "" || wr.hostID == "" {
		return fmt.Errorf("Both of apikey and hostID are needed to report result to Mackerel")
	}
	lastRe, err := re.loadLastResult()
	if err != nil {
		// resultFile something went wrong.
		// It may be no permission, broken json, not a normal file, and so on.
		// Though it rough, try to delete as workaround
		err := os.RemoveAll(re.resultFile())
		if err != nil {
			return err
		}
	}
	if lastRe == nil || !lastRe.Success || !re.Success {
		return wr.doReport(re)
	}
	return nil
}

func (wr *wrap) buildMsg(re *result) string {
	msg := re.Msg
	if re.Memo != "" {
		msg += "\nMemo: " + re.Memo
	}
	msg += "\n% " + strings.Join(re.Cmd, " ")
	if wr.verbose {
		msg += "\n" + re.Output
	}
	const messageLengthLimit = 1024
	runes := []rune(msg)
	if len(runes) > messageLengthLimit {
		msg = string(runes[0:messageLengthLimit])
	}
	return msg
}

func (wr *wrap) doReport(re *result) error {
	checkSt := mackerel.CheckStatusOK
	if !re.Success {
		if wr.warning {
			checkSt = mackerel.CheckStatusWarning
		} else {
			checkSt = mackerel.CheckStatusCritical
		}
	}
	crs := &mackerel.CheckReports{
		Reports: []*mackerel.CheckReport{
			{
				Source:     mackerel.NewCheckSourceHost(wr.hostID),
				Name:       re.checkName(),
				Status:     checkSt,
				OccurredAt: time.Now().Unix(),
				Message:    wr.buildMsg(re),
			},
		},
	}
	cli, err := mackerel.NewClientWithOptions(wr.apikey, wr.apibase, false)
	if err != nil {
		return err
	}
	return retry.Retry(3, time.Second*3, func() error {
		return cli.PostCheckReports(crs)
	})
}
