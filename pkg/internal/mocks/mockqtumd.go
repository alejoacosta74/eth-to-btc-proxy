package mocks

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"

	utils "github.com/alejoacosta74/qproxy/pkg/internal/testutils"
	"github.com/qtumproject/btcd/btcjson"
)

type MockQtumd struct {
	*httptest.Server
}

// NewMockQtumd creates a new mock qtumd server with the given handler.
// Params:
//   - responses: a map of method names to json responses
func NewMockQtumd(responses map[string]string) *MockQtumd {
	handlerFunc := responseHandler(responses)
	handler := http.HandlerFunc(handlerFunc)
	return &MockQtumd{
		Server: httptest.NewServer(handler),
	}
}

// ResponseHandler creates a handler function that for a given JSON RPC method
// returns the corresponding response. The handler function is used to create
// a mock qtumd server.
func responseHandler(responses map[string]string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// read request body
		var jsonreq btcjson.Request
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&jsonreq)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// find response for method
		method := jsonreq.Method
		response, ok := responses[method]
		if !ok {
			http.Error(w, "method not found: "+method, http.StatusInternalServerError)
			return
		}
		// write response
		w.Header().Set("Content-Type", "application/json")
		result := []byte(response)
		jsonResp := utils.NewJSONRPCResponse(1, result, nil)
		resp, err := json.Marshal(jsonResp)
		if err != nil {
			panic(err)
		}
		w.Write(resp)
	}
}
