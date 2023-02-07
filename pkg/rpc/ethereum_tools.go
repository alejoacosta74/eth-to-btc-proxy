package rpc

import (
	"bytes"
	"encoding/hex"
	"fmt"

	rpctypes "github.com/alejoacosta74/rpc-proxy/pkg/types"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	q "github.com/qtumproject/qtool/lib/common"
)

// VerifyTxSignature() verifies the signature of a signed ethereum transaction
// against a public key.
func VerifyTxSignature(tx *types.Transaction, pubKeyHex string) bool {
	//1. Get the signature from the signed transaction
	v, r, s := tx.RawSignatureValues()
	signature := make([]byte, 65)
	copy(signature[32-len(r.Bytes()):32], r.Bytes())
	copy(signature[64-len(s.Bytes()):64], s.Bytes())
	signature[64] = byte(v.Uint64() - 27)

	fmt.Printf("signature: %x, length: %d\n", signature, len(signature))

	//2. Recreate the raw transaction from the signed transaction
	var rawTx = types.NewTransaction(tx.Nonce(), *tx.To(), tx.Value(), tx.Gas(), tx.GasPrice(), tx.Data())

	//3. RLP encode the raw transaction
	var buf bytes.Buffer
	if err := rlp.Encode(&buf, rawTx); err != nil {
		fmt.Printf("error encoding raw tx: %v\n", err)
		return false
	}
	//4. Kecccak256 hash the raw transaction
	rawTxHashed := crypto.Keccak256Hash(buf.Bytes())

	//5. Append ethereum prefix to the hashed raw transaction
	digest := accounts.TextHash(rawTxHashed.Bytes())

	//6. Recover the public key from the signature and message digest
	recoveredPubKey, err := crypto.SigToPub(digest, signature)
	if err != nil {
		fmt.Printf("error recovering pubkey: %v\n", err)
		return false
	}

	fmt.Printf("==>recovered public key: %x\n", crypto.FromECDSAPub(recoveredPubKey))

	recoveredAddress := crypto.PubkeyToAddress(*recoveredPubKey)
	fmt.Printf("==>recovered Address: %s\n", recoveredAddress.String())

	//7. Convert the recovered public key to a hex string
	recoveredHex := hex.EncodeToString(crypto.FromECDSAPub(recoveredPubKey))

	//8. Compare the recovered public key hex string with the expected public key hex string
	return recoveredHex == pubKeyHex
}

// decodeRawTX decodes a raw ethereum RLP encoded transaction
func decodeRawTx(rawtx string) (*rpctypes.RawTransaction, error) {
	rawtx = q.RemoveHexPrefix(rawtx)

	rawtxBytes, err := hex.DecodeString(rawtx)
	if err != nil {
		return nil, err
	}

	// TODO: replace with go-ethereum types
	var tx rpctypes.RlpDecodedRawTransaction
	err = rlp.DecodeBytes(rawtxBytes, &tx)
	if err != nil {
		return nil, err
	}

	decodedTx, err := tx.Unmarshal()
	if err != nil {
		return nil, err
	}

	return decodedTx, nil
}
