package testdata

import (
	"github.com/qtumproject/btcd/btcjson"
)

var UnspentUTXO = []btcjson.ListUnspentResult{
	{
		TxID:    "06d48766c315b30b08cdcc2c34dc062ee4b3885dfc2414d5940c25225e058c2b",
		Vout:    1,
		Address: "qeVQ5JF6idPcrg1u9M3pCryXeebpj3Tbpk",
		//  Account: ,
		ScriptPubKey: "76a914e599be870c63d68a00a5019906d258a4ba5d1bac88ac",
		//  "RedeemScript": ,
		Amount:        0.01,
		Confirmations: 9,
		Spendable:     true,
	},
	{
		TxID:    "31f7d9e8d9d81439afec283f82a5156968a210492d15a4a69d3fe233f79efbbf",
		Vout:    1,
		Address: "qeVQ5JF6idPcrg1u9M3pCryXeebpj3Tbpk",
		//  Account: ,
		ScriptPubKey: "76a914e599be870c63d68a00a5019906d258a4ba5d1bac88ac",
		// RedeemScript:
		Amount:        0.01,
		Confirmations: 10,
		Spendable:     true,
	},
	{
		TxID:    "31f7d9e8d9d81439afec283f82a5156968a210492d15a4a69d3fe233f79efbbf",
		Vout:    0,
		Address: "qeVQ5JF6idPcrg1u9M3pCryXeebpj3Tbpk",
		//  Account: ,
		ScriptPubKey: "76a914e599be870c63d68a00a5019906d258a4ba5d1bac88ac",
		// RedeemScript:
		Amount:        1,
		Confirmations: 10,
		Spendable:     true,
	},
}

var GetAddressInfoResult = &btcjson.GetAddressInfoResult{
	IsMine:      true,
	IsWatchOnly: false,
}
