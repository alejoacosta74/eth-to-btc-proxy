package testutils

import (
	"fmt"
	"os"

	"github.com/alejoacosta74/gologger"
	"github.com/alejoacosta74/rpc-proxy/pkg/log"

	"github.com/qtumproject/btcd/chaincfg"
)

const (
	Network  = "testnet"
	QtumURL  = "http://localhost:13881"
	QtumUser = "qtum"
	QtumPass = "qtum"
)

var (
	ChainCfg *chaincfg.Params
	Debug    bool
)

func init() {
	// Read environment DEBUG variable
	environ := os.Environ()
	for _, env := range environ {
		if env == "DEBUG=true" {
			Debug = true
		}
	}

	SetLogger()
	SetNetworkConfig()
	if Debug {
		fmt.Printf("\nSetting logger with level debug: %+v\n", Debug)
		fmt.Printf("Setting testing network config to name: %+v and network: %+v \n", ChainCfg.Name, ChainCfg.Net)
	}
}

func SetLogger() {
	var logger *gologger.Logger
	if Debug {
		logger, _ = gologger.NewLogger(gologger.WithDebugLevel(true))
	} else {
		logger, _ = gologger.NewLogger(gologger.WithNullLogger())
	}
	log.SetLogger(logger)
}

func SetNetworkConfig() {
	switch Network {
	case "mainnet":
		ChainCfg = &chaincfg.QtumMainnetParams
	case "testnet":
		ChainCfg = &chaincfg.QtumTestnetParams
	case "regtest":
		ChainCfg = &chaincfg.QtumRegtestParams
	default:
		panic("wrong network")
	}

}

// GetTestingParams returns the current config params for unit tests
func GetTestingParams() (Network string, QtumURL string, QtumUser string, QtumPass string, cfg *chaincfg.Params) {
	cfg = ChainCfg
	return Network, QtumURL, QtumUser, QtumPass, cfg
}

func GetNetworkParams() *chaincfg.Params {
	return ChainCfg
}
