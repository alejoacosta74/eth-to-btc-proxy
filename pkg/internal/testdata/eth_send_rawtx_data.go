package testdata

import (
	"github.com/alejoacosta74/rpc-proxy/pkg/types"
)

const RawTx = "f86d8202b28477359400825208944592d8f8d7b001e72cb26a73e4fa1806a51ac79d880de0b6b3a7640000802ca05924bde7ef10aa88db9c66dd4f5fb16b46dff2319b9968be983118b57bb50562a001b24b31010004f13d9a26b320845257a6cfc2bf819a3d55e3fc86263c5f0772"

var DecodedRawTx = types.RawTransaction{
	From:     "0x96216849c49358B10257cb55b28eA603c874b05E",
	Gas:      "0x5208",
	GasPrice: "0x77359400",
	Hash:     "0xc429e5f128387d224ba8bed6885e86525e14bfdc2eb24b5e9c3351a1176fd81f",
	Input:    "0x",
	Nonce:    "0x2b2",
	To:       "0x4592d8f8d7b001e72cb26a73e4fa1806a51ac79d",
	Value:    "0xde0b6b3a7640000",
	V:        "0x2c",
	R:        "0x5924bde7ef10aa88db9c66dd4f5fb16b46dff2319b9968be983118b57bb50562",
	S:        "0x1b24b31010004f13d9a26b320845257a6cfc2bf819a3d55e3fc86263c5f0772",
}
