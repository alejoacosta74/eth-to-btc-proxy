package qtum

import (
	"github.com/alejoacosta74/rpc-proxy/pkg/wallet"
	"github.com/qtumproject/btcd/btcjson"
	"github.com/qtumproject/btcd/chaincfg/chainhash"
	"github.com/qtumproject/btcd/wire"
)

// interface for qtum rpc client
type Iqcli interface {

	// ImportAddressRescan imports the passed public address.
	//
	// When rescan is true, the block history is scanned for transactions
	// addressed to provided address.
	ImportAddressRescan(address string, account string, rescan bool) error

	// FindSpendableUTXO returns a list of spendable UTXOs for the given address
	// with a total amount greater than the given amount
	FindSpendableUTXO(address string, amount float64) ([]btcjson.ListUnspentResult, error)

	// GetAddressInfo returns information about the given qtum address.
	GetAddressInfo(address string) (*btcjson.GetAddressInfoResult, error)

	// PrepareRawTransaction creates a qtum/btc raw transaction using the given unspent outputs
	// to create inputs, and creates resulting outputs
	PrepareRawTransaction(unspent []btcjson.ListUnspentResult, sender, receiver string, amount float64) (*wire.MsgTx, error)

	// SendRawTransaction submits the encoded transaction to the server
	// which will then relay it to the network.
	SendRawTransaction(tx *wire.MsgTx, allowHighFees bool) (*chainhash.Hash, error)

	// SignRawTX signs the given raw transaction off-line using the given unspent outputs to create
	// signatures for the inputs.
	//
	// The transaction is not sent to the network.
	SignRawTX(tx *wire.MsgTx, unspent []btcjson.ListUnspentResult, w *wallet.QtumWallet) error

	// DecodeRawTransaction returns information about a transaction given its serialized bytes.
	DecodeRawTransaction(serializedTx []byte) (*btcjson.TxRawResult, error)

	// EstimateFee provides an estimated fee in bitcoins per kilobyte.
	EstimateSmartFee(confTarget int64, mode *btcjson.EstimateSmartFeeMode) (*btcjson.EstimateSmartFeeResult, error)
}
