package qtum

import (
	"context"
	"io"
	"os"
	"strings"

	"github.com/alejoacosta74/qproxy/pkg/log"
	"github.com/btcsuite/btclog"
	"github.com/pkg/errors"
	"github.com/qtumproject/btcd/chaincfg"
	"github.com/qtumproject/btcd/rpcclient"
)

// QtumClient is a wrapper for the btcd rpcclient used to connect to a Qtum Node
type QtumClient struct {
	*rpcclient.Client
	cfg *chaincfg.Params
}

func NewQtumClient(host, user, pass, network string) (*QtumClient, error) {
	log.With("module", "qcli").Tracef("Creating new qtum client for network: %s and host: %s", network, host)
	host = strings.TrimPrefix(host, "http://")
	// Connect to local bitcoin core RPC server using HTTP POST method.
	connCfg := &rpcclient.ConnConfig{
		Host:         host,
		User:         user,
		Pass:         pass,
		HTTPPostMode: true, // Bitcoin core only supports HTTP POST mode
		DisableTLS:   true, // Bitcoin core does not provide TLS by default
	}
	// Notice the notification parameter is nil since notifications are
	// not supported in HTTP POST mode.
	var loggerOutput io.Writer
	if log.IsDebug() {
		loggerOutput = os.Stdout
	} else {
		loggerOutput = io.Discard
	}
	backend := btclog.NewBackend(loggerOutput)
	rpcclient.UseLogger(backend.Logger("qtum"))
	qclient, err := rpcclient.New(connCfg, nil)
	if err != nil {
		return nil, err
	}
	// TODO: research
	// if log.IsDebug() {
	// 	_, err := qclient.DebugLevel("qtum=trace")
	// 	if err != nil {
	// 		log.With("module", "qtum").Debugf("Error setting debug level for qtum client: %+v", err)
	// 		return nil, err
	// 	}
	// }

	qcli := QtumClient{
		qclient,
		nil,
	}

	cfg, err := qcli.determineNetworkParams(network)
	if err != nil {
		return nil, errors.Wrapf(err, "Error determining network params for network: %s", network)
	}
	//TODO: Check if the network is the same as the one configured in the node
	// cfg, err := qcli.checkNetworkConfig()
	// if err != nil {
	// 	return nil, err
	// }
	qcli.cfg = cfg

	if qcli.Disconnected() {
		return nil, errors.New("Qtum client is disconnected")
	} else {
		log.With("module", "qcli").Tracef("Qtum client is connected")
	}

	return &qcli, nil
}

func (q *QtumClient) Stop(ctx context.Context) error {
	chErr := make(chan error)
	go func() {
		q.Shutdown()
		q.WaitForShutdown()
		chErr <- nil
	}()
	go func() {
		<-ctx.Done()
		chErr <- ctx.Err()
	}()
	return <-chErr
}

// checkNetworkConfig() returns the network info from the qtum node
func (q *QtumClient) checkNetworkConfig() (*chaincfg.Params, error) {
	info, err := q.GetInfo()
	if err != nil {
		return nil, errors.Wrapf(err, "Error getting info from qtum node")
	}
	log.With("module", "qcli").Debugf("Qtum network testnet: %+v", info.TestNet)
	if info.TestNet {
		return &chaincfg.QtumTestnetParams, nil
	} else {
		return &chaincfg.QtumMainnetParams, nil
	}
}

func (q *QtumClient) determineNetworkParams(network string) (*chaincfg.Params, error) {
	if network == "testnet" || network == "regtest" {
		return &chaincfg.QtumTestnetParams, nil
	} else {
		return &chaincfg.QtumMainnetParams, nil
	}
}
