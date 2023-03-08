package rpc

import (
	"github.com/alejoacosta74/qproxy/pkg/log"
)

// GetTransactionCount implements the eth_getTransactionCount JSON-RPC call.
//
// Returns the number of transactions sent from an address to be used
// to calculate the nonce field.
func (api *EthAPI) GetTransactionCount(address string, block string) (string, error) {
	log.With("method", "getTransactionCount").Debugf("GetTransactionCount called with address: %s, block: %s", address, block)

	return "0x0", nil
}
