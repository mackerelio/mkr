package main

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/mackerelio/mackerel-client-go"
	"github.com/mackerelio/mkr/logger"
	"gopkg.in/yaml.v2"
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

type graphsConfig struct {
	Title           string             `yaml:"title"`
	URLPath         string             `yaml:"url_path"`
	GraphType       string             `yaml:"graph_type"`
	Height          int                `yaml:"height"`
	Width           int                `yaml:"width"`
	HostGraphFormat []*hostGraphFormat `yaml:"host_graphs"`
	GraphFormat     []*graphFormat     `yaml:"graphs"`
}

type hostGraphFormat struct {
	Headline   string   `yaml:"headline"`
	HostIDs    []string `yaml:"host_ids"`
	GraphNames []string `yaml:"graph_names"`
	Period     string   `yaml:"period"`
}
type graphFormat struct {
	Headline    string      `yaml:"headline"`
	ColumnCount int         `yaml:"column_count"`
	GraphDefs   []*graphDef `yaml:"graph_def"`
}
type graphDef struct {
	HostID      string `yaml:"host_id"`
	ServiceName string `yaml:"service_name"`
	RoleName    string `yaml:"role_name"`
	Query       string `yaml:"query"`
	GraphName   string `yaml:"graph_name"`
	Period      string `yaml:"period"`
	Stacked     bool   `yaml:"stacked"`
	Simplified  bool   `yaml:"simplified"`
}

func (g graphDef) isHostGraph() bool {
	return g.HostID != ""
}
func (g graphDef) isServiceGraph() bool {
	return g.ServiceName != "" && g.RoleName == ""
}
func (g graphDef) isRoleGraph() bool {
	return g.ServiceName != "" && g.RoleName != ""
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

type serviceGraph struct {
	ServiceName string
	GraphType   string
	Graph       string
	Period      string
	height      int
	width       int
}

func (s serviceGraph) getURL(orgName string, isImage bool) string {
	extension := ""
	if isImage {
		extension = ".png"
	}
	u, _ := url.Parse(fmt.Sprintf("https://mackerel.io/embed/orgs/%s/services/%s"+extension, orgName, s.ServiceName))
	param := url.Values{}
	param.Add("graph", s.Graph)
	param.Add("period", s.Period)
	u.RawQuery = param.Encode()
	return u.String()
}
func (s serviceGraph) generateGraphString(orgName string) string {
	if s.GraphType == "iframe" {
		return makeIframeTag(orgName, s)
	}
	return makeImageMarkdown(orgName, s)
}
func (s serviceGraph) getHeight() int {
	return s.height
}
func (s serviceGraph) getWidth() int {
	return s.width
}

type roleGraph struct {
	ServiceName string
	RoleName    string
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
	u, _ := url.Parse(fmt.Sprintf("https://mackerel.io/embed/orgs/%s/services/%s/%s"+extension, orgName, r.ServiceName, r.RoleName))
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
	return fmt.Sprintf("[![graph](%s)](%s)", g.getURL(orgName, true), g.getURL(orgName, true))
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

	yml := graphsConfig{}
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
	if yml.GraphType == "" {
		yml.GraphType = "iframe"
	}
	if yml.GraphType != "iframe" && yml.GraphType != "image" {
		logger.Log("error", "graph_type should be 'iframe' or 'image'.")
		os.Exit(1)
	}
	if yml.Height == 0 {
		yml.Height = 200
	}
	if yml.Width == 0 {
		yml.Width = 400
	}

	if yml.HostGraphFormat != nil && yml.GraphFormat != nil {
		logger.Log("error", "you cannot specify both 'graphs' and host_graphs'.")
	}

	var markdown string
	for _, h := range yml.HostGraphFormat {
		markdown += generateHostGraphsMarkdown(org.Name, h, yml.GraphType, yml.Height, yml.Width, client)
	}
	for _, g := range yml.GraphFormat {
		markdown += generateGraphsMarkdown(org.Name, g, yml.GraphType, yml.Height, yml.Width)
	}

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

func generateHostGraphsMarkdown(orgName string, hostGraphs *hostGraphFormat, graphType string, height int, width int, client *mackerel.Client) string {

	var markdown string

	if hostGraphs.Headline != "" {
		markdown += fmt.Sprintf("## %s\n", hostGraphs.Headline)
	}
	markdown += generateHostGraphsTableHeader(hostGraphs.HostIDs, client)

	if hostGraphs.Period == "" {
		hostGraphs.Period = "1h"
	}

	for _, graphName := range hostGraphs.GraphNames {
		currentColumnCount := 0
		for _, hostID := range hostGraphs.HostIDs {
			h := &hostGraph{
				hostID,
				graphType,
				graphName,
				hostGraphs.Period,
				height,
				width,
			}
			markdown = appendMarkdown(markdown, h.generateGraphString(orgName), len(hostGraphs.HostIDs))
			currentColumnCount++
			if currentColumnCount >= len(hostGraphs.HostIDs) {
				markdown += "|\n"
				currentColumnCount = 0
			}
		}
	}

	return markdown
}

func generateHostGraphsTableHeader(hostIDs []string, client *mackerel.Client) string {
	var header string
	for _, hostID := range hostIDs {
		host, err := client.FindHost(hostID)
		logger.DieIf(err)

		var hostName string
		if host.DisplayName != "" {
			hostName = host.DisplayName
		} else {
			hostName = host.Name
		}

		header += "|" + hostName
	}

	header += "|\n" + generateGraphTableHeader(len(hostIDs))

	return header
}

func generateGraphsMarkdown(orgName string, graphs *graphFormat, graphType string, height int, width int) string {

	var markdown string
	if graphs.ColumnCount == 0 {
		graphs.ColumnCount = 1
	}
	currentColumnCount := 0

	if graphs.Headline != "" {
		markdown += fmt.Sprintf("## %s\n", graphs.Headline)
	}
	markdown += generateGraphTableHeader(graphs.ColumnCount)

	for i, gd := range graphs.GraphDefs {
		var graphDefCount = 0
		if gd.isHostGraph() {
			graphDefCount++
		}
		if gd.isServiceGraph() {
			graphDefCount++
		}
		if gd.isRoleGraph() {
			graphDefCount++
		}
		if gd.isExpressionGraph() {
			graphDefCount++
		}
		if graphDefCount != 1 {
			logger.Log("error", "at least one between hostId, service_name and query is required.")
			os.Exit(1)
		}

		if gd.Period == "" {
			gd.Period = "1h"
		}

		if gd.isHostGraph() {
			if gd.GraphName == "" {
				logger.Log("error", "graph_name is required for host graph.")
				os.Exit(1)
			}

			h := &hostGraph{
				gd.HostID,
				graphType,
				gd.GraphName,
				gd.Period,
				height,
				width,
			}
			markdown = appendMarkdown(markdown, h.generateGraphString(orgName), graphs.ColumnCount)
		}

		if gd.isServiceGraph() {
			if gd.GraphName == "" {
				logger.Log("error", "graph_name is required for service graph.")
				os.Exit(1)
			}

			h := &serviceGraph{
				gd.ServiceName,
				graphType,
				gd.GraphName,
				gd.Period,
				height,
				width,
			}
			markdown = appendMarkdown(markdown, h.generateGraphString(orgName), graphs.ColumnCount)
		}

		if gd.isRoleGraph() {
			if gd.GraphName == "" {
				logger.Log("error", "graph_name is required for role graph.")
				os.Exit(1)
			}

			r := &roleGraph{
				gd.ServiceName,
				gd.RoleName,
				graphType,
				gd.GraphName,
				gd.Period,
				gd.Stacked,
				gd.Simplified,
				height,
				width,
			}
			markdown = appendMarkdown(markdown, r.generateGraphString(orgName), graphs.ColumnCount)
		}

		if gd.isExpressionGraph() {
			e := &expressionGraph{
				gd.Query,
				graphType,
				gd.Period,
				height,
				width,
			}
			markdown = appendMarkdown(markdown, e.generateGraphString(orgName), graphs.ColumnCount)
		}

		currentColumnCount++
		if currentColumnCount >= graphs.ColumnCount || i >= len(graphs.GraphDefs)-1 {
			markdown += "|\n"
			currentColumnCount = 0
		}
	}

	return markdown
}

func generateGraphTableHeader(confColumnCount int) string {
	return strings.Repeat("|:-:", confColumnCount) + "|\n"
}

func appendMarkdown(markdown string, addItem string, confColumnCount int) string {
	return markdown + "|" + addItem
}
