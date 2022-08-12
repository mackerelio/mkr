package format

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/itchyny/gojq"
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
func PrettyPrintJSON(outStream io.Writer, src interface{}, query string) error {
	if query == "" {
		_, err := fmt.Fprintln(outStream, JSONMarshalIndent(src, "", "    "))
		return err
	}

	return FilterJSON(outStream, src, query)
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

func FilterJSON(outStream io.Writer, src interface{}, queryStr string) error {
	query, err := gojq.Parse(queryStr)
	if err != nil {
		return err
	}

	var dst interface{}
	jsonObj, err := json.Marshal(src)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonObj, &dst)
	if err != nil {
		return err
	}

	iter := query.Run(dst)
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}

		if err, ok := v.(error); ok {
			return err
		}

		if text, e := jsonScalarToString(v); e == nil {
			_, err := fmt.Fprintln(outStream, text)
			if err != nil {
				return err
			}
		} else {
			_, err = fmt.Fprintln(outStream, JSONMarshalIndent(v, "", "    "))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// copy from https://github.com/cli/cli/blob/d21d388b8dc10c8f04187c3afa6e0b44f0977c65/pkg/export/template.go#L110-L127
func jsonScalarToString(input interface{}) (string, error) {
	switch tt := input.(type) {
	case string:
		return tt, nil
	case float64:
		if math.Trunc(tt) == tt {
			return strconv.FormatFloat(tt, 'f', 0, 64), nil
		} else {
			return strconv.FormatFloat(tt, 'f', 2, 64), nil
		}
	case nil:
		return "", nil
	case bool:
		return fmt.Sprintf("%v", tt), nil
	default:
		return "", fmt.Errorf("cannot convert type to string: %v", tt)
	}
}
