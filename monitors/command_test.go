package monitors

import (
	"os"
	"testing"

	"github.com/mackerelio/mackerel-client-go"
)

func TestIsSameMonitor(t *testing.T) {
	a := &mackerel.MonitorConnectivity{ID: "12345", Name: "foo", Type: "connectivity"}
	b := &mackerel.MonitorConnectivity{Name: "foo", Type: "connectivity"}

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
	t.Run("connectivitiy", func(t *testing.T) {
		a := &mackerel.MonitorConnectivity{ID: "12345", Name: "foo", Type: "connectivity"}

		ret, err := validateRules([](mackerel.Monitor){a}, "test monitor")
		if ret != true {
			t.Errorf("should validate the rule: %s", err.Error())
		}
	})

	t.Run("valid anomalyDetection", func(t *testing.T) {
		a := &mackerel.MonitorAnomalyDetection{ID: "12345", Name: "anomaly", Type: "anomalyDetection", WarningSensitivity: "sensitive", Scopes: []string{"MyService: MyRole"}}

		ret, err := validateRules([](mackerel.Monitor){a}, "anomaly detection monitor")
		if ret != true {
			t.Errorf("should validate the rule: %s", err.Error())
		}
	})

	t.Run("invalid anomalyDetection", func(t *testing.T) {
		a := &mackerel.MonitorAnomalyDetection{ID: "12345", Name: "anomaly", Type: "anomalyDetection", WarningSensitivity: "sensitive"}

		ret, err := validateRules([](mackerel.Monitor){a}, "anomaly detection monitor")
		if ret == true || err == nil {
			t.Error("should invalidate the rule")
		}
	})

	t.Run("valid query monitoring rule", func(t *testing.T) {
		a := &mackerel.MonitorQuery{Name: "name", Type: "query", Query: "http.monitor.count", Operator: "<"}

		ret, err := validateRules([](mackerel.Monitor){a}, "query monitor")
		if !ret {
			t.Errorf("should validate the rule: %v", err)
		}
	})

	t.Run("invalid query monitoring rule", func(t *testing.T) {
		a := &mackerel.MonitorQuery{Name: "name", Type: "query", Operator: "<"}

		ret, err := validateRules([](mackerel.Monitor){a}, "query monitor")
		if ret == true || err == nil {
			t.Error("should invalidate the rule")
		}
	})
}

func pfloat64(x float64) *float64 {
	return &x
}

func TestDiffMonitors(t *testing.T) {
	const want = ` {
   "headers": [
   ],
   "name": "foo",
-  "responseTimeCritical": 1000,
   "service": "bar",
   "type": "external",
   "url": "http://example.com"
 },`
	a := &mackerel.MonitorExternalHTTP{ID: "12345", Name: "foo", Type: "external", URL: "http://example.com", Service: "bar", ResponseTimeCritical: pfloat64(1000), Headers: []mackerel.HeaderField{}}
	b := &mackerel.MonitorExternalHTTP{ID: "12345", Name: "foo", Type: "external", URL: "http://example.com", Service: "bar", Headers: []mackerel.HeaderField{}}
	if got := diffMonitor(a, b); got != want {
		t.Errorf("diffMonitor: got\n%s\nwant \n%s", got, want)
	}
}

func TestMonitorSaveRules(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "")
	if err != nil {
		t.Errorf("should not raise error: %v", err)
	}
	a := &mackerel.MonitorExternalHTTP{
		ID:                   "12345",
		Name:                 "foo",
		Type:                 "external",
		URL:                  "http://example.com",
		Service:              "bar",
		ResponseTimeCritical: pfloat64(1000),
		Headers:              []mackerel.HeaderField{},
	}
	if err := monitorSaveRules([]mackerel.Monitor{a}, tmpFile.Name()); err != nil {
		t.Fatal(err)
	}

	byt, _ := os.ReadFile(tmpFile.Name())
	content := string(byt)
	expected := `{
    "monitors": [
        {
            "id": "12345",
            "name": "foo",
            "type": "external",
            "url": "http://example.com",
            "service": "bar",
            "responseTimeCritical": 1000,
            "headers": []
        }
    ]
}
`
	if content != expected {
		t.Errorf("content should be:\n %s, but:\n %s", expected, content)
	}
}

func TestStringifyMonitor(t *testing.T) {
	a := &mackerel.MonitorConnectivity{ID: "12345", Name: "foo", Type: "connectivity"}
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
	a := &mackerel.MonitorConnectivity{
		ID:   "12345",
		Name: "foo",
		Type: "connectivity",
	}
	b := &mackerel.MonitorConnectivity{
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

	c := &mackerel.MonitorConnectivity{
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

	d := &mackerel.MonitorConnectivity{
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

func TestMonitorLoadRulesWithBOM(t *testing.T) {
	// XXX: t.TempDir is better, but it will cause "TempDir RemoveAll cleanup: remove C:\...\monitors.json: The process cannot access the file because it is being used by another process." error on Windows
	tmpFile, err := os.CreateTemp("", "")
	if err != nil {
		t.Errorf("should not raise error: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	json := `{"monitors": []}`

	_, err = tmpFile.WriteString(json)
	if err != nil {
		t.Errorf("should not raise error: %v", err)
	}
	_, err = monitorLoadRules(tmpFile.Name())
	if err != nil {
		t.Error("should accept JSON content no BOM")
	}

	utf8bom := "\xef\xbb\xbf"
	_, err = tmpFile.Seek(0, 0)
	if err != nil {
		t.Errorf("should not raise error: %v", err)
	}
	_, err = tmpFile.WriteString(utf8bom + json)
	if err != nil {
		t.Errorf("should not raise error: %v", err)
	}
	_, err = monitorLoadRules(tmpFile.Name())
	if err != nil {
		t.Error("should accept JSON content with BOM")
	}
}
