package qtum

import (
	"encoding/hex"
	"fmt"

	"github.com/alejoacosta74/rpc-proxy/pkg/log"
	"github.com/alejoacosta74/rpc-proxy/pkg/wallet"
	"github.com/pkg/errors"
	"github.com/qtumproject/btcd/btcjson"
	"github.com/qtumproject/btcd/btcutil"
	"github.com/qtumproject/btcd/chaincfg/chainhash"
	"github.com/qtumproject/btcd/txscript"
	"github.com/qtumproject/btcd/wire"
	"github.com/shopspring/decimal"
)

const (
	// DefaultGasPrice is the default gas price used for transactions
	DefaultGasPrice = 100000
	// Qtum is the number of satoshis in 1 Qtum
	Qtum = 100000000
	// Precision digits to use for decimal operations with Qtum amounts
	PrecisionExp = -8
)

// findCoinStake returns the transaction hash of the coinstake transaction in the block
//
// Used to calculate the effective gas fee
func (q *QtumClient) findCoinStakeTx(block *btcjson.GetBlockVerboseResult) (*chainhash.Hash, error) {

	for _, tx := range block.Tx {
		txHash, _ := chainhash.NewHashFromStr(tx)
		// GetTxOut returns the transaction output info if it's unspent and nil, otherwise.
		txOut, err := q.GetTxOut(txHash, 0, true)
		if err != nil {
			return nil, err
		}
		if txOut != nil && txOut.Coinstake {
			return txHash, nil
		}
	}
	return nil, nil
}

// CalculateEffectiveGasFee returns the effective gas fee of a transaction
func (q *QtumClient) CalculateEffectiveGasFee(vins []btcjson.Vin, block *btcjson.GetBlockVerboseResult) (int64, error) {
	if len(vins) == 0 {
		return 0, fmt.Errorf("no vins found in tx")
	}
	to := vins[0].Address

	// Get the coinbase tx of the block
	coinStakeTxHash, err := q.findCoinStakeTx(block)
	if err != nil {
		return 0, err
	}
	if coinStakeTxHash == nil {
		return 0, nil
	}
	coinStakeTx, err := q.GetRawTransactionVerbose(coinStakeTxHash)
	log.Debugf("coinbaseTx found: %#v", coinStakeTx.Hash)
	if err != nil {
		return 0, err
	}

	// Get the effective fee
	for _, vout := range coinStakeTx.Vout {
		if len(vout.ScriptPubKey.Addresses) > 0 && vout.ScriptPubKey.Addresses[0] == to {
			return vout.AmountSatoshi, nil
		}
	}
	return 0, nil
}

// GetTransactionIndex seeks a tx hash within the block and returns the index of the transaction
func (q *QtumClient) GetTransactionIndex(txHash string, block *btcjson.GetBlockVerboseResult) (int, error) {

	for i, tx := range block.Tx {
		if tx == txHash {
			return i, nil
		}
	}
	return 0, fmt.Errorf("tx not found in block")
}

// FindSpendableUTXO returns a list of spendable UTXOs for the given address
//
// Params:
//   - addr: the address to search for UTXOs in base58 format
func (q *QtumClient) FindSpendableUTXO(addr string) ([]btcjson.ListUnspentResult, error) {

	log.With("module", "qtum").Tracef("Searching unspent utxos for address %s: ", addr)
	address, err := btcutil.DecodeAddress(addr, q.cfg)
	if err != nil {
		log.With("module", "qtum").Tracef("Error decoding address: %s, error: %+v", addr, err)
		return nil, errors.Wrapf(err, "Error decoding address: %s", addr)
	}

	addresses := []btcutil.Address{address}
	unspent, err := q.ListUnspentMinMaxAddresses(0, 9999999, addresses)
	if err != nil {
		return nil, errors.Wrapf(err, "Error listing unspent utxos for address: %s", address.EncodeAddress())
	}
	return unspent, nil
}

// BuildUnsignedQtumTx creates a qtum/btc raw transaction using the given parameters
// to create inputs and resulting outputs.
// Returns a wire.MsgTx ready to be signed and sent to the network
//
// Params:
//   - unspent: a list of unspent outputs to use as inputs
//   - sender: the sender address in base58 format
//   - receiver: the receiver address in base58 format
//   - amount: the amount to send in Qtum
func (q *QtumClient) BuildUnsignedQtumTx(unspent []btcjson.ListUnspentResult, sender, receiver string, amount float64) (*wire.MsgTx, error) {

	//1. Create new empty transaction
	tx := wire.NewMsgTx(wire.TxVersion)

	//2. Get addresses and amounts
	if q.cfg == nil {
		return nil, errors.New("Network parameters not set in qtum client")
	}
	senderAddr, err := btcutil.DecodeAddress(sender, q.cfg)
	if err != nil {
		return nil, errors.Wrapf(err, "Error decoding sender address: %s", sender)
	}
	receiverAddr, err := btcutil.DecodeAddress(receiver, q.cfg)
	if err != nil {
		return nil, errors.Wrapf(err, "Error decoding receiver address: %s", receiver)
	}

	//3. Create inputs
	for _, utxo := range unspent {
		hash, err := chainhash.NewHashFromStr(utxo.TxID)
		if err != nil {
			return nil, errors.Wrapf(err, "Error creating chainhash: %s", utxo.TxID)
		}
		outPoint := wire.NewOutPoint(hash, utxo.Vout)
		txIn := wire.NewTxIn(outPoint, nil, nil)
		tx.AddTxIn(txIn)
	}

	//4. Create outputs

	// create receiver output
	receiverAmount, err := btcutil.NewAmount(amount)
	if err != nil {
		return nil, errors.Wrapf(err, "Error converting amount %v", amount)
	}
	receiverScript, err := txscript.PayToAddrScript(receiverAddr)
	if err != nil {
		return nil, errors.Wrapf(err, "Error creating receiver script: %s", receiverAddr.EncodeAddress())
	}
	log.With("module", "qtum").Tracef("receiverScript: %x", receiverScript)
	txOut := wire.NewTxOut(int64(receiverAmount), receiverScript)
	tx.AddTxOut(txOut)

	// TODO: research querying of gas fee dynamically

	change := calculateChange(unspent, amount)
	// create change output
	if change.GreaterThan(decimal.NewFromFloat(0)) {
		changeF, _ := change.Float64()
		changeAmount, err := btcutil.NewAmount(changeF)
		if err != nil {
			return nil, errors.Wrapf(err, "Error converting change amount %v", change)
		}
		senderScript, err := txscript.PayToAddrScript(senderAddr)
		if err != nil {
			return nil, errors.Wrapf(err, "Error creating sender script: %s", senderAddr.EncodeAddress())
		}
		log.With("module", "qtum").Tracef("senderScript: %x", senderScript)
		txOut = wire.NewTxOut(int64(changeAmount), senderScript)
		tx.AddTxOut(txOut)
	}

	return tx, nil
}

// SignRawTX signs the given raw transaction off-line using the given unspent outputs to create
// signatures for the inputs.
//
// The transaction is not sent to the network.
//
// Params:
//   - tx: the transaction to sign
//   - unspent: the list of outputs referenced by the inputs to sign
//   - w: the wallet to use for signing
func (q *QtumClient) SignRawTX(tx *wire.MsgTx, unspent []btcjson.ListUnspentResult, w wallet.IQtumWallet) error {

	// Sign inputs
	for i, txin := range tx.TxIn {
		//check unspent length
		if i >= len(unspent) {
			return errors.New("unspent length is less than txin length")
		}
		// Take the address from the input and get the private key from the wallet
		privKey, err := w.GetPrivateKey(unspent[i].Address)
		if err != nil {
			return errors.Wrapf(err, "Error getting private key for address: %s", unspent[i].Address)
		}
		scriptPubKey, err := hex.DecodeString(unspent[i].ScriptPubKey)
		if err != nil {
			return errors.Wrapf(err, "Error decoding scriptPubKey: %s", unspent[i].ScriptPubKey)
		}

		var sigScript []byte
		// check what type of scriptPubKey it is
		pKscriptType := txscript.GetScriptClass(scriptPubKey)
		switch pKscriptType {
		case txscript.PubKeyHashTy:
			log.With("module", "qtum").Tracef("scriptPubKey is P2PKH")
			sigScript, err = txscript.SignatureScript(
				tx,
				i,
				scriptPubKey,
				txscript.SigHashAll,
				privKey,
				true,
			)
			if err != nil {
				return errors.Wrapf(err, "error signing input type P2PKH with index %d", i)
			}
		case txscript.ScriptHashTy:
			// create scriptSig for P2SH
			log.With("module", "qtum").Tracef("scriptPubKey is P2SH")
			return errors.New("P2SH is not supported")
		case txscript.WitnessV0PubKeyHashTy:
			log.With("module", "qtum").Tracef("scriptPubKey is P2WPKH")
			return errors.New("P2WPKH is not supported")
		case txscript.WitnessV0ScriptHashTy:
			log.With("module", "qtum").Tracef("scriptPubKey is P2WSH")
			return errors.New("P2WSH is not supported")
		case txscript.PubKeyTy:
			log.With("module", "qtum").Tracef("scriptPubKey is P2PK")
			// create scriptSig for P2PK
			// sigHash, err := txscript.CalcSignatureHash(scriptPubKey, txscript.SigHashAll, tx, i)
			// if err != nil {
			// 	return errors.Wrapf(err, "Error calculating signature hash for input %d", i)
			// }
			signature, err := txscript.RawTxInSignature(tx, i, scriptPubKey, txscript.SigHashAll, privKey)
			if err != nil {
				return errors.Wrapf(err, "Error creating signature for input %d", i)
			}
			sigScript, err = txscript.NewScriptBuilder().AddData(signature).Script()
			if err != nil {
				return errors.Wrapf(err, "Error creating scriptSig for input %d", i)
			}

		case txscript.WitnessUnknownTy:
			// create scriptSig for P2PKH
			return errors.New("WitnessUnknown type is not supported")
		case txscript.NonStandardTy:
			// create scriptSig for non-standard
			log.With("module", "qtum").Tracef("scriptPubKey is non-standard")
			return errors.New("non-standard type is not supported")
		}
		/*
			sigScript, err := txscript.SignatureScript(
				tx,
				i,
				scriptPubKey,
				txscript.SigHashAll,
				privKey,
				true,
			)
			if err != nil {
				return errors.Wrapf(err, "Error signing input %d", i)
			}
		*/

		txin.SignatureScript = sigScript
	}

	return nil
}

// sumUTXO sums the amount of all unspent outputs in the given list
func sumUTXO(list []btcjson.ListUnspentResult) float64 {
	var sum float64
	for _, utxo := range list {
		sum += utxo.Amount
	}
	return sum
}

// calculateChange calculates the change amount due to the sender
func calculateChange(unspent []btcjson.ListUnspentResult, amount float64) decimal.Decimal {
	// TODO: research querying of gas fee dynamically
	// Estimate gas fee
	// txSize := int64(tx.SerializeSize() / 1000)
	// log.With("module", "qtum").Tracef("Estimated tx size: %v KB", txSize)
	// gas, err := q.EstimateFee(txSize)
	// gas, err := q.EstimateSmartFee(6, nil)
	// if err != nil {
	// 	return nil, errors.Wrapf(err, "Error estimating gas fee")
	// }
	// log.With("module", "qtum").Tracef("Estimated gas fee: %+v", gas)
	// if gas.FeeRate == nil {
	// 	*gas.FeeRate = 100000
	// }

	// ! hardcoded gas fee in satoshis
	gasSatoshis := DefaultGasPrice
	gas := decimal.NewFromFloatWithExponent(float64(gasSatoshis)/Qtum, PrecisionExp)
	utxoTotalAmount := decimal.NewFromFloatWithExponent(sumUTXO(unspent), PrecisionExp)
	amountToSend := decimal.NewFromFloatWithExponent(amount, PrecisionExp)
	change := utxoTotalAmount.Sub(amountToSend).Sub(gas)
	log.With("module", "qtum").Debugf("total amount value: %v, total utxo value: %v, gasFee value %v, change value: %v", amount, utxoTotalAmount, gas, change)
	return change

}
