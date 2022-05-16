package command

import (
	"fmt"

	"github.com/owncloud/ocis/v2/extensions/app-provider/pkg/command"
	"github.com/owncloud/ocis/v2/ocis-pkg/config"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/parser"
	"github.com/owncloud/ocis/v2/ocis/pkg/command/helper"
	"github.com/owncloud/ocis/v2/ocis/pkg/register"
	"github.com/urfave/cli/v2"
)

// AppProviderCommand is the entrypoint for the app provider command.
func AppProviderCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     cfg.AppProvider.Service.Name,
		Usage:    helper.SubcommandDescription(cfg.AppProvider.Service.Name),
		Category: "extensions",
		Before: func(c *cli.Context) error {
			if err := parser.ParseConfig(cfg); err != nil {
				fmt.Printf("%v", err)
			}
			cfg.AppProvider.Commons = cfg.Commons
			return nil
		},
		Subcommands: command.GetCommands(cfg.AppProvider),
	}
}

func init() {
	register.AddCommand(AppProviderCommand)
}
