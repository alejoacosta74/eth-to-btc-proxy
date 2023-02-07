package rpc

import "github.com/alejoacosta74/rpc-proxy/pkg/log"

// Version implements the net_version JSON-RPC call.
//
// Returns the current network protocol version.
//
// No params are required.
func (api *NetAPI) Version() (string, error) {
	// TODO: implement
	log.With("method", "net_version").Debugf("Net_version called. Returning hardwired 8995")
	return "8995", nil
}
