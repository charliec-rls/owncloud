package command

import (
	"context"
	"fmt"
	"time"

	"github.com/urfave/cli/v2"

	"github.com/owncloud/ocis/v2/ocis-pkg/config/configlog"
	"github.com/owncloud/ocis/v2/ocis-pkg/service/grpc"
	searchsvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/search/v0"
	"github.com/owncloud/ocis/v2/services/search/pkg/config"
	"github.com/owncloud/ocis/v2/services/search/pkg/config/parser"
)

// Index is the entrypoint for the server command.
func Index(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "index",
		Usage:    "index the files for one or more users",
		Category: "index management",
		Aliases:  []string{"i"},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "space",
				Aliases:  []string{"s"},
				Required: true,
				Usage:    "the id of the space to travers and index the files of",
			},
			&cli.StringFlag{
				Name:     "user",
				Aliases:  []string{"u"},
				Required: true,
				Usage:    "the id of the user that shall be used to access the files",
			},
		},
		Before: func(c *cli.Context) error {
			return configlog.ReturnFatal(parser.ParseConfig(cfg))
		},
		Action: func(c *cli.Context) error {
			client := searchsvc.NewSearchProviderService("com.owncloud.api.search", grpc.Client(grpc.RequestTimeout(60*time.Second)))
			_, err := client.IndexSpace(context.Background(), &searchsvc.IndexSpaceRequest{
				SpaceId: c.String("space"),
				UserId:  c.String("user"),
			})
			if err != nil {
				fmt.Println("failed to index space: " + err.Error())
				return err
			}
			return nil
		},
	}
}
