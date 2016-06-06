package main

import (
	"fmt"
	"io/ioutil"
	"net/url"
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
    A custom dashboard is registered from a yaml file.
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
func (e expressionGraph) getPermalink(orgName string) string {
	u, _ := url.Parse(fmt.Sprintf("https://mackerel.io/orgs/%s/advanced-graph", orgName))
	param := url.Values{}
	param.Add("query", e.Query)
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

func doGenerateDashboards(c *cli.Context) error {
	conffile := c.GlobalString("conf")

	isStdout := c.Bool("print")

	argFilePath := c.Args()
	if len(argFilePath) < 1 {
		cli.ShowCommandHelp(c, "generate")
		return cli.NewExitError("specify a yaml file.", 1)
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
		return cli.NewExitError("title is required in yaml.", 1)
	}
	if yml.URLPath == "" {
		return cli.NewExitError("url_path is required in yaml.", 1)
	}
	if yml.GraphType == "" {
		yml.GraphType = "iframe"
	}
	if yml.GraphType != "iframe" && yml.GraphType != "image" {
		return cli.NewExitError("graph_type should be 'iframe' or 'image'.", 1)
	}
	if yml.Height == 0 {
		yml.Height = 200
	}
	if yml.Width == 0 {
		yml.Width = 400
	}

	if yml.HostGraphFormat != nil && yml.GraphFormat != nil {
		return cli.NewExitError("you cannot specify both 'graphs' and host_graphs'.", 1)
	}

	var markdown string
	for _, h := range yml.HostGraphFormat {
		mdf := generateHostGraphsMarkdownFactory(h, yml.GraphType, yml.Height, yml.Width, client)
		markdown += mdf.generate(org.Name)
	}
	for _, g := range yml.GraphFormat {
		mdf, err := generateGraphsMarkdownFactory(g, yml.GraphType, yml.Height, yml.Width)
		if err != nil {
			return err
		}
		markdown += mdf.generate(org.Name)
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

func generateHostGraphsMarkdownFactory(hostGraphs *hostGraphFormat, graphType string, height int, width int, client *mackerel.Client) *markdownFactory {

	tableHeader := generateHostGraphsTableHeader(hostGraphs.HostIDs, client)

	if hostGraphs.Period == "" {
		hostGraphs.Period = "1h"
	}

	var baseGraphs []baseGraph
	for _, graphName := range hostGraphs.GraphNames {
		for _, hostID := range hostGraphs.HostIDs {
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
		ColumnCount: len(hostGraphs.HostIDs),
	}
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

	header += "|\n" + generateAlignmentLine(len(hostIDs))

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
