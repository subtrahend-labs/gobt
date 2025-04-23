//go:build integration
// +build integration

package extrinsics

import (
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/stretchr/testify/require"
	"github.com/subtrahend-labs/gobt/testutils"
)

func TestAdminUtilsModuleExtrinsics(t *testing.T) {
	t.Run("SudoSetNetworkRateLimit", func(t *testing.T) {
		setup(t)
		defer teardown(t)

		rateLimit := types.NewU64(1000)
		sudoCall, err := SudoSetNetworkRateLimitCall(env.Client, rateLimit)
		require.NoError(t, err, "Failed to create sudo_set_network_rate_limit ext")
		ext, err := NewSudoExt(env.Client, &sudoCall)
		testutils.SignAndSubmit(t, env.Client, ext, alice.coldkey.Keypair, uint32(alice.coldkey.AccInfo.Nonce))
		updateUserInfo(t, &alice)
	})
}
