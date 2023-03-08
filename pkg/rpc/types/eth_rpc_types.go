package types

import (
	"encoding/json"
	"math/big"

	// local "github.com/alejoacosta74/qproxy/pkg/types"
	"github.com/ethereum/go-ethereum/core/types"
)

// This file contains the types used for the RPC calls

// RPC Method: eth_getTransactionByHash
type Eth_GetTransactionByHashResponse struct {
	BlockHash        string `json:"blockHash"`
	BlockNumber      string `json:"blockNumber"`
	From             string `json:"from"`
	Gas              string `json:"gas"`
	GasPrice         string `json:"gasPrice"`
	Hash             string `json:"hash"`
	Input            string `json:"input"`
	Nonce            string `json:"nonce"`
	To               string `json:"to"`
	TransactionIndex string `json:"transactionIndex"`
	Value            string `json:"value"`
	V                string `json:"v"`
	R                string `json:"r"`
	S                string `json:"s"`
}

// RPC Method: eth_sendRawTransaction
type Eth_SendRawTransactionRequest struct {
	RawTx string `json:"rawtx"`
}

type Eth_SendRawTransactionResponse struct {
	Hash string `json:"hash"`
}

type RawTransaction struct {
	From             string `json:"from"`
	Gas              string `json:"gas"`
	GasPrice         string `json:"gasPrice"`
	Hash             string `json:"hash"`
	Input            string `json:"input"`
	Nonce            string `json:"nonce"`
	To               string `json:"to"`
	TransactionIndex string `json:"transactionIndex"`
	Value            string `json:"value"`
	V                string `json:"v"`
	R                string `json:"r"`
	S                string `json:"s"`
}

type TransactionHash struct {
	Hash string `json:"hash"`
}

type RlpDecodedRawTransaction struct {
	types.Transaction
}

// Unmarshal decodes the RLP encoded raw transaction into a RawTransaction type
func (tx *RlpDecodedRawTransaction) Unmarshal() (*RawTransaction, error) {

	var rawTx RawTransaction
	txJSONBytes, err := tx.MarshalJSON()
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(txJSONBytes, &rawTx)
	if err != nil {
		return nil, err
	}
	msg, err := tx.AsMessage(types.NewEIP155Signer(tx.ChainId()), big.NewInt(1))
	if err != nil {
		return nil, err
	}

	rawTx.From = msg.From().String()

	return &rawTx, nil

}

// RPC Method: eth_gasPrice
type Eth_GasPriceRequest string
type Eth_GasPriceResponse string

// RPC Method: net_version
type Net_VersionRequest string
type Net_VersionResponse string

// RPC Method: eth_getTransactionCount
type Eth_GetTransactionCountResponse struct {
	Count uint64 `json:"count"`
}
