package testdata

type QtumGetTxOutMap map[string]map[uint32][]byte

var QtumGetTxOutMapSample = QtumGetTxOutMap{
	"ebcb8609e1c7cbdd69471447d771c96e866ec8c3a54afd587c0d0cf1e66ffadc": {
		0: []byte(`{
				"bestblock": "eef79e3e1f8abceb219d399974e2b56b4af4428c6f4712d06c6088b7548453f3",
				"confirmations": 1520483,
				"value": 0.00000000,
				"scriptPubKey": {
				  "asm": "",
				  "hex": "",
				  "type": "nonstandard"
				},
				"coinbase": false,
				"coinstake": true
			      }`),
		1: []byte(`{}`),
		2: []byte(`{}`),
	},
	"06be5aa06f2373b0f4a5497e36de0398412723d19bfe11b3e1083933f1435623": {
		0: []byte(`{
			"bestblock": "0745ff853be05c13c4a9b659c26dfabf0fba5f10b6466c5c721796212b966878",
			"confirmations": 1522721,
			"value": 0.00000000,
			"scriptPubKey": {
			  "asm": "",
			  "hex": "",
			  "type": "nonstandard"
			},
			"coinbase": true,
			"coinstake": false
		      }`),
		1: []byte(`{}`),
		2: []byte(`{}`),
	},
	"d10175aedbdd0cf7550b53d3bbff80ebd389120aae4b7aae93f1c31e96eb8e21": {
		0: []byte(`{} `),
		1: []byte(`{}`),
		2: []byte(`{}`),
	},
	"0425fa39feed4cd6c93998159901095c147f8b0043823067dc1d25dabf950ac9": {
		0: []byte(`{} `),
		1: []byte(`{
			"bestblock": "1bb6fc364cef914c4dea5e056b2d2577a42b563a74589f52795d7c2b4f42ebb4",
			"confirmations": 1522738,
			"value": 0.00000000,
			"scriptPubKey": {
			  "asm": "1 93594441cb5de8b497ad8467d55412c2a0ef3659 6a4730440220396b30b7a2f2af482e585473b7575dd2f989f3f3d7cdee55fa34e93f23d5254d022055326cdcab38c58dc3e65c458bfb656cca8340f59534c00ad98b4d4d3303f459012103379c39b6fb2c705db608f98a8fc064f94c66faf894996ca88595487f9ef04a6e OP_SENDER 4 250000 40 -191784509 0000000000000000000000000000000000000086 OP_CALL",
			  "hex": "01011493594441cb5de8b497ad8467d55412c2a0ef36594c6b6a4730440220396b30b7a2f2af482e585473b7575dd2f989f3f3d7cdee55fa34e93f23d5254d022055326cdcab38c58dc3e65c458bfb656cca8340f59534c00ad98b4d4d3303f459012103379c39b6fb2c705db608f98a8fc064f94c66faf894996ca88595487f9ef04a6ec401040390d0030128043d666e8b140000000000000000000000000000000000000086c2",
			  "type": "call_sender"
			},
			"coinbase": false,
			"coinstake": false
		      }`),
		2: []byte(`{}`),
	},
}
