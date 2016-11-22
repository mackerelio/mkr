package main

import (
	"io/ioutil"
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
	const want = ` {
   "name": "foo",
-  "responseTimeCritical": 1000,
   "service": "bar",
   "type": "external",
   "url": "http://example.com"
 },`
	a := &mkr.Monitor{ID: "12345", Name: "foo", Type: "external", URL: "http://example.com", Service: "bar", ResponseTimeCritical: 1000}
	b := &mkr.Monitor{ID: "12345", Name: "foo", Type: "external", URL: "http://example.com", Service: "bar"}
	if got := diffMonitor(a, b); got != want {
		t.Errorf("diffMonitor: got\n%s\nwant \n%s", got, want)
	}
}

func TestMonitorSaveRules(t *testing.T) {
	tmpFile, err := ioutil.TempFile("", "")
	if err != nil {
		t.Errorf("should not raise error: %v", err)
	}
	a := &mkr.Monitor{
		ID:                   "12345",
		Name:                 "foo",
		Type:                 "external",
		URL:                  "http://example.com",
		Service:              "bar",
		ResponseTimeCritical: 1000,
	}
	monitorSaveRules([]*mkr.Monitor{a}, tmpFile.Name())

	byt, _ := ioutil.ReadFile(tmpFile.Name())
	content := string(byt)
	expected := `{
    "monitors": [
        {
            "id": "12345",
            "name": "foo",
            "type": "external",
            "url": "http://example.com",
            "service": "bar",
            "responseTimeCritical": 1000
        }
    ]
}
`
	if content != expected {
		t.Errorf("content should be:\n %s, but:\n %s", expected, content)
	}
}

func TestStringifyMonitor(t *testing.T) {
	a := &mkr.Monitor{ID: "12345", Name: "foo", Type: "connectivity"}
	expected := `+{
+  "id": "12345",
+  "name": "foo",
+  "type": "connectivity"
+},`

	r := stringifyMonitor(a, "+")
	if r != expected {
		t.Errorf("stringifyMonitor should be:\n%s\nbut:\n%s", expected, r)
	}
}

func TestDiffMonitorsWithScopes(t *testing.T) {
	a := &mkr.Monitor{
		ID:   "12345",
		Name: "foo",
		Type: "connectivity",
	}
	b := &mkr.Monitor{
		ID:     "12345",
		Name:   "foo",
		Type:   "connectivity",
		Scopes: []string{"sss: notebook"},
	}
	diff := diffMonitor(a, b)
	expected := ` {
   "name": "foo",
   "type": "connectivity"
+  "scopes": [
+    "sss: notebook"
+  ]
 },`
	if diff != expected {
		// t.Error(debugdiff.Diff(expected, diff))
		t.Errorf("expected:\n%s\n, output:\n%s\n", expected, diff)
	}

	diff = diffMonitor(b, a)
	expected = ` {
   "name": "foo",
-  "scopes": [
-    "sss: notebook"
-  ],
   "type": "connectivity"
 },`
	if diff != expected {
		t.Errorf("expected:\n%s\n, output:\n%s\n", expected, diff)
	}

	c := &mkr.Monitor{
		ID:     "12345",
		Name:   "foo",
		Type:   "connectivity",
		Scopes: []string{"sss: notebook", "ttt: notebook"},
	}

	diff = diffMonitor(b, c)
	expected = ` {
   "name": "foo",
   "scopes": [
     "sss: notebook"
+    "ttt: notebook"
   ],
   "type": "connectivity"
 },`
	if diff != expected {
		t.Errorf("expected:\n%s\n, output:\n%s\n", expected, diff)
	}

	d := &mkr.Monitor{
		ID:     "12345",
		Name:   "foo",
		Type:   "connectivity",
		Scopes: []string{"ttt: notebook"},
	}
	diff = diffMonitor(b, d)
	expected = ` {
   "name": "foo",
   "scopes": [
-    "sss: notebook"
+    "ttt: notebook"
   ],
   "type": "connectivity"
 },`
	if diff != expected {
		t.Errorf("expected:\n%s\n, output:\n%s\n", expected, diff)
	}
}
