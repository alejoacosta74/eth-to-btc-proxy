package rpc

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/alejoacosta74/rpc-proxy/pkg/log"
	"github.com/alejoacosta74/rpc-proxy/pkg/qtum"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"
	"github.com/qtumproject/btcd/chaincfg"
	"github.com/qtumproject/btcd/wire"
)

type RPCService struct {
	*rpc.Server
}

// NewRPCService creates a new RPC service instance with no registered handlers.
func NewRPCService() *RPCService {
	return &RPCService{rpc.NewServer()}
}

func NewEthereumRPCService(network string, qcli qtum.Iqcli) (*rpc.Server, error) {
	service := rpc.NewServer()
	api := NewAPI(context.Background(), qcli)
	cfg, err := getNetworkConfig(network)
	if err != nil {
		return nil, err
	}
	api.SetNetworkParams(cfg)
	ethAPI := (*EthAPI)(api)
	err = service.RegisterName("eth", ethAPI)
	if err != nil {
		return nil, errors.Wrap(err, "error registering eth namespace")
	}
	personalAPI := (*PersonalAPI)(api)
	err = service.RegisterName("personal", personalAPI)
	if err != nil {
		return nil, errors.Wrap(err, "error registering personal namespace")
	}

	netAPI := (*NetAPI)(api)
	err = service.RegisterName("net", netAPI)
	if err != nil {
		return nil, errors.Wrap(err, "error registering net namespace")
	}
	return service, nil
}

type API struct {
	ctx  context.Context
	qcli qtum.Iqcli
	cfg  *chaincfg.Params
}

func NewAPI(ctx context.Context, qcli qtum.Iqcli) *API {
	return &API{
		ctx:  ctx,
		qcli: qcli,
	}
}

// SetNetworkParams sets the network params like "testnet" or "mainnet" for the API
func (api *API) SetNetworkParams(cfg *chaincfg.Params) {
	api.cfg = cfg
}

type NetAPI API
type EthAPI API
type PersonalAPI API

// printQtumDecodedTX prints a decoded QTUM transaction
func (api *EthAPI) printQtumDecodedTX(qtumTx *wire.MsgTx, msg string) {
	var buf bytes.Buffer
	err := qtumTx.Serialize(&buf)
	if err != nil {
		log.With("method", "sendrawtx").Debugf(err.Error())
		return
	}
	decoded, err := api.qcli.DecodeRawTransaction(buf.Bytes())
	if err != nil {
		log.With("method", "sendrawtx").Debugf(err.Error())
		return
	}
	decodedBytes, err := json.Marshal(decoded)
	if err != nil {
		log.With("method", "sendrawtx").Debugf(err.Error())
		return
	}
	logPretty(msg, decodedBytes)
}

func getNetworkConfig(network string) (*chaincfg.Params, error) {
	switch network {
	case "regtest":
		return &chaincfg.QtumTestnetParams, nil
	case "testnet":
		return &chaincfg.QtumTestnetParams, nil
	case "mainnet":
		return &chaincfg.QtumMainnetParams, nil
	default:
		return nil, errors.New("Invalid network: " + network)
	}
}
