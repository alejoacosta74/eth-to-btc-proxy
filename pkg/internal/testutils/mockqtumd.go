package testutils

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"

	"github.com/alejoacosta74/rpc-proxy/pkg/internal/testdata"
	"github.com/alejoacosta74/rpc-proxy/pkg/log"

	"github.com/qtumproject/btcd/btcjson"
	"github.com/qtumproject/btcd/chaincfg/chainhash"
)

// supported RPC methods for mock qtumd
const (
	GETRAWTRANSACTION     = "getrawtransaction"
	GETBLOCK              = "getblock"
	GETTXOUT              = "gettxout"
	GETTRANSACTIONRECEIPT = "gettransactionreceipt"
	GETTRANSACTION        = "gettransaction"
	GETADDRESSINFO        = "getaddressinfo"
	IMPORTADDRESS         = "importaddress"
	GETINFO               = "getinfo"
	LISTUNSPENT           = "listunspent"
	SENDRAWTRANSACTION    = "sendrawtransaction"
)

var (
	utxos       []btcjson.ListUnspentResult // result from 'listunspent' call
	hash        chainhash.Hash
	addressInfo btcjson.GetAddressInfoResult
	infoWallet  btcjson.InfoWalletResult
)

var qtumHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	var jsonreq btcjson.Request
	decoder := json.NewDecoder(r.Body)
	var err error
	err = decoder.Decode(&jsonreq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	method := jsonreq.Method
	log.With("module", "mockqtumd").Debugf("received method: %s", method)
	var response interface{}

	switch method {
	case "getrawtransaction":
		var txid string
		err = json.Unmarshal(jsonreq.Params[0], &txid)
		if err != nil {
			break
		}
		var qRawTxResult btcjson.TxRawResult
		jsonRawTxResult, ok := testdata.QtumGetRawTxMapSample[txid]

		if !ok {
			http.Error(w, "tx hash not found by mock qtumd", http.StatusInternalServerError)
			return
		}
		err = json.Unmarshal(jsonRawTxResult, &qRawTxResult)
		if err != nil {
			break
		}
		response = qRawTxResult
	case "getblock":
		var blockId string
		err := json.Unmarshal(jsonreq.Params[0], &blockId)
		if err != nil {
			break
		}
		var qBlockResult btcjson.GetBlockVerboseResult
		blockOutput, ok := testdata.QtumGetBlockMapSample[blockId]
		if !ok {
			http.Error(w, "block not found", http.StatusInternalServerError)
			return
		}
		err = json.Unmarshal(blockOutput, &qBlockResult)
		if err != nil {
			break
		}
		response = qBlockResult
	case "gettxout":
		var txid string
		err := json.Unmarshal(jsonreq.Params[0], &txid)
		if err != nil {
			break
		}
		var vout uint32
		err = json.Unmarshal(jsonreq.Params[1], &vout)
		if err != nil {
			break
		}
		var qTxOutResult btcjson.GetTxOutResult
		txOutOutput, ok := testdata.QtumGetTxOutMapSample[txid][vout]
		if !ok {
			http.Error(w, "txout not found", http.StatusInternalServerError)
			return
		}
		err = json.Unmarshal(txOutOutput, &qTxOutResult)
		if err != nil {
			break
		}
		response = qTxOutResult
	case "gettransactionreceipt":
	case "getaddressinfo":
		response = addressInfo
	case "importaddress":
	case "getinfo":
		response = infoWallet
	case "listunspent":
		response = utxos
	case "sendrawtransaction":
		response = testdata.SendRawTXResult
	}
	log.With("module", "mockqtumd").Tracef("mockqtumd sending back response: %+v", response)

	var jsonResp *btcjson.Response
	if err != nil {
		jsonErr := btcjson.NewRPCError(btcjson.ErrRPCInvalidParameter, err.Error())
		jsonResp, _ = btcjson.NewResponse(btcjson.RpcVersion2, 1, nil, jsonErr)
	} else {
		responseBytes, _ := json.Marshal(response)
		jsonResp, _ = btcjson.NewResponse(btcjson.RpcVersion2, 1, responseBytes, nil)
	}

	jsonRespBytes, err := json.Marshal(jsonResp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonRespBytes)
})

func loadTestData(testname string) error {
	// read files named qtum*golden.json from testdata folder
	filenames := MOCKQTUMD_PREFIX + "*" + MOCKQTUMD_SUFFIX + ".json"
	path := filepath.Join(DATA_PATH, testname, filenames)
	paths, err := filepath.Glob(path)
	if err != nil {
		return err
	}

	// read each file and unmarshal to the corresponding struct
	var output interface{}
	for _, path := range paths {
		methodname := strings.TrimSuffix(filepath.Base(path), MOCKQTUMD_SUFFIX+".json")
		methodname = strings.TrimPrefix(methodname, MOCKQTUMD_PREFIX)

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		switch methodname {
		case GETADDRESSINFO:
			output = addressInfo
		case IMPORTADDRESS:
		case GETINFO:
			output = infoWallet
		case LISTUNSPENT:
			output = utxos
		case SENDRAWTRANSACTION:
			output = hash
		}

		if methodname == LISTUNSPENT {
			err = json.Unmarshal(content, &utxos)
		} else {
			err = json.Unmarshal(content, &output)
		}
		if err != nil {
			return err
		}

	}
	return nil

}

func NewMockQtumRpcServer(testname string) *httptest.Server {
	if testname == "" {
		panic("testname is empty")
	}
	err := loadTestData(testname)
	if err != nil {
		panic(err)
	}
	log.With("module", "mockqtumd").Debugf("starting new mockqtumd server")
	return httptest.NewServer(qtumHandler)
}
