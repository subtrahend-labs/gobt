//go:build integration
// +build integration

package extrinsics

import (
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/stretchr/testify/require"
	"github.com/subtrahend-labs/gobt/testutils"
	"github.com/subtrahend-labs/gobt/typetools"
)

func TestAdminUtilsModuleExtrinsics(t *testing.T) {
	t.Parallel()
	t.Run("SudoSetNetworkRateLimit", func(t *testing.T) {
		t.Parallel()
		env := setup(t)

		rateLimit := types.NewU64(1000)
		sudoCall, err := SudoSetNetworkRateLimitCall(env.Client, rateLimit)
		require.NoError(t, err, "Failed to create sudo_set_network_rate_limit ext")
		ext, _ := NewSudoExt(env.Client, &sudoCall)
		testutils.SignAndSubmit(t, env.Client, ext, env.Alice.Coldkey.Keypair, uint32(env.Alice.Coldkey.AccInfo.Nonce))
		updateUserInfo(t, &env.Alice, env, false)
	})

	t.Run("SudoSetDefaultTake", func(t *testing.T) {
		t.Parallel()
		env := setup(t)

		defaultTake := types.NewU16(1000)
		sudoCall, _ := SudoSetDefaultTakeCall(env.Client, defaultTake)
		ext, _ := NewSudoExt(env.Client, &sudoCall)
		testutils.SignAndSubmit(t, env.Client, ext, env.Alice.Coldkey.Keypair, uint32(env.Alice.Coldkey.AccInfo.Nonce))
		updateUserInfo(t, &env.Alice, env, false)

		storageKey, _ := types.CreateStorageKey(env.Client.Meta, "SubtensorModule", "MaxDelegateTake")

		var newDefaultTake types.U16
		env.Client.Api.RPC.State.GetStorageLatest(storageKey, &newDefaultTake)

		require.Equal(t, defaultTake, newDefaultTake, "Default take was not updated correctly")
	})

	t.Run("SudoSetTxRateLimit", func(t *testing.T) {
		t.Parallel()
		env := setup(t)

		txRateLimit := types.NewU64(1000)
		sudoCall, _ := SudoSetTxRateLimitCall(env.Client, txRateLimit)
		ext, _ := NewSudoExt(env.Client, &sudoCall)
		testutils.SignAndSubmit(t, env.Client, ext, env.Alice.Coldkey.Keypair, uint32(env.Alice.Coldkey.AccInfo.Nonce))
		updateUserInfo(t, &env.Alice, env, false)

		storageKey, _ := types.CreateStorageKey(env.Client.Meta, "SubtensorModule", "TxRateLimit")
		var newTxRateLimit types.U64
		env.Client.Api.RPC.State.GetStorageLatest(storageKey, &newTxRateLimit)

		require.Equal(t, txRateLimit, newTxRateLimit, "Tx rate limit was not updated correctly")
	})

	t.Run("SudoSetServingRateLimit", func(t *testing.T) {
		t.Parallel()
		env := setup(t)

		servingRateLimit := types.NewU64(1000)
		sudoCall, _ := SudoSetServingRateLimitCall(env.Client, 4, servingRateLimit)
		ext, _ := NewSudoExt(env.Client, &sudoCall)
		testutils.SignAndSubmit(t, env.Client, ext, env.Alice.Coldkey.Keypair, uint32(env.Alice.Coldkey.AccInfo.Nonce))
		updateUserInfo(t, &env.Alice, env, false)

		storageKey, _ := types.CreateStorageKey(env.Client.Meta, "SubtensorModule", "ServingRateLimit", typetools.Uint16ToBytes(uint16(4)))
		var newServingRateLimit types.U64
		env.Client.Api.RPC.State.GetStorageLatest(storageKey, &newServingRateLimit)

		require.Equal(t, servingRateLimit, newServingRateLimit, "Serving rate limit was not updated correctly")
	})

}
