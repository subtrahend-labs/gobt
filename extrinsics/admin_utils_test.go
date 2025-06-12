//go:build integration
// +build integration

package extrinsics

import (
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
		updateUserInfo(t, &env.Alice, env, false)
	})

	t.Run("SudoSetMaxDifficulty", func(t *testing.T) {
		t.Parallel()
		env := setup(t)

		netuid := uint16(0)
		max_difficulty := types.NewU64(1000)

		metagraph, err := runtime.GetMetagraph(env.Client, netuid, nil)
		require.NoError(t, err, "Failed to get initial metagraph")
		initialMaxDifficulty := metagraph.MaxDifficulty

		var sudoCall types.Call
		sudoCall, err = SudoSetMaxDifficultyCall(env.Client, netuid, max_difficulty)
		require.NoError(t, err, "Failed to create sudo_set_min_difficulty call")

		ext, err := NewSudoExt(env.Client, &sudoCall)
		require.NoError(t, err, "Failed to create sudo extrinsic")

		testutils.SignAndSubmit(t, env.Client, ext, env.Alice.Coldkey.Keypair, uint32(env.Alice.Coldkey.AccInfo.Nonce))

		metagraph, err = runtime.GetMetagraph(env.Client, netuid, nil)
		require.NoError(t, err, "Failed to get updated metagraph")

		newMaxDifficulty := types.NewU64(uint64(metagraph.MaxDifficulty.Int64()))
		require.Equal(t, max_difficulty, newMaxDifficulty, "Max difficulty was not set correctly")

		require.NotEqual(t, initialMaxDifficulty, newMaxDifficulty, "Max difficulty did not change")
	})

	t.Run("SudoSetDifficulty", func(t *testing.T) {
		t.Parallel()
		env := setup(t)

		netuid := uint16(0)
		default_difficulty := types.NewU64(0)

		metagraph, err := runtime.GetMetagraph(env.Client, netuid, nil)
		require.NoError(t, err, "Failed to get initial metagraph")
		initialDifficulty := metagraph.Difficulty

		var sudoCall types.Call
		sudoCall, err = SudoSetDifficultyCall(env.Client, netuid, default_difficulty)
		require.NoError(t, err, "Failed to create sudo_set_difficulty call")

		ext, err := NewSudoExt(env.Client, &sudoCall)
		require.NoError(t, err, "Failed to create sudo extrinsic")

		testutils.SignAndSubmit(t, env.Client, ext, env.Alice.Coldkey.Keypair, uint32(env.Alice.Coldkey.AccInfo.Nonce))

		metagraph, err = runtime.GetMetagraph(env.Client, netuid, nil)
		require.NoError(t, err, "Failed to get updated metagraph")

		newDifficulty := types.NewU64(uint64(metagraph.Difficulty.Int64()))
		require.Equal(t, default_difficulty, newDifficulty, "Difficulty was not set correctly")

		require.NotEqual(t, initialDifficulty, metagraph.Difficulty, "Difficulty did not change")
	})

	t.Run("SudoSetMinDifficulty", func(t *testing.T) {
		t.Parallel()
		env := setup(t)

		netuid := uint16(0)
		min_difficulty := types.NewU64(100)

		metagraph, err := runtime.GetMetagraph(env.Client, netuid, nil)
		require.NoError(t, err, "Failed to get initial metagraph")
		initialMinDifficulty := metagraph.MinDifficulty

		var sudoCall types.Call
		sudoCall, err = SudoSetMinDifficultyCall(env.Client, netuid, min_difficulty)
		require.NoError(t, err, "Failed to create sudo_set_min_difficulty call")

		ext, err := NewSudoExt(env.Client, &sudoCall)
		require.NoError(t, err, "Failed to create sudo extrinsic")

		testutils.SignAndSubmit(t, env.Client, ext, env.Alice.Coldkey.Keypair, uint32(env.Alice.Coldkey.AccInfo.Nonce))

		metagraph, err = runtime.GetMetagraph(env.Client, netuid, nil)
		require.NoError(t, err, "Failed to get updated metagraph")

		newMinDifficulty := types.NewU64(uint64(metagraph.MinDifficulty.Int64()))
		require.Equal(t, min_difficulty, newMinDifficulty, "Min difficulty was not set correctly")

		require.NotEqual(t, initialMinDifficulty, newMinDifficulty, "Min difficulty did not change")
	})
	t.Run("SudoSetWeightsVersionKey", func(t *testing.T) {
		t.Parallel()
		env := setup(t)

		netuid := uint16(0)
		weights_version_key := types.NewU64(1)

		metagraph, err := runtime.GetMetagraph(env.Client, netuid, nil)
		require.NoError(t, err, "Failed to get initial metagraph")
		initialVersion := metagraph.WeightsVersion

		sudoCall, err := SudoSetWeightsVersionKeyCall(env.Client, netuid, weights_version_key)
		require.NoError(t, err, "Failed to create sudo_set_weights_version_key call")

		ext, err := NewSudoExt(env.Client, &sudoCall)
		require.NoError(t, err, "Failed to create sudo extrinsic")

		testutils.SignAndSubmit(t, env.Client, ext, env.Alice.Coldkey.Keypair, uint32(env.Alice.Coldkey.AccInfo.Nonce))

		metagraph, err = runtime.GetMetagraph(env.Client, netuid, nil)
		require.NoError(t, err, "Failed to get updated metagraph")

		newVersion := types.NewU64(uint64(metagraph.WeightsVersion.Int64()))
		require.Equal(t, weights_version_key, newVersion, "Weights version key was not set correctly")

		require.NotEqual(t, initialVersion, newVersion, "Weights version key did not change")
	})

	t.Run("SudoSetTempo", func(t *testing.T) {
		t.Parallel()
		env := setup(t)

		netuid := uint16(0)
		tempo := uint16(1)

		metagraph, err := runtime.GetMetagraph(env.Client, netuid, nil)
		require.NoError(t, err, "Failed to get initial metagraph")
		initialTempo := metagraph.Tempo

		sudoCall, err := SudoSetTempoCall(env.Client, netuid, tempo)
		require.NoError(t, err, "Failed to create sudo_set_tempo call")

		ext, err := NewSudoExt(env.Client, &sudoCall)
		require.NoError(t, err, "Failed to create sudo extrinsic")

		testutils.SignAndSubmit(t, env.Client, ext, env.Alice.Coldkey.Keypair, uint32(env.Alice.Coldkey.AccInfo.Nonce))

		metagraph, err = runtime.GetMetagraph(env.Client, netuid, nil)
		require.NoError(t, err, "Failed to get updated metagraph")

		newTempo := uint16(metagraph.Tempo.Int64())
		require.Equal(t, tempo, newTempo, "Tempo was not set correctly")
		require.NotEqual(t, initialTempo, metagraph.Tempo, "Tempo did not change")
	})

	t.Run("SudoSetTotalIssuance", func(t *testing.T) {
		t.Parallel()
		env := setup(t)

		netuid := uint16(0)
		total_issuance := types.NewU64(1000)

		metagraph, err := runtime.GetMetagraph(env.Client, netuid, nil)
		require.NoError(t, err, "Failed to get initial metagraph")
		initialTotalIssuance := metagraph.TotalIssuance

		sudoCall, err := SudoSetTotalIssuanceCall(env.Client, netuid, total_issuance)
		require.NoError(t, err, "Failed to create sudo_set_total_issuance call")

		ext, err := NewSudoExt(env.Client, &sudoCall)
		require.NoError(t, err, "Failed to create sudo extrinsic")

		testutils.SignAndSubmit(t, env.Client, ext, env.Alice.Coldkey.Keypair, uint32(env.Alice.Coldkey.AccInfo.Nonce))

		metagraph, err = runtime.GetMetagraph(env.Client, netuid, nil)
		require.NoError(t, err, "Failed to get updated metagraph")

		newTotalIssuance := types.NewU64(uint64(metagraph.TotalIssuance.Int64()))
		require.Equal(t, total_issuance, newTotalIssuance, "Total issuance was not set correctly")
		require.NotEqual(t, initialTotalIssuance, newTotalIssuance, "Total issuance did not change")
	})

}
