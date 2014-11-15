package main

import (
	"os"

	"github.com/mackerelio/mackerel-agent/command"
	"github.com/mackerelio/mackerel-agent/config"
)

func LoadApikeyFromConfig() string {
	conf, err := config.LoadConfig(config.DefaultConfig.Conffile)
	if err != nil {
		return ""
	}
	return conf.Apikey
}

func LoadApikeyFromConfigOrEnv() string {
	if apiKey := os.Getenv("MACKEREL_APIKEY"); apiKey != "" {
		return apiKey
	}
	key := LoadApikeyFromConfig()
	return key
}

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
