package command

import (
	"github.com/owncloud/ocis/extensions/graph-explorer/pkg/command"
	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/urfave/cli/v2"
)

// GraphExplorerCommand is the entrypoint for the graph-explorer command.
func GraphExplorerCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     cfg.GraphExplorer.Service.Name,
		Usage:    subcommandDescription(cfg.GraphExplorer.Service.Name),
		Category: "extensions",
		Subcommands: command.GetCommands(cfg.GraphExplorer),
	}
}

func init() {
	register.AddCommand(GraphExplorerCommand)
}
