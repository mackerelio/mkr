package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/mackerelio/mkr/logger"
)

// HostFormat defines output json structure.
type HostFormat struct {
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
func PrettyPrintJSON(src interface{}) {
	data, err := json.MarshalIndent(src, "", "    ")
	logger.DieIf(err)
	fmt.Fprintln(os.Stdout, string(data))
}
