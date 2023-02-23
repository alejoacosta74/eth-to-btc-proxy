package rpc

import (
	"bytes"
	"encoding/hex"
	"fmt"

	"github.com/alejoacosta74/rpc-proxy/pkg/log"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/pkg/errors"
	q "github.com/qtumproject/qtool/lib/common"
)

// VerifyTxSignature() verifies the signature of a signed ethereum transaction
// against a public key.
func VerifyTxSignature(tx *types.Transaction, pubKeyHex string) bool {

	//1. Get the signature and signer from the signed transaction
	signature, signer := getSignatureAndSigner(tx)

	//2. Recreate the raw transaction from the signed transaction
	var rawTx = types.NewTransaction(tx.Nonce(), *tx.To(), tx.Value(), tx.Gas(), tx.GasPrice(), tx.Data())

	//3. RLP encode the raw transaction
	var buf bytes.Buffer
	if err := rlp.Encode(&buf, rawTx); err != nil {
		fmt.Printf("error encoding raw tx: %v\n", err)
		return false
	}
	//4. Hash the raw transaction
	rawTxHashed := signer.Hash(rawTx)
	digest := rawTxHashed.Bytes()

	//6. Recover the public key from the signature and message digest
	recoveredPubKey, err := crypto.SigToPub(digest, signature)
	if err != nil {
		log.With("method", "_eth_sendRawtx").Debugf("error recovering public key: %v", err)
		return false
	}

	log.With("method", "_eth_sendRawtx").Debugf("recovered public key: %x", crypto.FromECDSAPub(recoveredPubKey))

	recoveredAddress := crypto.PubkeyToAddress(*recoveredPubKey)
	log.With("method", "_eth_sendRawtx").Debugf("recovered address: %s", recoveredAddress.String())

	//7. Convert the recovered public key to a hex string
	recoveredPubKeyHex := hex.EncodeToString(crypto.FromECDSAPub(recoveredPubKey))

	//8. Compare the recovered public key hex string with the expected public key hex string
	return recoveredPubKeyHex == pubKeyHex
}

// decodeRawTX decodes a raw ethereum RLP encoded transaction
// and returns a go-ethereum types.Transaction
func decodeRawTx(rawtx string) (*types.Transaction, error) {
	rawtx = q.RemoveHexPrefix(rawtx)

	rawtxBytes, err := hex.DecodeString(rawtx)
	if err != nil {
		return nil, err
	}

	var tx = &types.Transaction{}

	err = rlp.DecodeBytes(rawtxBytes, tx)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

// getFromAddress returns the signer address from the public key
// of a signed ethereum transaction
func getFromAddress(tx *types.Transaction) (*common.Address, error) {

	signature, signer := getSignatureAndSigner(tx)

	recoveredPubKey, err := crypto.SigToPub(signer.Hash(tx).Bytes(), signature)
	if err != nil {
		return nil, errors.Wrap(err, "error recovering pubkey")
	}

	recoveredAddress := crypto.PubkeyToAddress(*recoveredPubKey)

	return &recoveredAddress, nil

}

// getSignatureAndSigner returns the signature and signer from
// a signed ethereum transaction
func getSignatureAndSigner(tx *types.Transaction) (signature []byte, signer types.Signer) {
	v, r, s := tx.RawSignatureValues()

	signature = make([]byte, 65)
	copy(signature[32-len(r.Bytes()):32], r.Bytes())
	copy(signature[64-len(s.Bytes()):64], s.Bytes())

	if tx.Protected() {
		signer = types.NewEIP155Signer(tx.ChainId())
		signature[64] = byte(v.Uint64() - 35 - 2*tx.ChainId().Uint64())
	} else {
		signer = types.HomesteadSigner{}
		signature[64] = byte(v.Uint64() - 27)
	}

	return signature, signer
}
