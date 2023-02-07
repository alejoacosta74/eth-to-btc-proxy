package wallet

import (
	"github.com/alejoacosta74/rpc-proxy/pkg/log"
	"github.com/pkg/errors"
	"github.com/qtumproject/btcd/chaincfg"
)

type Wallets map[string]*QtumWallet

var wallets = make(Wallets)

func GetWallets() Wallets {
	return wallets
}

func (ws Wallets) DeleteWallet(address string, passphrase string) error {
	if ws[address] == nil {
		return errors.New("Wallet not found for address: " + address)
	}

	log.With("module", "wallet").Debugf("Deleting wallet for address: %s", address)
	delete(ws, address)
	log.With("module", "wallet").Debugf("Succesfully deleted wallet for address %s", address)
	return nil
}

// NewWallet creates a new wallet for the given private key
func (ws Wallets) NewWallet(privKeyStr string, cfg *chaincfg.Params) (*QtumWallet, error) {
	// verify that the wallet does not exist for the eth address
	address, err := PrivKeyToEthAddress(privKeyStr)
	if err != nil {
		return nil, errors.Wrapf(err, "Error creating wallet for private: %s", privKeyStr)
	}
	if wallets[address.String()] != nil {
		return nil, errors.New("Wallet already exists for private key. Address: " + address.String())
	}

	// Create Qtum wallet
	w, err := NewQtumWallet(privKeyStr, cfg)
	if err != nil {
		return nil, errors.Wrapf(err, "Error creating wallet for private key: %s", privKeyStr)
	}
	w.SetEthereumAddress(address)

	qtumAddr, err := w.GetAddress()
	if err != nil {
		return nil, errors.Wrapf(err, "Error getting qtum address for private key: %s", privKeyStr)
	}
	log.With("module", "wallet").Debugf("Created wallet for eth addr: %s and qtum addr: %s", address, qtumAddr)

	ws[address.String()] = w
	return w, nil
}

// SeekWallet returns the Qtum wallet associated to the given ethereum  address. If the wallet
// is not found, an error is returned.
func (ws Wallets) SeekWallet(address string) (*QtumWallet, error) {
	if ws[address] == nil {
		return nil, errors.New("Wallet not found for address: " + address)
	}
	return wallets[address], nil
}
