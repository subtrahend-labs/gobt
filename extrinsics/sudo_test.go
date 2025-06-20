//go:build integration
// +build integration

package extrinsics

import (
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/stretchr/testify/require"
	"github.com/subtrahend-labs/gobt/testutils"
)

func TestSudoModuleExtrinsics(t *testing.T) {
	t.Parallel()

	t.Run("NewSudoCall", func(t *testing.T) {
		t.Parallel()
		env := setup(t)

		innerCall, _ := SudoSetNetworkRateLimitCall(env.Client, types.NewU64(100))

		sudoCall, _ := NewSudoCall(env.Client, &innerCall)
		require.NotEmpty(t, sudoCall, "Sudo call should not be empty")
	})

	t.Run("NewSudoExt", func(t *testing.T) {
		t.Parallel()
		env := setup(t)

		innerCall, _ := SudoSetNetworkRateLimitCall(env.Client, types.NewU64(150))

		sudoExt, _ := NewSudoExt(env.Client, &innerCall)
		require.NotNil(t, sudoExt, "Sudo extrinsic should not be nil")

		testutils.SignAndSubmit(
			t,
			env.Client,
			sudoExt,
			env.Alice.Coldkey.Keypair,
			uint32(env.Alice.Coldkey.AccInfo.Nonce),
		)

		t.Logf("Successfully executed sudo call with network rate limit: 150")
	})

	t.Run("SudoUncheckedWeight", func(t *testing.T) {
		t.Parallel()
		env := setup(t)

		innerCall, err := SudoSetNetworkRateLimitCall(env.Client, types.NewU64(100))
		require.NoError(t, err, "Failed to create inner call")

		customWeight := types.Weight{
			RefTime:   types.NewUCompactFromUInt(1000000000),
			ProofSize: types.NewUCompactFromUInt(65536),
		}

		sudoUncheckedCall, err := NewSudoUncheckedWeightCall(env.Client, &innerCall, customWeight)
		require.NotEmpty(t, sudoUncheckedCall, "Sudo unchecked weight call should not be empty")

		sudoUncheckedExt, err := NewSudoUncheckedWeightExt(env.Client, &innerCall, customWeight)
		require.NotNil(t, sudoUncheckedExt, "Sudo unchecked weight extrinsic should not be nil")

		testutils.SignAndSubmit(
			t,
			env.Client,
			sudoUncheckedExt,
			env.Alice.Coldkey.Keypair,
			uint32(env.Alice.Coldkey.AccInfo.Nonce),
		)

		t.Logf("Successfully executed sudo_unchecked_weight with custom weight: RefTime=%v, ProofSize=%v",
			customWeight.RefTime, customWeight.ProofSize)
	})

}
