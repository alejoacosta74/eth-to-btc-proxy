package qtum

import (
	"testing"

	"github.com/alejoacosta74/rpc-proxy/pkg/internal/mocks"
	utils "github.com/alejoacosta74/rpc-proxy/pkg/internal/testutils"
	"github.com/stretchr/testify/assert"
)

func TestVerifyAddress(t *testing.T) {
	// create mock responses
	const getwalletinfoJSON = `{
		"walletname": "example",
		"walletversion": 169900,
		"format": "bdb",
		"balance": 0.00000000,
		"stake": 0.00000000,
		"unconfirmed_balance": 0.00000000,
		"immature_balance": 0.00000000,
		"txcount": 0,
		"keypoololdest": 1676593362,
		"keypoolsize": 1000,
		"hdseedid": "c3b3fbfcb004572f5108368dd76665b9491a1b39",
		"keypoolsize_hd_internal": 1000,
		"paytxfee": 0.00000000,
		"private_keys_enabled": true,
		"avoid_reuse": false,
		"scanning": false,
		"descriptors": false
	      }`

	const getaddressinfoJSON = `{
		"address": "qUbxboqjBRp96j3La8D1RYkyqx5uQbJPoW",
		"scriptPubKey": "76a9147926223070547d2d15b2ef5e7383e541c338ffe988ac",
		"ismine": false,
		"solvable": false,
		"iswatchonly": false,
		"isscript": false,
		"iswitness": false,
		"ischange": false,
		"labels": [
		]
	}`

	const createwalletJSON = `{
		"name": "example",
		"warning": ""
	      }`

	var responses = map[string]string{
		"getwalletinfo":  getwalletinfoJSON,
		"getaddressinfo": getaddressinfoJSON,
		"createwallet":   createwalletJSON,
		"importaddress":  "null",
	}

	// create mock qtumd server
	mockQtumd := mocks.NewMockQtumd(responses)
	defer mockQtumd.Close()

	// create qtum client
	qcli, err := NewQtumClient(mockQtumd.URL, "qtum", "qtumpass", cfg.Net.String())
	utils.HandleFatalError(t, err)

	// create tests
	var tests = []struct {
		name    string
		address string
		wantErr bool
	}{
		{
			name:    "valid address",
			address: "qUbxboqjBRp96j3La8D1RYkyqx5uQbJPoW",
			wantErr: false,
		},
		// TODO: add test cases with wallet not found and address not found
	}

	// run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			err = qcli.VerifyAddress(tt.address)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}

}
