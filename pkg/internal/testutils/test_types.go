package testutils

import (
	"github.com/qtumproject/btcd/btcjson"
)

type TestDescription struct {
	TestName        string `json:"test_name"`
	TestDescription string `json:"test_description"`
	TxType          string `json:"tx_type"`
	QtumExplorerUrl string `json:"qtumexplorer"`
}
type TestEthGetRawTx struct {
	TxHash string                           `json:"tx_hash"`
	Result Eth_GetTransactionByHashResponse `json:"result"`
}

type TestQtumGetBlock struct {
	BlockHash string                        `json:"block_hash"`
	Result    btcjson.GetBlockVerboseResult `json:"result"`
}

type TestQtumGetRawTx struct {
	TxHash string                `json:"tx_hash"`
	Result []btcjson.TxRawResult `json:"result"`
}

type TestQtumGetTxOut struct {
	TxHash string                   `json:"tx_hash"`
	Result []btcjson.GetTxOutResult `json:"result"`
}

type EthGetTxByHashTestData struct {
	EthGetRawTx  TestEthGetRawTx  `json:"ethgetrawtx"`
	QtumGetBlock TestQtumGetBlock `json:"qtumgetblock"`
	QtumGetRawTx TestQtumGetRawTx `json:"qtumgetrawtx"`
	QtumGetTxOut TestQtumGetTxOut `json:"qtumgettxout"`
}

type EthGetTxByHashTest struct {
	Description TestDescription        `json:"description"`
	Data        EthGetTxByHashTestData `json:"data"`
}

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
