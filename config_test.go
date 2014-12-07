package main

import (
	"os"

	"github.com/mackerelio/mackerel-agent/config"
)

import (
	"testing"
)

func TestLoadApikeyFromConfig(t *testing.T) {
	config.DefaultConfig.Conffile = "test/mackerel-agent.conf"

	apiKey := LoadApikeyFromConfig()

	if apiKey != "123456ABCD" {
		t.Error("should be 123456ABCD")
	}
}

func TestLoadApikeyFromConfigOrEnv(t *testing.T) {
	os.Setenv("MACKEREL_APIKEY", "")

	config.DefaultConfig.Conffile = "test/mackerel-agent.conf"

	apiKey := LoadApikeyFromEnvOrConfig()

	if apiKey != "123456ABCD" {
		t.Error("should be 123456ABCD")
	}

	os.Setenv("MACKEREL_APIKEY", "ENV123456ABCD")

	apiKey = LoadApikeyFromEnvOrConfig()

	if apiKey != "ENV123456ABCD" {
		t.Error("should be ENV123456ABCD")
	}

	os.Setenv("MACKEREL_APIKEY", "")
}

func TestLoadHostIDFromConfig(t *testing.T) {
	config.DefaultConfig.Conffile = "test/mackerel-agent.conf"

	hostID := LoadHostIDFromConfig()

	if hostID == "" {
		t.Error("should not be empty")
	}

	if hostID != "9876ABCD" {
		t.Error("should be 9876ABCD")
	}
}
