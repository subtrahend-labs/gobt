//go:build integration
// +build integration

package extrinsics

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerateCommit(t *testing.T) {
	env := setup(t)
	defer env.Teardown()

	// Test parameters
	uids := []uint16{0, 1, 2}
	vals := []uint16{100, 200, 300}
	versionKey := uint64(843000)
	tempo := uint64(360)
	currentBlock := uint64(1000)
	netuid := uint16(1)
	revealEpochs := uint64(10)
	blockTime := 1.0
	hotkey := env.Bob.Hotkey.Keypair.PublicKey

	// Test GenerateCommit
	commitBytes, revealRound, err := GenerateCommit(
		uids,
		vals,
		versionKey,
		tempo,
		currentBlock,
		netuid,
		revealEpochs,
		blockTime,
		hotkey,
	)

	require.NoError(t, err, "Failed to generate commit")
	require.NotEmpty(t, commitBytes, "Generated commit should not be empty")
	require.Greater(t, revealRound, currentBlock, "Reveal round should be greater than current block")
}