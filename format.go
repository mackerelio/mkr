package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/mackerelio/mkr/utils"
)

type HostFormat struct {
	Id            string            `json:"id,omitempty"`
	Name          string            `json:"name,omitempty"`
	Status        string            `json:"status,omitempty"`
	Memo          string            `json:"memo,omitempty"`
	RoleFullnames []string          `json:"roleFullnames,omitempty"`
	IsRetired     bool              `json:"isRetired"` // 'omitempty' regard boolean 'false' as empty.
	CreatedAt     string            `json:"createdAt,omitempty"`
	IpAddresses   map[string]string `json:"ipAddresses,omitempty"`
}

func PrettyPrintJson(src interface{}) {
	data, err := json.MarshalIndent(src, "", "    ")
	utils.DieIf(err)
	fmt.Fprintln(os.Stdout, string(data))
}
