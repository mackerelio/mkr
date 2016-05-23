package main

import (
	"testing"
)

func TestHostIFrameGraph(t *testing.T) {
	h := &hostGraph{
		"hostid",
		"iframe",
		"loadavg5",
		"30m",
		200,
		600,
	}

	actual := h.generateGraphString("orgname")
	expected := `<iframe src="https://mackerel.io/embed/orgs/orgname/hosts/hostid?graph=loadavg5&period=30m" height="200" width="600" frameborder="0"></iframe>`

	if actual != expected {
		t.Errorf("output should be:\n%s\nbut:\n%s", expected, actual)
	}
}

func TestServiceIFrameGraph(t *testing.T) {
	r := &roleGraph{
		"hoge",
		"api",
		"iframe",
		"cpu",
		"6h",
		true,
		true,
		200,
		600,
	}

	actual := r.generateGraphString("orgname")
	expected := `<iframe src="https://mackerel.io/embed/orgs/orgname/services/hoge/api?graph=cpu&period=6h&simplified=true&stacked=true" height="200" width="600" frameborder="0"></iframe>`

	if actual != expected {
		t.Errorf("output should be:\n%s\nbut:\n%s", expected, actual)
	}
}

func TestExpressionIFrameGraph(t *testing.T) {
	e := &expressionGraph{
		"max(roleSlots('hoge:api','loadavg5'))",
		"iframe",
		"6h",
		200,
		600,
	}

	actual := e.generateGraphString("orgname")
	expected := `<iframe src="https://mackerel.io/embed/orgs/orgname/advanced-graph?period=6h&query=max%28roleSlots%28%27hoge%3Aapi%27%2C%27loadavg5%27%29%29" height="200" width="600" frameborder="0"></iframe>`

	if actual != expected {
		t.Errorf("output should be:\n%s\nbut:\n%s", expected, actual)
	}
}

func TestHostImageGraph(t *testing.T) {
	h := &hostGraph{
		"hostid",
		"image",
		"loadavg5",
		"30m",
		200,
		600,
	}

	actual := h.generateGraphString("orgname")
	expected := `[![graph](https://mackerel.io/embed/orgs/orgname/hosts/hostid.png?graph=loadavg5&period=30m)]()`

	if actual != expected {
		t.Errorf("output should be:\n%s\nbut:\n%s", expected, actual)
	}
}

func TestServiceImageGraph(t *testing.T) {
	r := &roleGraph{
		"hoge",
		"api",
		"image",
		"cpu",
		"6h",
		true,
		true,
		200,
		600,
	}

	actual := r.generateGraphString("orgname")
	expected := `[![graph](https://mackerel.io/embed/orgs/orgname/services/hoge/api.png?graph=cpu&period=6h&simplified=true&stacked=true)]()`

	if actual != expected {
		t.Errorf("output should be:\n%s\nbut:\n%s", expected, actual)
	}
}

func TestExpressionImageGraph(t *testing.T) {
	e := &expressionGraph{
		"max(roleSlots('hoge:api','loadavg5'))",
		"image",
		"6h",
		200,
		600,
	}

	actual := e.generateGraphString("orgname")
	expected := `[![graph](https://mackerel.io/embed/orgs/orgname/advanced-graph.png?period=6h&query=max%28roleSlots%28%27hoge%3Aapi%27%2C%27loadavg5%27%29%29)]()`

	if actual != expected {
		t.Errorf("output should be:\n%s\nbut:\n%s", expected, actual)
	}
}

func TestGenerateMarkDown(t *testing.T) {
	arr := []*graphDef{
		&graphDef{
			ServiceName: "hoge",
			RollName:    "api",
			Graph:       "cpu",
		},
		&graphDef{
			HostID:      "abcde",
			Graph:       "cpu",
		},
		&graphDef{
			Query:       "max(roleSlots('hoge:api','loadavg5'))",
			Graph:       "cpu",
		},
	}

	actual := generateMarkDown("orgname", arr, 2)
	expected := "|:-:|:-:|\n" +
	            `|<iframe src="https://mackerel.io/embed/orgs/orgname/services/hoge/api?graph=cpu&period=1h&simplified=false&stacked=false" height="200" width="400" frameborder="0"></iframe>|<iframe src="https://mackerel.io/embed/orgs/orgname/hosts/abcde?graph=cpu&period=1h" height="200" width="400" frameborder="0"></iframe>|` + "\n" +
		    `|<iframe src="https://mackerel.io/embed/orgs/orgname/advanced-graph?period=1h&query=max%28roleSlots%28%27hoge%3Aapi%27%2C%27loadavg5%27%29%29" height="200" width="400" frameborder="0"></iframe>|` + "\n"


	if actual != expected {
		t.Errorf("output should be:\n%s\nbut:\n%s", expected, actual)
	}
}
