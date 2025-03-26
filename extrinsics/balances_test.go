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
		aliceInitialBalance := new(big.Int).SetUint64(uint64(aliceInitialInfo.Data.Free))
		aliceFinalBalance := new(big.Int).SetUint64(uint64(aliceFinalInfo.Data.Free))
		actualAliceDiff := new(big.Int).Sub(aliceInitialBalance, aliceFinalBalance)
		assert.GreaterOrEqual(t, actualAliceDiff.Cmp(amount), 0,
			"Alice balance didn't decrease by at least %v: initial=%v, final=%v, diff=%v",
			amount, aliceInitialBalance, aliceFinalBalance, actualAliceDiff)
		bobInitialBalance := new(big.Int).SetUint64(uint64(bobInitialInfo.Data.Free))
		bobFinalBalance := new(big.Int).SetUint64(uint64(bobFinalInfo.Data.Free))
		actualBobDiff := new(big.Int).Sub(bobFinalBalance, bobInitialBalance)
		assert.Equal(t, 0, actualBobDiff.Cmp(amount),
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
