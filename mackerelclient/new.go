package mackerelclient

import (
	"os"

	mkr "github.com/mackerelio/mackerel-client-go"
	"github.com/mackerelio/mkr/logger"
	cli "gopkg.in/urfave/cli.v1"
)

// NewFromContext returns mackerel client from cli.Context
func NewFromContext(c *cli.Context) *mkr.Client {
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

	mackerel, err := mkr.NewClientWithOptions(apiKey, apiBase, os.Getenv("DEBUG") != "")
	logger.DieIf(err)

	return mackerel
}
