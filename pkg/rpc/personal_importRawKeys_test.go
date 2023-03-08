package rpc

import (
	"net/http"
	"net/http/httptest"
	"testing"

	utils "github.com/alejoacosta74/qproxy/pkg/internal/testutils"
	"github.com/alejoacosta74/qproxy/pkg/wallet"
	"github.com/stretchr/testify/assert"
)

const (
	// Test data
	ETHEREUM_ADDR     string = "0x96216849c49358B10257cb55b28eA603c874b05E"
	PRIVATEKEY_HEX    string = "fad9c8855b740a0b7ed4c221dbad0f33a83a49cad6b3fe8d5817ac83d38b6a19"
	WALLET_PASSPHRASE string = "test123"
)

func TestImportRawKey(t *testing.T) {
	assert := assert.New(t)
	rpcservice, err := getPersonalRPCService()
	utils.HandleFatalError(t, err)
	testserver := httptest.NewServer(rpcservice)
	defer testserver.Close()

	client := &http.Client{}

	t.Run("Verify derived address from private key", func(t *testing.T) {
		httpReq, err := createRPCRequest(testserver.URL, "personal_importRawKey", PRIVATEKEY_HEX, WALLET_PASSPHRASE)
		utils.HandleFatalError(t, err)

		httpResp, err := client.Do(httpReq)
		utils.HandleFatalError(t, err)

		var got string
		err = utils.ReadJSONResult(httpResp, &got)
		assert.NoError(err)
		want := ETHEREUM_ADDR
		assert.Equal(want, got)
	})
	t.Run("Verify wallet is created for address", func(t *testing.T) {
		ws := wallet.GetWallets()
		assert.NotNil(ws)
		w, err := ws.SeekWallet(ETHEREUM_ADDR)
		assert.NoError(err)
		assert.NotNil(w)
	})

}
