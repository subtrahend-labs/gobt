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

type Key struct {
	Keypair signature.KeyringPair
	Address types.MultiAddress
	AccID   *types.AccountID
	AccInfo *storage.AccountInfo
}

type User struct {
	username string
	coldkey  Key
	hotkey   Key
}

func setup(t *testing.T) {
	// Alice coldkey setup
	keyringAliceCold := signature.TestKeyringPairAlice
	aliceColdAccID, err := types.NewAccountID(keyringAliceCold.PublicKey)
	require.NoError(t, err, "Failed to get Alice coldkey account ID")
	aliceColdAddress, err := types.NewMultiAddressFromAccountID(keyringAliceCold.PublicKey)
	require.NoError(t, err, "Failed to create Alice coldkey multi address")
	aliceColdInfo, err := storage.GetAccountInfo(env.Client, keyringAliceCold.PublicKey)
	require.NoError(t, err, "Failed to get Alice coldkey account info")

	// Alice hotkey setup
	keyringAliceHot, err := signature.KeyringPairFromSecret("//Alice//hot", 42)
	require.NoError(t, err, "Failed to create Alice hotkey")
	aliceHotAccID, err := types.NewAccountID(keyringAliceHot.PublicKey)
	require.NoError(t, err, "Failed to get Alice hotkey account ID")
	aliceHotAddress, err := types.NewMultiAddressFromAccountID(keyringAliceHot.PublicKey)
	require.NoError(t, err, "Failed to create Alice hotkey multi address")
	// aliceHotInfo, err := storage.GetAccountInfo(env.Client, keyringAliceHot.PublicKey)
	// require.NoError(t, err, fmt.Sprintf("Failed to get Alice hotkey account info: %v", err))

	alice = User{
		username: "Alice",
		coldkey: Key{
			Keypair: keyringAliceCold,
			Address: aliceColdAddress,
			AccID:   aliceColdAccID,
			AccInfo: aliceColdInfo,
		},
		hotkey: Key{
			Keypair: keyringAliceHot,
			Address: aliceHotAddress,
			AccID:   aliceHotAccID,
			AccInfo: nil,
		},
	}

	// Bob coldkey setup
	keyringBobCold, err := signature.KeyringPairFromSecret("//Bob", 42)
	require.NoError(t, err, "Failed to create Bob coldkey")
	bobColdAccID, err := types.NewAccountID(keyringBobCold.PublicKey)
	require.NoError(t, err, "Failed to create Bob coldkey account ID")
	bobColdAddress, err := types.NewMultiAddressFromAccountID(keyringBobCold.PublicKey)
	require.NoError(t, err, "Failed to create Bob coldkey multi address")
	bobColdInfo, err := storage.GetAccountInfo(env.Client, keyringBobCold.PublicKey)
	require.NoError(t, err, "Failed to get Bob coldkey balance")

	// bob hotkey setup
	keyringBobHot, err := signature.KeyringPairFromSecret("//Bob//Hot", 42)
	require.NoError(t, err, "Failed to create Bob hotkey")
	bobHotAccID, err := types.NewAccountID(keyringBobHot.PublicKey)
	require.NoError(t, err, "Failed to get Bob hotkey account ID")
	bobHotAddress, err := types.NewMultiAddressFromAccountID(keyringBobHot.PublicKey)
	require.NoError(t, err, "Failed to create Bob hotkey multi address")
	// bobHotInfo, err := storage.GetAccountInfo(env.Client, keyringBobHot.PublicKey)
	// require.NoError(t, err, "Failed to get Bob hotkey account info")

	bob = User{
		username: "Bob",
		coldkey: Key{
			Keypair: keyringBobCold,
			Address: bobColdAddress,
			AccID:   bobColdAccID,
			AccInfo: bobColdInfo,
		},
		hotkey: Key{
			Keypair: keyringBobHot,
			Address: bobHotAddress,
			AccID:   bobHotAccID,
			AccInfo: nil,
		},
	}

	// Charlie coldkey setup
	keyringCharlieCold, err := signature.KeyringPairFromSecret("//Charlie", 0)
	require.NoError(t, err, "Failed to create Charlie coldkey")
	charlieColdAccID, err := types.NewAccountID(keyringCharlieCold.PublicKey)
	require.NoError(t, err, "Failed to create Charlie coldkey account ID")
	charlieColdAddress, err := types.NewMultiAddressFromAccountID(keyringCharlieCold.PublicKey)
	require.NoError(t, err, "Failed to create Charlie coldkey multi address")
	charlieColdInfo, err := storage.GetAccountInfo(env.Client, keyringCharlieCold.PublicKey)
	require.NoError(t, err, "Failed to get Charlie coldkey balance")

	// charlie hotkey setup
	keyringCharlieHot, err := signature.KeyringPairFromSecret("//CharlieHot"+fmt.Sprintf("%d", 44), 0)
	require.NoError(t, err, "Failed to create Charlie hotkey")
	charlieHotAccID, err := types.NewAccountID(keyringCharlieHot.PublicKey)
	require.NoError(t, err, "Failed to get Charlie hotkey account ID")
	charlieHotAddress, err := types.NewMultiAddressFromAccountID(keyringCharlieHot.PublicKey)
	require.NoError(t, err, "Failed to create Charlie hotkey multi address")
	// charlieHotInfo, err := storage.GetAccountInfo(env.Client, keyringCharlieHot.PublicKey)
	// require.NoError(t, err, "Failed to get Charlie hotkey account info")

	charlie = User{
		username: "Charlie",
		coldkey: Key{
			Keypair: keyringCharlieCold,
			Address: charlieColdAddress,
			AccID:   charlieColdAccID,
			AccInfo: charlieColdInfo,
		},
		hotkey: Key{
			Keypair: keyringCharlieHot,
			Address: charlieHotAddress,
			AccID:   charlieHotAccID,
			AccInfo: nil,
		},
	}
}

func updateUserInfo(t *testing.T, u *User) {
	info, err := storage.GetAccountInfo(env.Client, u.coldkey.Keypair.PublicKey)
	require.NoError(t, err, "Failed to get %s coldkey balance", u.username)
	u.coldkey.AccInfo = info
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
	aliceNonce := uint32(alice.coldkey.AccInfo.Nonce)
	bobCall, err := types.NewCall(env.Client.Meta, "Balances.force_set_balance", bob.coldkey.Address, initialBalanceUCompact)
	require.NoError(t, err, "Failed to create Bob balance call")
	extSudo, err := NewSudoExt(env.Client, &bobCall)
	require.NoError(t, err, "Failed to create sudo extrinsic")
	testutils.SignAndSubmit(t, env.Client, extSudo, alice.coldkey.Keypair, aliceNonce)

	charlieCall, err := types.NewCall(env.Client.Meta, "Balances.force_set_balance", charlie.coldkey.Address, initialBalanceUCompact)
	require.NoError(t, err, "Failed to create Charlie balance call")
	extSudo, err = NewSudoExt(env.Client, &charlieCall)
	require.NoError(t, err, "Failed to create sudo extrinsic")
	testutils.SignAndSubmit(t, env.Client, extSudo, alice.coldkey.Keypair, aliceNonce+1)

	updateUserInfo(t, &bob)
	updateUserInfo(t, &charlie)

	assert.Equal(t, initialBalanceU64, bob.coldkey.AccInfo.Data.Free, "Bob balance reset failed")
	assert.Equal(t, initialBalanceU64, charlie.coldkey.AccInfo.Data.Free, "Charlie balance reset failed")
}
