package cmd

import (
	"fmt"

	"ociswrapper/common"
	ocis "ociswrapper/ocis"
	ocisConfig "ociswrapper/ocis/config"
	wrapper "ociswrapper/wrapper"
	wrapperConfig "ociswrapper/wrapper/config"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ociswrapper",
	Short: "ociswrapper is a wrapper for oCIS server",
	Run: func(cmd *cobra.Command, args []string) {
		if err := cmd.Help(); err != nil {
			fmt.Printf("error executing help command: %v\n", err)
		}
	},
}

func serveCmd() *cobra.Command {
	serveCmd := &cobra.Command{
		Use:   "serve",
		Short: "Starts the server",
		Run: func(cmd *cobra.Command, args []string) {
			common.Wg.Add(2)

			// set configs
			ocisConfig.Set("bin", cmd.Flag("bin").Value.String())
			ocisConfig.Set("url", cmd.Flag("url").Value.String())
			ocisConfig.Set("retry", cmd.Flag("retry").Value.String())

			go ocis.Start(nil)
			go wrapper.Start(cmd.Flag("port").Value.String())
		},
	}

	// serve command args
	serveCmd.Flags().SortFlags = false
	serveCmd.Flags().StringP("bin", "", ocisConfig.Get("bin"), "Full oCIS binary path")
	serveCmd.Flags().StringP("url", "", ocisConfig.Get("url"), "oCIS server url")
	serveCmd.Flags().StringP("retry", "", ocisConfig.Get("retry"), "Number of retries to start oCIS server")
	serveCmd.Flags().StringP("port", "p", wrapperConfig.Get("port"), "Wrapper API server port")

	return serveCmd
}

// Execute executes the command
func Execute() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	rootCmd.AddCommand(serveCmd())
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("error executing command: %v\n", err)
	}
}
