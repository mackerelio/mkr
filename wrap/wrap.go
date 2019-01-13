package wrap

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"

	"github.com/Songmu/retry"
	"github.com/Songmu/wrapcommander"
	mackerel "github.com/mackerelio/mackerel-client-go"
	"github.com/mackerelio/mkr/logger"
	"golang.org/x/sync/errgroup"
	cli "gopkg.in/urfave/cli.v1"
)

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

func (wr *wrap) run() error {
	re := wr.runCmd()
	if err := wr.report(re); err != nil {
		logger.Logf("error", "failed to post following report to Mackerel: %s\n%s",
			err, re.buildMsg(wr.verbose))
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

func (wr *wrap) doReport(re *result) error {
	checkSt := mackerel.CheckStatusOK
	if !re.Success {
		if wr.warning {
			checkSt = mackerel.CheckStatusWarning
		} else {
			checkSt = mackerel.CheckStatusCritical
		}
	}
	payload := &mackerel.CheckReports{
		Reports: []*mackerel.CheckReport{
			{
				Source:     mackerel.NewCheckSourceHost(wr.hostID),
				Name:       re.checkName(),
				Status:     checkSt,
				OccurredAt: time.Now().Unix(),
				Message:    re.buildMsg(wr.verbose),
			},
		},
	}
	mcli, err := mackerel.NewClientWithOptions(wr.apikey, wr.apibase, false)
	if err != nil {
		return err
	}
	return retry.Retry(3, time.Second*3, func() error {
		return mcli.PostCheckReports(payload)
	})
}
