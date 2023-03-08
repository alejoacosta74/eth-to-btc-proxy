package qtum

import (
	"fmt"

	"github.com/alejoacosta74/qproxy/pkg/log"
	"github.com/pkg/errors"
	"github.com/qtumproject/btcd/btcjson"
)

var walletExists = false

var errorWalletNotFound = btcjson.NewRPCError(btcjson.ErrRPCWalletNotFound, "No wallet is loaded. Load a wallet using loadwallet or create a new one with createwallet. (Note: A default wallet is no longer automatically created)")

// VerifyAddress checks if the address is known to the node's wallet.
// If not, it will import the address and rescan the blockchain seeking transactions
// related to this address.
func (q *QtumClient) VerifyAddress(address string) error {
	// check the node's wallet exists
	if !walletExists {
		log.With("module", "qtum").Debugf("Verifying node wallet...")
		err := q.verifyNodeWallet()
		if err != nil {
			return errors.Wrap(err, "Error verifying node wallet")
		}
		walletExists = true
	}
	result, err := q.GetAddressInfo(address)
	if err != nil {
		return errors.Wrap(err, "Error getting info for address: "+address)
	}
	log.With("module", "qtum").Tracef("Address info result: %+v", result)
	if !result.IsWatchOnly && !result.IsMine {
		log.With("module", "qtum").Debugf("Address %s not found in wallet. Importing it...", address)
		err := q.ImportAddressRescan(address, "", true)
		if err != nil {
			return errors.Wrap(err, "Error importing address: "+address)
		}
		log.With("module", "qtum").Debugf("Address imported: %+v", address)
	}
	return nil
}

// VerifyNodeWallet checks that the node's wallet exists and if not, it will create it.
func (q *QtumClient) verifyNodeWallet() error {
	walletInfo, err := q.GetWalletInfo()
	if err != nil {
		// check if the error is because the node's wallet was not found
		if fmt.Sprintf("%v", err) == fmt.Sprintf("%v", errorWalletNotFound) {
			log.With("module", "qtum").Debugf("Wallet not found. Creating it...")
			result, err := q.CreateWallet("wallet")
			if err != nil {
				return errors.Wrap(err, "Error creating wallet")
			}
			log.With("module", "qtum").Debugf("Wallet created: %+v", result)
			return nil
		}
	}
	log.With("module", "qtum").Tracef("Wallet info: %+v", walletInfo)
	return nil
}
