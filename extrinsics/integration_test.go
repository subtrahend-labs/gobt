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

func setup(t *testing.T) *testutils.TestEnv {
	env, err := testutils.Setup()
	if err != nil {
		assert.FailNow(t, "Failed creating test setup")
	}
	env.InitialBalanceUint64 = 1000000000000000
	env.InitialBalanceU64 = types.NewU64(env.InitialBalanceUint64)
	env.InitialBalanceUCompact = types.NewUCompactFromUInt(env.InitialBalanceUint64)
	env.ZeroUint64 = 0
	env.ZeroU64 = types.NewU64(0)
	env.ZeroUCompact = types.NewUCompactFromUInt(env.ZeroUint64)

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

	env.Alice = testutils.User{
		Username: "Alice",
		Coldkey: testutils.Key{
			Keypair: keyringAliceCold,
			Address: aliceColdAddress,
			AccID:   aliceColdAccID,
			AccInfo: aliceColdInfo,
		},
		Hotkey: testutils.Key{
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

	env.Bob = testutils.User{
		Username: "Bob",
		Coldkey: testutils.Key{
			Keypair: keyringBobCold,
			Address: bobColdAddress,
			AccID:   bobColdAccID,
			AccInfo: bobColdInfo,
		},
		Hotkey: testutils.Key{
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

	env.Charlie = testutils.User{
		Username: "Charlie",
		Coldkey: testutils.Key{
			Keypair: keyringCharlieCold,
			Address: charlieColdAddress,
			AccID:   charlieColdAccID,
			AccInfo: charlieColdInfo,
		},
		Hotkey: testutils.Key{
			Keypair: keyringCharlieHot,
			Address: charlieHotAddress,
			AccID:   charlieHotAccID,
			AccInfo: nil,
		},
	}
	return env
}

func setupSubnet(t *testing.T, env *testutils.TestEnv) {
	sudoCall, err := SudoSetNetworkRateLimitCall(env.Client, types.NewU64(0))
	require.NoError(t, err, "Failed to create sudo_set_network_rate_limit ext")
	ext, err := NewSudoExt(env.Client, &sudoCall)
	testutils.SignAndSubmit(t, env.Client, ext, env.Alice.Coldkey.Keypair, uint32(env.Alice.Coldkey.AccInfo.Nonce))
	require.NoError(t, err, "Failed to create root_register ext")
	updateUserInfo(t, &env.Alice, env, false)

	ext, err = RegisterNetworkExt(env.Client, *env.Bob.Hotkey.AccID)
	require.NoError(t, err, "Failed to create register_network ext")
	testutils.SignAndSubmit(t, env.Client, ext, env.Bob.Coldkey.Keypair, uint32(env.Bob.Coldkey.AccInfo.Nonce))
	updateUserInfo(t, &env.Bob, env, false)
}

func updateUserInfo(t *testing.T, u *testutils.User, env *testutils.TestEnv, doHot bool) {
	infoCold, err := storage.GetAccountInfo(env.Client, u.Coldkey.Keypair.PublicKey)
	require.NoError(t, err, "Failed to get %s coldkey balance", u.Username)
	u.Coldkey.AccInfo = infoCold
	if doHot {
		infoHot, err := storage.GetAccountInfo(env.Client, u.Hotkey.Keypair.PublicKey)
		require.NoError(t, err, "Failed to get %s hotkey balance", u.Username)
		u.Hotkey.AccInfo = infoHot
	}
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
