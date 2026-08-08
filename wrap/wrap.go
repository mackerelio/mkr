package wrap

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/Songmu/retry"
	"github.com/Songmu/wrapcommander"
	"github.com/mackerelio/mackerel-client-go"
	"github.com/mackerelio/mkr/logger"
	"github.com/urfave/cli/v3"
)

type wrap struct {
	name                 string
	detail               bool
	note                 string
	warning              bool
	autoClose            bool
	notificationInterval time.Duration
	hostID               string
	apibase              string
	apikey               string
	cmd                  []string

	outStream, errStream io.Writer
}

func (wr *wrap) run(ctx context.Context) error {
	re := wr.runCmd(ctx)
	if err := wr.report(context.Background(), re); err != nil {
		msg, _ := re.buildMsg(wr.detail)
		logger.Logf("error", "failed to post following report to Mackerel: %s\n%s", err, msg)
	}
	if !re.Success {
		return cli.Exit(re.Msg, re.ExitCode)
	}
	return nil
}

type syncBuffer struct {
	bytes.Buffer
	mu sync.Mutex
}

func (b *syncBuffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.Buffer.Write(p)
}

func (wr *wrap) runCmd(ctx context.Context) *result {
	cmd := exec.CommandContext(ctx, wr.cmd[0], wr.cmd[1:]...)
	re := &result{
		Cmd:  wr.cmd,
		Name: wr.name,
		Note: wr.note,
	}

	bufMerged := &syncBuffer{}
	cmd.Stdout = io.MultiWriter(wr.outStream, bufMerged)
	cmd.Stderr = io.MultiWriter(wr.errStream, bufMerged)

	if err := cmd.Start(); err != nil {
		return re.errorEnd("command invocation failed with follwing error: %s", err)
	}

	cmdErr := cmd.Wait()
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
	re.Output = bufMerged.String()

	re.Success = re.ExitCode == 0
	return re
}

func (wr *wrap) report(ctx context.Context, re *result) error {
	if wr.autoClose {
		defer func() {
			err := re.saveResult()
			if err != nil {
				logger.Logf("error", "failed to save result: %s", err)
			}
		}()
	}

	if wr.apikey == "" || wr.hostID == "" {
		return fmt.Errorf("Both of apikey and hostID are needed to report result to Mackerel") // nolint
	}
	var lastRe *result
	if wr.autoClose {
		var err error
		lastRe, err = re.loadLastResult()
		if err != nil {
			// resultFile something went wrong.
			// It may be no permission, broken json, not a normal file, and so on.
			// Though it is rough, try to delete as workaround.
			err := os.RemoveAll(re.resultFile())
			if err != nil {
				return err
			}
		}
	}
	if !re.Success || wr.autoClose && (lastRe == nil || !lastRe.Success) {
		return wr.doReport(ctx, re)
	}
	return nil
}

func (wr *wrap) doReport(ctx context.Context, re *result) error {
	checkSt := mackerel.CheckStatusOK
	if !re.Success {
		if wr.warning {
			checkSt = mackerel.CheckStatusWarning
		} else {
			checkSt = mackerel.CheckStatusCritical
		}
	}
	niInMinutes := uint(wr.notificationInterval.Minutes())
	if 0 < niInMinutes && niInMinutes < 10 {
		niInMinutes = 10
	}

	msg, err := re.buildMsg(wr.detail)
	if err != nil {
		return err
	}
	payload := &mackerel.CheckReports{
		Reports: []*mackerel.CheckReport{
			{
				Source:               mackerel.NewCheckSourceHost(wr.hostID),
				Name:                 re.checkName(),
				Status:               checkSt,
				OccurredAt:           time.Now().Unix(),
				Message:              msg,
				NotificationInterval: niInMinutes,
			},
		},
	}
	mcli, err := mackerel.NewClientWithOptions(wr.apikey, wr.apibase, false)
	if err != nil {
		return err
	}
	return retry.Retry(3, time.Second*3, func() error {
		return mcli.PostCheckReportsContext(ctx, payload)
	})
}
