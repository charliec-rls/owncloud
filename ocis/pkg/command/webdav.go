package command

import (
	"github.com/owncloud/ocis/extensions/webdav/pkg/command"
	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/urfave/cli/v2"
)

// WebDAVCommand is the entrypoint for the webdav command.
func WebDAVCommand(cfg *config.Config) *cli.Command {

	return &cli.Command{
		Name:        cfg.WebDAV.Service.Name,
		Usage:       subcommandDescription(cfg.WebDAV.Service.Name),
		Category:    "extensions",
		Subcommands: command.GetCommands(cfg.WebDAV),
	}
}

func init() {
	register.AddCommand(WebDAVCommand)
}
