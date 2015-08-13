package main

import (
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
