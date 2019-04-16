package wrap

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/Songmu/wrapcommander"
)

type result struct {
	Cmd        []string
	Name, Note string

	Output   string `json:"-"`
	ExitCode int
	Signaled bool

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

const msgTplText = `{{.Msg}}
{{- if ne .Note "" }}
Note: {{.Note}}{{end}}
% {{.Command}}
{{- if .Detail }}
{{.Output}}{{end}}`

var msgTpl *template.Template

func init() {
	msgTpl = template.Must(template.New("msg").Parse(msgTplText))
}

func (re *result) buildMsg(detail bool) string {
	s := struct {
		Msg, Note, Command, Output string
		Detail                     bool
	}{
		re.Msg, re.Note, strings.Join(re.Cmd, " "), re.Output,
		detail,
	}
	buf := &bytes.Buffer{}
	template.Must(msgTpl.Clone()).Execute(buf, s)
	const messageLengthLimit = 1024 // https://mackerel.io/api-docs/entry/check-monitoring#post
	return truncate(buf.String(), messageLengthLimit, "\n...\n")
}

func truncate(src string, limit int, sep string) string {
	rs := []rune(src)
	if len(rs) <= limit {
		return src
	}
	seprs := []rune(sep)
	var i int
	rest := limit - len(seprs)
	if 0 < rest {
		i += rest / 2
		rest -= rest / 2
	}
	i += copy(rs[i:], seprs)
	if 0 < rest {
		copy(rs[i:], rs[len(rs)-rest:])
	}
	return string(rs[:limit])
}
