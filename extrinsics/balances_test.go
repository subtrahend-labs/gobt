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
	t.Parallel()
	t.Run("TransferAllowDeath", func(t *testing.T) {
		env := setup(t)

		amountU64 := uint64(100000000)
		bobInitial := uint64(env.Bob.Coldkey.AccInfo.Data.Free)
		charlieInitial := uint64(env.Charlie.Coldkey.AccInfo.Data.Free)
		ext, err := TransferAllowDeathExt(env.Client, env.Charlie.Coldkey.Address, types.NewUCompact(new(big.Int).SetUint64(amountU64)))
		require.NoError(t, err, "Failed to create extrinsic")
		testutils.SignAndSubmit(t, env.Client, ext, env.Bob.Coldkey.Keypair, uint32(env.Bob.Coldkey.AccInfo.Nonce))

		updateUserInfo(t, &env.Bob, env)
		updateUserInfo(t, &env.Charlie, env)

		bobFinal := uint64(env.Bob.Coldkey.AccInfo.Data.Free)
		charlieFinal := uint64(env.Charlie.Coldkey.AccInfo.Data.Free)
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
		env := setup(t)

		amountU64 := uint64(100000000)
		bobInitial := uint64(env.Bob.Coldkey.AccInfo.Data.Free)
		charlieInitial := uint64(env.Charlie.Coldkey.AccInfo.Data.Free)
		ext, err := TransferKeepAliveExt(env.Client, env.Charlie.Coldkey.Address, types.NewUCompact(new(big.Int).SetUint64(amountU64)))
		require.NoError(t, err, "Failed to create extrinsic")
		testutils.SignAndSubmit(t, env.Client, ext, env.Bob.Coldkey.Keypair, uint32(env.Bob.Coldkey.AccInfo.Nonce))

		updateUserInfo(t, &env.Bob, env)
		updateUserInfo(t, &env.Charlie, env)

		bobFinal := uint64(env.Bob.Coldkey.AccInfo.Data.Free)
		charlieFinal := uint64(env.Charlie.Coldkey.AccInfo.Data.Free)
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
		env := setup(t)

		amountU64 := uint64(100000000)
		source := env.Bob.Coldkey.Address
		recipient := env.Charlie.Coldkey.Address
		bobInitial := uint64(env.Bob.Coldkey.AccInfo.Data.Free)
		charlieInitial := uint64(env.Charlie.Coldkey.AccInfo.Data.Free)

		forceTransferCall, err := ForceTransferCall(env.Client, source, recipient, types.NewUCompact(new(big.Int).SetUint64(amountU64)))
		require.NoError(t, err, "Failed to create Call")
		ext, err := NewSudoExt(env.Client, &forceTransferCall)
		require.NoError(t, err, "Failed to create extrinsic")
		testutils.SignAndSubmit(t, env.Client, ext, env.Alice.Coldkey.Keypair, uint32(env.Alice.Coldkey.AccInfo.Nonce))

		updateUserInfo(t, &env.Bob, env)
		updateUserInfo(t, &env.Charlie, env)

		bobFinal := uint64(env.Bob.Coldkey.AccInfo.Data.Free)
		charlieFinal := uint64(env.Charlie.Coldkey.AccInfo.Data.Free)

		bobDiff := bobInitial - bobFinal
		charlieDiff := charlieFinal - charlieInitial
		assert.GreaterOrEqual(t, bobDiff, amountU64, "Bob balance didn't decrease by at least %v: initial=%v, final=%v, diff=%v")
		assert.Equal(t, amountU64, charlieDiff, "Charlie balance didn't increase by %v: initial=%v, final=%v, diff=%v")
	})

	// t.Run("TransferAll", func(t *testing.T) {
	// 	setup(t)
	// 	defer teardown(t)

	// 	bobInitial := uint64(bob.coldkey.AccInfo.Data.Free)
	// 	charlieInitial := uint64(charlie.coldkey.AccInfo.Data.Free)
	// 	ext, err := TransferAllExt(env.Client, charlie.address, false)
	// 	require.NoError(t, err, "Failed to create extrinsic")
	// 	testutils.SignAndSubmit(t, env.Client, ext, bob.keyring, uint32(bob.coldkey.AccInfo.Nonce))

	// 	updateUserInfo(t, &bob)
	// 	updateUserInfo(t, &charlie)

	// 	bobFinal := uint64(bob.coldkey.AccInfo.Data.Free)
	// 	charlieFinal := uint64(charlie.coldkey.AccInfo.Data.Free)
	// 	bobDiff := bobInitial - bobFinal
	// 	assert.Equal(t, bobDiff, bobInitial,
	// 		"Bob balance didn't decrease by %v: initial=%v, final=%v, diff=%v",
	// 		bobInitial, bobInitial, bobFinal, bobDiff)

	// 	//charlieDiff := charlieFinal - charlieInitial
	// })
}
