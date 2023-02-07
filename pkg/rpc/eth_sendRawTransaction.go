package rpc

import (
	"github.com/alejoacosta74/rpc-proxy/pkg/log"
	"github.com/alejoacosta74/rpc-proxy/pkg/types"
	"github.com/alejoacosta74/rpc-proxy/pkg/wallet"
	"github.com/pkg/errors"
	qcommon "github.com/qtumproject/qtool/lib/common"
	qtool "github.com/qtumproject/qtool/lib/tools"
)

// SendRawTransactionRequest implements the eth_sendRawTransaction JSON-RPC call.
//
// Receives an ethereum signed transaction, decodes it and
// creates a new transaction signed with the stored private key
func (api *EthAPI) SendRawTransaction(rawtx string) (*types.Eth_SendRawTransactionResponse, error) {
	log.With("method", "sendrawtx").Debugf("SendRawTransaction called with rawtx: %+v", rawtx)

	// Decode raw transaction
	decodedTx, err := decodeRawTx(rawtx)
	if err != nil {
		log.With("method", "sendrawtx").Debugf(err.Error())
		return nil, errors.Wrapf(err, "Error decoding raw transaction: %s", rawtx)
	}
	log.With("module", "eth_sendRawTransaction").Tracef("Decoded transaction: %+v", *decodedTx)

	// Load wallet for sender eth hex address
	ws := wallet.GetWallets()
	w, err := ws.SeekWallet(decodedTx.From)
	if err != nil {
		log.With("method", "sendrawtx").Debugf(err.Error())
		return nil, errors.Wrapf(err, "Error loading wallet for address: %s", decodedTx.From)
	}
	// Get address in qtum/btc format
	addr, err := w.GetAddress()
	log.With("module", "eth_sendRawTransaction").Debugf("Wallet returned address encoded: %s", addr)

	if err != nil {
		log.With("method", "sendrawtx").Debugf(err.Error())
		return nil, errors.Wrapf(err, "Error getting qtum address for wallet: %s", decodedTx.From)
	}

	// Convert value in wei to amount in qtum
	amount, err := qcommon.ConvertWeiToQtum(decodedTx.Value)
	if err != nil {
		log.With("method", "sendrawtx").Debugf(err.Error())
		return nil, errors.Wrapf(err, "Error converting amount to float: %s", decodedTx.Value)
	}
	log.With("method", "sendrawtx").Debugf("Amount in wei: %v,  amount in Qtum: %f", decodedTx.Value, amount)

	// Find spendable UTXO for sender address and amount
	unspent, err := api.qcli.FindSpendableUTXO(addr, amount)
	if err != nil {
		log.With("method", "sendrawtx").Debugf(err.Error())
		return nil, errors.Wrapf(err, "Error finding spendable UTXO for address: %s", addr)
	}
	log.With("method", "sendrawtx").Debugf("Found %d utxos to spent", len(unspent))

	receiver, err := qtool.AddressHexToBase58(decodedTx.To, api.cfg)
	if err != nil {
		log.With("method", "sendrawtx").Debugf(err.Error())
		return nil, errors.Wrapf(err, "Error converting receiver address to base58: %s", decodedTx.To)
	}

	// Create qtum transaction
	log.With("method", "sendrawtx").Debugf("Receiver address: %s", receiver)
	qtumTx, err := api.qcli.PrepareRawTransaction(unspent, addr, receiver, amount)
	if err != nil {
		log.With("method", "sendrawtx").Debugf(err.Error())
		return nil, errors.Wrapf(err, "Error preparing transaction")
	}

	if log.IsDebug() {
		api.printQtumDecodedTX(qtumTx, "Decoded unsigned qtum tx")
	}

	// Sign qtum transaction
	err = api.qcli.SignRawTX(qtumTx, unspent, w)
	if err != nil {
		log.With("method", "sendrawtx").Debugf(err.Error())
		return nil, errors.Wrapf(err, "Error signing transaction")
	}

	if log.IsDebug() {
		api.printQtumDecodedTX(qtumTx, "Decoded signed qtum tx")
	}

	// Send qtum raw transaction
	txid, err := api.qcli.SendRawTransaction(qtumTx, true)
	if err != nil {
		log.With("method", "sendrawtx").Debugf(err.Error())
		return nil, errors.Wrapf(err, "Error sending transaction")
	}
	log.With("method", "sendrawtx").Debugf("Transaction sent with txid: %s", txid.String())

	return &types.Eth_SendRawTransactionResponse{
		Hash: "0x" + txid.String(),
	}, nil
}
