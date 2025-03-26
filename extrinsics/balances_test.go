//go:build integration
// +build integration

package extrinsics

import (
	"fmt"
	"math/big"
	"os"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/subtrahend-labs/gobt/storage"
	"github.com/subtrahend-labs/gobt/testutils"
	"github.com/vedhavyas/go-subkey/v2"
)

var env *testutils.TestEnv

func TestMain(m *testing.M) {
	var err error
	env, err = testutils.Setup()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Test setup failed: %v\n", err)
		os.Exit(1)
	}
	defer env.Teardown()

	os.Exit(m.Run())
}

func resetState(t *testing.T) {
	keyringAlice := signature.TestKeyringPairAlice
	aliceAccID, err := types.NewAccountID(keyringAlice.PublicKey)
	require.NoError(t, err, "Failed to create Alice account ID")
	keyringBob, err := signature.KeyringPairFromSecret("//Bob", 0)
	bobAccID, err := types.NewAccountID(keyringBob.PublicKey)
	require.NoError(t, err, "Failed to create Bob account ID")

	const initialBalance = 1_000_000
	bigZero := types.U128{big.NewInt(0)}
	aliceBalance := types.U128{big.NewInt(1_000_000_000_000)}
	bobBalance := types.U128{big.NewInt(0)}

	// Use root (or a sudo key if configured) to force-set balances
	sudoKey := signature.TestKeyringPairAlice // Assuming Alice has sudo in dev mode
	extAlice := NewForceSetBalance(env.Client, *aliceAccID, aliceBalance, types.U128{big.NewInt(0)})
	testutils.SignAndSubmit(t, env.Client, extAlice, sudoKey, 0) // Nonce might need adjustment

	extBob := NewForceSetBalance(env.Client, *bobAccID, bobBalance, types.U128{big.NewInt(0)})
	testutils.SignAndSubmit(t, env.Client, extBob, sudoKey, 0) // Nonce might need adjustment

	// Verify reset
	aliceInfo, err := storage.GetAccountInfo(env.Client, keyringAlice.PublicKey)
	require.NoError(t, err, "Failed to get Alice balance after reset")
	assert.Equal(t, uint64(initialBalance), uint64(aliceInfo.Data.Free), "Alice balance reset failed")

	bobInfo, err := storage.GetAccountInfo(env.Client, keyringBob.PublicKey)
	require.NoError(t, err, "Failed to get Bob balance after reset")
	assert.Equal(t, uint64(initialBalance), uint64(bobInfo.Data.Free), "Bob balance reset failed")
}

func TestBalanceModuleExtrinsics(t *testing.T) {
	t.Run("TransferAllowDeath", func(t *testing.T) {

		keyringAlice := signature.TestKeyringPairAlice
		keyringBob, err := signature.KeyringPairFromSecret("//Bob", 0)
		require.NoError(t, err, "Failed to create Bob key")

		aliceInitialInfo, err := storage.GetAccountInfo(env.Client, keyringAlice.PublicKey)
		require.NoError(t, err, "Failed to get Alice balance")
		bobInitialInfo, err := storage.GetAccountInfo(env.Client, keyringBob.PublicKey)
		require.NoError(t, err, "Failed to get Bob balance")

		amountU64 := uint64(100000000)
		amount := new(big.Int).SetUint64(amountU64)
		bobMultiAddress, err := types.NewMultiAddressFromAccountID(keyringBob.PublicKey)
		require.NoError(t, err, "Failed to create Bob multi address")

		ext := NewTransferAllowDeath(env.Client, bobMultiAddress, amount)
		testutils.SignAndSubmit(t, env.Client, ext, keyringAlice, uint64(aliceInitialInfo.Nonce))

		aliceFinalInfo, err := storage.GetAccountInfo(env.Client, keyringAlice.PublicKey)
		require.NoError(t, err, "Failed to get Alice final balance")
		bobFinalInfo, err := storage.GetAccountInfo(env.Client, keyringBob.PublicKey)
		require.NoError(t, err, "Failed to get Bob final balance")

		aliceInitialBalance := uint64(aliceInitialInfo.Data.Free)
		aliceFinalBalance := uint64(aliceFinalInfo.Data.Free)
		actualAliceDiff := aliceInitialBalance - aliceFinalBalance
		assert.GreaterOrEqual(t, actualAliceDiff, amountU64,
			"Alice balance didn't decrease by at least %v: initial=%v, final=%v, diff=%v",
			amount, aliceInitialBalance, aliceFinalBalance, actualAliceDiff)

		bobInitialBalance := uint64(bobInitialInfo.Data.Free)
		bobFinalBalance := uint64(bobFinalInfo.Data.Free)
		actualBobDiff := bobFinalBalance - bobInitialBalance
		assert.Equal(t, amountU64, actualBobDiff,
			"Bob balance didn't increase by %v: initial=%v, final=%v, diff=%v",
			amount, bobInitialBalance, bobFinalBalance, actualBobDiff)
	})

	t.Run("TransferKeepAlive", func(t *testing.T) {
		keyringAlice := signature.TestKeyringPairAlice
		keyringBob, err := signature.KeyringPairFromSecret("//Bob", 0)
		require.NoError(t, err, "Failed to create Bob key")

		aliceInitialInfo, err := storage.GetAccountInfo(env.Client, keyringAlice.PublicKey)
		require.NoError(t, err, "Failed to get Alice balance")

		bobInitialInfo, err := storage.GetAccountInfo(env.Client, keyringBob.PublicKey)
		require.NoError(t, err, "Failed to get Bob balance")

		amountU64 := uint64(100000000)
		amount := new(big.Int).SetUint64(amountU64)
		bobMultiAddress, err := types.NewMultiAddressFromAccountID(keyringBob.PublicKey)
		require.NoError(t, err, "Failed to create Bob multi address")

		ext := NewTransferKeepAlive(env.Client, bobMultiAddress, amount)
		testutils.SignAndSubmit(t, env.Client, ext, keyringAlice, uint64(aliceInitialInfo.Nonce))

		aliceFinalInfo, err := storage.GetAccountInfo(env.Client, keyringAlice.PublicKey)
		require.NoError(t, err, "Failed to get Alice final balance")
		bobFinalInfo, err := storage.GetAccountInfo(env.Client, keyringBob.PublicKey)
		require.NoError(t, err, "Failed to get Bob final balance")

		aliceInitialBalance := uint64(aliceInitialInfo.Data.Free)
		aliceFinalBalance := uint64(aliceFinalInfo.Data.Free)
		actualAliceDiff := aliceInitialBalance - aliceFinalBalance
		bobInitialBalance := uint64(bobInitialInfo.Data.Free)
		bobFinalBalance := uint64(bobFinalInfo.Data.Free)
		actualBobDiff := bobFinalBalance - bobInitialBalance

		assert.GreaterOrEqual(t, actualAliceDiff, amountU64,
			"Alice balance didn't decrease by at least %v: initial=%v, final=%v, diff=%v",
			amount, aliceInitialBalance, aliceFinalBalance, actualAliceDiff)

		assert.Equal(t, amountU64, actualBobDiff,
			"Bob balance didn't increase by %v: initial=%v, final=%v, diff=%v",
			amount, bobInitialBalance, bobFinalBalance, actualBobDiff)
	})
}
