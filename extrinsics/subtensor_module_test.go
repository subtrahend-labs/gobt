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
		updateUserInfo(t, &env.Bob, env)
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
		updateUserInfo(t, &env.Bob, env)
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
		ext, err := BurnedRegisterExt(env.Client, *env.Bob.Hotkey.AccID, netuid)
		require.NoError(t, err, "Failed to create burned_register ext")
		testutils.SignAndSubmit(t, env.Client, ext, env.Bob.Coldkey.Keypair, uint32(env.Bob.Coldkey.AccInfo.Nonce))
		updateUserInfo(t, &env.Bob, env)

		// Now test ServeAxon with Bob's hotkey
		version := types.NewU32(1)

		// 127.0.0.1 in hex (0x7F000001) converted to big.Int
		ipInt, _ := new(big.Int).SetString("2130706433", 10) // decimal representation of 127.0.0.1
		ip := types.NewU128(*ipInt)

		port := types.NewU16(8080)
		ipType := types.NewU8(4)   // IPv4
		protocol := types.NewU8(0) // HTTP
		placeholder1 := types.NewU8(0)
		placeholder2 := types.NewU8(0)

		// Create an empty certificate for this test
		certificate := []byte{}

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
			certificate,
		)
		require.NoError(t, err, "Failed to create serve_axon ext")

		// Sign and submit the extrinsic
		testutils.SignAndSubmit(
			t,
			env.Client,
			serveAxonExt,
			env.Bob.Coldkey.Keypair,
			uint32(env.Bob.Coldkey.AccInfo.Nonce),
		)

		// Update user info after transaction
		updateUserInfo(t, &env.Bob, env)
	})

	// t.Run("SetWeights", func(t *testing.T) {
	// setup(t)
	// defer teardown(t)

	// netuid := types.NewU16(4)
	// uids := []types.U16{types.NewU16(1), types.NewU16(2)}
	// weights := []types.U16{types.NewU16(10), types.NewU16(20)}
	// versionKey := types.NewU64(843000)

	// ext, err := SetWeightsExt(env.Client, netuid, uids, weights, versionKey)
	// require.NoError(t, err, "Failed to create set_weights ext")
	// testutils.SignAndSubmit(t, env.Client, ext, bob.keyring, uint32(bob.accountInfo.Nonce))
	// updateUserInfo(t, &bob)

	// })
}
