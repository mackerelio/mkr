package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/mackerelio/mackerel-client-go"
	"github.com/mackerelio/mkr/logger"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/url"
	"os"
	"strconv"
	"strings"
)

var commandDashboards = cli.Command{
	Name: "dashboards",
	Subcommands: []cli.Command{
		{
			Name:  "generate",
			Usage: "Generate custom dashboard",
			Description: `
    A custom dashboard is registered from a yaml file..
    Requests "POST /api/v0/dashboards". See https://mackerel.io/ja/api-docs/entry/dashboards#create.
`,
			Action: doGenerateDashboards,
			Flags: []cli.Flag{
				cli.BoolFlag{Name: "print, p", Usage: "markdown is output in standard output."},
			},
		},
	},
}

type configYAML struct {
	Title       string      `yaml:"title"`
	URLPath     string      `yaml:"url_path"`
	ColumnCount int         `yaml:"column_count"`
	Graphs      []*graphDef `yaml:"graphs"`
}
type graphDef struct {
	HostID      string `yaml:"host_id"`
	ServiceName string `yaml:"service_name"`
	RollName    string `yaml:"roll_name"`
	Query       string `yaml:"query"`
	Graph       string `yaml:"graph"`
	GraphType   string `yaml:"graph_type"`
	Period      string `yaml:"period"`
	Stacked     bool   `yaml:"stacked"`
	Simplified  bool   `yaml:"simplified"`
	Height      int    `yaml:"height"`
	Width       int    `yaml:"width"`
}

func (g graphDef) isHostGraph() bool {
	return g.HostID != ""
}
func (g graphDef) isRollGraph() bool {
	return g.ServiceName != "" && g.RollName != ""
}
func (g graphDef) isExpressionGraph() bool {
	return g.Query != ""
}

type baseGraph interface {
	getURL(string, bool) string
	getHeight() int
	getWidth() int
}

type hostGraph struct {
	HostID    string
	GraphType string
	Graph     string
	Period    string
	height    int
	width     int
}

func (h hostGraph) getURL(orgName string, isImage bool) string {
	extension := ""
	if isImage {
		extension = ".png"
	}
	u, _ := url.Parse(fmt.Sprintf("https://mackerel.io/embed/orgs/%s/hosts/%s"+extension, orgName, h.HostID))
	param := url.Values{}
	param.Add("graph", h.Graph)
	param.Add("period", h.Period)
	u.RawQuery = param.Encode()
	return u.String()
}
func (h hostGraph) generateGraphString(orgName string) string {
	if h.GraphType == "iframe" {
		return makeIframeTag(orgName, h)
	}
	return makeImageMarkdown(orgName, h)
}
func (h hostGraph) getHeight() int {
	return h.height
}
func (h hostGraph) getWidth() int {
	return h.width
}

type roleGraph struct {
	ServiceName string
	RollName    string
	GraphType   string
	Graph       string
	Period      string
	Stacked     bool
	Simplified  bool
	height      int
	width       int
}

func (r roleGraph) getURL(orgName string, isImage bool) string {
	extension := ""
	if isImage {
		extension = ".png"
	}
	u, _ := url.Parse(fmt.Sprintf("https://mackerel.io/embed/orgs/%s/services/%s/%s"+extension, orgName, r.ServiceName, r.RollName))
	param := url.Values{}
	param.Add("graph", r.Graph)
	param.Add("stacked", strconv.FormatBool(r.Stacked))
	param.Add("simplified", strconv.FormatBool(r.Simplified))
	param.Add("period", r.Period)
	u.RawQuery = param.Encode()
	return u.String()
}
func (r roleGraph) generateGraphString(orgName string) string {
	if r.GraphType == "iframe" {
		return makeIframeTag(orgName, r)
	}
	return makeImageMarkdown(orgName, r)
}
func (r roleGraph) getHeight() int {
	return r.height
}
func (r roleGraph) getWidth() int {
	return r.width
}

type expressionGraph struct {
	Query     string
	GraphType string
	Period    string
	height    int
	width     int
}

func (e expressionGraph) getURL(orgName string, isImage bool) string {
	extension := ""
	if isImage {
		extension = ".png"
	}
	u, _ := url.Parse(fmt.Sprintf("https://mackerel.io/embed/orgs/%s/advanced-graph"+extension, orgName))
	param := url.Values{}
	param.Add("query", e.Query)
	param.Add("period", e.Period)
	u.RawQuery = param.Encode()
	return u.String()
}
func (e expressionGraph) generateGraphString(orgName string) string {
	if e.GraphType == "iframe" {
		return makeIframeTag(orgName, e)
	}
	return makeImageMarkdown(orgName, e)
}
func (e expressionGraph) getHeight() int {
	return e.height
}
func (e expressionGraph) getWidth() int {
	return e.width
}

func makeIframeTag(orgName string, g baseGraph) string {
	return fmt.Sprintf(`<iframe src="%s" height="%d" width="%d" frameborder="0"></iframe>`, g.getURL(orgName, false), g.getHeight(), g.getWidth())
}

func makeImageMarkdown(orgName string, g baseGraph) string {
	return fmt.Sprintf("[![graph](%s)]()", g.getURL(orgName, true))
}

func doGenerateDashboards(c *cli.Context) error {
	conffile := c.GlobalString("conf")

	isStdout := c.Bool("print")

	argFilePath := c.Args()
	if len(argFilePath) < 1 {
		logger.Log("error", "at least one argumet is required.")
		cli.ShowCommandHelp(c, "generate")
		os.Exit(1)
	}

	buf, err := ioutil.ReadFile(argFilePath[0])
	logger.DieIf(err)

	yml := configYAML{}
	err = yaml.Unmarshal(buf, &yml)
	logger.DieIf(err)

	client := newMackerel(conffile)

	org, err := client.GetOrg()
	logger.DieIf(err)

	if yml.Title == "" {
		logger.Log("error", "title is required in yaml.")
		os.Exit(1)
	}
	if yml.URLPath == "" {
		logger.Log("error", "url_path is required in yaml.")
		os.Exit(1)
	}
	if yml.ColumnCount == 0 {
		yml.ColumnCount = 1
	}

	markdown := generateMarkDown(org.Name, yml.Graphs, yml.ColumnCount)

	if isStdout {
		fmt.Println(markdown)
	} else {
		updateDashboard := &mackerel.Dashboard{
			Title:        yml.Title,
			BodyMarkDown: markdown,
			URLPath:      yml.URLPath,
		}

		dashboards, fetchError := client.FindDashboards()
		logger.DieIf(fetchError)

		dashboardID := ""
		for _, ds := range dashboards {
			if ds.URLPath == yml.URLPath {
				dashboardID = ds.ID
			}
		}

		if dashboardID == "" {
			_, createError := client.CreateDashboard(updateDashboard)
			logger.DieIf(createError)
		} else {
			_, updateError := client.UpdateDashboard(dashboardID, updateDashboard)
			logger.DieIf(updateError)
		}
	}

	return nil
}

func generateMarkDown(orgName string, graphs []*graphDef, confColumnCount int) string {

	var markdown string
	var currentColumnCount = 0

	for i, g := range graphs {

		var graphDefCount = 0
		if g.isHostGraph() {
			graphDefCount++
		}
		if g.isRollGraph() {
			graphDefCount++
		}
		if g.isExpressionGraph() {
			graphDefCount++
		}
		if graphDefCount != 1 {
			logger.Log("error", "at least one between hostId, (service_name and roll_name) and query is required.")
			os.Exit(1)
		}

		if g.GraphType == "" {
			g.GraphType = "iframe"
		}
		if g.GraphType != "iframe" && g.GraphType != "image" {
			logger.Log("error", "graph_type should 'iframe' or 'image'.")
			os.Exit(1)
		}

		if g.Height == 0 {
			g.Height = 200
		}
		if g.Width == 0 {
			g.Width = 400
		}

		if g.Period == "" {
			g.Period = "1h"
		}

		if g.isHostGraph() {
			if g.Graph == "" {
				logger.Log("error", "graph is required for host graph.")
				os.Exit(1)
			}

			h := &hostGraph{
				g.HostID,
				g.GraphType,
				g.Graph,
				g.Period,
				g.Height,
				g.Width,
			}
			markdown = appendMarkdown(markdown, h.generateGraphString(orgName), confColumnCount)
		}

		if g.isRollGraph() {
			if g.Graph == "" {
				logger.Log("error", "graph is required for roll graph.")
				os.Exit(1)
			}

			r := &roleGraph{
				g.ServiceName,
				g.RollName,
				g.GraphType,
				g.Graph,
				g.Period,
				g.Stacked,
				g.Simplified,
				g.Height,
				g.Width,
			}
			markdown = appendMarkdown(markdown, r.generateGraphString(orgName), confColumnCount)
		}

		if g.isExpressionGraph() {
			e := &expressionGraph{
				g.Query,
				g.GraphType,
				g.Period,
				g.Height,
				g.Width,
			}
			markdown = appendMarkdown(markdown, e.generateGraphString(orgName), confColumnCount)
		}

		currentColumnCount++
		if currentColumnCount >= confColumnCount || i >= len(graphs)-1 {
			if strings.HasPrefix(markdown, "|") {
				markdown += "|"
			}
			markdown += "\n"
			currentColumnCount = 0
		}
	}

	return generateTableHeader(confColumnCount) + markdown
}

func appendMarkdown(markdown string, addItem string, confColumnCount int) string {
	if confColumnCount == 1 {
		return markdown + addItem
	}
	return markdown + "|" + addItem
}

func generateTableHeader(confColumnCount int) string {
	header := ""
	if confColumnCount > 1 {
		for i := 0; i < confColumnCount; i++ {
			header += "|:-:"
		}
		return header + "|\n"
	}
	return header
}
