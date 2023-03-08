package rpc

import (
	"github.com/alejoacosta74/qproxy/pkg/log"
	"github.com/alejoacosta74/qproxy/pkg/wallet"
	"github.com/pkg/errors"
)

// RPC Method: personal_importRawKey
// Imports the given unencrypted private key (hex string) into
// the wallet
//
// Returns the ethereum address of the new account.
func (api *PersonalAPI) ImportRawKey(keydata string, passphrase string) (string, error) {
	log.With("method", "importrawkey").Debugf("ImportRawKey called with req.KeyData: %s, req.Passphrase: %s", keydata, passphrase)

	// TODO: implement btcd wallets for persisten storage

	ws := wallet.GetWallets()

	w, err := ws.NewWallet(keydata, api.cfg)
	if err != nil {
		return "", errors.Wrapf(err, "Error importing raw key: %s", keydata)
	}

	return w.GetEthereumAddress().String(), nil

}
