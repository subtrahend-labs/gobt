//go:build integration
// +build integration

package extrinsics

import (
	"fmt"
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
