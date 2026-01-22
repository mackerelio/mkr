package dashboards

import (
	"encoding/json"
	"fmt"
	"os"

	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"

	"github.com/mackerelio/mackerel-client-go"
	"github.com/mackerelio/mkr/format"
	"github.com/mackerelio/mkr/logger"
	"github.com/mackerelio/mkr/mackerelclient"
	"github.com/urfave/cli/v3"
)

var Command = &cli.Command{
	Name:  "dashboards",
	Usage: "Manipulate custom dashboards",
	Description: `
    Manipulate custom dashboards. With no subcommand specified, this will show all dashboards. See https://mackerel.io/docs/entry/advanced/cli
`,
	Action: doListDashboards,
	Subcommands: []*cli.Command{
		{
			Name:      "pull",
			Usage:     "Pull custom dashboards",
			ArgsUsage: "--id <id>",
			Description: `
	Pull custom dashboards from Mackerel server and output these to local files.
`,
			Action: doPullDashboard,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "id",
					Usage: "dashboard ID to pull (optional, if not specified, pulls all dashboards)",
				},
			},
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
				&cli.StringFlag{
					Name:    "file-path",
					Aliases: []string{"F"},
					Usage:   "read dashboard from the file",
				},
			},
		},

		// urfave/cli will call default Action when running undefined subcommands,
		// So we leave command defintion to cause error when executing `mkr dashboards generate`.
		{
			Name:   "generate",
			Hidden: true,
			Action: func(c *cli.Context) error {
				logger.Log("error", "`mkr dashboards generate` command has been obsolete")
				os.Exit(1)
				return nil
			},
		},
	},
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

	var dashboards []*mackerel.Dashboard
	if id := c.String("id"); id != "" {
		dashboard, err := client.FindDashboard(id)
		logger.DieIf(err)
		dashboards = append(dashboards, dashboard)
	} else {
		var err error
		dashboards, err = client.FindDashboards()
		logger.DieIf(err)
	}

	for _, d := range dashboards {
		filename := fmt.Sprintf("dashboard-%s.json", d.ID)
		file, err := os.Create(filename)
		logger.DieIf(err)
		_, err = file.WriteString(format.JSONMarshalIndent(d, "", "    "))
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
	fallback := unicode.UTF8.NewDecoder()
	r := transform.NewReader(src, unicode.BOMOverride(fallback))

	dec := json.NewDecoder(r)
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
