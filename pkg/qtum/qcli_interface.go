package qtum

import (
	"github.com/alejoacosta74/qproxy/pkg/wallet"
	// "github.com/btcsuite/btcutil"
	"github.com/qtumproject/btcd/btcjson"
	"github.com/qtumproject/btcd/btcutil"
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
	FindSpendableUTXO(address string) ([]btcjson.ListUnspentResult, error)

	// GetAddressInfo returns information about the given qtum address.
	GetAddressInfo(address string) (*btcjson.GetAddressInfoResult, error)

	// BuildUnsignedQtumTx creates a qtum/btc raw transaction using the given unspent outputs
	// to create inputs, and creates resulting outputs
	BuildUnsignedQtumTx(unspent []btcjson.ListUnspentResult, sender, receiver string, amount float64) (*wire.MsgTx, error)

	// SendRawTransaction submits the encoded transaction to the server
	// which will then relay it to the network.
	SendRawTransaction(tx *wire.MsgTx, allowHighFees bool) (*chainhash.Hash, error)

	// SignRawTX signs the given raw transaction off-line using the given unspent outputs to create
	// signatures for the inputs.
	//
	// The transaction is not sent to the network.
	SignRawTX(tx *wire.MsgTx, unspent []btcjson.ListUnspentResult, wallet wallet.IQtumWallet) error

	// DecodeRawTransaction returns information about a transaction given its serialized bytes.
	DecodeRawTransaction(serializedTx []byte) (*btcjson.TxRawResult, error)

	// EstimateFee provides an estimated fee in bitcoins per kilobyte.
	EstimateSmartFee(confTarget int64, mode *btcjson.EstimateSmartFeeMode) (*btcjson.EstimateSmartFeeResult, error)

	// VerifyAddress checks if the address is known to the node's wallet.
	// If not, it will import the address and rescan the blockchain seeking transactions
	// related to this address.
	VerifyAddress(address string) error

	GetBalance(account string) (btcutil.Amount, error)
}
