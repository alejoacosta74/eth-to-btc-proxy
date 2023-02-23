package wallet

import (
	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/qtumproject/btcd/btcutil"
	"github.com/qtumproject/btcd/chaincfg"
	qtool "github.com/qtumproject/qtool/lib/tools"
)

// QtumWallet represents a wallet for the Qtum blockchain
//
// wif is the private key
// cfg is the chain configuration
// ethereumAddr is the ethereum address associated with the wallet
type QtumWallet struct {
	wif          *btcutil.WIF
	cfg          *chaincfg.Params
	ethereumAddr *common.Address
}

func NewQtumWallet(privKey string, cfg *chaincfg.Params) (*QtumWallet, error) {
	result, err := qtool.ConvertPrivateKeyToWIF(privKey, cfg.Net.String())
	if err != nil {
		return nil, errors.Wrap(err, "Failed to convert private key to WIF")
	}
	wif, err := btcutil.DecodeWIF(result.WIF)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to decode WIF")
	}
	return &QtumWallet{wif: wif, cfg: cfg}, nil
}

// Returns the private key for the given address provided it
// can be generated from the wallet's private key
func (w *QtumWallet) GetPrivateKey(address string) (*secp256k1.PrivateKey, error) {
	// Confirm the address can be generated from the private key
	generatedAddr, err := btcutil.NewAddressPubKey(w.wif.PrivKey.PubKey().SerializeCompressed(), w.cfg)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to generate address from private key")
	}
	if generatedAddr.EncodeAddress() != address {
		return nil, errors.Wrapf(err, "Address mismatch: %s != %s", generatedAddr.EncodeAddress(), address)
	}
	return w.wif.PrivKey, nil
}

func (w *QtumWallet) SetEthereumAddress(address *common.Address) {
	w.ethereumAddr = address
}

func (w *QtumWallet) GetEthereumAddress() *common.Address {
	return w.ethereumAddr
}

/*
func (w *QtumWallet) GetAddressPubKey() (btcutil.Address, error) {
	addrPubKey, err := btcutil.NewAddressPubKey(w.wif.PrivKey.PubKey().SerializeCompressed(), w.cfg)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to generate address pubkey from private key")
	}
	return addrPubKey, nil
}
*/

// GetQtumAddress returns the address associated with the wallet's private key in base58 format
func (w *QtumWallet) GetQtumAddress() (string, error) {
	addrPubKey, err := btcutil.NewAddressPubKey(w.wif.PrivKey.PubKey().SerializeCompressed(), w.cfg)
	if err != nil {
		return "", errors.Wrap(err, "Failed to generate address pubkey from private key")
	}
	return addrPubKey.EncodeAddress(), nil
}
