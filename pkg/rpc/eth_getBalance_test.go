package rpc

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/alejoacosta74/rpc-proxy/pkg/internal/mocks"
	"github.com/qtumproject/btcd/btcjson"
)

func TestGetBalance(t *testing.T) {

	// create a new mock qtum client
	mockQcli := mocks.NewMockQCli()

	// create mock responses
	var defaultListUnspentResponse []btcjson.ListUnspentResult
	_ = json.Unmarshal([]byte(mocks.DefaultListUnspentResponseJSON), &defaultListUnspentResponse)

	var emptyListUnspentResponse = []btcjson.ListUnspentResult{}

	tests := []struct {
		name     string
		response []btcjson.ListUnspentResult
		want     string
		wantErr  bool
	}{
		{
			name:     "GetBalance with default LinstUnspent response",
			response: defaultListUnspentResponse,
			// test balance is 60.000 qtum
			// 60.000 qtum = 60.000 eth
			// 60.000 eth = 60.000 * 10^18 wei
			// 60.000 * 10^18 wei = 60000000000000000000000
			// 60000000000000000000000 wei = 0xCB49B44BA602D800000
			want:    strings.ToLower("0xCB49B44BA602D800000"),
			wantErr: false,
		},
		{
			name:     "GetBalance with empty LinstUnspent response",
			response: emptyListUnspentResponse,
			want:     "0x0",
			wantErr:  false,
		},
	}

	address := "0x96216849c49358B10257cb55b28eA603c874b05E"

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockQcli.FindSpendableUTXOResult = tt.response
			api := NewAPI(context.Background(), mockQcli)
			api.SetNetworkParams(cfg)
			ethAPI := (*EthAPI)(api)

			got, err := ethAPI.GetBalance(address, 0)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBalance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetBalance() got = %v, want %v", got, tt.want)
			}
		})
	}

}
