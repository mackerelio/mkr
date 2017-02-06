package main

import (
	"os"

	mkr "github.com/mackerelio/mackerel-client-go"
	"github.com/mackerelio/mkr/logger"
	"gopkg.in/urfave/cli.v1"
)

var commandAnnotations = cli.Command{
	Name: "annotations",
	Subcommands: []cli.Command{
		{
			Name:        "create",
			Usage:       "create annotation",
			Description: "Creates an annotation.",
			Action:      doAnnotationsCreate,
			Flags: []cli.Flag{
				cli.StringFlag{Name: "title", Value: "", Usage: "Title for annotation."},
				cli.StringFlag{Name: "description", Value: "", Usage: "Description for annotation."},
				cli.IntFlag{Name: "from"},
				cli.IntFlag{Name: "to"},
				cli.StringFlag{Name: "service, s", Value: "", Usage: "Service name for annotation."},
				cli.StringSliceFlag{
					Name:  "role, r",
					Value: &cli.StringSlice{},
					Usage: "Roles for annotation. Multiple choices are allowed.",
				},
			},
		},
		{
			Name:        "list",
			Usage:       "list annotations",
			Description: "Shows annotations by service name and duration(from and to).",
			Action:      doAnnotationsList,
			Flags: []cli.Flag{
				cli.StringFlag{Name: "service, s", Value: "", Usage: "Service name for annotation."},
				cli.IntFlag{Name: "from"},
				cli.IntFlag{Name: "to"},
			},
		},
		{
			Name:        "update",
			Usage:       "update annotation",
			Description: "Updates an annotation.",
			Action:      doAnnotationsUpdate,
			Flags: []cli.Flag{
				cli.StringFlag{Name: "id", Value: "", Usage: "Annotation ID."},
				cli.StringFlag{Name: "service, s", Value: "", Usage: "Service name for annotation."},
				cli.StringFlag{Name: "title", Value: "", Usage: "Title for annotation."},
				cli.StringFlag{Name: "description", Value: "", Usage: "Description for annotation."},
				cli.IntFlag{Name: "from"},
				cli.IntFlag{Name: "to"},
				cli.StringSliceFlag{
					Name:  "role, r",
					Value: &cli.StringSlice{},
					Usage: "Roles for annotation. Multiple choices are allowed.",
				},
			},
		},
		{
			Name:        "delete",
			Usage:       "delete annotation",
			Description: "Delete graph annotation by annotation id.",
			Action:      doAnnotationsDelete,
			Flags: []cli.Flag{
				cli.StringFlag{Name: "id", Value: "", Usage: "Reason of closing alert."},
			},
		},
	},
}

func doAnnotationsCreate(c *cli.Context) error {
	title := c.String("title")
	description := c.String("description")
	from := c.Int64("from")
	to := c.Int64("to")
	service := c.String("service")
	roles := c.StringSlice("role")

	if service == "" || from == 0 || to == 0 {
		_ = cli.ShowCommandHelp(c, "create")
		os.Exit(1)
	}

	client := newMackerelFromContext(c)
	err := client.CreateGraphAnnotation(&mkr.GraphAnnotation{
		Title:       title,
		Description: description,
		From:        from,
		To:          to,
		Service:     service,
		Roles:       roles,
	})
	logger.DieIf(err)
	return nil
}

func doAnnotationsList(c *cli.Context) error {
	service := c.String("service")
	from := c.Int64("from")
	to := c.Int64("to")

	if service == "" || from == 0 || to == 0 {
		_ = cli.ShowCommandHelp(c, "list")
		os.Exit(1)
	}

	client := newMackerelFromContext(c)
	annotations, err := client.FindGraphAnnotations(service, from, to)
	logger.DieIf(err)
	PrettyPrintJSON(annotations)
	return nil
}

func doAnnotationsUpdate(c *cli.Context) error {
	annotationID := c.String("id")
	title := c.String("title")
	description := c.String("description")
	from := c.Int64("from")
	to := c.Int64("to")
	service := c.String("service")
	roles := c.StringSlice("role")

	if service == "" || from == 0 || to == 0 {
		_ = cli.ShowCommandHelp(c, "update")
		os.Exit(1)
	}

	client := newMackerelFromContext(c)
	annotation, err := client.UpdateGraphAnnotation(annotationID, &mkr.GraphAnnotation{
		Title:       title,
		Description: description,
		From:        from,
		To:          to,
		Service:     service,
		Roles:       roles,
	})
	logger.DieIf(err)
	PrettyPrintJSON(annotation)
	return nil
}

func doAnnotationsDelete(c *cli.Context) error {
	annotationID := c.String("id")

	if annotationID == "" {
		_ = cli.ShowCommandHelp(c, "delete")
		os.Exit(1)
	}

	client := newMackerelFromContext(c)
	annotation, err := client.DeleteGraphAnnotation(annotationID)
	logger.DieIf(err)
	PrettyPrintJSON(annotation)
	return nil
}
