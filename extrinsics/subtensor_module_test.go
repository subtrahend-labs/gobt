//go:build integration
// +build integration

package extrinsics

import (
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	//	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/subtrahend-labs/gobt/testutils"
)

func TestSubtensorModuleExtrinsics(t *testing.T) {

	t.Run("SetWeights", func(t *testing.T) {
		setup(t)
		defer teardown(t)

		netuid := types.NewU16(0)
		uids := []types.U16{types.NewU16(1), types.NewU16(2)}
		weights := []types.U16{types.NewU16(10), types.NewU16(20)}
		versionKey := types.NewU64(843000)

		ext, err := SetWeightsExt(env.Client, netuid, uids, weights, versionKey)
		require.NoError(t, err, "Failed to create set_weights ext")
		testutils.SignAndSubmit(t, env.Client, ext, bob.keyring, uint32(bob.accountInfo.Nonce))
		updateUserInfo(t, &bob)

	})
}
