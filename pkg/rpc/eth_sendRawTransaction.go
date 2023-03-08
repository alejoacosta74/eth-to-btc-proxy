package rpc

import (
	"fmt"

	"github.com/alejoacosta74/qproxy/pkg/log"
	rpctypes "github.com/alejoacosta74/qproxy/pkg/rpc/types"
	"github.com/alejoacosta74/qproxy/pkg/wallet"
	"github.com/qtumproject/btcd/btcjson"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	qcommon "github.com/qtumproject/qtool/lib/common"
	qtool "github.com/qtumproject/qtool/lib/tools"
)

// TODO: replace return type with string

// SendRawTransactionRequest implements the eth_sendRawTransaction JSON-RPC call.
//
// Receives an ethereum signed transaction, decodes it and
// creates a new transaction signed with the stored private key
func (api *EthAPI) SendRawTransaction(rawtx string) (*rpctypes.Eth_SendRawTransactionResponse, error) {
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
	sender, err := getFromAddress(decodedTx)
	if err != nil {
		log.With("method", "sendrawtx").Debugf(err.Error())
		return nil, errors.Wrapf(err, "Error getting sender address from raw transaction: %s", rawtx)
	}
	w, err := ws.SeekWallet(sender.String())
	if err != nil {
		log.With("method", "sendrawtx").Debugf(err.Error())
		return nil, errors.Wrapf(err, "Error loading wallet for address: %s", sender.String())
	}
	// Get address in qtum/btc format
	addr, err := w.GetQtumAddress()
	log.With("module", "eth_sendRawTransaction").Debugf("Wallet returned address encoded: %s", addr)

	if err != nil {
		log.With("method", "sendrawtx").Debugf(err.Error())
		return nil, errors.Wrapf(err, "Error getting qtum address for wallet: %s", sender.String())
	}

	// Convert value in wei to amount in qtum
	weiAmount := decodedTx.Value
	amount, err := qcommon.ConvertWeiToQtum(hexutil.EncodeBig(weiAmount()))
	if err != nil {
		log.With("method", "sendrawtx").Debugf(err.Error())
		return nil, errors.Wrapf(err, "Error converting amount to float: %v", decodedTx.Value().Int64())
	}
	log.With("method", "sendrawtx").Debugf("Amount in wei: %v,  amount in Qtum: %f", decodedTx.Value().Int64(), amount)

	// ensure the address is known to the node's wallet
	err = api.qcli.VerifyAddress(addr)
	if err != nil {
		return nil, errors.Wrapf(err, "Error verifying address: %s", addr)
	}

	// Find spendable UTXO for sender address and amount
	unspent, err := api.qcli.FindSpendableUTXO(addr)
	if err != nil {
		log.With("method", "sendrawtx").Debugf(err.Error())
		return nil, errors.Wrapf(err, "Error finding spendable UTXO for address: %s", addr)
	}
	spendable, err := getUTXOtoSpend(unspent, amount)
	if err != nil {
		log.With("method", "sendrawtx").Debugf(err.Error())
		return nil, errors.Wrapf(err, "Error getting UTXO to spend for address: %s", addr)
	}
	log.With("method", "sendrawtx").Debugf("Found %d utxos to spent", len(spendable))

	if log.IsDebug() {
		fmt.Printf("UTXO to spend:\n")
		fmt.Printf("scriptPubKey: %s\n", spendable[0].ScriptPubKey)
		fmt.Printf("redeemscript: %s\n", spendable[0].RedeemScript)
		fmt.Printf("amount: %f\n", spendable[0].Amount)
		fmt.Printf("address: %s\n", spendable[0].Address)
		fmt.Printf("txid: %s\n", spendable[0].TxID)
		fmt.Printf("vout: %d\n", spendable[0].Vout)
		fmt.Printf("confirmations: %d\n", spendable[0].Confirmations)
		fmt.Printf("spendable: %t\n", spendable[0].Spendable)

	}

	// Convert receiver address to base58
	receiver, err := qtool.AddressHexToBase58(decodedTx.To().String(), api.cfg)
	if err != nil {
		log.With("method", "sendrawtx").Debugf(err.Error())
		return nil, errors.Wrapf(err, "Error converting receiver address to base58: %s", decodedTx.To().String())
	}

	// Create qtum transaction
	log.With("method", "sendrawtx").Debugf("Receiver address: %s", receiver)
	qtumTx, err := api.qcli.BuildUnsignedQtumTx(spendable, addr, receiver, amount)
	if err != nil {
		log.With("method", "sendrawtx").Debugf(err.Error())
		return nil, errors.Wrapf(err, "Error preparing transaction")
	}

	if log.IsDebug() {
		api.printQtumDecodedTX(qtumTx, "Decoded unsigned qtum tx")
	}

	// Sign qtum transaction
	err = api.qcli.SignRawTX(qtumTx, spendable, w)
	if err != nil {
		log.With("method", "sendrawtx").Debugf(err.Error())
		return nil, errors.Wrapf(err, "Error signing transaction")
	}

	if log.IsDebug() {
		api.printQtumDecodedTX(qtumTx, "Decoded signed qtum tx")
	}

	// Send qtum raw transaction
	qtumHash, err := api.qcli.SendRawTransaction(qtumTx, true)
	if err != nil {
		log.With("method", "sendrawtx").Debugf(err.Error())
		return nil, errors.Wrapf(err, "Error sending transaction")
	}
	log.With("method", "sendrawtx").Debugf("Transaction sent with txid: %s", qtumHash.String())

	return &rpctypes.Eth_SendRawTransactionResponse{
		Hash: decodedTx.Hash().String(),
	}, nil
}

// getUTXOtoSpend receives a list of unspent UTXOs and returns
// a list of UTXOs that can be used to spend the amount requested
func getUTXOtoSpend(unspent []btcjson.ListUnspentResult, amount float64) ([]btcjson.ListUnspentResult, error) {
	var utxos []btcjson.ListUnspentResult
	var total float64
	for _, utxo := range unspent {
		if utxo.Confirmations < 6 {
			continue
		}
		utxos = append(utxos, utxo)
		total += utxo.Amount
		if total >= amount {
			break
		}
	}
	if total < amount {
		return nil, errors.New("Not enough funds to spend")
	}
	return utxos, nil
}
