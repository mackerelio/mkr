package main

type HostFormat struct {
	Id            string   `json:"id,omitempty"`
	Name          string   `json:"name,omitempty"`
	Status        string   `json:"status,omitempty"`
	Memo          string   `json:"memo,omitempty"`
	RoleFullnames []string `json:"roleFullnames,omitempty"`
	IsRetired     bool     `json:"isRetired,omitempty"`
	CreatedAt     string   `json:"createdAt,omitempty"`
}
