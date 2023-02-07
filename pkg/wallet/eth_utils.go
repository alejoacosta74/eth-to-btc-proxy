package wallet

import (
	"crypto/ecdsa"

	"github.com/pkg/errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// PrivKeyToAddress converts a private hex key string to an ethereum address
func PrivKeyToEthAddress(privKeyStr string) (*common.Address, error) {
	privateKey, err := crypto.HexToECDSA(privKeyStr)
	if err != nil {
		return nil, errors.New("Error converting private key to ECDSA:" + err.Error())
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("Type assertion failed for private key")

	}
	address := crypto.PubkeyToAddress(*publicKeyECDSA)
	return &address, nil
}
