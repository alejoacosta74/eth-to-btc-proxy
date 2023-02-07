package rpc

import (
	"net/http"
	"net/http/httptest"
	"testing"

	utils "github.com/alejoacosta74/rpc-proxy/pkg/internal/testutils"
	"github.com/alejoacosta74/rpc-proxy/pkg/qtum"
	"github.com/alejoacosta74/rpc-proxy/pkg/types"
	"github.com/alejoacosta74/rpc-proxy/pkg/wallet"
	"github.com/stretchr/testify/assert"
)

var cfg = utils.GetNetworkParams()

func TestDecodeRawTx(t *testing.T) {
	type inputData struct {
		Rawtx string
	}
	var input = new(inputData)

	var want = new(types.RawTransaction)

	var fn utils.TestingFn = func(t *testing.T, testname string, i int, inputPath string, goldenPath []string, inputIntfc interface{}, wantIntfc interface{}) {
		t.Run(testname, func(t *testing.T) {
			want := wantIntfc.(*types.RawTransaction)
			input := inputIntfc.(*inputData)

			got, err := decodeRawTx(input.Rawtx)
			assert.NoError(t, err)
			assert.Equal(t, want, got)
		})
	}

	utils.RunTests(t, fn, "rpc", "decoderawtransaction", input, want)

}
func TestSendRawTx(t *testing.T) {

	type inputData struct {
		Rawtx   string
		PrivKey string
		WIF     string
		Address string
	}

	type wantData struct {
		QtumHash string
		Error    bool
	}

	var fn utils.TestingFn = func(t *testing.T, testname string, i int, inputPath string, goldenPath []string, inputIntfc interface{}, wantIntfc interface{}) {
		t.Run(testname, func(t *testing.T) {
			want := wantIntfc.(*wantData)
			input := inputIntfc.(*inputData)

			mockQtumd := utils.NewMockQtumRpcServer(testname)
			defer mockQtumd.Close()

			// Create a new qtum client
			qcli, err := qtum.NewQtumClient(mockQtumd.URL, "qtum", "qtumpass", cfg.Net.String())
			handleFatalError(t, err)
			// Create a rpc service handler and a httptest server
			rpcservice, err := getETHRPCService(cfg, qcli)
			handleFatalError(t, err)
			testserver := httptest.NewServer(rpcservice)
			defer testserver.Close()

			// Create a http client
			client := &http.Client{}

			// Get wallets to use for running the test
			ws := wallet.GetWallets()

			// create wallet for test
			_, err = ws.NewWallet(input.PrivKey, cfg)
			utils.HandleFatalError(t, err)
			httpReq, err := createRPCRequest(testserver.URL, "eth_sendRawTransaction", input.Rawtx)
			handleFatalError(t, err)

			httpResp, err := client.Do(httpReq)
			handleFatalError(t, err)

			var got types.Eth_SendRawTransactionResponse
			err = utils.ReadJSONResult(httpResp, &got)
			handleFatalError(t, err)

			assert.Equal(t, want.QtumHash, got.Hash)
		})

	}

	var input = new(inputData)
	var want = new(wantData)

	utils.RunTests(t, fn, "rpc", "sendrawtransaction", input, want)
}
