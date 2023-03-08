package mocks

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/alejoacosta74/qproxy/pkg/wallet"
	"github.com/qtumproject/btcd/btcjson"
	"github.com/qtumproject/btcd/btcutil"
	"github.com/qtumproject/btcd/chaincfg/chainhash"
	"github.com/qtumproject/btcd/wire"
)

func NewMockQtumClient(rawTxResult *btcjson.TxRawResult, blockResult *btcjson.GetBlockVerboseResult) *MockQcli {
	return &MockQcli{
		RawTxResult: rawTxResult,
		BlockResult: blockResult,
	}
}

func NewMockQCli() *MockQcli {
	qcli := &MockQcli{}
	qcli.loadDefaultResponses()
	return qcli
}

type MockQcli struct {
	RawTxResult               *btcjson.TxRawResult           // Mock response for GetRawTransactionVerbose
	BlockResult               *btcjson.GetBlockVerboseResult // Mock response for GetBlockVerbose
	AddressResult             *btcjson.GetAddressInfoResult  // Mock response for GetAddressInfo
	BuildUnsignedQtumTxResult *wire.MsgTx                    // Mock response for BuildUnsignedQtumTx
	FindSpendableUTXOResult   []btcjson.ListUnspentResult    // Mock response for FindSpendableUTXO
	SendRawTransactionResult  *chainhash.Hash                // Mock response for SendRawTransaction
	DefaultResponses          map[string]interface{}
}

// Interface methods

func (q *MockQcli) DecodeRawTransaction(serializedTx []byte) (*btcjson.TxRawResult, error) {
	return nil, nil
}

func (q *MockQcli) EstimateSmartFee(confTarget int64, mode *btcjson.EstimateSmartFeeMode) (*btcjson.EstimateSmartFeeResult, error) {
	return nil, nil
}

func (q *MockQcli) GetRawTransactionVerbose(txHash *chainhash.Hash) (*btcjson.TxRawResult, error) {
	if txHash.String() == q.RawTxResult.Txid {
		return q.RawTxResult, nil
	}
	return nil, errors.New("tx not found")
}

func (q *MockQcli) GetBlockVerbose(blockHash *chainhash.Hash) (*btcjson.GetBlockVerboseResult, error) {
	if blockHash.String() == q.BlockResult.Hash {
		return q.BlockResult, nil
	}
	return nil, errors.New("block not found")
}

func (q *MockQcli) ListUnspentMinMaxAddresses(minConf int, maxConf int, addrs []btcutil.Address) ([]btcjson.ListUnspentResult, error) {
	return nil, nil
}

func (q *MockQcli) GetAddressInfo(address string) (*btcjson.GetAddressInfoResult, error) {
	return q.AddressResult, nil
}

func (q *MockQcli) ImportAddressRescan(address string, account string, rescan bool) error {
	return nil
}

// Mocked methods

func (q *MockQcli) VerifyAddress(address string) error {
	return nil
}

func (q *MockQcli) FindSpendableUTXO(address string) ([]btcjson.ListUnspentResult, error) {
	return q.FindSpendableUTXOResult, nil
}

func (q *MockQcli) BuildUnsignedQtumTx(unspent []btcjson.ListUnspentResult, sender, receiver string, amount float64) (*wire.MsgTx, error) {
	return nil, nil
}

func (q *MockQcli) SignRawTX(tx *wire.MsgTx, unspent []btcjson.ListUnspentResult, w wallet.IQtumWallet) error {
	return nil
}

func (q *MockQcli) SendRawTransaction(tx *wire.MsgTx, allowHighFees bool) (*chainhash.Hash, error) {
	return q.SendRawTransactionResult, nil
}

func (q *MockQcli) GetBalance(account string) (btcutil.Amount, error) {
	return 0, nil
}

// Mockqcli default responses

// Default response for FindSpendableUTXO()
var DefaultListUnspentResponseJSON string = `[
	{
	  "txid": "bbe399eebaf12849cb306af8218460061223baa8cb76216358dd68429c921500",
	  "vout": 0,
	  "address": "qUbxboqjBRp96j3La8D1RYkyqx5uQbJPoW",
	  "label": "",
	  "scriptPubKey": "76a9147926223070547d2d15b2ef5e7383e541c338ffe988ac",
	  "amount": 20000.00000000,
	  "confirmations": 2763,
	  "spendable": true,
	  "solvable": false,
	  "safe": true
	},
	{
	  "txid": "8225bd905c83553ebba2bb80887608ea5cc315a3d34e5d1283359b6ffa862e00",
	  "vout": 0,
	  "address": "qUbxboqjBRp96j3La8D1RYkyqx5uQbJPoW",
	  "label": "",
	  "scriptPubKey": "76a9147926223070547d2d15b2ef5e7383e541c338ffe988ac",
	  "amount": 20000.00000000,
	  "confirmations": 2468,
	  "spendable": true,
	  "solvable": false,
	  "safe": true
	},
	{
	  "txid": "50db899a5e3eb817381d82719327720d408daaff6b55a9e5878786b0d44a5f00",
	  "vout": 0,
	  "address": "qUbxboqjBRp96j3La8D1RYkyqx5uQbJPoW",
	  "label": "",
	  "scriptPubKey": "76a9147926223070547d2d15b2ef5e7383e541c338ffe988ac",
	  "amount": 20000.00000000,
	  "confirmations": 2213,
	  "spendable": false,
	  "solvable": false,
	  "safe": true
	}
      ]`

var listUnspentResponse []btcjson.ListUnspentResult

// Default response for GetAddressInfo()
var getAddressInfoResponseJSON string = ` {
	"address": "qUbxboqjBRp96j3La8D1RYkyqx5uQbJPoW",
	"scriptPubKey": "76a9147926223070547d2d15b2ef5e7383e541c338ffe988ac",
	"ismine": false,
	"solvable": false,
	"iswatchonly": false,
	"isscript": false,
	"iswitness": false,
	"ischange": false,
	"labels": [
	]
}`

// Default response for BuildUnsignedQtumTx()
var buildUnsignedQtumTxResponseHEX string = "01000000020015929c4268dd58632176cba8ba231206608421f86a30cb4928f1baee99e3bb0000000000ffffffff0015929c4268dd58632176cba8ba231206608421f86a30cb4928f1baee99e3bb0000000000ffffffff0170dc5d54a30300001976a914e599be870c63d68a00a5019906d258a4ba5d1bac88ac00000000"

var getAddressInfoResponse btcjson.GetAddressInfoResult

// Default response for SendRawTransaction()
var hashHex string = "1dbf40139b6038d5f19b43c592b33a5ad3fe55494e6407712de55cff6b2938da"

func (q *MockQcli) loadDefaultResponses() {

	// Set default response for FindSpendableUTXO()
	err := json.Unmarshal([]byte(DefaultListUnspentResponseJSON), &listUnspentResponse)
	if err != nil {
		fmt.Printf("Error unmarshalling unspent response: %v\n", err)
		panic(err)
	}
	q.FindSpendableUTXOResult = listUnspentResponse

	// Set default response for GetAddressInfo()
	err = json.Unmarshal([]byte(getAddressInfoResponseJSON), &getAddressInfoResponse)
	if err != nil {
		fmt.Printf("Error unmarshalling getAddressInfo: %v\n", err)
		panic(err)
	}
	q.AddressResult = &getAddressInfoResponse

	// Set default response for BuildUnsignedQtumTx()
	unsigned, err := hex.DecodeString(buildUnsignedQtumTxResponseHEX)
	if err != nil {
		fmt.Printf("Error decoding unsigned tx: %v\n", err)
		panic(err)
	}
	buf := bytes.NewBuffer(unsigned)
	unsignedTx := new(wire.MsgTx)
	err = unsignedTx.Deserialize(buf)
	if err != nil {
		fmt.Printf("Error deserializing unsigned tx: %v\n", err)
		panic(err)
	}

	// Set default response for SendRawTransaction()
	hash, err := chainhash.NewHashFromStr(hashHex)
	if err != nil {
		fmt.Printf("Error converting hash: %v\n", err)
		panic(err)
	}
	q.SendRawTransactionResult = hash

}
