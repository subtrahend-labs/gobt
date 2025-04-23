//go:build integration
// +build integration

package extrinsics

import (
	"fmt"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/stretchr/testify/require"
	"github.com/subtrahend-labs/gobt/testutils"
)

func TestSubtensorModuleExtrinsics(t *testing.T) {
	t.Run("RootRegister", func(t *testing.T) {
		setup(t)
		defer teardown(t)

		ext, err := RootRegisterExt(env.Client, *bob.hotkey.AccID)
		require.NoError(t, err, "Failed to create root_register ext")
		testutils.SignAndSubmit(t, env.Client, ext, bob.coldkey.Keypair, uint32(bob.coldkey.AccInfo.Nonce))
		updateUserInfo(t, &bob)
	})

	t.Skip("Register", func(t *testing.T) {
	})

	t.Run("BurnedRegister", func(t *testing.T) {
		setup(t)
		defer teardown(t)

		netuid := types.NewU16(1)

		// Use BurnedRegisterExt instead of RegisterNetworkExt
		ext, err := BurnedRegisterExt(env.Client, *bob.hotkey.AccID, netuid)
		require.NoError(t, err, "Failed to create burned_register ext")
		testutils.SignAndSubmit(t, env.Client, ext, bob.coldkey.Keypair, uint32(bob.coldkey.AccInfo.Nonce))
		updateUserInfo(t, &bob)
	})

	t.Run("RegisterNetwork", func(t *testing.T) {
		setup(t)
		defer teardown(t)

		ext, err := RootRegisterExt(env.Client, *bob.hotkey.AccID)
		require.NoError(t, err, "Failed to create root_register ext")
		testutils.SignAndSubmit(t, env.Client, ext, bob.coldkey.Keypair, uint32(bob.coldkey.AccInfo.Nonce))
		fmt.Println("Bob's nonce after root_register:", bob.coldkey.AccInfo.Nonce)

		sudoCall, err := SudoSetNetworkRateLimitCall(env.Client, types.NewU64(0))
		require.NoError(t, err, "Failed to create sudo_set_network_rate_limit ext")
		ext, err = NewSudoExt(env.Client, &sudoCall)
		testutils.SignAndSubmit(t, env.Client, ext, alice.coldkey.Keypair, uint32(alice.coldkey.AccInfo.Nonce))
		require.NoError(t, err, "Failed to create root_register ext")
		fmt.Println("Will I ever make progress")

		ext, err = RegisterNetworkExt(env.Client, *bob.hotkey.AccID)
		require.NoError(t, err, "Failed to create register_network ext")
		testutils.SignAndSubmit(t, env.Client, ext, bob.coldkey.Keypair, uint32(bob.coldkey.AccInfo.Nonce+1))
		fmt.Println("Here we are again on my own")
	})

	// t.Run("SetWeights", func(t *testing.T) {
	// 	setup(t)
	// 	defer teardown(t)

	// 	netuid := types.NewU16(4)
	// 	uids := []types.U16{types.NewU16(1), types.NewU16(2)}
	// 	weights := []types.U16{types.NewU16(10), types.NewU16(20)}
	// 	versionKey := types.NewU64(843000)

	// 	ext, err := SetWeightsExt(env.Client, netuid, uids, weights, versionKey)
	// 	require.NoError(t, err, "Failed to create set_weights ext")
	// 	testutils.SignAndSubmit(t, env.Client, ext, bob.keyring, uint32(bob.accountInfo.Nonce))
	// 	updateUserInfo(t, &bob)

	// })
}
