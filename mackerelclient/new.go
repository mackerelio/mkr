package mackerelclient

import (
	"fmt"
	"os"

	"github.com/urfave/cli"

	"github.com/mackerelio/mackerel-agent/config"
	"github.com/mackerelio/mackerel-client-go"

	"github.com/mackerelio/mkr/logger"
)

// New returns new mackerel client
func New(conffile, apibase string) (Client, error) {
	apikey := os.Getenv("MACKEREL_APIKEY")
	var conf *config.Config
	if apikey == "" {
		var err error
		conf, err = config.LoadConfig(conffile)
		if err != nil {
			return nil, err
		}
		apikey = conf.Apikey
	}
	if apikey == "" {
		return nil, fmt.Errorf("no mackerel apikeys are specified from MACKEREL_APIKEY or config")
	}
	if apibase == "" {
		if conf == nil {
			conf, _ = config.LoadConfig(conffile)
			if conf == nil {
				conf = config.DefaultConfig
			}
		}
		apibase = conf.Apibase
	}
	return mackerel.NewClientWithOptions(apikey, apibase, os.Getenv("DEBUG") != "")
}

// NewFromContext returns mackerel client from cli.Context
func NewFromContext(c *cli.Context) *mackerel.Client {
	confFile := c.GlobalString("conf")
	apiBase := c.GlobalString("apibase")
	apiKey := LoadApikeyFromEnvOrConfig(confFile)
	if apiKey == "" {
		logger.Log("error", `
    MACKEREL_APIKEY environment variable is not set. (Try "export MACKEREL_APIKEY='<Your apikey>'")
`)
		os.Exit(1)
	}

	if apiBase == "" {
		apiBase = LoadApibaseFromConfigWithFallback(confFile)
	}

	client, err := mackerel.NewClientWithOptions(apiKey, apiBase, os.Getenv("DEBUG") != "")
	logger.DieIf(err)

	return client
}
