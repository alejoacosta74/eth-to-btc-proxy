package mocks

import (
	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"github.com/ethereum/go-ethereum/common"
)

type MockQtumWallet struct {
}

func (w *MockQtumWallet) GetPrivateKey(address string) (*secp256k1.PrivateKey, error) {
	return nil, nil
}

func (w *MockQtumWallet) SetEthereumAddress(address *common.Address) {
}

func (w *MockQtumWallet) GetEthereumAddress() *common.Address {
	return nil
}

func (w *MockQtumWallet) GetQtumAddress() (string, error) {
	return "", nil
}
