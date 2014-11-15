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

	apiKey := LoadApikeyFromConfigOrEnv()

	if apiKey != "123456ABCD" {
		t.Error("should be 123456ABCD")
	}

	os.Setenv("MACKEREL_APIKEY", "ENV123456ABCD")

	apiKey = LoadApikeyFromConfigOrEnv()

	if apiKey != "ENV123456ABCD" {
		t.Error("should be ENV123456ABCD")
	}

	os.Setenv("MACKEREL_APIKEY", "")
}

func TestLoadHostIdFromConfig(t *testing.T) {
	config.DefaultConfig.Conffile = "test/mackerel-agent.conf"

	hostId := LoadHostIdFromConfig()

	if hostId == "" {
		t.Error("should not empty")
	}

	if hostId != "9876ABCD" {
		t.Error("should be 9876ABCD")
	}
}
