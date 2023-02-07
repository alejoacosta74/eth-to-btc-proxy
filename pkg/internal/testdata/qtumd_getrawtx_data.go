package testdata

type QtumGetRawTxMap map[string][]byte

var QtumGetRawTxMapSample = QtumGetRawTxMap{
	"ebcb8609e1c7cbdd69471447d771c96e866ec8c3a54afd587c0d0cf1e66ffadc": []byte(`{
		"txid": "ebcb8609e1c7cbdd69471447d771c96e866ec8c3a54afd587c0d0cf1e66ffadc",
		"hash": "ebcb8609e1c7cbdd69471447d771c96e866ec8c3a54afd587c0d0cf1e66ffadc",
		"version": 2,
		"size": 210,
		"vsize": 210,
		"weight": 840,
		"locktime": 0,
		"vin": [
		  {
		    "txid": "2485b395f94cbbddee483ba6468dd4e938b68d56019d532a09791d5c1fcadfa1",
		    "vout": 2,
		    "scriptSig": {
		      "asm": "3044022000b1ce0031d545cbeec68fe0bd8f332230e25e3479cd9e2cf7a18853a61dea6a022002e286710275756066e01a07068727bf334792bbfe588e9b820c204a76b47fc5[ALL]",
		      "hex": "473044022000b1ce0031d545cbeec68fe0bd8f332230e25e3479cd9e2cf7a18853a61dea6a022002e286710275756066e01a07068727bf334792bbfe588e9b820c204a76b47fc501"
		    },
		    "value": 100.42384175,
		    "valueSat": 10042384175,
		    "address": "Qc6iYCZWn4BauKXGYirRG8pMtgdHMk2dzn",
		    "sequence": 4294967295
		  }
		],
		"vout": [
		  {
		    "value": 0.00000000,
		    "valueSat": 0,
		    "n": 0,
		    "scriptPubKey": {
		      "asm": "",
		      "hex": "",
		      "type": "nonstandard"
		    }
		  },
		  {
		    "value": 104.43827170,
		    "valueSat": 10443827170,
		    "n": 1,
		    "scriptPubKey": {
		      "asm": "02e61c13b9bd4f5ca51d79576f97d98f3700faa43a743cb2bea64a97a200d8cba8 OP_CHECKSIG",
		      "hex": "2102e61c13b9bd4f5ca51d79576f97d98f3700faa43a743cb2bea64a97a200d8cba8ac",
		      "type": "pubkey"
		    },
		    "spentTxId": "8309121300d2da4455df574a81810a8cc168b9d77da60b47048e47edf965876b",
		    "spentIndex": 430,
		    "spentHeight": 832994
		  },
		  {
		    "value": 0.08801920,
		    "valueSat": 8801920,
		    "n": 2,
		    "scriptPubKey": {
		      "asm": "OP_DUP OP_HASH160 93594441cb5de8b497ad8467d55412c2a0ef3659 OP_EQUALVERIFY OP_CHECKSIG",
		      "hex": "76a91493594441cb5de8b497ad8467d55412c2a0ef365988ac",
		      "address": "Qa36NrNdFgr4XeMxKdZeSZ1FGCdSNLmqXh",
		      "type": "pubkeyhash",
		      "addresses": [
			"Qa36NrNdFgr4XeMxKdZeSZ1FGCdSNLmqXh"
		      ],
		      "reqSigs": 1
		    },
		    "spentTxId": "65618168f5a5cef6db33869388e356df91d215187358df2272fab59ad5f1dc70",
		    "spentIndex": 1,
		    "spentHeight": 792374
		  }
		],
		"hex": "0200000001a1dfca1f5c1d79092a539d01568db638e9d48d46a63b48eeddbb4cf995b385240200000048473044022000b1ce0031d545cbeec68fe0bd8f332230e25e3479cd9e2cf7a18853a61dea6a022002e286710275756066e01a07068727bf334792bbfe588e9b820c204a76b47fc501ffffffff03000000000000000000e227806e02000000232102e61c13b9bd4f5ca51d79576f97d98f3700faa43a743cb2bea64a97a200d8cba8ac804e8600000000001976a91493594441cb5de8b497ad8467d55412c2a0ef365988ac00000000",
		"blockhash": "d0d972f12350c6d0242c3ea8079634f79251322b0d2d4589031feb4f4eb4e51c",
		"confirmations": 1520487,
		"time": 1605390336,
		"blocktime": 1605390336,
		"height": 732721
	      }`),

	"0425fa39feed4cd6c93998159901095c147f8b0043823067dc1d25dabf950ac9": []byte(`{
		"txid": "0425fa39feed4cd6c93998159901095c147f8b0043823067dc1d25dabf950ac9",
		"hash": "0425fa39feed4cd6c93998159901095c147f8b0043823067dc1d25dabf950ac9",
		"version": 2,
		"size": 368,
		"vsize": 368,
		"weight": 1472,
		"locktime": 732720,
		"vin": [
		  {
		    "txid": "d10175aedbdd0cf7550b53d3bbff80ebd389120aae4b7aae93f1c31e96eb8e21",
		    "vout": 0,
		    "scriptSig": {
		      "asm": "304402205072650c6aafdd27c6c54732dda88481530c82eba4d5b45b8b6a3d5226dfdacb02205038e5d509a5c0442092845a855f8c3ab409846af306afbfdaf9e3a7ac3a95fe[ALL] 03379c39b6fb2c705db608f98a8fc064f94c66faf894996ca88595487f9ef04a6e",
		      "hex": "47304402205072650c6aafdd27c6c54732dda88481530c82eba4d5b45b8b6a3d5226dfdacb02205038e5d509a5c0442092845a855f8c3ab409846af306afbfdaf9e3a7ac3a95fe012103379c39b6fb2c705db608f98a8fc064f94c66faf894996ca88595487f9ef04a6e"
		    },
		    "value": 2.00000000,
		    "valueSat": 200000000,
		    "address": "Qa36NrNdFgr4XeMxKdZeSZ1FGCdSNLmqXh",
		    "sequence": 4294967294
		  }
		],
		"vout": [
		  {
		    "value": 1.89868837,
		    "valueSat": 189868837,
		    "n": 0,
		    "scriptPubKey": {
		      "asm": "OP_DUP OP_HASH160 93594441cb5de8b497ad8467d55412c2a0ef3659 OP_EQUALVERIFY OP_CHECKSIG",
		      "hex": "76a91493594441cb5de8b497ad8467d55412c2a0ef365988ac",
		      "address": "Qa36NrNdFgr4XeMxKdZeSZ1FGCdSNLmqXh",
		      "type": "pubkeyhash",
		      "addresses": [
			"Qa36NrNdFgr4XeMxKdZeSZ1FGCdSNLmqXh"
		      ],
		      "reqSigs": 1
		    },
		    "spentTxId": "e588141d2646fa6f1fd865ae141df476f6687e36d2f90e2e38caeb483fe5dbfb",
		    "spentIndex": 4,
		    "spentHeight": 792374
		  },
		  {
		    "value": 0.00000000,
		    "valueSat": 0,
		    "n": 1,
		    "scriptPubKey": {
		      "asm": "1 93594441cb5de8b497ad8467d55412c2a0ef3659 6a4730440220396b30b7a2f2af482e585473b7575dd2f989f3f3d7cdee55fa34e93f23d5254d022055326cdcab38c58dc3e65c458bfb656cca8340f59534c00ad98b4d4d3303f459012103379c39b6fb2c705db608f98a8fc064f94c66faf894996ca88595487f9ef04a6e OP_SENDER 4 250000 40 -191784509 0000000000000000000000000000000000000086 OP_CALL",
		      "hex": "01011493594441cb5de8b497ad8467d55412c2a0ef36594c6b6a4730440220396b30b7a2f2af482e585473b7575dd2f989f3f3d7cdee55fa34e93f23d5254d022055326cdcab38c58dc3e65c458bfb656cca8340f59534c00ad98b4d4d3303f459012103379c39b6fb2c705db608f98a8fc064f94c66faf894996ca88595487f9ef04a6ec401040390d0030128043d666e8b140000000000000000000000000000000000000086c2",
		      "type": "call_sender"
		    }
		  }
		],
		"hex": "0200000001218eeb961ec3f193ae7a4bae0a1289d3eb80ffbbd3530b55f70cdddbae7501d1000000006a47304402205072650c6aafdd27c6c54732dda88481530c82eba4d5b45b8b6a3d5226dfdacb02205038e5d509a5c0442092845a855f8c3ab409846af306afbfdaf9e3a7ac3a95fe012103379c39b6fb2c705db608f98a8fc064f94c66faf894996ca88595487f9ef04a6efeffffff02252b510b000000001976a91493594441cb5de8b497ad8467d55412c2a0ef365988ac0000000000000000a801011493594441cb5de8b497ad8467d55412c2a0ef36594c6b6a4730440220396b30b7a2f2af482e585473b7575dd2f989f3f3d7cdee55fa34e93f23d5254d022055326cdcab38c58dc3e65c458bfb656cca8340f59534c00ad98b4d4d3303f459012103379c39b6fb2c705db608f98a8fc064f94c66faf894996ca88595487f9ef04a6ec401040390d0030128043d666e8b140000000000000000000000000000000000000086c2302e0b00",
		"blockhash": "d0d972f12350c6d0242c3ea8079634f79251322b0d2d4589031feb4f4eb4e51c",
		"confirmations": 1515793,
		"time": 1605390336,
		"blocktime": 1605390336,
		"height": 732721
	      }`),
	"06be5aa06f2373b0f4a5497e36de0398412723d19bfe11b3e1083933f1435623": []byte(`{
		"txid": "06be5aa06f2373b0f4a5497e36de0398412723d19bfe11b3e1083933f1435623",
		"hash": "180ef3478f6cad5027349d2f45e35eacb5238ab54211ec5c087873ebf68a9e5f",
		"version": 2,
		"size": 148,
		"vsize": 121,
		"weight": 484,
		"locktime": 0,
		"vin": [
		  {
		    "coinbase": "03312e0b00",
		    "sequence": 4294967295
		  }
		],
		"vout": [
		  {
		    "value": 0.00000000,
		    "valueSat": 0,
		    "n": 0,
		    "scriptPubKey": {
		      "asm": "",
		      "hex": "",
		      "type": "nonstandard"
		    }
		  },
		  {
		    "value": 0.00000000,
		    "valueSat": 0,
		    "n": 1,
		    "scriptPubKey": {
		      "asm": "OP_RETURN aa21a9ed1841d96476b566dfd40b40c8981d8b0e031cfcb7a1a45bf312e67fa824d20849",
		      "hex": "6a24aa21a9ed1841d96476b566dfd40b40c8981d8b0e031cfcb7a1a45bf312e67fa824d20849",
		      "type": "nulldata"
		    }
		  }
		],
		"hex": "020000000001010000000000000000000000000000000000000000000000000000000000000000ffffffff0503312e0b00ffffffff020000000000000000000000000000000000266a24aa21a9ed1841d96476b566dfd40b40c8981d8b0e031cfcb7a1a45bf312e67fa824d208490120000000000000000000000000000000000000000000000000000000000000000000000000",
		"blockhash": "d0d972f12350c6d0242c3ea8079634f79251322b0d2d4589031feb4f4eb4e51c",
		"confirmations": 1522745,
		"time": 1605390336,
		"blocktime": 1605390336,
		"height": 732721
	      }`),
}
