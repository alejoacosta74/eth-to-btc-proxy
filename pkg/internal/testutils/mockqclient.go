package testutils

import (
	"errors"

	"github.com/alejoacosta74/rpc-proxy/pkg/wallet"
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

type MockQcli struct {
	RawTxResult       *btcjson.TxRawResult
	BlockResult       *btcjson.GetBlockVerboseResult
	ListUnspentResult []btcjson.ListUnspentResult
	AddressResult     *btcjson.GetAddressInfoResult
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
	return nil, errors.New("Block not found")
}

func (q *MockQcli) ListUnspentMinMaxAddresses(minConf int, maxConf int, addrs []btcutil.Address) ([]btcjson.ListUnspentResult, error) {
	return q.ListUnspentResult, nil
}

func (q *MockQcli) GetAddressInfo(address string) (*btcjson.GetAddressInfoResult, error) {
	return q.AddressResult, nil
}

func (q *MockQcli) ImportAddressRescan(address string, account string, rescan bool) error {
	return nil
}

func (q *MockQcli) FindSpendableUTXO(address string, amount float64) ([]btcjson.ListUnspentResult, error) {
	return q.ListUnspentResult, nil
}

func (q *MockQcli) PrepareRawTransaction(unspent []btcjson.ListUnspentResult, sender, receiver string, amount float64) (*wire.MsgTx, error) {
	return nil, nil
}

func (q *MockQcli) SignRawTX(tx *wire.MsgTx, unspent []btcjson.ListUnspentResult, w *wallet.QtumWallet) error {
	return nil
}

func (q *MockQcli) SendRawTransaction(tx *wire.MsgTx, allowHighFees bool) (*chainhash.Hash, error) {
	return nil, nil
}
