package wrap

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/Songmu/wrapcommander"
)

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

func normalizeName(str string) string {
	return reg.ReplaceAllString(strings.TrimSpace(str), "_")
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

func (re *result) buildMsg(verbose bool) string {
	msg := re.Msg
	if re.Memo != "" {
		msg += "\nMemo: " + re.Memo
	}
	msg += "\n% " + strings.Join(re.Cmd, " ")
	if verbose {
		msg += "\n" + re.Output
	}
	const messageLengthLimit = 1024
	runes := []rune(msg)
	if len(runes) > messageLengthLimit {
		msg = string(runes[0:messageLengthLimit])
	}
	return msg
}
