package rpc

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/alejoacosta74/qproxy/pkg/internal/mocks"
	utils "github.com/alejoacosta74/qproxy/pkg/internal/testutils"
	"github.com/alejoacosta74/qproxy/pkg/wallet"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/qtumproject/btcd/btcjson"
)

var cfg = utils.GetNetworkParams()

func TestSendRawTx(t *testing.T) {
	const PRIVATEKEY = "00821d8c8a3627adc68aa4034fea953b2f5da553fab312db3fa274240bd49f35"

	// create a new mock qtum client
	mockQcli := mocks.NewMockQCli()

	// create a new eth api
	api := NewAPI(context.Background(), mockQcli)
	api.SetNetworkParams(cfg)
	ethAPI := (*EthAPI)(api)

	// create a new wallet
	ws := wallet.GetWallets()
	_, err := ws.NewWallet(PRIVATEKEY, cfg)
	utils.HandleFatalError(t, err)

	// create a raw ethereum transaction
	tx := newEthereumTx()
	signer := types.HomesteadSigner{}
	signedTx, err := signEthereumTx(tx, signer, PRIVATEKEY)
	if err != nil {
		t.Fatalf("error signing transaction: %v", err)
	}
	buf := new(bytes.Buffer)
	if err := signedTx.EncodeRLP(buf); err != nil {
		t.Fatalf("error encoding transaction: %v", err)
	}
	rlpEncodedTx := buf.Bytes()
	rlpEncodedTxHex := hex.EncodeToString(rlpEncodedTx)

	// call the eth_sendRawTransaction method
	got, err := ethAPI.SendRawTransaction(rlpEncodedTxHex)
	if err != nil {
		t.Fatalf("error calling eth_sendRawTransaction: %v", err)
	}
	want := signedTx.Hash().Hex()
	if got.Hash != want {
		t.Fatalf("got %v, want %v", got, want)
	}

}

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
var gasSatoshis int = 10000
var gas float64 = float64(gasSatoshis) / 100000000

func TestGetUTXOtoSpend(t *testing.T) {

	err := json.Unmarshal([]byte(listUnspentResponseJSON), &listUnspentResponse)
	utils.HandleFatalError(t, err)

	// Define and run tests
	tests := []struct {
		name    string
		address string
		amount  float64
		want    []btcjson.ListUnspentResult
		wantErr bool
	}{
		{
			name:    "Test 1",
			address: "qUbxboqjBRp96j3La8D1RYkyqx5uQbJPoW",
			amount:  20000,
			want:    listUnspentResponse[0:1],
			wantErr: false,
		},
		{
			name:    "Test 2",
			address: "qUbxboqjBRp96j3La8D1RYkyqx5uQbJPoW",
			amount:  40000,
			want:    listUnspentResponse[0:2],
			wantErr: false,
		},
		{
			name:    "Test 3",
			address: "qUbxboqjBRp96j3La8D1RYkyqx5uQbJPoW",
			amount:  40001,
			want:    listUnspentResponse[0:3],
			wantErr: false,
		},
		{
			name:    "Test 4",
			address: "qUbxboqjBRp96j3La8D1RYkyqx5uQbJPoW",
			amount:  80001,
			want:    listUnspentResponse[0:3],
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getUTXOtoSpend(listUnspentResponse, tt.amount)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("qcli.FindSpendableUTXO() error = %v, wantErr %v", err, tt.wantErr)
					return
				} else {
					return
				}
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("QtumClient.FindSpendableUTXO() = \n%v\n ====>want \n%v\n", got, tt.want)
			}
		})
	}

}
