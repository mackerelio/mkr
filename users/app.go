package users

import (
	"fmt"
	"io"
	"text/template"
	"time"

	"github.com/mackerelio/mkr/format"
	"github.com/mackerelio/mkr/mackerelclient"
)

type appLogger interface {
	Log(string, string)
	Error(error)
}

type userApp struct {
	client    mackerelclient.Client
	logger    appLogger
	outStream io.Writer
	jqFilter  string
}

type findUsersParam struct {
	verbose bool
	format  string
}

// User defines output json structure.
type User struct {
	ID                      string   `json:"id,omitempty"`
	ScreenName              string   `json:"screenName,omitempty"`
	Email                   string   `json:"email,omitempty"`
	Authority               string   `json:"authority,omitempty"`
	IsInRegistrationProcess bool     `json:"isInRegistrationProcess"`
	IsMFAEnabled            bool     `json:"isMFAEnabled"`
	AuthenticationMethods   []string `json:"authenticationMethods,omitempty"`
	JoinedAt                string   `json:"joinedAt,omitempty"`
}

func (ua *userApp) findUsers(param findUsersParam) error {
	users, err := ua.client.FindUsers()
	if err != nil {
		return err
	}

	switch {
	case param.format != "" && ua.jqFilter != "":
		return fmt.Errorf("--format and --jq options are incompatible.")
	case param.format != "":
		t, err := template.New("format").Parse(param.format)
		if err != nil {
			return err
		}
		return t.Execute(ua.outStream, users)
	case param.verbose:
		return format.PrettyPrintJSON(ua.outStream, users, ua.jqFilter)
	default:
		usersFormat := make([]*User, 0)
		for _, user := range users {
			usersFormat = append(usersFormat, &User{
				ID:                      user.ID,
				ScreenName:              user.ScreenName,
				Email:                   user.Email,
				Authority:               user.Authority,
				IsInRegistrationProcess: user.IsInRegistrationProcess,
				IsMFAEnabled:            user.IsMFAEnabled,
				AuthenticationMethods:   user.AuthenticationMethods,
				JoinedAt:                format.ISO8601Extended(time.Unix(user.JoinedAt, 0)),
			})
		}
		return format.PrettyPrintJSON(ua.outStream, usersFormat, ua.jqFilter)
	}
}
