package main

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/mackerelio/mackerel-agent/config"
)

const idFileName = "id"

func idFilePath(root string) string {
	return filepath.Join(root, idFileName)
}

func loadHostID(root string) (string, error) {
	content, err := ioutil.ReadFile(idFilePath(root))
	if err != nil {
		return "", err
	}
	return string(content), nil
}

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
	hostID, err := loadHostID(conf.Root)
	if err != nil {
		return ""
	}
	return hostID
}
