package testdata

type QtumGetBlockMap map[string][]byte

var QtumGetBlockMapSample = QtumGetBlockMap{
	"d0d972f12350c6d0242c3ea8079634f79251322b0d2d4589031feb4f4eb4e51c": []byte(`{
		"hash": "d0d972f12350c6d0242c3ea8079634f79251322b0d2d4589031feb4f4eb4e51c",
		"confirmations": 1515780,
		"height": 732721,
		"version": 536870912,
		"versionHex": "20000000",
		"merkleroot": "09e61ca9c94a02f01732980839eb92b25a7e793e01e187f696c7809e379659ba",
		"time": 1605390336,
		"mediantime": 1605389552,
		"nonce": 0,
		"bits": "1a055bff",
		"difficulty": 3130403.781442982,
		"chainwork": "0000000000000000000000000000000000000000000001f36ab2de690c1b38a0",
		"nTx": 4,
		"hashStateRoot": "cef6d1715ad99591e04f749c6420aa0cda5f31caa4f84fcf5bea38b06b658d0e",
		"hashUTXORoot": "ae921e57ea715053820bd30a8cb96eb0118e46ea7a091bdbe5f702614559bb0a",
		"prevoutStakeHash": "2485b395f94cbbddee483ba6468dd4e938b68d56019d532a09791d5c1fcadfa1",
		"prevoutStakeVoutN": 2,
		"previousblockhash": "4fc10c3cfa02102275fd6eac36a1a32ee20eb65e6c3e1b2be7227b3723a000ff",
		"nextblockhash": "3f1bb8785bf08bb9adaaded261c54873e880495f2b70688c9bdf38b9dce546e8",
		"flags": "proof-of-stake",
		"proofhash": "0000000000000000000000000000000000000000000000000000000000000000",
		"modifier": "09637b326d5243e6e085c282e5d391832387a29976ab920a6fafa1ee58e7f02e",
		"signature": "20d8b54d0c295bc5bf3c1e3763cb3a529aa0478cb28baf48e3b2f0c9c3ddfb869d1f07f90711d8d5fdd86a2b9071561799552ddc02564310d2d368bb7880e5210c",
		"strippedsize": 1162,
		"size": 1198,
		"weight": 4684,
		"tx": [
		  "06be5aa06f2373b0f4a5497e36de0398412723d19bfe11b3e1083933f1435623",
		  "ebcb8609e1c7cbdd69471447d771c96e866ec8c3a54afd587c0d0cf1e66ffadc",
		  "d10175aedbdd0cf7550b53d3bbff80ebd389120aae4b7aae93f1c31e96eb8e21",
		  "0425fa39feed4cd6c93998159901095c147f8b0043823067dc1d25dabf950ac9"
		]
	      }`),
}
