//go:build integration
// +build integration

package extrinsics

import (
	"math/big"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/subtrahend-labs/gobt/testutils"
)

func TestBalanceModuleExtrinsics(t *testing.T) {
	t.Run("TransferAllowDeath", func(t *testing.T) {
		setup(t)
		defer teardown(t)

		amountU64 := uint64(100000000)
		bobInitial := uint64(bob.accountInfo.Data.Free)
		charlieInitial := uint64(charlie.accountInfo.Data.Free)
		ext, err := TransferAllowDeathExt(env.Client, charlie.address, types.NewUCompact(new(big.Int).SetUint64(amountU64)))
		require.NoError(t, err, "Failed to create extrinsic")
		testutils.SignAndSubmit(t, env.Client, ext, bob.keyring, uint32(bob.accountInfo.Nonce))

		updateUserInfo(t, &bob)
		updateUserInfo(t, &charlie)

		bobFinal := uint64(bob.accountInfo.Data.Free)
		charlieFinal := uint64(charlie.accountInfo.Data.Free)
		bobDiff := bobInitial - bobFinal
		assert.GreaterOrEqual(t, bobDiff, amountU64,
			"Bob balance didn't decrease by at least %v: initial=%v, final=%v, diff=%v",
			amountU64, bobInitial, bobFinal, bobDiff)

		charlieDiff := charlieFinal - charlieInitial
		assert.Equal(t, amountU64, charlieDiff,
			"Charlie balance didn't increase by %v: initial=%v, final=%v, diff=%v",
			amountU64, charlieInitial, charlieFinal, charlieDiff)
	})

	t.Run("TransferKeepAlive", func(t *testing.T) {
		setup(t)
		defer teardown(t)

		amountU64 := uint64(100000000)
		bobInitial := uint64(bob.accountInfo.Data.Free)
		charlieInitial := uint64(charlie.accountInfo.Data.Free)
		ext, err := TransferKeepAliveExt(env.Client, charlie.address, types.NewUCompact(new(big.Int).SetUint64(amountU64)))
		require.NoError(t, err, "Failed to create extrinsic")
		testutils.SignAndSubmit(t, env.Client, ext, bob.keyring, uint32(bob.accountInfo.Nonce))

		updateUserInfo(t, &bob)
		updateUserInfo(t, &charlie)

		bobFinal := uint64(bob.accountInfo.Data.Free)
		charlieFinal := uint64(charlie.accountInfo.Data.Free)
		bobDiff := bobInitial - bobFinal
		assert.GreaterOrEqual(t, bobDiff, amountU64,
			"Bob balance didn't decrease by at least %v: initial=%v, final=%v, diff=%v",
			amountU64, bobInitial, bobFinal, bobDiff)

		charlieDiff := charlieFinal - charlieInitial
		assert.Equal(t, amountU64, charlieDiff,
			"Charlie balance didn't increase by %v: initial=%v, final=%v, diff=%v",
			amountU64, charlieInitial, charlieFinal, charlieDiff)
	})

	t.Run("ForceTransfer", func(t *testing.T) {
		setup(t)
		defer teardown(t)

		amountU64 := uint64(100000000)
		source := bob.address
		recipient := charlie.address
		bobInitial := uint64(bob.accountInfo.Data.Free)
		charlieInitial := uint64(charlie.accountInfo.Data.Free)

		forceTransferCall, err := ForceTransferCall(env.Client, source, recipient, types.NewUCompact(new(big.Int).SetUint64(amountU64)))
		require.NoError(t, err, "Failed to create Call")
		ext := NewSudo(env.Client, &forceTransferCall)
		testutils.SignAndSubmit(t, env.Client, ext, alice.keyring, uint32(alice.accountInfo.Nonce))

		updateUserInfo(t, &bob)
		updateUserInfo(t, &charlie)

		bobFinal := uint64(bob.accountInfo.Data.Free)
		charlieFinal := uint64(charlie.accountInfo.Data.Free)

		bobDiff := bobInitial - bobFinal
		charlieDiff := charlieFinal - charlieInitial
		assert.GreaterOrEqual(t, bobDiff, amountU64, "Bob balance didn't decrease by at least %v: initial=%v, final=%v, diff=%v")
		assert.Equal(t, amountU64, charlieDiff, "Charlie balance didn't increase by %v: initial=%v, final=%v, diff=%v")

	})

	// t.Run("TransferAll", func(t *testing.T) {
	// 	setup(t)
	// 	defer teardown(t)

	// 	bobInitial := uint64(bob.accountInfo.Data.Free)
	// 	charlieInitial := uint64(charlie.accountInfo.Data.Free)
	// 	ext, err := TransferAllExt(env.Client, charlie.address, false)
	// 	require.NoError(t, err, "Failed to create extrinsic")
	// 	testutils.SignAndSubmit(t, env.Client, ext, bob.keyring, uint32(bob.accountInfo.Nonce))

	// 	updateUserInfo(t, &bob)
	// 	updateUserInfo(t, &charlie)

	// 	bobFinal := uint64(bob.accountInfo.Data.Free)
	// 	charlieFinal := uint64(charlie.accountInfo.Data.Free)
	// 	bobDiff := bobInitial - bobFinal
	// 	assert.Equal(t, bobDiff, bobInitial,
	// 		"Bob balance didn't decrease by %v: initial=%v, final=%v, diff=%v",
	// 		bobInitial, bobInitial, bobFinal, bobDiff)

	// 	//charlieDiff := charlieFinal - charlieInitial
	// })

	t.Run("SetWeights", func(t *testing.T) {
		setup(t)
		defer teardown(t)

		netuid := types.NewU16(0)
		dests := []types.U16{types.NewU16(1), types.NewU16(2)}
		weights := []types.U16{types.NewU16(10), types.NewU16(20)}
		versionKey := types.NewU64(843000)

		call, err := SetWeightsCall(env.Client, netuid, dests, weights, versionKey)
		require.NoError(t, err, "Failed to create call")
		ext := NewSudo(env.Client, &call)
		testutils.SignAndSubmit(t, env.Client, ext, alice.keyring, uint32(alice.accountInfo.Nonce))

		updateUserInfo(t, &bob)
		updateUserInfo(t, &charlie)

		assert.Equal(t, initialBalanceU64, bob.accountInfo.Data.Free,
			"Bob balance didn't change: initial=%v, final=%v", initialBalanceU64, bob.accountInfo.Data.Free)
		assert.Equal(t, initialBalanceU64, charlie.accountInfo.Data.Free,
			"Charlie balance didn't change: initial=%v, final=%v", initialBalanceU64, charlie.accountInfo.Data.Free)
	})
}
