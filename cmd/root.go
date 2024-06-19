package cmd

import (
	"os"

	"github.com/mikuta0407/tcpproxy/proxy"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use: "tcp-proxy",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}

	proxy.Proxy(configFile)
}

var configFile string

func init() {
	rootCmd.Flags().StringVarP(&configFile, "config", "c", "", "config file")
}
