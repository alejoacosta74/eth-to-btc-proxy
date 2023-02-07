package qtum

import (
	"testing"

	data "github.com/alejoacosta74/rpc-proxy/pkg/internal/testdata"
	utils "github.com/alejoacosta74/rpc-proxy/pkg/internal/testutils"
	"github.com/stretchr/testify/assert"
)

func TestVerifyAddress(t *testing.T) {
	var cfg = utils.GetNetworkParams()
	assert := assert.New(t)
	mockQtumd := utils.NewMockQtumRpcServer("")
	defer mockQtumd.Close()
	qcli, err := NewQtumClient(mockQtumd.URL, "user", "pass", cfg.Net.String())
	utils.HandleFatalError(t, err)

	tests := []struct {
		name    string
		address string
		want    bool
	}{
		{
			name:    "Verify address",
			address: data.AddressB58,
			want:    true,
		},
		{
			name:    "Verify address (bis)",
			address: data.AddressB58,
			want:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err = qcli.VerifyAddress(tt.address)
			assert.NoError(err)
		})
	}
}
