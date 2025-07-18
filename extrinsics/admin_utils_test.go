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

	t.Run("SudoSetDefaultTake", func(t *testing.T) {
		t.Parallel()
		env := setup(t)

		defaultTake := types.NewU16(1000)
		sudoCall, _ := SudoSetDefaultTakeCall(env.Client, defaultTake)
		ext, _ := NewSudoExt(env.Client, &sudoCall)
		testutils.SignAndSubmit(t, env.Client, ext, env.Alice.Coldkey.Keypair, uint32(env.Alice.Coldkey.AccInfo.Nonce))

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

		storageKey, _ := types.CreateStorageKey(env.Client.Meta, "SubtensorModule", "ServingRateLimit", typetools.Uint16ToBytes(uint16(4)))
		var newServingRateLimit types.U64
		env.Client.Api.RPC.State.GetStorageLatest(storageKey, &newServingRateLimit)

		require.Equal(t, servingRateLimit, newServingRateLimit, "Serving rate limit was not updated correctly")
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
	t.Run("SudoSetAdjustmentInterval", func(t *testing.T) {
		t.Parallel()
		env := setup(t)

		adjustmentInterval := types.NewU16(1000)
		sudoCall, _ := SudoSetAdjustmentIntervalCall(env.Client, 0, adjustmentInterval)
		ext, _ := NewSudoExt(env.Client, &sudoCall)
		testutils.SignAndSubmit(t, env.Client, ext, env.Alice.Coldkey.Keypair, uint32(env.Alice.Coldkey.AccInfo.Nonce))

		storageKey, _ := types.CreateStorageKey(env.Client.Meta, "SubtensorModule", "AdjustmentInterval", typetools.Uint16ToBytes(uint16(0)))
		var newAdjustmentInterval types.U16
		env.Client.Api.RPC.State.GetStorageLatest(storageKey, &newAdjustmentInterval)

		require.Equal(t, adjustmentInterval, newAdjustmentInterval, "Adjustment interval was not updated correctly")
	})

	t.Run("SudoSetAdjustmentAlpha", func(t *testing.T) {
		t.Parallel()
		env := setup(t)

		adjustmentAlpha := types.NewU64(1000)
		sudoCall, _ := SudoSetAdjustmentAlphaCall(env.Client, 0, adjustmentAlpha)
		ext, _ := NewSudoExt(env.Client, &sudoCall)
		testutils.SignAndSubmit(t, env.Client, ext, env.Alice.Coldkey.Keypair, uint32(env.Alice.Coldkey.AccInfo.Nonce))

		storageKey, _ := types.CreateStorageKey(env.Client.Meta, "SubtensorModule", "AdjustmentAlpha", typetools.Uint16ToBytes(uint16(0)))
		var newAdjustmentAlpha types.U64
		env.Client.Api.RPC.State.GetStorageLatest(storageKey, &newAdjustmentAlpha)

		require.Equal(t, adjustmentAlpha, newAdjustmentAlpha, "Adjustment alpha was not updated correctly")
	})

	t.Run("SudoSetMaxWeightLimit", func(t *testing.T) {
		t.Parallel()
		env := setup(t)

		maxWeightLimit := types.NewU16(1000)
		sudoCall, _ := SudoSetMaxWeightLimitCall(env.Client, 0, maxWeightLimit)
		ext, _ := NewSudoExt(env.Client, &sudoCall)
		testutils.SignAndSubmit(t, env.Client, ext, env.Alice.Coldkey.Keypair, uint32(env.Alice.Coldkey.AccInfo.Nonce))

		storageKey, _ := types.CreateStorageKey(env.Client.Meta, "SubtensorModule", "MaxWeightsLimit", typetools.Uint16ToBytes(uint16(0)))
		var newMaxWeightLimit types.U16
		env.Client.Api.RPC.State.GetStorageLatest(storageKey, &newMaxWeightLimit)

		require.Equal(t, maxWeightLimit, newMaxWeightLimit, "Max weight limit was not updated correctly")
	})

	t.Run("SudoSetImmunityPeriod", func(t *testing.T) {
		t.Parallel()
		env := setup(t)

		immunityPeriod := types.NewU16(1000)
		sudoCall, _ := SudoSetImmunityPeriodCall(env.Client, 0, immunityPeriod)
		ext, _ := NewSudoExt(env.Client, &sudoCall)
		testutils.SignAndSubmit(t, env.Client, ext, env.Alice.Coldkey.Keypair, uint32(env.Alice.Coldkey.AccInfo.Nonce))

		storageKey, _ := types.CreateStorageKey(env.Client.Meta, "SubtensorModule", "ImmunityPeriod", typetools.Uint16ToBytes(uint16(0)))
		var newImmunityPeriod types.U16
		env.Client.Api.RPC.State.GetStorageLatest(storageKey, &newImmunityPeriod)

		require.Equal(t, immunityPeriod, newImmunityPeriod, "Immunity period was not updated correctly")
	})

	t.Run("SudoSetMinAllowedWeights", func(t *testing.T) {
		t.Parallel()
		env := setup(t)

		minAllowedWeights := types.NewU16(1000)
		sudoCall, _ := SudoSetMinAllowedWeightsCall(env.Client, 0, minAllowedWeights)
		ext, _ := NewSudoExt(env.Client, &sudoCall)
		testutils.SignAndSubmit(t, env.Client, ext, env.Alice.Coldkey.Keypair, uint32(env.Alice.Coldkey.AccInfo.Nonce))

		storageKey, _ := types.CreateStorageKey(env.Client.Meta, "SubtensorModule", "MinAllowedWeights", typetools.Uint16ToBytes(uint16(0)))
		var newMinAllowedWeights types.U16
		env.Client.Api.RPC.State.GetStorageLatest(storageKey, &newMinAllowedWeights)

		require.Equal(t, minAllowedWeights, newMinAllowedWeights, "Min allowed weights was not updated correctly")
	})
	t.Run("SudoSetMaxAllowedUids", func(t *testing.T) {
		t.Parallel()
		env := setup(t)

		maxAllowedUids := types.NewU16(1000)
		sudoCall, _ := SudoSetMaxAllowedUidsCall(env.Client, 0, maxAllowedUids)
		ext, _ := NewSudoExt(env.Client, &sudoCall)
		testutils.SignAndSubmit(t, env.Client, ext, env.Alice.Coldkey.Keypair, uint32(env.Alice.Coldkey.AccInfo.Nonce))

		storageKey, _ := types.CreateStorageKey(env.Client.Meta, "SubtensorModule", "MaxAllowedUids", typetools.Uint16ToBytes(uint16(0)))
		var newMaxAllowedUids types.U16
		env.Client.Api.RPC.State.GetStorageLatest(storageKey, &newMaxAllowedUids)

		require.Equal(t, maxAllowedUids, newMaxAllowedUids, "Max allowed uids was not updated correctly")
	})

	t.Run("SudoSetKappa", func(t *testing.T) {
		t.Parallel()
		env := setup(t)

		kappa := types.NewU16(1000)
		sudoCall, _ := SudoSetKappaCall(env.Client, 0, kappa)
		ext, _ := NewSudoExt(env.Client, &sudoCall)
		testutils.SignAndSubmit(t, env.Client, ext, env.Alice.Coldkey.Keypair, uint32(env.Alice.Coldkey.AccInfo.Nonce))

		storageKey, _ := types.CreateStorageKey(env.Client.Meta, "SubtensorModule", "Kappa", typetools.Uint16ToBytes(uint16(0)))
		var newKappa types.U16
		env.Client.Api.RPC.State.GetStorageLatest(storageKey, &newKappa)

		require.Equal(t, kappa, newKappa, "Kappa was not updated correctly")
	})

	t.Run("SudoSetRho", func(t *testing.T) {
		t.Parallel()
		env := setup(t)

		rho := types.NewU16(1000)
		sudoCall, _ := SudoSetRhoCall(env.Client, 0, rho)
		ext, _ := NewSudoExt(env.Client, &sudoCall)
		testutils.SignAndSubmit(t, env.Client, ext, env.Alice.Coldkey.Keypair, uint32(env.Alice.Coldkey.AccInfo.Nonce))

		storageKey, _ := types.CreateStorageKey(env.Client.Meta, "SubtensorModule", "Rho", typetools.Uint16ToBytes(uint16(0)))
		var newRho types.U16
		env.Client.Api.RPC.State.GetStorageLatest(storageKey, &newRho)

		require.Equal(t, rho, newRho, "Rho was not updated correctly")

	})

	t.Run("SudoSetActivityCutoff", func(t *testing.T) {
		t.Parallel()
		env := setup(t)

		activityCutoff := types.NewU16(1000)
		sudoCall, _ := SudoSetActivityCutoffCall(env.Client, 0, activityCutoff)
		ext, _ := NewSudoExt(env.Client, &sudoCall)
		testutils.SignAndSubmit(t, env.Client, ext, env.Alice.Coldkey.Keypair, uint32(env.Alice.Coldkey.AccInfo.Nonce))

		storageKey, _ := types.CreateStorageKey(env.Client.Meta, "SubtensorModule", "ActivityCutoff", typetools.Uint16ToBytes(uint16(0)))
		var newActivityCutoff types.U16
		env.Client.Api.RPC.State.GetStorageLatest(storageKey, &newActivityCutoff)

		require.Equal(t, activityCutoff, newActivityCutoff, "Activity cutoff was not updated correctly")
	})

	t.Run("SudoSetNetworkRegistrationAllowed", func(t *testing.T) {
		t.Parallel()
		env := setup(t)

		registrationAllowed := true
		sudoCall, _ := SudoSetNetworkRegistrationAllowedCall(env.Client, 0, registrationAllowed)
		ext, _ := NewSudoExt(env.Client, &sudoCall)
		testutils.SignAndSubmit(t, env.Client, ext, env.Alice.Coldkey.Keypair, uint32(env.Alice.Coldkey.AccInfo.Nonce))

		storageKey, _ := types.CreateStorageKey(env.Client.Meta, "SubtensorModule", "NetworkRegistrationAllowed", typetools.Uint16ToBytes(uint16(0)))
		var newNetworkRegistrationAllowed bool
		env.Client.Api.RPC.State.GetStorageLatest(storageKey, &newNetworkRegistrationAllowed)

		require.Equal(t, registrationAllowed, newNetworkRegistrationAllowed, "Network registration allowed was not updated correctly")
	})

	t.Run("SudoSetNetworkPowRegistrationAllowed", func(t *testing.T) {
		t.Parallel()
		env := setup(t)

		registrationAllowed := true
		sudoCall, _ := SudoSetNetworkPowRegistrationAllowedCall(env.Client, 0, registrationAllowed)
		ext, _ := NewSudoExt(env.Client, &sudoCall)
		testutils.SignAndSubmit(t, env.Client, ext, env.Alice.Coldkey.Keypair, uint32(env.Alice.Coldkey.AccInfo.Nonce))

		storageKey, _ := types.CreateStorageKey(env.Client.Meta, "SubtensorModule", "NetworkPowRegistrationAllowed", typetools.Uint16ToBytes(uint16(0)))
		var newNetworkPowRegistrationAllowed bool
		env.Client.Api.RPC.State.GetStorageLatest(storageKey, &newNetworkPowRegistrationAllowed)

		require.Equal(t, registrationAllowed, newNetworkPowRegistrationAllowed, "Network pow registration allowed was not updated correctly")
	})

	t.Run("SudoSetTargetRegistrationsPerInterval", func(t *testing.T) {
		t.Parallel()
		env := setup(t)

		targetRegistrationsPerInterval := types.NewU16(1000)
		sudoCall, _ := SudoSetTargetRegistrationsPerIntervalCall(env.Client, 0, targetRegistrationsPerInterval)
		ext, _ := NewSudoExt(env.Client, &sudoCall)
		testutils.SignAndSubmit(t, env.Client, ext, env.Alice.Coldkey.Keypair, uint32(env.Alice.Coldkey.AccInfo.Nonce))

		storageKey, _ := types.CreateStorageKey(env.Client.Meta, "SubtensorModule", "TargetRegistrationsPerInterval", typetools.Uint16ToBytes(uint16(0)))
		var newTargetRegistrationsPerInterval types.U16
		env.Client.Api.RPC.State.GetStorageLatest(storageKey, &newTargetRegistrationsPerInterval)

		require.Equal(t, targetRegistrationsPerInterval, newTargetRegistrationsPerInterval, "Target registrations per interval was not updated correctly")
	})

	t.Run("SudoSetMinBurn", func(t *testing.T) {
		t.Parallel()
		env := setup(t)

		minBurn := types.NewU64(1000)
		sudoCall, _ := SudoSetMinBurnCall(env.Client, 0, minBurn)
		ext, _ := NewSudoExt(env.Client, &sudoCall)
		testutils.SignAndSubmit(t, env.Client, ext, env.Alice.Coldkey.Keypair, uint32(env.Alice.Coldkey.AccInfo.Nonce))

		storageKey, _ := types.CreateStorageKey(env.Client.Meta, "SubtensorModule", "MinBurn", typetools.Uint16ToBytes(uint16(0)))
		var newMinBurn types.U64
		env.Client.Api.RPC.State.GetStorageLatest(storageKey, &newMinBurn)

		require.Equal(t, minBurn, newMinBurn, "Min burn was not updated correctly")
	})

	t.Run("SudoSetMaxBurn", func(t *testing.T) {
		t.Parallel()
		env := setup(t)

		maxBurn := types.NewU64(1000)
		sudoCall, _ := SudoSetMaxBurnCall(env.Client, 0, maxBurn)
		ext, _ := NewSudoExt(env.Client, &sudoCall)
		testutils.SignAndSubmit(t, env.Client, ext, env.Alice.Coldkey.Keypair, uint32(env.Alice.Coldkey.AccInfo.Nonce))

		storageKey, _ := types.CreateStorageKey(env.Client.Meta, "SubtensorModule", "MaxBurn", typetools.Uint16ToBytes(uint16(0)))
		var newMaxBurn types.U64
		env.Client.Api.RPC.State.GetStorageLatest(storageKey, &newMaxBurn)

		require.Equal(t, maxBurn, newMaxBurn, "Max burn was not updated correctly")
	})

	t.Run("SudoSetDifficulty", func(t *testing.T) {
		t.Parallel()
		env := setup(t)

		difficulty := types.NewU64(1000)
		sudoCall, _ := SudoSetDifficultyCall(env.Client, 0, difficulty)
		ext, _ := NewSudoExt(env.Client, &sudoCall)
		testutils.SignAndSubmit(t, env.Client, ext, env.Alice.Coldkey.Keypair, uint32(env.Alice.Coldkey.AccInfo.Nonce))

		storageKey, _ := types.CreateStorageKey(env.Client.Meta, "SubtensorModule", "Difficulty", typetools.Uint16ToBytes(uint16(0)))
		var newDifficulty types.U64
		env.Client.Api.RPC.State.GetStorageLatest(storageKey, &newDifficulty)

		require.Equal(t, difficulty, newDifficulty, "Difficulty was not updated correctly")
	})

	t.Run("SudoSetMaxAllowedValidators", func(t *testing.T) {
		t.Parallel()
		env := setup(t)

		//Rust code has check that max allowed validators is less than or equal to max allowed uids
		//So we increase max allowed uids first
		maxAllowedUids := types.NewU16(1000)
		sudoCall, _ := SudoSetMaxAllowedUidsCall(env.Client, 0, maxAllowedUids)
		ext, _ := NewSudoExt(env.Client, &sudoCall)
		testutils.SignAndSubmit(t, env.Client, ext, env.Alice.Coldkey.Keypair, uint32(env.Alice.Coldkey.AccInfo.Nonce))
		updateUserInfo(t, &env.Alice, env, false) //updated so that nonce is updated for next transaction

		maxAllowedValidators := types.NewU16(10)
		sudoCall, _ = SudoSetMaxAllowedValidatorsCall(env.Client, 0, maxAllowedValidators)
		ext, _ = NewSudoExt(env.Client, &sudoCall)
		testutils.SignAndSubmit(t, env.Client, ext, env.Alice.Coldkey.Keypair, uint32(env.Alice.Coldkey.AccInfo.Nonce))

		storageKey, _ := types.CreateStorageKey(env.Client.Meta, "SubtensorModule", "MaxAllowedValidators", typetools.Uint16ToBytes(uint16(0)))
		var newMaxAllowedValidators types.U16
		env.Client.Api.RPC.State.GetStorageLatest(storageKey, &newMaxAllowedValidators)

		require.Equal(t, maxAllowedValidators, newMaxAllowedValidators, "Max allowed validators was not updated correctly")
	})

	t.Run("SudoSetBondsMovingAverage", func(t *testing.T) {
		t.Parallel()
		env := setup(t)

		bondsMovingAverage := types.NewU64(1000)
		sudoCall, _ := SudoSetBondsMovingAverageCall(env.Client, 0, bondsMovingAverage)
		ext, _ := NewSudoExt(env.Client, &sudoCall)
		testutils.SignAndSubmit(t, env.Client, ext, env.Alice.Coldkey.Keypair, uint32(env.Alice.Coldkey.AccInfo.Nonce))

		storageKey, _ := types.CreateStorageKey(env.Client.Meta, "SubtensorModule", "BondsMovingAverage", typetools.Uint16ToBytes(uint16(0)))
		var newBondsMovingAverage types.U64
		env.Client.Api.RPC.State.GetStorageLatest(storageKey, &newBondsMovingAverage)

		require.Equal(t, bondsMovingAverage, newBondsMovingAverage, "Bonds moving average was not updated correctly")
	})

	t.Run("SudoSetBondsPenalty", func(t *testing.T) {
		t.Parallel()
		env := setup(t)

		bondsPenalty := types.NewU16(1000)
		sudoCall, _ := SudoSetBondsPenaltyCall(env.Client, 0, bondsPenalty)
		ext, _ := NewSudoExt(env.Client, &sudoCall)
		testutils.SignAndSubmit(t, env.Client, ext, env.Alice.Coldkey.Keypair, uint32(env.Alice.Coldkey.AccInfo.Nonce))

		storageKey, _ := types.CreateStorageKey(env.Client.Meta, "SubtensorModule", "BondsPenalty", typetools.Uint16ToBytes(uint16(0)))
		var newBondsPenalty types.U16
		env.Client.Api.RPC.State.GetStorageLatest(storageKey, &newBondsPenalty)

		require.Equal(t, bondsPenalty, newBondsPenalty, "Bonds penalty was not updated correctly")
	})

	t.Run("SudoSetMaxRegistrationsPerBlock", func(t *testing.T) {
		t.Parallel()
		env := setup(t)

		maxRegistrationsPerBlock := types.NewU16(1000)
		sudoCall, _ := SudoSetMaxRegistrationsPerBlockCall(env.Client, 0, maxRegistrationsPerBlock)
		ext, _ := NewSudoExt(env.Client, &sudoCall)
		testutils.SignAndSubmit(t, env.Client, ext, env.Alice.Coldkey.Keypair, uint32(env.Alice.Coldkey.AccInfo.Nonce))

		storageKey, _ := types.CreateStorageKey(env.Client.Meta, "SubtensorModule", "MaxRegistrationsPerBlock", typetools.Uint16ToBytes(uint16(0)))
		var newMaxRegistrationsPerBlock types.U16
		env.Client.Api.RPC.State.GetStorageLatest(storageKey, &newMaxRegistrationsPerBlock)

		require.Equal(t, maxRegistrationsPerBlock, newMaxRegistrationsPerBlock, "Max registrations per block was not updated correctly")
	})

	t.Run("SudoSetSubnetOwnerCut", func(t *testing.T) {
		t.Parallel()
		env := setup(t)

		subnetOwnerCut := types.NewU16(1000)
		sudoCall, _ := SudoSetSubnetOwnerCutCall(env.Client, subnetOwnerCut)
		ext, _ := NewSudoExt(env.Client, &sudoCall)
		testutils.SignAndSubmit(t, env.Client, ext, env.Alice.Coldkey.Keypair, uint32(env.Alice.Coldkey.AccInfo.Nonce))

		storageKey, _ := types.CreateStorageKey(env.Client.Meta, "SubtensorModule", "SubnetOwnerCut")
		var newSubnetOwnerCut types.U16
		env.Client.Api.RPC.State.GetStorageLatest(storageKey, &newSubnetOwnerCut)
		require.Equal(t, subnetOwnerCut, newSubnetOwnerCut, "Subnet owner cut was not updated correctly")

	})

}
