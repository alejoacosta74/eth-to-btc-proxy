package rpc

import (
	"net/http"
	"net/http/httptest"
	"testing"

	data "github.com/alejoacosta74/rpc-proxy/pkg/internal/testdata"
	utils "github.com/alejoacosta74/rpc-proxy/pkg/internal/testutils"
	"github.com/alejoacosta74/rpc-proxy/pkg/wallet"
	"github.com/stretchr/testify/assert"
)

func TestImportRawKey(t *testing.T) {
	assert := assert.New(t)
	rpcservice, err := getPersonalRPCService()
	handleFatalError(t, err)
	testserver := httptest.NewServer(rpcservice)
	defer testserver.Close()

	client := &http.Client{}

	t.Run("Verify derived address from private key", func(t *testing.T) {
		httpReq, err := createRPCRequest(testserver.URL, "personal_importRawKey", data.PrivKeyString, data.KeyStorePassphrase)
		handleFatalError(t, err)

		httpResp, err := client.Do(httpReq)
		handleFatalError(t, err)

		var got string
		err = utils.ReadJSONResult(httpResp, &got)
		assert.NoError(err)
		want := data.EthereumAddress
		assert.Equal(want, got)
	})
	t.Run("Verify wallet is created for address", func(t *testing.T) {
		ws := wallet.GetWallets()
		assert.NotNil(ws)
		w, err := ws.SeekWallet(data.EthereumAddress)
		assert.NoError(err)
		assert.NotNil(w)
	})

}
