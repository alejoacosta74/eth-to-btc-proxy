package rpc

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// logPretty prints out a JSON byte array as a JSON indented format
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
