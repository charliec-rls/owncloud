package command

import (
	"os"

	"github.com/owncloud/ocis/v2/ocis-pkg/clihelper"
	"github.com/owncloud/ocis/v2/ocis-pkg/config"
	"github.com/owncloud/ocis/v2/ocis/pkg/register"
	"github.com/urfave/cli/v2"
)

// Execute is the entry point for the ocis command.
func Execute() error {
	cfg := config.DefaultConfig()

	app := clihelper.DefaultApp(&cli.App{
		Name:  "ocis",
		Usage: "ownCloud Infinite Scale",
	})

	for _, fn := range register.Commands {
		app.Commands = append(
			app.Commands,
			fn(cfg),
		)
	}

	cli.HelpFlag = &cli.BoolFlag{
		Name:  "help,h",
		Usage: "Show the help",
	}

	return app.Run(os.Args)
}
