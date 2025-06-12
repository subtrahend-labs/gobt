//go:build integration
// +build integration

package extrinsics

import (
	"math/big"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/stretchr/testify/require"
	"github.com/subtrahend-labs/gobt/runtime"
	"github.com/subtrahend-labs/gobt/testutils"
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
	})

	t.Run("SudoToggleEvmPrecompile", func(t *testing.T) {
		t.Parallel()
		env := setup(t)

		call, err := SudoToggleEvmPrecompileCall(env.Client, types.U8(1), types.Bool(true))
		require.NoError(t, err, "Failed to create sudo_toggle_evm_precompile call")

		sudoExt, err := NewSudoExt(env.Client, &call)
		require.NoError(t, err, "Failed to create sudo ext")

		testutils.SignAndSubmit(t, env.Client, sudoExt, env.Alice.Coldkey.Keypair, uint32(env.Alice.Coldkey.AccInfo.Nonce))
	})

	t.Run("SudoSetSubnetMovingAlpha", func(t *testing.T) {
		t.Parallel()
		env := setup(t)

		alpha := runtime.I96F32{
			Bits: types.NewU128(*big.NewInt(1000000000)),
		}

		ext, err := SudoSetSubnetMovingAlphaExt(env.Client, alpha)
		require.NoError(t, err)

		sudoExt, err := NewSudoExt(env.Client, &ext.Method)
		require.NoError(t, err)

		blockHash := testutils.SignAndSubmit(t, env.Client, sudoExt, env.Alice.Coldkey.Keypair, uint32(env.Alice.Coldkey.AccInfo.Nonce))
		require.NotNil(t, blockHash)

		meta, _ := env.Client.Api.RPC.State.GetMetadataLatest()
		storageKey, _ := types.CreateStorageKey(meta, "SubtensorModule", "SubnetMovingAlpha")
		var result runtime.I96F32
		env.Client.Api.RPC.State.GetStorageLatest(storageKey, &result)

		require.Equal(t, alpha.Bits, result.Bits)
		t.Logf("SubnetMovingAlpha updated successfully: %v", result.Bits)
	})

	t.Run("SudoSetSubnetOwnerHotkey", func(t *testing.T) {
		t.Parallel()
		env := setup(t)
		setupSubnet(t, env)

		ext, err := SudoSetSubnetOwnerHotkeyExt(env.Client, types.U16(1), *env.Bob.Hotkey.AccID)
		require.NoError(t, err)

		sudoExt, err := NewSudoExt(env.Client, &ext.Method)
		require.NoError(t, err)

		blockHash := testutils.SignAndSubmit(t, env.Client, sudoExt, env.Alice.Coldkey.Keypair, uint32(env.Alice.Coldkey.AccInfo.Nonce))
		require.NotNil(t, blockHash, "Transaction should be included in block")

		t.Logf("SudoSetSubnetOwnerHotkey transaction included in block: %x", blockHash)
	})

	t.Run("SudoSetEmaPriceHalvingPeriod", func(t *testing.T) {
		t.Parallel()
		env := setup(t)
		setupSubnet(t, env)

		ext, err := SudoSetEmaPriceHalvingPeriodExt(env.Client, types.U16(1), types.U64(7200))
		require.NoError(t, err, "Failed to create sudo_set_ema_price_halving_period ext")

		sudoExt, err := NewSudoExt(env.Client, &ext.Method)
		require.NoError(t, err, "Failed to create sudo ext")

		blockHash := testutils.SignAndSubmit(t, env.Client, sudoExt, env.Alice.Coldkey.Keypair, uint32(env.Alice.Coldkey.AccInfo.Nonce))
		require.NotNil(t, blockHash, "Transaction should be included in block")

		t.Logf("SudoSetEmaPriceHalvingPeriod transaction included in block: %x", blockHash)
	})

}
