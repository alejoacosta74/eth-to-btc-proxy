package qtum

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"math"
	"testing"

	utils "github.com/alejoacosta74/rpc-proxy/pkg/internal/testutils"
	"github.com/alejoacosta74/rpc-proxy/pkg/wallet"
	"github.com/qtumproject/btcd/btcec/v2"
	"github.com/qtumproject/btcd/btcjson"
	"github.com/qtumproject/btcd/txscript"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

const (
	SENDER_ADDR    = "qUbxboqjBRp96j3La8D1RYkyqx5uQbJPoW"
	SENDER_PRIVKEY = "00821d8c8a3627adc68aa4034fea953b2f5da553fab312db3fa274240bd49f35"
	RECEIVER_ADDR  = "qeVQ5JF6idPcrg1u9M3pCryXeebpj3Tbpk"
)

// Mock responses from qtumd
const listUnspentResponseJSON string = `[
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
var cfg = utils.GetNetworkParams()
var gasSatoshis int = 10000
var gas float64 = float64(gasSatoshis) / 100000000

func TestBuildUnsignedQtumTx(t *testing.T) {
	err := json.Unmarshal([]byte(listUnspentResponseJSON), &listUnspentResponse)
	utils.HandleFatalError(t, err)

	var unspent1 []btcjson.ListUnspentResult
	unspent1 = append(unspent1, listUnspentResponse[0])

	var unspent2 []btcjson.ListUnspentResult
	unspent2 = append(unspent2, listUnspentResponse[1])

	unspent3 := append(unspent1, unspent2...)

	// Create qtum client
	qcli, _ := NewQtumClient("", "qtum", "qtumpass", cfg.Net.String())

	// Define and run tests
	tests := []struct {
		name         string
		unspent      []btcjson.ListUnspentResult
		outputAmount float64
		wantChange   float64
		wantOutputs  int
		wantErr      bool
	}{
		{
			name:         "test1: one input, amount < input.amount",
			unspent:      unspent1,
			outputAmount: 10000.1,
			wantChange:   9999.9 - gas,
			wantOutputs:  2,
			wantErr:      false,
		},
		{
			name:         "test2: one input, amount = input.amount-gas",
			unspent:      unspent1,
			outputAmount: 20000.0000000 - gas,
			wantChange:   0,
			wantOutputs:  1,
			wantErr:      false,
		},
		{
			name:         "test3: two inputs, amount < inputs.amount",
			unspent:      unspent3,
			outputAmount: 30000.3,
			wantChange:   9999.7 - gas,
			wantOutputs:  2,
			wantErr:      false,
		},
		{
			name:         "test4: two inputs, amount = inputs.amount-gas",
			unspent:      unspent3,
			outputAmount: 40000 - gas,
			wantChange:   0,
			wantOutputs:  1,
			wantErr:      false,
		},
		{
			name:         "test5: two inputs, amount with decimals",
			unspent:      unspent3,
			outputAmount: 30000.8,
			wantChange:   9999.2 - gas,
			wantOutputs:  2,
			wantErr:      false,
		},
		{
			name:         "test6: two inputs, amount == inputs.amount",
			unspent:      unspent3,
			outputAmount: 40000.0000000,
			wantChange:   0,
			wantOutputs:  1,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			unsignedTx, err := qcli.BuildUnsignedQtumTx(tt.unspent, SENDER_ADDR, RECEIVER_ADDR, tt.outputAmount)
			// check error
			if err != nil {
				if tt.wantErr {
					return
				}
				t.Fatal(err)
			}

			// check receiver address is correct in tx output
			scriptBytes := unsignedTx.TxOut[0].PkScript
			script, err := txscript.ParsePkScript(scriptBytes)
			if err != nil {
				t.Fatal(err)
			}
			receiver, err := script.Address(cfg)
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, RECEIVER_ADDR, receiver.String())

			// check receiver amount is correct in tx output
			receiverSatoshis := unsignedTx.TxOut[0].Value
			receiverAmount := decimal.NewFromFloatWithExponent(float64(receiverSatoshis)/Qtum, PrecisionExp)
			if !receiverAmount.Equals(decimal.NewFromFloatWithExponent(tt.outputAmount, PrecisionExp)) {
				t.Errorf("receiver amount is not correct, want %v, got %v", tt.outputAmount, receiverAmount)
			}

			// check number of outputs is correct
			assert.Equal(t, tt.wantOutputs, len(unsignedTx.TxOut))

			// check change amount is correct
			if tt.wantOutputs > 1 {
				changeSatoshis := unsignedTx.TxOut[1].Value
				changeQtum := float64(changeSatoshis) / Qtum
				if !utils.AreEqual(changeQtum, tt.wantChange) {
					t.Errorf("change amount is not correct, want %v, got %v", tt.wantChange, changeQtum)
				}

			}
		})
	}

}

func round(num float64, decimalPlaces int) float64 {
	shift := math.Pow(10, float64(decimalPlaces))
	return math.Round(num*shift) / shift
}

func TestSignRawTx(t *testing.T) {

	// Mocked unspent inputs to build unsigned tx
	listUnspent := []btcjson.ListUnspentResult{}
	err := json.Unmarshal([]byte(listUnspentResponseJSON), &listUnspent)
	utils.HandleFatalError(t, err)

	inputs := listUnspent[0:2]

	// Create qtum client
	qcli, _ := NewQtumClient("", "qtum", "qtumpass", cfg.Net.String())

	// Build unsigned tx
	amount := 10000.1
	tx, err := qcli.BuildUnsignedQtumTx(inputs, SENDER_ADDR, RECEIVER_ADDR, amount)
	utils.HandleFatalError(t, err)

	// Mocked wallet
	wallet, err := wallet.NewQtumWallet(SENDER_PRIVKEY, cfg)
	utils.HandleFatalError(t, err)

	// Sign tx
	err = qcli.SignRawTX(tx, inputs, wallet)
	utils.HandleFatalError(t, err)

	// Check all inputs are signed by sender
	for i, txin := range tx.TxIn {
		// check sigScript is not empty
		sigScript := txin.SignatureScript
		if len(sigScript) == 0 {
			t.Fatal("tx is not signed")
		}
		// Check input is signed by sender

		// 1. Get pubkey from sigScript
		signaturePubKey := txin.SignatureScript[len(sigScript)-33 : len(sigScript)]

		// 2. Get pubkey from sender private key
		privKeyBytes, _ := hex.DecodeString(SENDER_PRIVKEY)
		_, senderPubKey := btcec.PrivKeyFromBytes(privKeyBytes)

		// 3. Check pubkey is the same
		if !bytes.Equal(signaturePubKey, senderPubKey.SerializeCompressed()) {
			t.Fatalf("tx %d is not signed by sender", i)
		}
		// TODO: add signature verification for each input
		// 1. Get signature from sigScript
		// 2. Get hash from tx
		// 3. Verify signature
	}

}
