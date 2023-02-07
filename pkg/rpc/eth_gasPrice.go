package rpc

import "github.com/alejoacosta74/rpc-proxy/pkg/log"

// GasPrice implements the eth_GasPrice JSON-RPC call.
//
// Returns the current price per gas in wei.
//
// No params are required.
func (api *EthAPI) GasPrice() (string, error) {
	// TODO: implement
	log.With("method", "gasPrice").Debugf("GasPrice called. Returning hardwired 0x9184e72a000")
	return "0x9184e72a000", nil
}
