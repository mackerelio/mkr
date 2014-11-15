package main

import (
	"github.com/mackerelio/mackerel-agent/config"
)

import (
	"testing"
)

func TestLoadHostIdFromConfig(t *testing.T) {
	config.DefaultConfig.Conffile = "test/mackerel-agent.conf"

	hostId := LoadHostIdFromConfig()
	if hostId == "" {
		t.Error("should not empty")
	}

	if hostId != "9876ABCD" {
		debug(hostId)
		t.Error("should be 9876ABCD")
	}
}
