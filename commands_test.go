package main

import (
	"strings"
	"testing"

	cli "gopkg.in/urfave/cli.v1"
)

func TestCommands_requirements(t *testing.T) {
	var cs, subcs []cli.Command
	for _, c := range Commands {
		if len(c.Subcommands) == 0 {
			cs = append(cs, c)
		} else {
			for _, sc := range c.Subcommands {
				cs = append(cs, sc)
			}
			subcs = append(subcs, c)
		}
	}
	for _, c := range cs {
		if !strings.HasPrefix(c.Description, "\n    ") {
			t.Errorf("%s: cli.Command.Description should start with '\\n    ', got:\n%s", c.Name, c.Description)
		}
		if !strings.HasSuffix(c.Description, "\n") {
			t.Errorf("%s: cli.Command.Description should end with '\\n', got:\n%s", c.Name, c.Description)
		}
		if len(c.Flags) > 0 && c.ArgsUsage == "" {
			t.Errorf("%s: cli.Command.ArgsUsage should not be empty. Describe flag options.", c.Name)
		}
	}
	for _, sc := range subcs {
		if sc.Action == nil {
			if sc.Description == "" && sc.Usage == "" {
				t.Errorf("%s: Neither .Description nor .Usage should be empty", sc.Name)

			}
		} else if sc.Description == "" {
			t.Errorf("%s: cli.Command.Description should not be empty", sc.Name)
		}
	}
}
