package rpc

import (
	"bytes"
	"context"
	"encoding/hex"
	"testing"

	"github.com/alejoacosta74/rpc-proxy/pkg/internal/mocks"
	utils "github.com/alejoacosta74/rpc-proxy/pkg/internal/testutils"
	"github.com/alejoacosta74/rpc-proxy/pkg/wallet"
	"github.com/ethereum/go-ethereum/core/types"
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
