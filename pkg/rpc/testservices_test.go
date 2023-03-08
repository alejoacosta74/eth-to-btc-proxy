package rpc

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/big"
	"net/http"

	utils "github.com/alejoacosta74/qproxy/pkg/internal/testutils"
	"github.com/alejoacosta74/qproxy/pkg/log"
	"github.com/alejoacosta74/qproxy/pkg/qtum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/qtumproject/btcd/chaincfg"
)

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

// newEthereumTx() creates a new unsigned ethereum transaction
// with default parameters.
func newEthereumTx() *types.Transaction {
	nonce := uint64(0)
	gasPrice := big.NewInt(20000000000)
	gasLimit := uint64(21000)
	toAddress := common.HexToAddress("0x71517f86711b4bff4d789ad6fee9a58d8af1c6bb")
	amount := big.NewInt(1000000)
	tx := types.NewTransaction(nonce, toAddress, amount, gasLimit, gasPrice, nil)
	return tx
}

// signEthereumTx() signs an ethereum transaction with a private key.
func signEthereumTx(tx *types.Transaction, signer types.Signer, privKeyHex string) (*types.Transaction, error) {
	privKey, err := crypto.HexToECDSA(privKeyHex)
	if err != nil {
		return nil, err
	}
	signedTx, err := types.SignTx(tx, signer, privKey)
	if err != nil {
		return nil, err
	}
	return signedTx, nil
}

// privToPubKey converts a ECDSA private key hex string to a public key.
func privToPubKey(privKeyHex string) (string, error) {
	privKey, err := crypto.HexToECDSA(privKeyHex)
	if err != nil {
		return "", err
	}
	pubKey := privKey.Public()
	pubKeyECDSA, ok := pubKey.(*ecdsa.PublicKey)
	if !ok {
		return "", fmt.Errorf("error casting public key to ECDSA")
	}
	pubKeyBytes := crypto.FromECDSAPub(pubKeyECDSA)
	return hex.EncodeToString(pubKeyBytes), nil
}
