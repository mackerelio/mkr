package main

import (
	"github.com/mackerelio/mackerel-agent/command"
	"github.com/mackerelio/mackerel-agent/config"
)

func LoadHostIdFromConfig() string {
	conf, err := config.LoadConfig(config.DefaultConfig.Conffile)
	if err != nil {
		return ""
	}
	hostId, err := command.LoadHostId(conf.Root)
	if err != nil {
		return ""
	}
	return hostId
}
