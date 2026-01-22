package annotations

import (
	"context"
	"io"
	"os"

	"github.com/mackerelio/mackerel-client-go"
	"github.com/mackerelio/mkr/format"
	"github.com/mackerelio/mkr/jq"
	"github.com/mackerelio/mkr/logger"
	"github.com/mackerelio/mkr/mackerelclient"
	"github.com/urfave/cli/v3"
)

var Command = &cli.Command{
	Name:  "annotations",
	Usage: "Manipulate graph annotations",
	Description: `
    Manipulate graph annotations. Requests APIs under "/api/v0/graph-annotations".
    See https://mackerel.io/api-docs/entry/graph-annotations .
`,
	Commands: []*cli.Command{
		{
			Name:      "create",
			Usage:     "create a graph annotation",
			ArgsUsage: "--title <title> [--description <description>] [--description-file <file-path>] --from <from> --to <to> --service|-s <service> [--role|-r <role>]",
			Description: `
    Creates a graph annotation.
`,
			Action: doAnnotationsCreate,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "title",
					Usage: "Title for annotation",
				},
				&cli.StringFlag{
					Name:  "description",
					Usage: "Description for annotation",
				},
				&cli.StringFlag{
					Name:  "description-file",
					Usage: `Read description text for annotation from file (use "-" to read from stdin)`,
				},
				&cli.IntFlag{
					Name:  "from",
					Usage: "Starting time (epoch seconds)",
				},
				&cli.IntFlag{
					Name:  "to",
					Usage: "Ending time (epoch seconds)",
				},
				&cli.StringFlag{
					Name:    "service",
					Aliases: []string{"s"},
					Usage:   "Service name for annotation",
				},
				&cli.StringSliceFlag{
					Name:    "role",
					Aliases: []string{"r"},
					Usage:   "Roles for annotation. Multiple choices are allowed",
				},
			},
		},
		{
			Name:      "list",
			Usage:     "list annotations",
			ArgsUsage: "--from <from> --to <to> --service|-s <service> [--jq <formula>]",
			Description: `
    Shows annotations by service name and duration (from and to)
`,
			Action: doAnnotationsList,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "service",
					Aliases: []string{"s"},
					Usage:   "Service name for annotation",
				},
				&cli.IntFlag{
					Name:  "from",
					Usage: "Starting time (epoch seconds)",
				},
				&cli.IntFlag{
					Name:  "to",
					Usage: "Ending time (epoch seconds)",
				},
				jq.CommandLineFlag,
			},
		},
		{
			Name:      "update",
			Usage:     "update annotation",
			ArgsUsage: "--id <id> [--title <title>] [--description <description>] [--description-file <file-path>] --from <from> --to <to> --service|-s <service> [--role|-r <role>]",
			Description: `
    Updates an annotation
`,
			Action: doAnnotationsUpdate,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "id",
					Usage: "Annotation ID.",
				},
				&cli.StringFlag{
					Name:    "service",
					Aliases: []string{"s"},
					Usage:   "Service name for annotation",
				},
				&cli.StringFlag{
					Name:  "title",
					Usage: "Title for annotation",
				},
				&cli.StringFlag{
					Name:  "description",
					Usage: "Description for annotation",
				},
				&cli.StringFlag{
					Name:  "description-file",
					Usage: `Read description text for annotation from file (use "-" to read from stdin)`,
				},
				&cli.IntFlag{
					Name:  "from",
					Usage: "Starting time (epoch seconds)",
				},
				&cli.IntFlag{
					Name:  "to",
					Usage: "Ending time (epoch seconds)",
				},
				&cli.StringSliceFlag{
					Name:    "role",
					Aliases: []string{"r"},
					Usage:   "Roles for annotation. Multiple choices are allowed",
				},
			},
		},
		{
			Name:      "delete",
			Usage:     "delete annotation",
			ArgsUsage: "--id <id>",
			Description: `
    Delete graph annotation by annotation id
`,
			Action: doAnnotationsDelete,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "id",
					Usage: "Graph annotation ID",
				},
			},
		},
	},
}

func doAnnotationsCreate(ctx context.Context, c *cli.Command) error {
	title := c.String("title")
	description := c.String("description")
	descriptionFile := c.String("description-file")
	from := c.Int64("from")
	to := c.Int64("to")
	service := c.String("service")
	roles := c.StringSlice("role")

	if title == "" {
		_ = cli.ShowCommandHelp(c, "create")
		return cli.Exit("`title` is a required field to create a graph annotation.", 1)
	}

	if service == "" {
		_ = cli.ShowCommandHelp(c, "create")
		return cli.Exit("`service` is a required field to create a graph annotation.", 1)
	}

	if from == 0 {
		_ = cli.ShowCommandHelp(c, "create")
		return cli.Exit("`from` is a required field to create a graph annotation.", 1)
	}

	if to == 0 {
		_ = cli.ShowCommandHelp(c, "create")
		return cli.Exit("`to` is a required field to create a graph annotation.", 1)
	}

	if description != "" && descriptionFile != "" {
		_ = cli.ShowCommandHelp(c, "create")
		return cli.Exit("specify one of `description` or `description-file`.", 1)
	}

	if descriptionFile != "" {
		var (
			b   []byte
			err error
		)
		if descriptionFile == "-" {
			b, err = io.ReadAll(os.Stdin)
		} else {
			b, err = os.ReadFile(descriptionFile)
		}
		logger.DieIf(err)
		description = string(b)
	}

	client := mackerelclient.NewFromContext(c)
	annotation, err := client.CreateGraphAnnotation(&mackerel.GraphAnnotation{
		Title:       title,
		Description: description,
		From:        from,
		To:          to,
		Service:     service,
		Roles:       roles,
	})
	logger.DieIf(err)
	err = format.PrettyPrintJSON(os.Stdout, annotation, "")
	logger.DieIf(err)
	return nil
}

func doAnnotationsList(ctx context.Context, c *cli.Command) error {
	service := c.String("service")
	from := c.Int64("from")
	to := c.Int64("to")

	if service == "" {
		_ = cli.ShowCommandHelp(c, "list")
		return cli.Exit("`service` is a required field to list graph annotations.", 1)
	}

	if from == 0 {
		_ = cli.ShowCommandHelp(c, "list")
		return cli.Exit("`from` is a required field to list graph annotations.", 1)
	}

	if to == 0 {
		_ = cli.ShowCommandHelp(c, "list")
		return cli.Exit("`to` is a required field to list graph annotations.", 1)
	}

	client := mackerelclient.NewFromContext(c)
	annotations, err := client.FindGraphAnnotations(service, from, to)
	logger.DieIf(err)
	err = format.PrettyPrintJSON(os.Stdout, annotations, c.String("jq"))
	logger.DieIf(err)
	return nil
}

func doAnnotationsUpdate(ctx context.Context, c *cli.Command) error {
	annotationID := c.String("id")
	title := c.String("title")
	description := c.String("description")
	descriptionFile := c.String("description-file")
	from := c.Int64("from")
	to := c.Int64("to")
	service := c.String("service")
	roles := c.StringSlice("role")

	if annotationID == "" {
		_ = cli.ShowCommandHelp(c, "update")
		return cli.Exit("`id` is a required field to delete a update annotation.", 1)
	}

	if service == "" {
		_ = cli.ShowCommandHelp(c, "update")
		return cli.Exit("`service` is a required field to update a graph annotation.", 1)
	}

	if from == 0 {
		_ = cli.ShowCommandHelp(c, "update")
		return cli.Exit("`from` is a required field to update a graph annotation.", 1)
	}

	if to == 0 {
		_ = cli.ShowCommandHelp(c, "update")
		return cli.Exit("`to` is a required field to update a graph annotation.", 1)
	}

	if description != "" && descriptionFile != "" {
		_ = cli.ShowCommandHelp(c, "create")
		return cli.Exit("specify one of `description` or `description-file`.", 1)
	}

	if descriptionFile != "" {
		var (
			b   []byte
			err error
		)
		if descriptionFile == "-" {
			b, err = io.ReadAll(os.Stdin)
		} else {
			b, err = os.ReadFile(descriptionFile)
		}
		logger.DieIf(err)
		description = string(b)
	}

	client := mackerelclient.NewFromContext(c)
	annotation, err := client.UpdateGraphAnnotation(annotationID, &mackerel.GraphAnnotation{
		Title:       title,
		Description: description,
		From:        from,
		To:          to,
		Service:     service,
		Roles:       roles,
	})
	logger.DieIf(err)
	err = format.PrettyPrintJSON(os.Stdout, annotation, "")
	logger.DieIf(err)
	return nil
}

func doAnnotationsDelete(ctx context.Context, c *cli.Command) error {
	annotationID := c.String("id")

	if annotationID == "" {
		_ = cli.ShowCommandHelp(c, "delete")
		return cli.Exit("`id` is a required field to delete a graph annotation.", 1)
	}

	client := mackerelclient.NewFromContext(c)
	annotation, err := client.DeleteGraphAnnotation(annotationID)
	logger.DieIf(err)
	err = format.PrettyPrintJSON(os.Stdout, annotation, "")
	logger.DieIf(err)
	return nil
}
