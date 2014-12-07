package main

import (
	"os"

	"github.com/mackerelio/mackerel-agent/command"
	"github.com/mackerelio/mackerel-agent/config"
)

// LoadApikeyFromConfig gets mackerel.io apikey from mackerel-agent.conf if it's installed mackerel-agent on localhost
func LoadApikeyFromConfig() string {
	conf, err := config.LoadConfig(config.DefaultConfig.Conffile)
	if err != nil {
		return ""
	}
	return conf.Apikey
}

// LoadApikeyFromEnvOrConfig is similar to LoadApikeyFromConfig. return MACKEREL_APIKEY environment value if defined MACKEREL_APIKEY 
func LoadApikeyFromEnvOrConfig() string {
	if apiKey := os.Getenv("MACKEREL_APIKEY"); apiKey != "" {
		return apiKey
	}
	key := LoadApikeyFromConfig()
	return key
}

// LoadHostIDFromConfig gets localhost's hostID from conf.Root (ex. /var/lib/mackerel/id) if it's installed mackerel-agent on localhost
func LoadHostIDFromConfig() string {
	conf, err := config.LoadConfig(config.DefaultConfig.Conffile)
	if err != nil {
		return ""
	}
	hostID, err := command.LoadHostId(conf.Root)
	if err != nil {
		return ""
	}
	return hostID
}
