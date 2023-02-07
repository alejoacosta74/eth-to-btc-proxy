package qtum

import (
	"testing"

	utils "github.com/alejoacosta74/rpc-proxy/pkg/internal/testutils"
	"github.com/alejoacosta74/rpc-proxy/pkg/types"
	"github.com/qtumproject/btcd/btcjson"
	qcommon "github.com/qtumproject/qtool/lib/common"
	qtool "github.com/qtumproject/qtool/lib/tools"
	"github.com/stretchr/testify/assert"
)

var cfg = utils.GetNetworkParams()

func TestFindSpendableUTXO(t *testing.T) {
	assert := assert.New(t)
	testnames := utils.ReadTestNames(t)
	var input = struct {
		Address string
		Amount  float64
	}{}

	var want = struct {
		Unspent []btcjson.ListUnspentResult
		Amount  float64
		Error   bool
	}{}
	for _, testname := range testnames {
		// Create a mock qtumd server
		mockQtumd := utils.NewMockQtumRpcServer(testname)
		defer mockQtumd.Close()

		// Create a qtum client
		qcli, err := NewQtumClient(mockQtumd.URL, "qtum", "qtumpass", cfg.Net.String())
		utils.HandleFatalError(t, err)

		// Get the input and golden test data files
		inputPaths, goldenPaths := utils.GetTestFilePaths(t, testname, "qtum", "listunspent")
		for i, inputPath := range inputPaths {
			t.Run(testname, func(t *testing.T) {

				// Load the input data for this particular test
				utils.LoadDataFromFile(t, inputPath, &input)

				// Call the method
				unspent, err := qcli.FindSpendableUTXO(input.Address, input.Amount)

				// Load the golden data for this particular test
				utils.LoadDataFromFile(t, goldenPaths[i], &want)

				if want.Error {
					assert.Error(err)
					return
				} else {
					sum := sumUTXO(unspent)
					assert.NoError(err)
					assert.Equal(want.Amount, sum)
					assert.Equal(want.Unspent, unspent)

				}
			})
		}
	}
}

func TestPrepareRawTransaction(t *testing.T) {
	assert := assert.New(t)
	testnames := utils.ReadTestNames(t)
	var input = struct {
		Unspent      []btcjson.ListUnspentResult
		DecodedEthTx types.RawTransaction
	}{}

	var want = struct {
		TxHash string
	}{}
	for _, testname := range testnames {
		inputPaths, goldenPaths := utils.GetTestFilePaths(t, testname, "qtum", "preparerawtx")
		for i, inputPath := range inputPaths {
			t.Run(testname, func(t *testing.T) {

				mockQtumd := utils.NewMockQtumRpcServer(testname)
				defer mockQtumd.Close()

				qcli, err := NewQtumClient(mockQtumd.URL, "user", "pass", cfg.Net.String())
				utils.HandleFatalError(t, err)

				// Load the input data for this particular test
				utils.LoadDataFromFile(t, inputPath, &input)
				unspent := input.Unspent
				decodedEthTx := input.DecodedEthTx
				sender, err := qtool.ConvertAddressHexToBase58(decodedEthTx.From, cfg.Name, cfg.Net.String())
				utils.HandleFatalError(t, err)
				receiver, err := qtool.ConvertAddressHexToBase58(decodedEthTx.To, cfg.Name, cfg.Net.String())
				utils.HandleFatalError(t, err)
				amount, err := qcommon.ConvertWeiToQtum(decodedEthTx.Value)
				utils.HandleFatalError(t, err)

				rawTx, err := qcli.PrepareRawTransaction(unspent, sender.String(), receiver.String(), amount)
				assert.NoError(err)

				utils.LoadDataFromFile(t, goldenPaths[i], &want)
				assert.Equal(want.TxHash, rawTx.TxHash().String())

			})

		}
	}

}
