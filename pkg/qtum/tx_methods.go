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
// with a total amount greater than the given amount
func (q *QtumClient) FindSpendableUTXO(addr string, amount float64) ([]btcjson.ListUnspentResult, error) {

	// ensure the address is known to the node's wallet
	err := q.VerifyAddress(addr)
	if err != nil {
		return nil, errors.Wrapf(err, "Error verifying address: %s", addr)
	}
	log.With("module", "qtum").Tracef("Searching unspent utxos for address %s with a total amount of %v: ", addr, amount)
	address, err := btcutil.DecodeAddress(addr, q.cfg)
	if err != nil {
		return nil, errors.Wrapf(err, "Error decoding address: %s", addr)
	}

	addresses := []btcutil.Address{address}
	unspent, err := q.ListUnspentMinMaxAddresses(0, 9999999, addresses)
	if err != nil {
		return nil, errors.Wrapf(err, "Error listing unspent utxos for address: %s", address.EncodeAddress())
	}
	var total float64
	// TODO: sort by amount?
	for i, utxo := range unspent {
		if utxo.Confirmations < 6 {
			continue
		}
		total += utxo.Amount
		if total >= amount {
			log.With("module", "qtum").Tracef("Unspent utxos returned: %+v", unspent)
			return unspent[:i+1], nil
		}
	}
	log.With("module", "qtum").Tracef("not enough UTXOs found: %+v", unspent)
	return unspent, errors.New("not enough UTXOs found")
}

// PrepareRawTransaction creates a qtum/btc raw transaction using the given parameters
// to create inputs and resulting outputs
func (q *QtumClient) PrepareRawTransaction(unspent []btcjson.ListUnspentResult, sender, receiver string, amount float64) (*wire.MsgTx, error) {

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
	// ! hardcoded gas fee
	gas := float64(100000)

	// Calculate change
	utxoTotalAmount := sumUTXO(unspent)
	change := utxoTotalAmount - amount - gas
	log.With("module", "qtum").Debugf("Total utxo value: %v, change value: %v", utxoTotalAmount, change)

	if change > 0 {
		changeAmount, err := btcutil.NewAmount(change)
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
func (q *QtumClient) SignRawTX(tx *wire.MsgTx, unspent []btcjson.ListUnspentResult, w *wallet.QtumWallet) error {

	// Sign inputs
	for i, utxo := range unspent {
		privKey, err := w.GetPrivateKey(utxo.Address)
		if err != nil {
			return errors.Wrapf(err, "Error getting private key for address: %s", utxo.Address)
		}
		scriptPubKey, err := hex.DecodeString(utxo.ScriptPubKey)
		if err != nil {
			return errors.Wrapf(err, "Error decoding scriptPubKey: %s", utxo.ScriptPubKey)
		}
		tx.TxIn[i].SignatureScript, err = txscript.SignatureScript(tx, i, scriptPubKey, txscript.SigHashAll, privKey, true)
		if err != nil {
			return errors.Wrapf(err, "Error signing input %d", i)
		}

	}
	return nil
}

func sumUTXO(list []btcjson.ListUnspentResult) float64 {
	var sum float64
	for _, utxo := range list {
		sum += utxo.Amount
	}
	return sum
}
