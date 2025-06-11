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

	t.Run("SudoSetMaxDifficulty", func(t *testing.T) {
		t.Parallel()
		env := setup(t)

		netuid := uint16(0)
		max_difficulty := types.NewU64(1000)
		default_difficulty := types.NewU64(0)

		sudoCall, err := SudoSetDifficultyCall(env.Client, netuid, default_difficulty)
		require.NoError(t, err, "Failed to create sudo_set_difficulty call")
		sudoCall, err = SudoSetMaxDifficultyCall(env.Client, netuid, max_difficulty)
		require.NoError(t, err, "Failed to create sudo_set_min_difficulty call")

		ext, err := NewSudoExt(env.Client, &sudoCall)
		require.NoError(t, err, "Failed to create sudo extrinsic")

		testutils.SignAndSubmit(t, env.Client, ext, env.Alice.Coldkey.Keypair, uint32(env.Alice.Coldkey.AccInfo.Nonce))
		updateUserInfo(t, &env.Alice, env, false)
	})

	t.Run("SudoSetDifficulty", func(t *testing.T) {
		t.Parallel()
		env := setup(t)

		netuid := uint16(0)
		default_difficulty := types.NewU64(0)

		sudoCall, err := SudoSetDifficultyCall(env.Client, netuid, default_difficulty)
		require.NoError(t, err, "Failed to create sudo_set_difficulty call")

		// Wrap in sudo extrinsic
		ext, err := NewSudoExt(env.Client, &sudoCall)
		require.NoError(t, err, "Failed to create sudo extrinsic")

		// Submit as root (usually Alice in testnets)
		testutils.SignAndSubmit(t, env.Client, ext, env.Alice.Coldkey.Keypair, uint32(env.Alice.Coldkey.AccInfo.Nonce))
		updateUserInfo(t, &env.Alice, env, false)
	})

	t.Run("SudoSetMinDifficulty", func(t *testing.T) {
		t.Parallel()
		env := setup(t)

		netuid := uint16(0)
		default_difficulty := types.NewU64(0)
		min_difficulty := types.NewU64(100)

		sudoCall, err := SudoSetDifficultyCall(env.Client, netuid, default_difficulty)
		require.NoError(t, err, "Failed to create sudo_set_difficulty call")
		sudoCall, err = SudoSetMinDifficultyCall(env.Client, netuid, min_difficulty)
		require.NoError(t, err, "Failed to create sudo_set_min_difficulty call")

		ext, err := NewSudoExt(env.Client, &sudoCall)
		require.NoError(t, err, "Failed to create sudo extrinsic")

		testutils.SignAndSubmit(t, env.Client, ext, env.Alice.Coldkey.Keypair, uint32(env.Alice.Coldkey.AccInfo.Nonce))
		updateUserInfo(t, &env.Alice, env, false)
	})

}
