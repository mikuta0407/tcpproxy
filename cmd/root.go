package cmd

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mikuta0407/tcpproxy/config"
	"github.com/mikuta0407/tcpproxy/proxy"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use: "tcpproxy",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}

	start()
}

var proxyConfig config.Config
var configFile string

func init() {

	rootCmd.Flags().StringVarP(&configFile, "config", "c", "", "config file")
	var err error
	proxyConfig, err = config.LoadConfigFile(configFile)
	if err != nil {
		panic(err)
	}

}

func start() {
	log.Println("Starting port forwarding tool...")

	cancelFuncs := make(map[string]context.CancelFunc)
	proxy := proxy.NewProxy(&proxyConfig, cancelFuncs)

	// SIGHUPシグナルをハンドリングするためのチャネルを設定
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGHUP)
	// シグナルを待つGoroutineを起動
	go func() {
		for {
			<-sigs
			log.Println("Received SIGHUP signal. Reloading configuration...")
			newProxyConfig, err := config.LoadConfigFile(configFile)
			if err != nil {
				panic(err)
			}
			proxy.StopAllProxies()
			proxy.ConfigMutex.Lock()
			proxy.ActiveConfig = &newProxyConfig
			proxy.ConfigMutex.Unlock()
			proxy.Start()
			log.Println("Configuration reloaded successfully.")
		}
	}()

	// プロキシ起動
	proxy.Start()

	// プログラムが終了しないように待機
	select {}

}
