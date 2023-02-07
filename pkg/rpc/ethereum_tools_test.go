package rpc

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
)

func TestVerifyTxSignature(t *testing.T) {
	assert := assert.New(t)
	var privKeyHex string = "fad9c8855b740a0b7ed4c221dbad0f33a83a49cad6b3fe8d5817ac83d38b6a19"
	privKey, _ := crypto.HexToECDSA(privKeyHex)
	address := crypto.PubkeyToAddress(privKey.PublicKey)
	fmt.Printf("public key: %x\n", crypto.FromECDSAPub(&privKey.PublicKey))
	fmt.Printf("address: %s\n", address.Hex())
	tx := NewTx()

	signer := types.HomesteadSigner{}

	signedTx, err := SignTx(tx, signer, privKeyHex)
	if err != nil {
		panic("error signing transaction: %v" + err.Error())
	}
	PrintRawTx(signedTx, "signed transaction")
	pubKeyHex, err := PrivToPubKey(privKeyHex)
	if err != nil {
		panic("error converting private key to public key: %v" + err.Error())
	}
	verify := VerifyTxSignature(signedTx, pubKeyHex)
	assert.True(verify)

}

// PrivToPubKey converts a ECDSA private key hex string to a public key.
func PrivToPubKey(privKeyHex string) (string, error) {
	privKey, err := crypto.HexToECDSA(privKeyHex)
	if err != nil {
		return "", err
	}
	pubKey := privKey.Public()
	pubKeyECDSA, ok := pubKey.(*ecdsa.PublicKey)
	if !ok {
		return "", fmt.Errorf("error casting public key to ECDSA")
	}
	pubKeyBytes := crypto.FromECDSAPub(pubKeyECDSA)
	return hex.EncodeToString(pubKeyBytes), nil
}

// NewTx() creates a new unsigned ethereum transaction with default parameters.
func NewTx() *types.Transaction {
	nonce := uint64(0)
	gasPrice := big.NewInt(20000000000)
	gasLimit := uint64(21000)
	toAddress := common.HexToAddress("0x71517f86711b4bff4d789ad6fee9a58d8af1c6bb")
	amount := big.NewInt(1000000)
	tx := types.NewTransaction(nonce, toAddress, amount, gasLimit, gasPrice, nil)
	return tx
}

// SignTx() signs an ethereum transaction with a private key.
func SignTx(tx *types.Transaction, signer types.Signer, privKeyHex string) (*types.Transaction, error) {
	privKey, err := crypto.HexToECDSA(privKeyHex)
	if err != nil {
		return nil, err
	}
	signedTx, err := types.SignTx(tx, signer, privKey)
	if err != nil {
		return nil, err
	}
	return signedTx, nil
}

func PrintRawTx(tx *types.Transaction, msg string) {
	rawTx, _ := tx.MarshalJSON()
	logPretty(msg, rawTx)
}
