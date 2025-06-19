//go:build integration
// +build integration

package extrinsics

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/stretchr/testify/require"
	"github.com/subtrahend-labs/gobt/storage"
	"github.com/subtrahend-labs/gobt/testutils"
)

func TestSubtensorModuleExtrinsics(t *testing.T) {
	t.Parallel()
	t.Run("RootRegister", func(t *testing.T) {
		t.Parallel()
		env := setup(t)

		ext, err := RootRegisterExt(env.Client, *env.Bob.Hotkey.AccID)
		require.NoError(t, err, "Failed to create root_register ext")
		testutils.SignAndSubmit(t, env.Client, ext, env.Bob.Coldkey.Keypair, uint32(env.Bob.Coldkey.AccInfo.Nonce))
		updateUserInfo(t, &env.Bob, env, false)
	})

	t.Run("BurnedRegister", func(t *testing.T) {
		t.Parallel()
		env := setup(t)
		setupSubnet(t, env)

		netuid := types.NewU16(1)

		// Use BurnedRegisterExt instead of RegisterNetworkExt
		ext, err := BurnedRegisterExt(env.Client, *env.Bob.Hotkey.AccID, netuid)
		require.NoError(t, err, "Failed to create burned_register ext")
		testutils.SignAndSubmit(t, env.Client, ext, env.Bob.Coldkey.Keypair, uint32(env.Bob.Coldkey.AccInfo.Nonce))
		updateUserInfo(t, &env.Bob, env, false)
	})

	t.Run("RegisterNetwork", func(t *testing.T) {
		t.Parallel()
		env := setup(t)

		ext, err := RootRegisterExt(env.Client, *env.Bob.Hotkey.AccID)
		require.NoError(t, err, "Failed to create root_register ext")
		testutils.SignAndSubmit(t, env.Client, ext, env.Bob.Coldkey.Keypair, uint32(env.Bob.Coldkey.AccInfo.Nonce))
		fmt.Println("Bob's nonce after root_register:", env.Bob.Coldkey.AccInfo.Nonce)

		sudoCall, err := SudoSetNetworkRateLimitCall(env.Client, types.NewU64(0))
		require.NoError(t, err, "Failed to create sudo_set_network_rate_limit ext")
		ext, err = NewSudoExt(env.Client, &sudoCall)
		testutils.SignAndSubmit(t, env.Client, ext, env.Alice.Coldkey.Keypair, uint32(env.Alice.Coldkey.AccInfo.Nonce))
		require.NoError(t, err, "Failed to create root_register ext")
		fmt.Println("Will I ever make progress")

		ext, err = RegisterNetworkExt(env.Client, *env.Bob.Hotkey.AccID)
		require.NoError(t, err, "Failed to create register_network ext")
		testutils.SignAndSubmit(t, env.Client, ext, env.Bob.Coldkey.Keypair, uint32(env.Bob.Coldkey.AccInfo.Nonce+1))
		fmt.Println("Here we are again on my own")
	})

	t.Run("ServeAxon", func(t *testing.T) {
		t.Parallel()
		env := setup(t)

		setupSubnet(t, env)

		netuid := types.NewU16(1)
		ext, err := RootRegisterExt(env.Client, *env.Bob.Hotkey.AccID)
		require.NoError(t, err, "Failed to create root_register ext")
		testutils.SignAndSubmit(t, env.Client, ext, env.Bob.Coldkey.Keypair, uint32(env.Bob.Coldkey.AccInfo.Nonce))
		updateUserInfo(t, &env.Bob, env, false)

		version := types.NewU32(0)

		ipInt, _ := new(big.Int).SetString("1676056785", 10)
		ip := types.NewU128(*ipInt)

		port := types.NewU16(8080)
		ipType := types.NewU8(4)
		protocol := types.NewU8(0)
		placeholder1 := types.NewU8(0)
		placeholder2 := types.NewU8(0)

		serveAxonExt, err := ServeAxonExt(
			env.Client,
			netuid,
			version,
			ip,
			port,
			ipType,
			protocol,
			placeholder1,
			placeholder2,
		)
		require.NoError(t, err, "Failed to create serve_axon ext")

		testutils.SignAndSubmit(
			t,
			env.Client,
			serveAxonExt,
			env.Bob.Hotkey.Keypair,
			uint32(0),
		)

	})

	t.Run("ServeAxonTLS", func(t *testing.T) {
		t.Parallel()
		env := setup(t)

		setupSubnet(t, env)

		netuid := types.NewU16(1)

		ext, err := RootRegisterExt(env.Client, *env.Bob.Hotkey.AccID)
		require.NoError(t, err, "Failed to create root_register ext")
		testutils.SignAndSubmit(t, env.Client, ext, env.Bob.Coldkey.Keypair, uint32(env.Bob.Coldkey.AccInfo.Nonce))
		updateUserInfo(t, &env.Bob, env, false)

		// Now test ServeAxon with Bob's hotkey
		version := types.NewU32(0)

		ipInt, _ := new(big.Int).SetString("1676056785", 10)
		ip := types.NewU128(*ipInt)

		port := types.NewU16(8080)
		ipType := types.NewU8(4)   // IPv4
		protocol := types.NewU8(0) // HTTP
		placeholder1 := types.NewU8(0)
		placeholder2 := types.NewU8(0)

		// Create and submit the ServeAxon extrinsic
		serveAxonExt, err := ServeAxonExt(
			env.Client,
			netuid,
			version,
			ip,
			port,
			ipType,
			protocol,
			placeholder1,
			placeholder2,
		)
		require.NoError(t, err, "Failed to create serve_axon ext")

		// Sign and submit the extrinsic
		testutils.SignAndSubmit(
			t,
			env.Client,
			serveAxonExt,
			env.Bob.Hotkey.Keypair,
			uint32(0),
		)

	})

	t.Run("AddStake", func(t *testing.T) {
		t.Parallel()
		env := setup(t)

		setupSubnet(t, env)

		netuid := types.NewU16(1)
		ext, err := RootRegisterExt(env.Client, *env.Bob.Hotkey.AccID)
		require.NoError(t, err, "Failed to create root_register ext")
		testutils.SignAndSubmit(t, env.Client, ext, env.Bob.Coldkey.Keypair, uint32(env.Bob.Coldkey.AccInfo.Nonce))
		updateUserInfo(t, &env.Bob, env, false)

		initialBalance := uint64(env.Bob.Coldkey.AccInfo.Data.Free)
		t.Logf("Bob's initial balance: %v TAO", initialBalance)

		amount_staked := types.NewU64(1000000000)

		addStakeExt, err := AddStakeExt(
			env.Client,
			*env.Bob.Hotkey.AccID,
			netuid,
			amount_staked,
		)

		testutils.SignAndSubmit(
			t,
			env.Client,
			addStakeExt,
			env.Bob.Coldkey.Keypair,
			uint32(env.Bob.Coldkey.AccInfo.Nonce),
		)

		finalInfo, err := storage.GetAccountInfo(env.Client, env.Bob.Coldkey.Keypair.PublicKey, nil)
		finalBalance := uint64(finalInfo.Data.Free)
		t.Logf("Bob's final balance: %v TAO", finalBalance)
		require.Equal(t, initialBalance-uint64(amount_staked), finalBalance, "Balance should decrease by staked amount")
	})

	t.Run("RemoveStake", func(t *testing.T) {
		t.Parallel()
		env := setup(t)

		setupSubnet(t, env)

		netuid := types.NewU16(1)

		// Register Bob to root network
		ext, err := RootRegisterExt(env.Client, *env.Bob.Hotkey.AccID)
		require.NoError(t, err, "Failed to create root_register ext")
		testutils.SignAndSubmit(t, env.Client, ext, env.Bob.Coldkey.Keypair, uint32(env.Bob.Coldkey.AccInfo.Nonce))
		updateUserInfo(t, &env.Bob, env, false)

		// Register Bob to the subnet (required for staking)
		burnedRegisterExt, err := BurnedRegisterExt(env.Client, *env.Bob.Hotkey.AccID, netuid)
		require.NoError(t, err, "Failed to create burned_register ext")
		testutils.SignAndSubmit(t, env.Client, burnedRegisterExt, env.Bob.Coldkey.Keypair, uint32(env.Bob.Coldkey.AccInfo.Nonce))
		updateUserInfo(t, &env.Bob, env, false)

		// Add some stake first
		amount_staked := types.NewU64(1000000000)
		addStakeExt, err := AddStakeExt(
			env.Client,
			*env.Bob.Hotkey.AccID,
			netuid,
			amount_staked,
		)
		require.NoError(t, err, "Failed to create add_stake ext")
		testutils.SignAndSubmit(
			t,
			env.Client,
			addStakeExt,
			env.Bob.Coldkey.Keypair,
			uint32(env.Bob.Coldkey.AccInfo.Nonce),
		)
		updateUserInfo(t, &env.Bob, env, false)

		initialBalance := uint64(env.Bob.Coldkey.AccInfo.Data.Free)
		t.Logf("Bob's initial balance: %v TAO", initialBalance)

		amount_unstaked := types.NewU64(500000000)

		removeStakeExt, err := RemoveStakeExt(
			env.Client,
			*env.Bob.Hotkey.AccID,
			netuid,
			amount_unstaked,
		)
		require.NoError(t, err, "Failed to create remove_stake ext")

		testutils.SignAndSubmit(
			t,
			env.Client,
			removeStakeExt,
			env.Bob.Coldkey.Keypair,
			uint32(env.Bob.Coldkey.AccInfo.Nonce),
		)

		finalInfo, err := storage.GetAccountInfo(env.Client, env.Bob.Coldkey.Keypair.PublicKey, nil)
		require.NoError(t, err, "Failed to get final account info")
		finalBalance := uint64(finalInfo.Data.Free)
		t.Logf("Bob's final balance: %v TAO", finalBalance)

		balanceIncrease := finalBalance - initialBalance
		require.GreaterOrEqual(t, balanceIncrease, uint64(amount_unstaked),
			"Balance should increase by at least the unstaked amount (may include emission rewards)")
	})

	t.Run("AddStakeLimit", func(t *testing.T) {
		t.Parallel()
		env := setup(t)
		setupSubnet(t, env)

		netuid := types.NewU16(1)
		ext, err := RootRegisterExt(env.Client, *env.Bob.Hotkey.AccID)
		require.NoError(t, err, "Failed to create root_register ext")
		testutils.SignAndSubmit(t, env.Client, ext, env.Bob.Coldkey.Keypair, uint32(env.Bob.Coldkey.AccInfo.Nonce))
		updateUserInfo(t, &env.Bob, env, false)

		initialBalance := uint64(env.Bob.Coldkey.AccInfo.Data.Free)
		t.Logf("Bob's initial balance: %v TAO", initialBalance)

		amount_staked := types.NewU64(10000000000)
		limit_price := types.NewU64(10000000000000000000)
		allow_partial := types.NewBool(true)

		addStakeLimitExt, err := AddStakeLimitExt(
			env.Client,
			*env.Bob.Hotkey.AccID,
			netuid,
			amount_staked,
			limit_price,
			allow_partial,
		)

		testutils.SignAndSubmit(
			t,
			env.Client,
			addStakeLimitExt,
			env.Bob.Coldkey.Keypair,
			uint32(env.Bob.Coldkey.AccInfo.Nonce),
		)

		finalInfo, err := storage.GetAccountInfo(env.Client, env.Bob.Coldkey.Keypair.PublicKey, nil)
		finalBalance := uint64(finalInfo.Data.Free)
		t.Logf("Bob's final balance: %v TAO", finalBalance)
		require.Equal(t, initialBalance-uint64(amount_staked), finalBalance, "Balance should not decrease by staked amount")
	})

	t.Run("RemoveStakeLimit", func(t *testing.T) {
		t.Parallel()
		env := setup(t)
		setupSubnet(t, env)

		netuid := types.NewU16(1)

		ext, err := RootRegisterExt(env.Client, *env.Bob.Hotkey.AccID)
		require.NoError(t, err, "Failed to create root_register ext")
		testutils.SignAndSubmit(t, env.Client, ext, env.Bob.Coldkey.Keypair, uint32(env.Bob.Coldkey.AccInfo.Nonce))
		updateUserInfo(t, &env.Bob, env, false)

		burnedRegisterExt, err := BurnedRegisterExt(env.Client, *env.Bob.Hotkey.AccID, netuid)
		require.NoError(t, err, "Failed to create burned_register ext")
		testutils.SignAndSubmit(t, env.Client, burnedRegisterExt, env.Bob.Coldkey.Keypair, uint32(env.Bob.Coldkey.AccInfo.Nonce))
		updateUserInfo(t, &env.Bob, env, false)

		amount_staked := types.NewU64(1000000000)
		addStakeExt, err := AddStakeExt(
			env.Client,
			*env.Bob.Hotkey.AccID,
			netuid,
			amount_staked,
		)
		require.NoError(t, err, "Failed to create add_stake ext")
		testutils.SignAndSubmit(
			t,
			env.Client,
			addStakeExt,
			env.Bob.Coldkey.Keypair,
			uint32(env.Bob.Coldkey.AccInfo.Nonce),
		)
		updateUserInfo(t, &env.Bob, env, false)

		initialBalance := uint64(env.Bob.Coldkey.AccInfo.Data.Free)
		t.Logf("Bob's initial balance: %v TAO", initialBalance)

		amount_unstaked := types.NewU64(500000000)

		unstake_limit_price := types.NewU64(1)
		allow_partial_unstake := types.NewBool(true)

		removeStakeLimitExt, err := RemoveStakeLimitExt(
			env.Client,
			*env.Bob.Hotkey.AccID,
			netuid,
			amount_unstaked,
			unstake_limit_price,
			allow_partial_unstake,
		)
		require.NoError(t, err, "Failed to create remove_stake_limit ext")

		testutils.SignAndSubmit(
			t,
			env.Client,
			removeStakeLimitExt,
			env.Bob.Coldkey.Keypair,
			uint32(env.Bob.Coldkey.AccInfo.Nonce),
		)

		finalInfo, err := storage.GetAccountInfo(env.Client, env.Bob.Coldkey.Keypair.PublicKey, nil)
		finalBalance := uint64(finalInfo.Data.Free)
		t.Logf("Bob's final balance: %v TAO", finalBalance)

		balanceIncrease := finalBalance - initialBalance
		require.GreaterOrEqual(t, balanceIncrease, uint64(amount_unstaked),
			"Balance should increase by at least the unstaked amount (may include emission rewards)")
	})

}
