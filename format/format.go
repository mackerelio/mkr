package format

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/mackerelio/mkr/logger"
)

// Host defines output json structure.
type Host struct {
	ID            string            `json:"id,omitempty"`
	Name          string            `json:"name,omitempty"`
	DisplayName   string            `json:"displayName,omitempty"`
	Status        string            `json:"status,omitempty"`
	Memo          string            `json:"memo,omitempty"`
	RoleFullnames []string          `json:"roleFullnames,omitempty"`
	IsRetired     bool              `json:"isRetired"` // 'omitempty' regard boolean 'false' as empty.
	CreatedAt     string            `json:"createdAt,omitempty"`
	IPAddresses   map[string]string `json:"ipAddresses,omitempty"`
}

// PrettyPrintJSON output indented json via stdout.
func PrettyPrintJSON(src interface{}) error {
	_, err := fmt.Fprintln(os.Stdout, JSONMarshalIndent(src, "", "    "))
	return err
}

// JSONMarshalIndent call json.MarshalIndent and replace encoded angle brackets
func JSONMarshalIndent(src interface{}, prefix, indent string) string {
	dataRaw, err := json.MarshalIndent(src, prefix, indent)
	logger.DieIf(err)
	return replaceAngleBrackets(string(dataRaw))
}

func replaceAngleBrackets(s string) string {
	s = strings.Replace(s, "\\u003c", "<", -1)
	return strings.Replace(s, "\\u003e", ">", -1)
}

// ISO8601Extended format
func ISO8601Extended(t time.Time) string {
	const layoutISO8601Exetnded = "2006-01-02T15:04:05-07:00"
	return t.Format(layoutISO8601Exetnded)
}
