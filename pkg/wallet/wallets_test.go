package wallet

import (
	"io"
	"os"
	"testing"

	utils "github.com/alejoacosta74/rpc-proxy/pkg/internal/testutils"
	"github.com/stretchr/testify/assert"
)

var cfg = utils.GetNetworkParams()

func TestNewWallet(t *testing.T) {
	assert := assert.New(t)

	// const privKeyStr = data.PrivKeyString
	const privKeyStr = "85cbc7b1adfe877051d746c3996a01c2bc3e7a6988490439b1f4b4c2b465322d"
	// const addressStr = data.EthereumAddress
	const addressStr = "0xA6d2799a4b465805421bd10247386a708F01DB03"
	const addressB58 = "qTQeBZsvBmmLevSu6cU3wGwyHeZdEp9Tkx"
	ws := GetWallets()

	t.Run("Test new wallet ", func(t *testing.T) {
		w, err := ws.NewWallet(privKeyStr, cfg)
		assert.Nil(err)
		assert.NotNil(w)
		ethAddress := w.GetEthereumAddress().String()
		assert.Equal(addressStr, ethAddress)
		qtumAddr, err := w.GetAddress()
		assert.Nil(err)
		assert.Equal(addressB58, qtumAddr)
	})
	t.Run("Test duplicated wallet", func(t *testing.T) {

		w, err := ws.NewWallet(privKeyStr, cfg)
		assert.NotNil(err)
		assert.Nil(w)
	})
	t.Run("Check existing wallets", func(t *testing.T) {
		assert.Equal(1, len(wallets))
	})
	t.Run("Delete wallet and keystore", func(t *testing.T) {
		err := ws.DeleteWallet(addressStr, "")
		assert.Nil(err)
		assert.Equal(0, len(wallets))
	})

}

func IsEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1) // Or f.Readdir(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err // Either not empty or error, suits both cases
}
