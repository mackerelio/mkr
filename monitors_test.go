package main

import (
	"strings"
	"testing"

	mkr "github.com/mackerelio/mackerel-client-go"
)

func TestIsSameMonitor(t *testing.T) {
	a := &mkr.Monitor{ID: "12345", Name: "foo", Type: "connectivity"}
	b := &mkr.Monitor{Name: "foo", Type: "connectivity"}

	_, ret := isSameMonitor(a, b, true)
	if ret != true {
		t.Error("should recognize same monitors")
	}

	_, ret = isSameMonitor(a, b, false)
	if ret == true {
		t.Error("should not recognize same monitors")
	}
}

func TestValidateRoles(t *testing.T) {
	a := &mkr.Monitor{ID: "12345", Name: "foo", Type: "connectivity"}

	ret, err := validateRules([](*mkr.Monitor){a}, "test monitor")
	if ret != true {
		t.Errorf("should validate the rule: %s", err.Error())
	}

}

func TestDiffMonitors(t *testing.T) {
	a := &mkr.Monitor{ID: "12345", Name: "foo", Type: "external", URL: "http://example.com", Service: "bar", ResponseTimeCritical: 1000}
	b := &mkr.Monitor{ID: "12345", Name: "foo", Type: "external", URL: "http://example.com", Service: "bar"}

	ret := diffMonitor(a, b)

	zz := []string{
		"  {",
		"   \"name\": \"foo\",",
		"   \"type\": \"external\",",
		"   \"url\": \"http://example.com\",",
		"   \"service\": \"bar\",",
		"-  \"responseTimeCritical\": 1000.000000,",
		"+  \"responseTimeCritical\": 0.000000,",
		"  },",
	}

	if ret != strings.Join(zz, "\n") {
		t.Errorf("should validate the rule: %s\n%s", ret, strings.Join(zz, "\n"))
	}

}
