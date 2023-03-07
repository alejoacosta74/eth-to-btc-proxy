/*
Copyright © 2022 Alejo Acosta
*/
package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alejoacosta74/gologger"

	"github.com/alejoacosta74/rpc-proxy/pkg/log"
	"github.com/alejoacosta74/rpc-proxy/pkg/qtum"
	"github.com/alejoacosta74/rpc-proxy/pkg/server"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "qproxy",
	Short: "A JSON RPC proxy that converts a eth tx to qtum tx ",
	Long: `Proxy server that converts an Ethereum signed transaction to a 
Qtum (Bitcoin) signed transaction and sends it for broadcasting to the Qtum node.
`,
	Run:               runQtumProxy,
	Args:              cobra.MaximumNArgs(0),
	PersistentPreRunE: runPreRunE,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var (
	address         string
	backendUrl      string
	qtumRpcEndPoint string
	qtumUser        string
	qtumPass        string
	network         string

	logger   *gologger.Logger
	cfgFile  string
	logLevel string
)

func init() {

	rootCmd.Flags().StringVarP(&address, "address", "a", ":8080", "Address to listen on")
	rootCmd.PersistentFlags().StringVarP(&backendUrl, "backend", "b", "http://127.0.0.1:7545", "Backend URL to proxy to")
	rootCmd.PersistentFlags().StringVarP(&qtumRpcEndPoint, "qtumrpc", "q", "127.0.0.1:3889", "Qtum RPC endpoint")
	rootCmd.PersistentFlags().StringVarP(&qtumUser, "user", "u", "qtum", "Qtum user")
	rootCmd.PersistentFlags().StringVarP(&qtumPass, "pass", "p", "qtum", "Qtum password")
	rootCmd.PersistentFlags().StringVarP(&network, "network", "n", "regtest", "Qtum network")

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.qproxy.env)")
	rootCmd.PersistentFlags().StringVarP(&logLevel, "loglevel", "l", "", "log level (trace, debug, info, warn, error, fatal, panic")

	cobra.MarkFlagFilename(rootCmd.PersistentFlags(), "config")
	viper.BindPFlag("loglevel", rootCmd.PersistentFlags().Lookup("loglevel"))
	viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))
	viper.BindPFlag("nework", rootCmd.PersistentFlags().Lookup("network"))

}

func runPreRunE(cmd *cobra.Command, args []string) error {
	err := loadConfig()
	if err != nil {
		return err
	}
	err = setLogger()
	if err != nil {
		return err
	}
	return nil

}

func runQtumProxy(cmd *cobra.Command, args []string) {

	// Create new Qtum RPC client
	qclient, err := qtum.NewQtumClient(qtumRpcEndPoint, qtumUser, qtumPass, network)
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}
	// Create new proxy server
	srv, err := server.NewServer(address, backendUrl, qclient, network)
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}
	go func() {

		if err := srv.Start(); err != nil {
			logger.Fatal(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// Accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt, syscall.SIGQUIT)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	err = srv.Stop(ctx)
	if err != nil {
		logger.WithField("module", "root").Fatal(err)
	}
	// stop qtum client
	logger.WithField("module", "root").Debug("Stopping Qtum client")
	err = qclient.Stop(ctx)
	if err != nil {
		logger.WithField("module", "root").Fatal(err)
	}

	logger.Println("shutting down")
	os.Exit(0)

}

func setLogger() error {
	var err error
	level := viper.GetString("loglevel")
	switch level {
	case "trace":
		logger, err = gologger.NewLogger(gologger.WithLevel(gologger.TraceLevel))
	case "debug":
		logger, err = gologger.NewLogger(gologger.WithLevel(gologger.DebugLevel))
	case "info":
		logger, err = gologger.NewLogger(gologger.WithLevel(gologger.InfoLevel))
	case "warn":
		logger, err = gologger.NewLogger(gologger.WithLevel(gologger.WarnLevel))
	case "error":
		logger, err = gologger.NewLogger(gologger.WithLevel(gologger.ErrorLevel))
	case "fatal":
		logger, err = gologger.NewLogger(gologger.WithLevel(gologger.FatalLevel))
	case "panic":
		logger, err = gologger.NewLogger(gologger.WithLevel(gologger.PanicLevel))
	default:
		logger, err = gologger.NewLogger(gologger.WithNullLogger())
	}
	if err != nil {
		return err
	}
	log.SetLogger(logger)
	return nil
}

// loadConfig loads the ethcli configuration from a
// file or from env variables
func loadConfig() error {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
		viper.SetConfigType("yaml")

		if err := viper.ReadInConfig(); err != nil {
			return fmt.Errorf("failed to read config file - %s", err)
		}
	} else {
		// Default location for config file is $HOME/.ethcli
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return errors.Wrap(err, "error getting user home directory")
		} else {
			configFile := homeDir + "/.ethcli/config.yaml"
			viper.SetConfigFile(configFile)

			if err := viper.ReadInConfig(); err != nil {
				return errors.Wrap(err, "failed to read config file")
			}
		}
	}

	viper.AutomaticEnv()
	return nil
}
