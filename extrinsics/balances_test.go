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

var (
	env                    *testutils.TestEnv
	alice                  User
	bob                    User
	initialBalanceUint64   uint64         = 1000000000000000
	initialBalanceU64      types.U64      = types.NewU64(initialBalanceUint64)
	initialBalanceUCompact types.UCompact = types.NewUCompactFromUInt(initialBalanceUint64)
	zeroUint64             uint64         = 0
	zeroU64                types.U64      = types.NewU64(0)
	zeroUCompact           types.UCompact = types.NewUCompactFromUInt(zeroUint64)
)

type User struct {
	username    string
	keyring     signature.KeyringPair
	address     types.MultiAddress
	accID       *types.AccountID
	accountInfo *storage.AccountInfo
}

func setup(t *testing.T) {
	keyringAlice := signature.TestKeyringPairAlice
	aliceAccID, err := types.NewAccountID(keyringAlice.PublicKey)
	require.NoError(t, err, "Failed to get Alice balance")
	multiAddress, err := types.NewMultiAddressFromAccountID(keyringAlice.PublicKey)
	require.NoError(t, err, "Failed to create Alice multi address")
	aliceInitialInfo, err := storage.GetAccountInfo(env.Client, keyringAlice.PublicKey)
	alice = User{
		username:    "Alice",
		keyring:     keyringAlice,
		accID:       aliceAccID,
		address:     multiAddress,
		accountInfo: aliceInitialInfo,
	}

	require.NoError(t, err, "Failed to create Alice account ID")

	keyringBob, err := signature.KeyringPairFromSecret("//Bob", 0)
	bobAccID, err := types.NewAccountID(keyringBob.PublicKey)
	require.NoError(t, err, "Failed to create Bob account ID")
	bobInitialInfo, err := storage.GetAccountInfo(env.Client, keyringBob.PublicKey)
	require.NoError(t, err, "Failed to get Bob balance")
	multiAddress, err = types.NewMultiAddressFromAccountID(keyringBob.PublicKey)
	require.NoError(t, err, "Failed to create Bob multi address")
	bob = User{
		username:    "Bob",
		keyring:     keyringBob,
		accID:       bobAccID,
		accountInfo: bobInitialInfo,
		address:     multiAddress,
	}
}

func updateUserInfo(t *testing.T, u *User) {
	info, err := storage.GetAccountInfo(env.Client, u.keyring.PublicKey)
	require.NoError(t, err, "Failed to get %s balance", u.username)
	u.accountInfo = info
}

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

func teardown(t *testing.T) {
	updateUserInfo(t, &alice)
	aliceNonce := uint32(alice.accountInfo.Nonce)
	bobCall, err := types.NewCall(env.Client.Meta, "Balances.force_set_balance", bob.address, initialBalanceUCompact)
	require.NoError(t, err, "Failed to create call")
	extSudo := NewSudo(env.Client, &bobCall)
	testutils.SignAndSubmit(t, env.Client, extSudo, alice.keyring, aliceNonce)

	extResetAlice, err := types.NewCall(env.Client.Meta, "Balances.force_set_balance", alice.address, initialBalanceUCompact)
	require.NoError(t, err, "Failed to create call")
	extSudo = NewSudo(env.Client, &extResetAlice)
	testutils.SignAndSubmit(t, env.Client, extSudo, alice.keyring, aliceNonce+1)

	updateUserInfo(t, &alice)
	updateUserInfo(t, &bob)
	assert.Equal(t, initialBalanceU64, bob.accountInfo.Data.Free, "Bob balance reset failed")
	assert.Equal(t, initialBalanceU64, alice.accountInfo.Data.Free, "Alice balance reset failed")
}

func TestBalanceModuleExtrinsics(t *testing.T) {
	t.Run("TransferAllowDeath", func(t *testing.T) {
		setup(t)
		defer teardown(t)

		amountU64 := uint64(100000000)
		aliceInitialBalance := uint64(alice.accountInfo.Data.Free)
		fmt.Println("Alice init")
		fmt.Println(aliceInitialBalance)
		fmt.Println("bob init")
		bobInitialBalance := uint64(bob.accountInfo.Data.Free)
		fmt.Println(bobInitialBalance)
		ext := NewTransferAllowDeath(env.Client, bob.address, types.NewUCompact(new(big.Int).SetUint64(amountU64)))
		testutils.SignAndSubmit(t, env.Client, ext, alice.keyring, uint32(alice.accountInfo.Nonce))

		updateUserInfo(t, &alice)
		updateUserInfo(t, &bob)
		aliceFinalBalance := uint64(alice.accountInfo.Data.Free)
		bobFinalBalance := uint64(bob.accountInfo.Data.Free)
		actualAliceDiff := aliceInitialBalance - aliceFinalBalance
		assert.GreaterOrEqual(t, actualAliceDiff, amountU64,
			"Alice balance didn't decrease by at least %v: initial=%v, final=%v, diff=%v",
			amountU64, aliceInitialBalance, aliceFinalBalance, actualAliceDiff)
		actualBobDiff := bobFinalBalance - bobInitialBalance
		assert.Equal(t, amountU64, actualBobDiff,
			"Bob balance didn't increase by %v: initial=%v, final=%v, diff=%v",
			amountU64, bobInitialBalance, bobFinalBalance, actualBobDiff)
	})

	//	t.Run("TransferKeepAlive", func(t *testing.T) {
	//		keyringAlice := signature.TestKeyringPairAlice
	//		keyringBob, err := signature.KeyringPairFromSecret("//Bob", 0)
	//		require.NoError(t, err, "Failed to create Bob key")
	//
	//		aliceInitialInfo, err := storage.GetAccountInfo(env.Client, keyringAlice.PublicKey)
	//		require.NoError(t, err, "Failed to get Alice balance")
	//
	//		bobInitialInfo, err := storage.GetAccountInfo(env.Client, keyringBob.PublicKey)
	//		require.NoError(t, err, "Failed to get Bob balance")
	//
	//		amountU64 := uint64(100000000)
	//		amount := new(big.Int).SetUint64(amountU64)
	//		bobMultiAddress, err := types.NewMultiAddressFromAccountID(keyringBob.PublicKey)
	//		require.NoError(t, err, "Failed to create Bob multi address")
	//
	//		ext := NewTransferKeepAlive(env.Client, bobMultiAddress, amount)
	//		testutils.SignAndSubmit(t, env.Client, ext, keyringAlice, uint64(aliceInitialInfo.Nonce))
	//
	//		aliceFinalInfo, err := storage.GetAccountInfo(env.Client, keyringAlice.PublicKey)
	//		require.NoError(t, err, "Failed to get Alice final balance")
	//		bobFinalInfo, err := storage.GetAccountInfo(env.Client, keyringBob.PublicKey)
	//		require.NoError(t, err, "Failed to get Bob final balance")
	//
	//		aliceInitialBalance := uint64(aliceInitialInfo.Data.Free)
	//		aliceFinalBalance := uint64(aliceFinalInfo.Data.Free)
	//		actualAliceDiff := aliceInitialBalance - aliceFinalBalance
	//		bobInitialBalance := uint64(bobInitialInfo.Data.Free)
	//		bobFinalBalance := uint64(bobFinalInfo.Data.Free)
	//		actualBobDiff := bobFinalBalance - bobInitialBalance
	//
	//		assert.GreaterOrEqual(t, actualAliceDiff, amountU64,
	//			"Alice balance didn't decrease by at least %v: initial=%v, final=%v, diff=%v",
	//			amount, aliceInitialBalance, aliceFinalBalance, actualAliceDiff)
	//
	//		assert.Equal(t, amountU64, actualBobDiff,
	//			"Bob balance didn't increase by %v: initial=%v, final=%v, diff=%v",
	//			amount, bobInitialBalance, bobFinalBalance, actualBobDiff)
	//	})
}
