package wallet

import (
	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"github.com/ethereum/go-ethereum/common"
)

type IQtumWallet interface {
	GetPrivateKey(address string) (*secp256k1.PrivateKey, error)
	SetEthereumAddress(address *common.Address)
	GetEthereumAddress() *common.Address
	GetQtumAddress() (string, error)
}
