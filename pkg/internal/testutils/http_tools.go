package testutils

import (
	"bytes"
	"io"
	"net/http"
)

// CreateHTTPRequest creates a new HTTP request with the given http method and url
func CreateHTTPRequest(httpMethod string, url string, buf *bytes.Buffer) (*http.Request, error) {
	req, err := http.NewRequest("POST", url, buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

// NewEchoHandler creates a handler that echoes the received body back to the caller
func NewEchoHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// echo back the request body
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Write(bodyBytes)
	})
}

// Creates a new JSON RPC request with the given methods and arguments/params
func CreateRPCRequest(url string, method string, args ...string) (*http.Request, error) {
	params := []string{}
	for _, arg := range args {
		arg = "\"" + arg + "\""
		params = append(params, arg)
	}
	jsonReq, err := CreateJSONRequest(method, params...)
	if err != nil {
		return nil, err
	}
	return CreateHTTPRequest("POST", url, bytes.NewBuffer(jsonReq))
}
