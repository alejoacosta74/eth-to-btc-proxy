package rpc

import (
	"bytes"
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/core/types"
)

var privKeyHex string = "fad9c8855b740a0b7ed4c221dbad0f33a83a49cad6b3fe8d5817ac83d38b6a19"

var tests = []struct {
	name   string
	signer types.Signer
	want   bool
}{
	{
		name:   "homestead",
		signer: types.HomesteadSigner{},
		want:   true,
	},
	{
		name:   "eip155 chain id 1",
		signer: types.NewEIP155Signer(big.NewInt(1)),
		want:   true,
	},
	{
		name:   "eip155 chain id 2",
		signer: types.NewEIP155Signer(big.NewInt(2)),
		want:   true,
	},
}

func TestVerifyTxSignature(t *testing.T) {
	// privKey, _ := crypto.HexToECDSA(privKeyHex)
	// address := crypto.PubkeyToAddress(privKey.PublicKey)
	// fmt.Printf("public key: %x\n", crypto.FromECDSAPub(&privKey.PublicKey))
	// fmt.Printf("address: %s\n", address.Hex())
	tx := newEthereumTx()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			signedTx, err := signEthereumTx(tx, tt.signer, privKeyHex)
			if err != nil {
				panic("error signing transaction: %v" + err.Error())
			}
			// PrintRawTx(signedTx, "signed transaction")
			pubKeyHex, err := privToPubKey(privKeyHex)
			if err != nil {
				panic("error converting private key to public key: %v" + err.Error())
			}
			verify := VerifyTxSignature(signedTx, pubKeyHex)
			if verify != tt.want {
				t.Fatalf("VerifyTxSignature() = %v, want %v", verify, tt.want)
			}
		})
	}

}

func TestDecodeRawTx(t *testing.T) {
	tx := newEthereumTx()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			signedTx, err := signEthereumTx(tx, tt.signer, privKeyHex)
			if err != nil {
				t.Fatalf("error signing transaction: %v", err)
			}

			buf := new(bytes.Buffer)
			if err := signedTx.EncodeRLP(buf); err != nil {
				t.Fatalf("error encoding transaction: %v", err)
			}
			rlpEncodedTx := buf.Bytes()
			rlpEncodedTxHex := hex.EncodeToString(rlpEncodedTx)

			decodedTx, err := decodeRawTx(rlpEncodedTxHex)
			if err != nil {
				t.Fatalf("error decoding transaction: %v", err)
			}

			want, err := signedTx.MarshalJSON()
			if err != nil {
				t.Fatalf("error marshalling transaction: %v", err)
			}
			got, err := decodedTx.MarshalJSON()
			if err != nil {
				t.Fatalf("error marshalling transaction: %v", err)
			}
			if !bytes.Equal(want, got) {
				t.Fatalf("decoded transaction does not match original transaction: want %s, got %s", want, got)
			}
		})
	}

}
