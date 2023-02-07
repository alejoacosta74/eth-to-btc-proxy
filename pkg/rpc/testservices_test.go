package rpc

import (
	"bytes"
	"context"
	"net/http"
	"testing"

	utils "github.com/alejoacosta74/rpc-proxy/pkg/internal/testutils"
	"github.com/alejoacosta74/rpc-proxy/pkg/log"
	"github.com/alejoacosta74/rpc-proxy/pkg/qtum"
	"github.com/qtumproject/btcd/chaincfg"
)

func handleFatalError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("Fatal error: %+v", err)
	}
}

// createRPCRequest is a helper function to create a new JSON RPC request
// with the given method and arguments.
func createRPCRequest(url string, method string, args ...string) (*http.Request, error) {
	params := []string{}
	for _, arg := range args {
		arg = `"` + arg + `"`
		params = append(params, arg)
	}
	log.With("module", "rpc").Debugf("Creating json request with method: %s and params: %+v", method, params)
	jsonReq, err := utils.CreateJSONRequest(method, params...)
	if err != nil {
		return nil, err
	}
	return utils.CreateHTTPRequest("POST", url, bytes.NewBuffer(jsonReq))
}

// getETHRPCService is a helper function for unit testing that creates
// a new RPC service for the eth_ namespace
// Returns an RPC server based on go-ethereum RPC server
func getETHRPCService(cfg *chaincfg.Params, qcli qtum.Iqcli) (*RPCService, error) {
	api := NewAPI(context.Background(), qcli)
	api.SetNetworkParams(cfg)
	ethAPI := (*EthAPI)(api)
	rpcservice := NewRPCService()
	err := rpcservice.RegisterName("eth", ethAPI)
	if err != nil {
		return nil, err
	}
	return rpcservice, nil

}

// getPersonalRPCService is a helper function for unit testing that creates
// a new RPC service for the personal_ namespace
// Returns an RPC server based on go-ethereum RPC server
func getPersonalRPCService() (*RPCService, error) {
	api := NewAPI(context.Background(), nil)
	cfg := utils.GetNetworkParams()
	api.SetNetworkParams(cfg)
	personalAPI := (*PersonalAPI)(api)
	rpcservice := NewRPCService()
	err := rpcservice.RegisterName("personal", personalAPI)
	if err != nil {
		return nil, err
	}
	return rpcservice, nil
}
