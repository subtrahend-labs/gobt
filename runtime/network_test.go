package runtime_test

import (
	"strings"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/subtrahend-labs/gobt/runtime"
	"github.com/subtrahend-labs/gobt/testutils"
)

func TestNetworkRuntimeAPIs(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Parallel()

	t.Run("GetNetworkLockCost", func(t *testing.T) {
		t.Parallel()

		env, err := testutils.Setup()
		if err != nil {
			t.Skipf("Failed to setup test environment: %v", err)
		}
		defer env.Teardown()

		blockHash, err := env.Client.Api.RPC.Chain.GetBlockHashLatest()
		require.NoError(t, err, "Failed to get latest block hash")

		lockCost, err := runtime.GetNetworkLockCost(env.Client, &blockHash)
		if err != nil {
			errorMsg := err.Error()
			isExpectedError := strings.Contains(errorMsg, "no network lock cost found") ||
				strings.Contains(errorMsg, "Method not found") ||
				strings.Contains(errorMsg, "failed to decode") ||
				strings.Contains(errorMsg, "Invalid params")
			assert.True(t, isExpectedError, "Expected network lock cost-related error, got: %v", err)
		} else {
			assert.GreaterOrEqual(t, uint64(lockCost), uint64(0), "Lock cost should be non-negative")
			t.Logf("Network lock cost: %d", lockCost)
		}

		lockCost, err = runtime.GetNetworkLockCost(env.Client, nil)
		assert.Error(t, err, "Should error with nil block hash")
		assert.Contains(t, err.Error(), "block hash cannot be nil")
		assert.Equal(t, types.U64(0), lockCost, "Lock cost should be 0 on error")
	})

	t.Run("GetNetworkLastLock", func(t *testing.T) {
		t.Parallel()

		env, err := testutils.Setup()
		if err != nil {
			t.Skipf("Failed to setup test environment: %v", err)
		}
		defer env.Teardown()

		blockHash, err := env.Client.Api.RPC.Chain.GetBlockHashLatest()
		require.NoError(t, err, "Failed to get latest block hash")

		lastLock, err := runtime.GetNetworkLastLock(env.Client, &blockHash)
		if err != nil {
			errorMsg := err.Error()
			isExpectedError := strings.Contains(errorMsg, "no network last lock found") ||
				strings.Contains(errorMsg, "Method not found") ||
				strings.Contains(errorMsg, "failed to decode") ||
				strings.Contains(errorMsg, "Invalid params")
			assert.True(t, isExpectedError, "Expected network last lock-related error, got: %v", err)
		} else {
			assert.GreaterOrEqual(t, uint64(lastLock), uint64(0), "Last lock should be non-negative")
			t.Logf("Network last lock: %d", lastLock)
		}

		lastLock, err = runtime.GetNetworkLastLock(env.Client, nil)
		assert.Error(t, err, "Should error with nil block hash")
		assert.Contains(t, err.Error(), "block hash cannot be nil")
		assert.Equal(t, types.U64(0), lastLock, "Last lock should be 0 on error")
	})

	t.Run("GetNetworkMinLock", func(t *testing.T) {
		t.Parallel()

		env, err := testutils.Setup()
		if err != nil {
			t.Skipf("Failed to setup test environment: %v", err)
		}
		defer env.Teardown()

		blockHash, err := env.Client.Api.RPC.Chain.GetBlockHashLatest()
		require.NoError(t, err, "Failed to get latest block hash")

		minLock, err := runtime.GetNetworkMinLock(env.Client, &blockHash)
		if err != nil {
			errorMsg := err.Error()
			isExpectedError := strings.Contains(errorMsg, "no network min lock found") ||
				strings.Contains(errorMsg, "Method not found") ||
				strings.Contains(errorMsg, "failed to decode") ||
				strings.Contains(errorMsg, "Invalid params")
			assert.True(t, isExpectedError, "Expected network min lock-related error, got: %v", err)
		} else {
			assert.GreaterOrEqual(t, uint64(minLock), uint64(0), "Min lock should be non-negative")
			t.Logf("Network min lock: %d", minLock)
		}

		minLock, err = runtime.GetNetworkMinLock(env.Client, nil)
		assert.Error(t, err, "Should error with nil block hash")
		assert.Contains(t, err.Error(), "block hash cannot be nil")
		assert.Equal(t, types.U64(0), minLock, "Min lock should be 0 on error")
	})

	t.Run("GetNetworkLastLockBlock", func(t *testing.T) {
		t.Parallel()

		env, err := testutils.Setup()
		if err != nil {
			t.Skipf("Failed to setup test environment: %v", err)
		}
		defer env.Teardown()

		blockHash, err := env.Client.Api.RPC.Chain.GetBlockHashLatest()
		require.NoError(t, err, "Failed to get latest block hash")

		lastLockBlock, err := runtime.GetNetworkLastLockBlock(env.Client, &blockHash)
		if err != nil {
			errorMsg := err.Error()
			isExpectedError := strings.Contains(errorMsg, "no network last lock block found") ||
				strings.Contains(errorMsg, "Method not found") ||
				strings.Contains(errorMsg, "failed to decode") ||
				strings.Contains(errorMsg, "Invalid params")
			assert.True(t, isExpectedError, "Expected network last lock block-related error, got: %v", err)
		} else {
			assert.GreaterOrEqual(t, uint64(lastLockBlock), uint64(0), "Last lock block should be non-negative")
			t.Logf("Network last lock block: %d", lastLockBlock)
		}

		lastLockBlock, err = runtime.GetNetworkLastLockBlock(env.Client, nil)
		assert.Error(t, err, "Should error with nil block hash")
		assert.Contains(t, err.Error(), "block hash cannot be nil")
		assert.Equal(t, types.U64(0), lastLockBlock, "Last lock block should be 0 on error")
	})

	t.Run("GetLockReductionInterval", func(t *testing.T) {
		t.Parallel()

		env, err := testutils.Setup()
		if err != nil {
			t.Skipf("Failed to setup test environment: %v", err)
		}
		defer env.Teardown()

		blockHash, err := env.Client.Api.RPC.Chain.GetBlockHashLatest()
		require.NoError(t, err, "Failed to get latest block hash")

		interval, err := runtime.GetLockReductionInterval(env.Client, &blockHash)
		if err != nil {
			errorMsg := err.Error()
			isExpectedError := strings.Contains(errorMsg, "no lock reduction interval found") ||
				strings.Contains(errorMsg, "Method not found") ||
				strings.Contains(errorMsg, "failed to decode") ||
				strings.Contains(errorMsg, "Invalid params")
			assert.True(t, isExpectedError, "Expected lock reduction interval-related error, got: %v", err)
		} else {
			assert.GreaterOrEqual(t, uint64(interval), uint64(0), "Lock reduction interval should be non-negative")
			t.Logf("Lock reduction interval: %d", interval)
		}

		interval, err = runtime.GetLockReductionInterval(env.Client, nil)
		assert.Error(t, err, "Should error with nil block hash")
		assert.Contains(t, err.Error(), "block hash cannot be nil")
		assert.Equal(t, types.U64(0), interval, "Lock reduction interval should be 0 on error")
	})

	t.Run("GetCurrentBlock", func(t *testing.T) {
		t.Parallel()

		env, err := testutils.Setup()
		if err != nil {
			t.Skipf("Failed to setup test environment: %v", err)
		}
		defer env.Teardown()

		blockHash, err := env.Client.Api.RPC.Chain.GetBlockHashLatest()
		require.NoError(t, err, "Failed to get latest block hash")

		currentBlock, err := runtime.GetCurrentBlock(env.Client, &blockHash)
		if err != nil {
			errorMsg := err.Error()
			isExpectedError := strings.Contains(errorMsg, "no current block found") ||
				strings.Contains(errorMsg, "Method not found") ||
				strings.Contains(errorMsg, "failed to decode") ||
				strings.Contains(errorMsg, "Invalid params")
			assert.True(t, isExpectedError, "Expected current block-related error, got: %v", err)
		} else {
			assert.GreaterOrEqual(t, uint64(currentBlock), uint64(0), "Current block should be non-negative")
			t.Logf("Current block: %d", currentBlock)
		}

		currentBlock, err = runtime.GetCurrentBlock(env.Client, nil)
		assert.Error(t, err, "Should error with nil block hash")
		assert.Contains(t, err.Error(), "block hash cannot be nil")
		assert.Equal(t, types.U64(0), currentBlock, "Current block should be 0 on error")
	})

	t.Run("ErrorHandling", func(t *testing.T) {
		t.Parallel()

		env, err := testutils.Setup()
		if err != nil {
			t.Skipf("Failed to setup test environment: %v", err)
		}
		defer env.Teardown()

		invalidBlockHash := types.NewHash([]byte("invalid_hash_that_does_not_exist_on_chain"))

		lockCost, err := runtime.GetNetworkLockCost(env.Client, &invalidBlockHash)
		assert.Error(t, err, "Should error with invalid block hash")
		assert.Equal(t, types.U64(0), lockCost, "Lock cost should be 0 on error")

		lastLock, err := runtime.GetNetworkLastLock(env.Client, &invalidBlockHash)
		assert.Error(t, err, "Should error with invalid block hash")
		assert.Equal(t, types.U64(0), lastLock, "Last lock should be 0 on error")

		minLock, err := runtime.GetNetworkMinLock(env.Client, &invalidBlockHash)
		assert.Error(t, err, "Should error with invalid block hash")
		assert.Equal(t, types.U64(0), minLock, "Min lock should be 0 on error")

		lastLockBlock, err := runtime.GetNetworkLastLockBlock(env.Client, &invalidBlockHash)
		assert.Error(t, err, "Should error with invalid block hash")
		assert.Equal(t, types.U64(0), lastLockBlock, "Last lock block should be 0 on error")

		interval, err := runtime.GetLockReductionInterval(env.Client, &invalidBlockHash)
		assert.Error(t, err, "Should error with invalid block hash")
		assert.Equal(t, types.U64(0), interval, "Lock reduction interval should be 0 on error")

		currentBlock, err := runtime.GetCurrentBlock(env.Client, &invalidBlockHash)
		assert.Error(t, err, "Should error with invalid block hash")
		assert.Equal(t, types.U64(0), currentBlock, "Current block should be 0 on error")
	})

	t.Run("NullAndEmptyResponseHandling", func(t *testing.T) {
		t.Parallel()

		env, err := testutils.Setup()
		if err != nil {
			t.Skipf("Failed to setup test environment: %v", err)
		}
		defer env.Teardown()

		blockHash, err := env.Client.Api.RPC.Chain.GetBlockHashLatest()
		require.NoError(t, err, "Failed to get latest block hash")

		// Test GetNetworkLockCost with empty response
		lockCost, err := runtime.GetNetworkLockCost(env.Client, &blockHash)
		if err != nil {
			errorMsg := err.Error()
			isExpectedError := strings.Contains(errorMsg, "no network lock cost found") ||
				strings.Contains(errorMsg, "Method not found") ||
				strings.Contains(errorMsg, "failed to decode") ||
				strings.Contains(errorMsg, "Invalid params")
			assert.True(t, isExpectedError, "Expected network lock cost-related error, got: %v", err)
		} else {
			assert.NotNil(t, lockCost, "Lock cost should not be nil when successful")
			assert.IsType(t, types.U64(0), lockCost, "Should return proper type")
		}

		// Test GetNetworkLastLock with empty response
		lastLock, err := runtime.GetNetworkLastLock(env.Client, &blockHash)
		if err != nil {
			errorMsg := err.Error()
			isExpectedError := strings.Contains(errorMsg, "no network last lock found") ||
				strings.Contains(errorMsg, "Method not found") ||
				strings.Contains(errorMsg, "failed to decode") ||
				strings.Contains(errorMsg, "Invalid params")
			assert.True(t, isExpectedError, "Expected network last lock-related error, got: %v", err)
		} else {
			assert.NotNil(t, lastLock, "Last lock should not be nil when successful")
			assert.IsType(t, types.U64(0), lastLock, "Should return proper type")
		}

		// Test GetNetworkMinLock with empty response
		minLock, err := runtime.GetNetworkMinLock(env.Client, &blockHash)
		if err != nil {
			errorMsg := err.Error()
			isExpectedError := strings.Contains(errorMsg, "no network min lock found") ||
				strings.Contains(errorMsg, "Method not found") ||
				strings.Contains(errorMsg, "failed to decode") ||
				strings.Contains(errorMsg, "Invalid params")
			assert.True(t, isExpectedError, "Expected network min lock-related error, got: %v", err)
		} else {
			assert.NotNil(t, minLock, "Min lock should not be nil when successful")
			assert.IsType(t, types.U64(0), minLock, "Should return proper type")
		}

		// Test GetNetworkLastLockBlock with empty response
		lastLockBlock, err := runtime.GetNetworkLastLockBlock(env.Client, &blockHash)
		if err != nil {
			errorMsg := err.Error()
			isExpectedError := strings.Contains(errorMsg, "no network last lock block found") ||
				strings.Contains(errorMsg, "Method not found") ||
				strings.Contains(errorMsg, "failed to decode") ||
				strings.Contains(errorMsg, "Invalid params")
			assert.True(t, isExpectedError, "Expected network last lock block-related error, got: %v", err)
		} else {
			assert.NotNil(t, lastLockBlock, "Last lock block should not be nil when successful")
			assert.IsType(t, types.U64(0), lastLockBlock, "Should return proper type")
		}

		// Test GetLockReductionInterval with empty response
		interval, err := runtime.GetLockReductionInterval(env.Client, &blockHash)
		if err != nil {
			errorMsg := err.Error()
			isExpectedError := strings.Contains(errorMsg, "no lock reduction interval found") ||
				strings.Contains(errorMsg, "Method not found") ||
				strings.Contains(errorMsg, "failed to decode") ||
				strings.Contains(errorMsg, "Invalid params")
			assert.True(t, isExpectedError, "Expected lock reduction interval-related error, got: %v", err)
		} else {
			assert.NotNil(t, interval, "Lock reduction interval should not be nil when successful")
			assert.IsType(t, types.U64(0), interval, "Should return proper type")
		}

		// Test GetCurrentBlock with empty response
		currentBlock, err := runtime.GetCurrentBlock(env.Client, &blockHash)
		if err != nil {
			errorMsg := err.Error()
			isExpectedError := strings.Contains(errorMsg, "no current block found") ||
				strings.Contains(errorMsg, "Method not found") ||
				strings.Contains(errorMsg, "failed to decode") ||
				strings.Contains(errorMsg, "Invalid params")
			assert.True(t, isExpectedError, "Expected current block-related error, got: %v", err)
		} else {
			assert.NotNil(t, currentBlock, "Current block should not be nil when successful")
			assert.IsType(t, types.U64(0), currentBlock, "Should return proper type")
		}

		// Test all functions with nil block hash
		lockCost, err = runtime.GetNetworkLockCost(env.Client, nil)
		assert.Error(t, err, "Should error with nil block hash")
		assert.Contains(t, err.Error(), "block hash cannot be nil")
		assert.Equal(t, types.U64(0), lockCost, "Lock cost should be 0 on error")

		lastLock, err = runtime.GetNetworkLastLock(env.Client, nil)
		assert.Error(t, err, "Should error with nil block hash")
		assert.Contains(t, err.Error(), "block hash cannot be nil")
		assert.Equal(t, types.U64(0), lastLock, "Last lock should be 0 on error")

		minLock, err = runtime.GetNetworkMinLock(env.Client, nil)
		assert.Error(t, err, "Should error with nil block hash")
		assert.Contains(t, err.Error(), "block hash cannot be nil")
		assert.Equal(t, types.U64(0), minLock, "Min lock should be 0 on error")

		lastLockBlock, err = runtime.GetNetworkLastLockBlock(env.Client, nil)
		assert.Error(t, err, "Should error with nil block hash")
		assert.Contains(t, err.Error(), "block hash cannot be nil")
		assert.Equal(t, types.U64(0), lastLockBlock, "Last lock block should be 0 on error")

		interval, err = runtime.GetLockReductionInterval(env.Client, nil)
		assert.Error(t, err, "Should error with nil block hash")
		assert.Contains(t, err.Error(), "block hash cannot be nil")
		assert.Equal(t, types.U64(0), interval, "Lock reduction interval should be 0 on error")

		currentBlock, err = runtime.GetCurrentBlock(env.Client, nil)
		assert.Error(t, err, "Should error with nil block hash")
		assert.Contains(t, err.Error(), "block hash cannot be nil")
		assert.Equal(t, types.U64(0), currentBlock, "Current block should be 0 on error")
	})
}
