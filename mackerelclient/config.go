package mackerelclient

import (
	"os"

	"github.com/mackerelio/mackerel-agent/config"
)

// LoadApibaseFromConfig gets mackerel api Base URL (usually https://api.mackerelio.com/) from mackerel-agent.conf if it's installed mackerel-agent on localhost
func LoadApibaseFromConfig(conffile string) string {
	conf, err := config.LoadConfig(conffile)
	if err != nil {
		return ""
	}
	return conf.Apibase
}

// LoadApibaseFromConfigWithFallback get mackerel api Base URL from mackerel-agent.conf,
// and fallbacks to default (https://api.mackerelio.com/) if not specified.
func LoadApibaseFromConfigWithFallback(conffile string) string {
	apiBase := LoadApibaseFromConfig(conffile)
	if apiBase == "" {
		return config.DefaultConfig.Apibase
	}
	return apiBase
}

// LoadApikeyFromConfig gets mackerel.io apikey from mackerel-agent.conf if it's installed mackerel-agent on localhost
func LoadApikeyFromConfig(conffile string) string {
	conf, err := config.LoadConfig(conffile)
	if err != nil {
		return ""
	}
	return conf.Apikey
}

// LoadApikeyFromEnvOrConfig is similar to LoadApikeyFromConfig. return MACKEREL_APIKEY environment value if defined MACKEREL_APIKEY
func LoadApikeyFromEnvOrConfig(conffile string) string {
	if apiKey := os.Getenv("MACKEREL_APIKEY"); apiKey != "" {
		return apiKey
	}
	key := LoadApikeyFromConfig(conffile)
	return key
}

// LoadHostIDFromConfig gets localhost's hostID from conf.Root (ex. /var/lib/mackerel/id) if it's installed mackerel-agent on localhost
func LoadHostIDFromConfig(conffile string) string {
	conf, err := config.LoadConfig(conffile)
	if err != nil {
		return ""
	}
	hostID, err := conf.LoadHostID()
	if err != nil {
		return ""
	}
	return hostID
}
