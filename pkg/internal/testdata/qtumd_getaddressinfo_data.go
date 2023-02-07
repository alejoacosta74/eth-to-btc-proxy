package testdata

import (
	"github.com/qtumproject/btcd/btcjson"
)

var AddressResult = btcjson.GetAddressInfoResult{
	IsMine:      true,
	IsWatchOnly: false,
}
