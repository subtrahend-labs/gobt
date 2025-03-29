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
	charlie                User
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

	keyringCharlie, err := signature.KeyringPairFromSecret("//Charlie", 0)
	charlieAccID, err := types.NewAccountID(keyringCharlie.PublicKey)
	require.NoError(t, err, "Failed to create Charlie account ID")
	charlieInitialInfo, err := storage.GetAccountInfo(env.Client, keyringCharlie.PublicKey)
	require.NoError(t, err, "Failed to get Charlie balance")
	multiAddress, err = types.NewMultiAddressFromAccountID(keyringCharlie.PublicKey)
	require.NoError(t, err, "Failed to create Charlie multi address")
	charlie = User{
		username:    "Charlie",
		keyring:     keyringCharlie,
		accID:       charlieAccID,
		accountInfo: charlieInitialInfo,
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

	charlieCall, err := types.NewCall(env.Client.Meta, "Balances.force_set_balance", charlie.address, initialBalanceUCompact)
	require.NoError(t, err, "Failed to create call")
	extSudo = NewSudo(env.Client, &charlieCall)
	testutils.SignAndSubmit(t, env.Client, extSudo, alice.keyring, aliceNonce+1)

	updateUserInfo(t, &bob)
	updateUserInfo(t, &charlie)

	assert.Equal(t, initialBalanceU64, bob.accountInfo.Data.Free, "Bob balance reset failed")
	assert.Equal(t, initialBalanceU64, charlie.accountInfo.Data.Free, "Charlie balance reset failed")
}

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

	t.Run("TransferAll", func(t *testing.T) {
		setup(t)
		defer teardown(t)

		bobInitial := uint64(bob.accountInfo.Data.Free)
		charlieInitial := uint64(charlie.accountInfo.Data.Free)
		ext, err := TransferAllExt(env.Client, charlie.address, false)
		require.NoError(t, err, "Failed to create extrinsic")
		testutils.SignAndSubmit(t, env.Client, ext, bob.keyring, uint32(bob.accountInfo.Nonce))

		updateUserInfo(t, &bob)
		updateUserInfo(t, &charlie)

		bobFinal := uint64(bob.accountInfo.Data.Free)
		charlieFinal := uint64(charlie.accountInfo.Data.Free)
		bobDiff := bobInitial - bobFinal
		assert.Equal(t, bobDiff, bobInitial,
			"Bob balance didn't decrease by %v: initial=%v, final=%v, diff=%v",
			bobInitial, bobInitial, bobFinal, bobDiff)

		//charlieDiff := charlieFinal - charlieInitial
	})
}
