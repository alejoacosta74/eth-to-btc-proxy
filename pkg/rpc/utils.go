package rpc

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/alejoacosta74/rpc-proxy/pkg/log"
	"github.com/qtumproject/btcd/wire"
)

// printQtumDecodedTX prints a decoded QTUM transaction
func (api *EthAPI) printQtumDecodedTX(qtumTx *wire.MsgTx, msg string) {
	var buf bytes.Buffer
	err := qtumTx.Serialize(&buf)
	if err != nil {
		log.With("method", "sendrawtx").Debugf(err.Error())
		return
	}
	decoded, err := api.qcli.DecodeRawTransaction(buf.Bytes())
	if err != nil {
		log.With("method", "sendrawtx").Debugf(err.Error())
		return
	}
	decodedBytes, err := json.Marshal(decoded)
	if err != nil {
		log.With("method", "sendrawtx").Debugf(err.Error())
		return
	}
	logPretty(msg, decodedBytes)
}

// logPretty logs a JSON string in a json indented format
func logPretty(msg string, output []byte) {
	if len(output) > 0 {
		var prettyJSON bytes.Buffer
		if err := json.Indent(&prettyJSON, output, "", "    "); err != nil {
			fmt.Printf("Error decoding JSON: %v\n", err)
		} else {
			fmt.Printf("\n%s :\n%s\n", msg, prettyJSON.String())

		}
	}
}
