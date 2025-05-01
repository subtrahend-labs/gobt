//go:build integration
// +build integration

package extrinsics

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/stretchr/testify/require"
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

		// First, set up a subnet
		setupSubnet(t, env)

		// Bob needs to register to the subnet first using BurnedRegister
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

		// Update user info after transaction
		updateUserInfo(t, &env.Bob, env, false)
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

		// Update user info after transaction
		updateUserInfo(t, &env.Bob, env, false)
	})

	t.Run("CommitCRV3Weights", func(t *testing.T) {
		t.Parallel()
		env := setup(t)
		setupSubnet(t, env)
		netuid := types.NewU16(1)

		ext, err := RootRegisterExt(env.Client, *env.Bob.Hotkey.AccID)
		require.NoError(t, err)
		testutils.SignAndSubmit(t, env.Client, ext, env.Bob.Coldkey.Keypair, uint32(env.Bob.Coldkey.AccInfo.Nonce))
		updateUserInfo(t, &env.Bob, env, false)

		ext, err = BurnedRegisterExt(env.Client, *env.Bob.Hotkey.AccID, netuid)
		require.NoError(t, err)
		testutils.SignAndSubmit(t, env.Client, ext, env.Bob.Coldkey.Keypair, uint32(env.Bob.Coldkey.AccInfo.Nonce))
		updateUserInfo(t, &env.Bob, env, false)

		// 3) pick a trivial self‐weight to commit & reveal
		uids := []types.U16{types.NewU16(1)}
		weights := []types.U16{types.NewU16(1)}

		// convert to raw uint16 slices for the FFI
		uidsRaw := make([]uint16, len(uids))
		valsRaw := make([]uint16, len(weights))
		for i := range uids {
			uidsRaw[i] = uint16(uids[i])
			valsRaw[i] = uint16(weights[i])
		}

		// 4) grab on‐chain params
		currentBlock, err := env.Client.GetBlockNumber()
		require.NoError(t, err)

		// these must match your network’s hyperparams;
		// adjust or fetch them programmatically if they differ
		versionKey := uint64(env.VersionKey)
		tempo := uint64(1)
		revealEpochs := uint64(1)
		blockTime := float64(6)

		// 5) call into Rust to get the encrypted commit + reveal round
		commitBytes, revealRound, err := GenerateCommit(
			uidsRaw, valsRaw,
			versionKey,
			tempo,
			currentBlock,
			uint16(netuid),
			revealEpochs,
			blockTime,
		)
		require.NoError(t, err)
		require.NotEmpty(t, commitBytes, "encrypted commit should not be empty")
		require.Greater(t, revealRound, currentBlock, "reveal round must be in the future")

		// 6) submit the commit_crv3_weights extrinsic
		commitExt, err := CommitCRV3WeightsExt(env.Client, netuid, types.Bytes(commitBytes), types.NewU64(revealRound))
		require.NoError(t, err)
		testutils.SignAndSubmit(t, env.Client, commitExt, env.Bob.Coldkey.Keypair, uint32(env.Bob.Coldkey.AccInfo.Nonce))
		updateUserInfo(t, &env.Bob, env, false)

		// 7) wait until we reach the reveal round
		for i := 0; i < 20; i++ {
			blk, err := env.Client.GetBlockNumber()
			require.NoError(t, err)
			if blk >= revealRound {
				break
			}
			time.Sleep(6 * time.Second)
		}

		// 8) submit the reveal_weights extrinsic
		revealExt, err := RevealWeightsExt(
			env.Client,
			netuid,
			uids,
			weights,
			types.NewU64(versionKey),
			types.NewU64(revealRound),
		)
		require.NoError(t, err)
		testutils.SignAndSubmit(t, env.Client, revealExt, env.Bob.Coldkey.Keypair, uint32(env.Bob.Coldkey.AccInfo.Nonce))
		updateUserInfo(t, &env.Bob, env, false)

	})

}
