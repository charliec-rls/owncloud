package command

import (
	"fmt"

	"github.com/owncloud/ocis/v2/ocis-pkg/config"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/parser"
	"github.com/owncloud/ocis/v2/ocis/pkg/register"
	"github.com/owncloud/ocis/v2/ocis/pkg/runtime"
	"github.com/urfave/cli/v2"
)

// Server is the entrypoint for the server command.
func Server(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "server",
		Usage:    "start a fullstack server (runtime and all extensions in supervised mode)",
		Category: "fullstack",
		Before: func(c *cli.Context) error {
			if err := parser.ParseConfig(cfg); err != nil {
				fmt.Printf("%v", err)
				return err
			}
			return nil
		},
		Action: func(c *cli.Context) error {
			r := runtime.New(cfg)
			return r.Start()
		},
	}
}

func init() {
	register.AddCommand(Server)
}
