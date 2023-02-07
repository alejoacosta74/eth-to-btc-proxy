package testutils

import (
	"encoding/json"
	"fmt"

	// "errors"
	"io"
	"net/http"
	"os"

	"github.com/pkg/errors"
)

// CreateJSONRequest creates a JSON RPC request with the given method and params
// and returns the JSON encoded request as a byte array.
func CreateJSONRequest(method string, paramsStr ...string) ([]byte, error) {
	params := make([]json.RawMessage, len(paramsStr))
	for i, param := range paramsStr {
		params[i] = json.RawMessage(param)
	}
	jsonReq := NewJSONRPCRequest(1, method, params)
	return json.Marshal(jsonReq)
}

// ReadJSONResult reads the JSON response from the given http response and
// unmarshals the result into the given result interface.
func ReadJSONResult(resp *http.Response, result interface{}) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.New("error reading body: " + err.Error())
	}
	var jsonResp *JSONRPCResponse
	err = json.Unmarshal(body, &jsonResp)
	if err != nil {
		return errors.New("error unmarshaling JSON response: " + err.Error())
	}
	if jsonResp.Error != nil {
		return errors.New("JSON Response with error: " + jsonResp.Error.Message + " (code: " + fmt.Sprint(jsonResp.Error.Code) + ")")
	}
	err = json.Unmarshal(jsonResp.Result, result)
	if err != nil {
		return errors.New("error unmarshaling JSON result: " + err.Error())
	}
	return nil
}

// Reads a JSON file and unmarshals it into the given object
func ReadJSONFromFile(filename string, obj interface{}) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	return decoder.Decode(obj)
}
