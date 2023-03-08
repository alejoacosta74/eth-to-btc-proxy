package rpc

import (
	"github.com/alejoacosta74/qproxy/pkg/log"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	qcommon "github.com/qtumproject/qtool/lib/common"
	qtool "github.com/qtumproject/qtool/lib/tools"
	"github.com/shopspring/decimal"
)

// GetBalance implements the eth_getBalance JSON-RPC call.
//
// Retrieves the balance of an ethereum address.
// If the block number is not specified, the "*latest*" block is used.
func (api *EthAPI) GetBalance(address string, blockNumber interface{}) (string, error) {

	// 1. Convert the address from hex to base58
	address = qcommon.RemoveHexPrefix(address)
	addrBase58, err := qtool.AddressHexToBase58(address, api.cfg)
	if err != nil {
		return "", errors.Wrapf(err, "Error converting address from hex to base58: %s", address)
	}
	log.With("method", "getbalance").Debugf("GetBalance called with address: %s (%s), blockNumber: %v", address, addrBase58, blockNumber)

	// 2. Verify the address is known to the node
	err = api.qcli.VerifyAddress(addrBase58)
	if err != nil {
		return "", errors.Wrapf(err, "Error verifying address: %s (hex %s)", addrBase58, address)
	}

	// 3. get a list of unspent outputs for the address
	unspent, err := api.qcli.FindSpendableUTXO(addrBase58)
	if err != nil {
		return "", errors.Wrapf(err, "Error getting unspent outputs for address: %s", addrBase58)
	}

	// 4. sum the amount of each output
	var balance float64
	for _, utxo := range unspent {
		balance += utxo.Amount
	}

	// 5. convert the amount from qtum to wei
	b := decimal.NewFromFloat(balance)
	balanceSat := qcommon.ConvertFromQtumToSatoshis(b)
	balanceWei := qcommon.ConvertFromSatoshiToWei(balanceSat.BigInt())

	// 6. Convert the amount from wei to hex
	hexBalance := hexutil.EncodeBig(balanceWei)

	log.With("method", "getbalance").Debugf("Address: %s (%s) - found %d unspent utxos with balance (hex): %s", address, addrBase58, len(unspent), hexBalance)

	return hexBalance, nil
}
