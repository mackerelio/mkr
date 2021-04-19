package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/mackerelio/mackerel-client-go"
	"github.com/mackerelio/mkr/format"
	"github.com/mackerelio/mkr/logger"
	"github.com/mackerelio/mkr/mackerelclient"
	"github.com/urfave/cli"
)

var commandDashboards = cli.Command{
	Name:  "dashboards",
	Usage: "Generating custom dashboards",
	Description: `
    Manipulate custom dashboards. With no subcommand specified, this will show all dashboards. See https://mackerel.io/docs/entry/advanced/cli
`,
	Action: doListDashboards,
	Subcommands: []cli.Command{
		{
			Name:  "pull",
			Usage: "Pull custom dashboards",
			Description: `
	Pull custom dashboards from Mackerel server and output these to local files.
`,
			Action: doPullDashboard,
		},
		{
			Name:      "push",
			Usage:     "Push custom dashboard",
			ArgsUsage: "--file-path | F <file>",
			Description: `
	Push custom dashboards to Mackerel server from a specified file.
	When "id" is defined in the file, updates the dashboard.
	Otherwise creates a new dashboard.
`,
			Action: doPushDashboard,
			Flags: []cli.Flag{
				cli.StringFlag{Name: "file-path, F", Usage: "read dashboard from the file"},
			},
		},
	},
}

type graphsConfig struct {
	ConfigVersion   string             `yaml:"config_version"`
	Title           string             `yaml:"title"`
	URLPath         string             `yaml:"url_path"`
	Format          string             `yaml:"format"`
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
	GraphTitle  string `yaml:"title"`
	Unit        string `yaml:"unit"`
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

func (g graphDef) getBaseGraph(graphType string, height int, width int) (baseGraph baseGraph, err error) {
	if g.isHostGraph() {
		if g.GraphName == "" {
			return nil, cli.NewExitError("graph_name is required for host graph.", 1)
		}

		return hostGraph{
			g.HostID,
			graphType,
			g.GraphName,
			g.Period,
			height,
			width,
		}, nil
	}

	if g.isServiceGraph() {
		if g.GraphName == "" {
			return nil, cli.NewExitError("graph_name is required for service graph.", 1)
		}

		return serviceGraph{
			g.ServiceName,
			graphType,
			g.GraphName,
			g.Period,
			height,
			width,
		}, nil
	}

	if g.isRoleGraph() {
		if g.GraphName == "" {
			return nil, cli.NewExitError("graph_name is required for role graph.", 1)
		}

		return roleGraph{
			g.ServiceName,
			g.RoleName,
			graphType,
			g.GraphName,
			g.Period,
			g.Stacked,
			g.Simplified,
			height,
			width,
		}, nil
	}

	if g.isExpressionGraph() {
		return expressionGraph{
			g.Query,
			graphType,
			g.GraphTitle,
			g.Unit,
			g.Period,
			height,
			width,
		}, nil
	}

	return nil, cli.NewExitError("either host_id, service_name or query should be specified.", 1)
}

type baseGraph interface {
	getURL(string, bool) string
	getPermalink(string) string
	getHeight() int
	getWidth() int
	generateGraphString(orgName string) string
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
func (h hostGraph) getPermalink(orgName string) string {
	u, _ := url.Parse(fmt.Sprintf("https://mackerel.io/orgs/%s/hosts/%s/-/graphs/%s", orgName, h.HostID, url.QueryEscape(h.Graph)))
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
func (s serviceGraph) getPermalink(orgName string) string {
	u, _ := url.Parse(fmt.Sprintf("https://mackerel.io/orgs/%s/services/%s/-/graphs", orgName, s.ServiceName))
	param := url.Values{}
	param.Add("name", s.Graph)
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
func (r roleGraph) getPermalink(orgName string) string {
	u, _ := url.Parse(fmt.Sprintf("https://mackerel.io/orgs/%s/services/%s/%s/-/graph", orgName, r.ServiceName, r.RoleName))
	param := url.Values{}
	param.Add("name", r.Graph)
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
	Title     string
	Unit      string
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
	param.Add("title", e.Title)
	param.Add("unit", e.Unit)
	u.RawQuery = param.Encode()
	return u.String()
}
func (e expressionGraph) getPermalink(orgName string) string {
	u, _ := url.Parse(fmt.Sprintf("https://mackerel.io/orgs/%s/advanced-graph", orgName))
	param := url.Values{}
	param.Add("query", e.Query)
	param.Add("title", e.Title)
	param.Add("unit", e.Unit)
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

type markdownFactory struct {
	Headline    string
	TableHeader string
	BaseGraphs  []baseGraph
	ColumnCount int
}

func (mdf markdownFactory) generate(orgName string) string {
	markdown := ""
	if mdf.Headline != "" {
		markdown += fmt.Sprintf("## %s\n", mdf.Headline)
	}

	markdown += mdf.TableHeader

	for i, g := range mdf.BaseGraphs {
		markdown += "|" + g.generateGraphString(orgName)
		if i%mdf.ColumnCount >= mdf.ColumnCount-1 || i >= len(mdf.BaseGraphs)-1 {
			markdown += "|\n"
		}
	}
	return markdown
}

func makeIframeTag(orgName string, g baseGraph) string {
	return fmt.Sprintf(`<iframe src="%s" height="%d" width="%d" frameborder="0"></iframe>`, g.getURL(orgName, false), g.getHeight(), g.getWidth())
}

func makeImageMarkdown(orgName string, g baseGraph) string {
	return fmt.Sprintf("[![graph](%s)](%s)", g.getURL(orgName, true), g.getPermalink(orgName))
}

func generateHostGraphsMarkdownFactory(hostGraphs *hostGraphFormat, graphType string, height int, width int) *markdownFactory {

	if hostGraphs.Period == "" {
		hostGraphs.Period = "1h"
	}

	var baseGraphs []baseGraph

	tableHeader := generateHostGraphsTableHeader(hostGraphs.GraphNames)

	for _, hostID := range hostGraphs.HostIDs {
		for _, graphName := range hostGraphs.GraphNames {
			baseGraphs = append(baseGraphs, hostGraph{
				hostID,
				graphType,
				graphName,
				hostGraphs.Period,
				height,
				width,
			})
		}
	}

	return &markdownFactory{
		Headline:    hostGraphs.Headline,
		TableHeader: tableHeader,
		BaseGraphs:  baseGraphs,
		ColumnCount: len(hostGraphs.GraphNames),
	}
}

func generateHostGraphsTableHeader(graphNames []string) string {
	var header string
	for _, graphName := range graphNames {
		header += "|" + graphName
	}

	header += "|\n" + generateAlignmentLine(len(graphNames))

	return header
}

func generateGraphsMarkdownFactory(graphs *graphFormat, graphType string, height int, width int) (mdf *markdownFactory, err error) {

	if graphs.ColumnCount == 0 {
		graphs.ColumnCount = 1
	}

	tableHeader := generateAlignmentLine(graphs.ColumnCount)

	var baseGraphs []baseGraph
	for _, gd := range graphs.GraphDefs {
		if gd.Period == "" {
			gd.Period = "1h"
		}

		bg, err := gd.getBaseGraph(graphType, height, width)
		if err != nil {
			return nil, err
		}

		baseGraphs = append(baseGraphs, bg)
	}

	return &markdownFactory{
		Headline:    graphs.Headline,
		TableHeader: tableHeader,
		BaseGraphs:  baseGraphs,
		ColumnCount: graphs.ColumnCount,
	}, nil
}

func generateAlignmentLine(count int) string {
	return strings.Repeat("|:-:", count) + "|\n"
}

func doListDashboards(c *cli.Context) error {
	client := mackerelclient.NewFromContext(c)

	dashboards, err := client.FindDashboards()
	logger.DieIf(err)

	fmt.Println(format.JSONMarshalIndent(dashboards, "", "    "))
	return nil
}

func doPullDashboard(c *cli.Context) error {
	client := mackerelclient.NewFromContext(c)

	dashboards, err := client.FindDashboards()
	logger.DieIf(err)
	for _, d := range dashboards {
		dashboard, err := client.FindDashboard(d.ID)
		logger.DieIf(err)
		filename := fmt.Sprintf("dashboard-%s.json", d.ID)
		file, err := os.Create(filename)
		logger.DieIf(err)
		_, err = file.WriteString(format.JSONMarshalIndent(dashboard, "", "    "))
		logger.DieIf(err)
		file.Close()
		logger.Log("info", fmt.Sprintf("Dashboard file is saved to '%s'(title:%s)", filename, d.Title))
	}
	return nil
}

func doPushDashboard(c *cli.Context) error {
	client := mackerelclient.NewFromContext(c)

	f := c.String("file-path")
	src, err := os.Open(f)
	logger.DieIf(err)

	dec := json.NewDecoder(src)
	var dashboard mackerel.Dashboard
	err = dec.Decode(&dashboard)
	logger.DieIf(err)
	if id := dashboard.ID; id != "" {
		_, err := client.FindDashboard(id)
		logger.DieIf(err)

		_, err = client.UpdateDashboard(id, &dashboard)
		logger.DieIf(err)
	} else {
		_, err := client.CreateDashboard(&dashboard)
		logger.DieIf(err)
	}
	return nil
}
