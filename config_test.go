package main

import "os"

import (
	"testing"
)

func TestLoadApibaseFromConfig(t *testing.T) {
	conffile := "test/mackerel-agent.conf"

	apiBase := LoadApibaseFromConfig(conffile)

	if apiBase != "https://example.com/" {
		t.Error("should be https://example.com/")
	}
}

func TestLoadApikeyFromConfig(t *testing.T) {
	conffile := "test/mackerel-agent.conf"

	apiKey := LoadApikeyFromConfig(conffile)

	if apiKey != "123456ABCD" {
		t.Error("should be 123456ABCD")
	}
}

func TestLoadApikeyFromConfigOrEnv(t *testing.T) {
	os.Setenv("MACKEREL_APIKEY", "")

	conffile := "test/mackerel-agent.conf"

	apiKey := LoadApikeyFromEnvOrConfig(conffile)

	if apiKey != "123456ABCD" {
		t.Error("should be 123456ABCD")
	}

	os.Setenv("MACKEREL_APIKEY", "ENV123456ABCD")

	apiKey = LoadApikeyFromEnvOrConfig(conffile)

	if apiKey != "ENV123456ABCD" {
		t.Error("should be ENV123456ABCD")
	}

	os.Setenv("MACKEREL_APIKEY", "")
}

func TestLoadHostIDFromConfig(t *testing.T) {
	conffile := "test/mackerel-agent.conf"

	hostID := LoadHostIDFromConfig(conffile)

	if hostID == "" {
		t.Error("should not be empty")
	}

	if hostID != "9876ABCD" {
		t.Error("should be 9876ABCD")
	}
}
