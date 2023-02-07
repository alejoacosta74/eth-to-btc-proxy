package testdata

import "github.com/qtumproject/btcd/btcjson"

var GetInfoResultMainnet *btcjson.InfoWalletResult = &btcjson.InfoWalletResult{
	ProtocolVersion: 70015,
	WalletVersion:   130000,
	Balance:         0,
	Blocks:          0,
	TimeOffset:      0,
	Connections:     0,
	Proxy:           "",
	Difficulty:      0,
	TestNet:         false,
}

var GetInfoResultTestnet *btcjson.InfoWalletResult = &btcjson.InfoWalletResult{
	ProtocolVersion: 70015,
	WalletVersion:   130000,
	Balance:         0,
	Blocks:          0,
	TimeOffset:      0,
	Connections:     0,
	Proxy:           "",
	Difficulty:      0,
	TestNet:         true,
}
