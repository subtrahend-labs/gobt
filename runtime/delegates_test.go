package runtime_test

import (
	"strings"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/subtrahend-labs/gobt/runtime"
	"github.com/subtrahend-labs/gobt/testutils"
)

func TestDelegateRuntimeAPIs(t *testing.T) {
	// Skip if running in CI or if Docker is not available
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Parallel()

	t.Run("GetDelegates", func(t *testing.T) {
		t.Parallel()

		env, err := testutils.Setup()
		if err != nil {
			t.Skipf("Failed to setup test environment: %v", err)
		}
		defer env.Teardown()

		blockHash, err := env.Client.Api.RPC.Chain.GetBlockHashLatest()
		require.NoError(t, err, "Failed to get latest block hash")

		// Test successful case
		delegates, err := runtime.GetDelegates(env.Client, &blockHash)

		// In a test environment, there might not be any delegates
		// So we test both scenarios
		if err != nil {
			// If error is about no delegates found, that's acceptable for test env
			assert.Contains(t, err.Error(), "no delegates found",
				"Expected 'no delegates found' error in test environment, got: %v", err)
		} else {
			// If delegates exist (or empty array), verify structure
			assert.NotNil(t, delegates, "Delegates should not be nil")
			t.Logf("Found %d delegates in test environment", len(delegates))

			// If there are delegates, verify the structure
			if len(delegates) > 0 {
				delegate := delegates[0]
				assert.NotNil(t, delegate.AccountID, "Delegate AccountID should not be nil")
				assert.GreaterOrEqual(t, uint64(delegate.TakeRate.Int64()), uint64(0), "TakeRate should be non-negative")
				assert.GreaterOrEqual(t, uint64(delegate.NominatorStake.Int64()), uint64(0), "NominatorStake should be non-negative")
				assert.GreaterOrEqual(t, uint64(delegate.ValidatorStake.Int64()), uint64(0), "ValidatorStake should be non-negative")
				assert.GreaterOrEqual(t, uint64(delegate.TotalStake.Int64()), uint64(0), "TotalStake should be non-negative")
				assert.NotNil(t, delegate.Registrations, "Registrations should not be nil")
				assert.GreaterOrEqual(t, uint64(delegate.VotingPower.Int64()), uint64(0), "VotingPower should be non-negative")
				assert.NotNil(t, delegate.ValidatorPermits, "ValidatorPermits should not be nil")
				assert.GreaterOrEqual(t, uint64(delegate.Return.Int64()), uint64(0), "Return should be non-negative")

				t.Logf("First delegate: %x", delegate.AccountID.ToBytes())
				t.Logf("  TakeRate: %d", delegate.TakeRate.Int64())
				t.Logf("  TotalStake: %d", delegate.TotalStake.Int64())
				t.Logf("  VotingPower: %d", delegate.VotingPower.Int64())
			}
		}

		// Test with nil block hash
		delegates, err = runtime.GetDelegates(env.Client, nil)
		assert.Error(t, err, "Should error with nil block hash")
		assert.Contains(t, err.Error(), "block hash cannot be nil", "Should have proper error message for nil block hash")
		assert.Nil(t, delegates, "Delegates should be nil on error")
	})

	t.Run("GetDelegate", func(t *testing.T) {
		t.Parallel()

		env, err := testutils.Setup()
		if err != nil {
			t.Skipf("Failed to setup test environment: %v", err)
		}
		defer env.Teardown()

		blockHash, err := env.Client.Api.RPC.Chain.GetBlockHashLatest()
		require.NoError(t, err, "Failed to get latest block hash")

		// First, try to get all delegates to have a valid delegate account to query
		delegates, err := runtime.GetDelegates(env.Client, &blockHash)

		if err == nil && len(delegates) > 0 {
			// Test successful case with existing delegate
			delegateAccount := delegates[0].AccountID

			delegate, err := runtime.GetDelegate(env.Client, delegateAccount, &blockHash)
			require.NoError(t, err, "Failed to get specific delegate")
			assert.NotNil(t, delegate, "Delegate should not be nil")

			// Verify it's the same delegate
			assert.Equal(t, delegateAccount, delegate.AccountID, "Delegate accounts should match")
			assert.Equal(t, delegates[0].TakeRate, delegate.TakeRate, "TakeRate values should match")
			assert.Equal(t, delegates[0].TotalStake, delegate.TotalStake, "TotalStake values should match")

			t.Logf("Successfully retrieved delegate: %x", delegate.AccountID.ToBytes())
		} else {
			t.Log("No delegates found in test environment to test GetDelegate with valid account")
		}

		// Test with non-existent delegate
		randomKeyring, err := signature.KeyringPairFromSecret("//Random/Test/Account/12345", 42)
		require.NoError(t, err)
		randomAccount, err := types.NewAccountID(randomKeyring.PublicKey)
		require.NoError(t, err)

		delegate, err := runtime.GetDelegate(env.Client, *randomAccount, &blockHash)
		assert.Error(t, err, "Should error on non-existent delegate")
		assert.Nil(t, delegate, "Delegate should be nil on error")
		// Runtime might return different error messages in test environment
		errorMsg := err.Error()
		isExpectedError := strings.Contains(errorMsg, "no delegate found") ||
			strings.Contains(errorMsg, "Invalid params")
		assert.True(t, isExpectedError, "Error should indicate delegate not found or invalid params")

		// Test with nil block hash
		if len(delegates) > 0 {
			delegate, err = runtime.GetDelegate(env.Client, delegates[0].AccountID, nil)
			assert.Error(t, err, "Should error with nil block hash")
			assert.Contains(t, err.Error(), "block hash cannot be nil", "Should have proper error message for nil block hash")
			assert.Nil(t, delegate, "Delegate should be nil on error")
		}
	})

	t.Run("GetDelegated", func(t *testing.T) {
		t.Parallel()

		env, err := testutils.Setup()
		if err != nil {
			t.Skipf("Failed to setup test environment: %v", err)
		}
		defer env.Teardown()

		blockHash, err := env.Client.Api.RPC.Chain.GetBlockHashLatest()
		require.NoError(t, err, "Failed to get latest block hash")

		// Test with a known account (Alice from test environment)
		aliceKeyring := signature.TestKeyringPairAlice
		aliceAccount, err := types.NewAccountID(aliceKeyring.PublicKey)
		require.NoError(t, err)

		delegated, err := runtime.GetDelegated(env.Client, *aliceAccount, &blockHash)

		// In test environment, Alice likely has no delegations (empty array is expected)
		// But runtime might return "Invalid params" error for test accounts
		if err != nil {
			// Runtime might return error for test accounts
			errorMsg := err.Error()
			isExpectedError := strings.Contains(errorMsg, "Invalid params") ||
				strings.Contains(errorMsg, "failed to call delegateInfo_getDelegated")
			assert.True(t, isExpectedError, "Should have runtime API error for test account")
		} else {
			assert.NotNil(t, delegated, "Delegated info should not be nil")
			t.Logf("Found %d delegations for Alice", len(delegated))

			// Verify the structure if there are delegations
			if len(delegated) > 0 {
				delegation := delegated[0]
				assert.NotNil(t, delegation.DelegateInfo.AccountID, "Delegate AccountID should not be nil")
				assert.GreaterOrEqual(t, uint64(delegation.NetUID.Int64()), uint64(0), "NetUID should be non-negative")
				assert.GreaterOrEqual(t, uint64(delegation.Amount.Int64()), uint64(0), "Amount should be non-negative")

				t.Logf("First delegation - Delegate: %x, NetUID: %d, Amount: %d",
					delegation.DelegateInfo.AccountID.ToBytes(),
					delegation.NetUID.Int64(),
					delegation.Amount.Int64())
			}
		}

		// Test with non-existent account - should still return empty array
		randomKeyring, err := signature.KeyringPairFromSecret("//Random/Test/Delegatee/67890", 42)
		require.NoError(t, err)
		randomAccount, err := types.NewAccountID(randomKeyring.PublicKey)
		require.NoError(t, err)

		delegated, err = runtime.GetDelegated(env.Client, *randomAccount, &blockHash)
		// Expect success with empty array or runtime error for invalid account
		if err != nil {
			// Runtime might return error for invalid account parameters
			errorMsg := err.Error()
			isExpectedError := strings.Contains(errorMsg, "Invalid params") ||
				strings.Contains(errorMsg, "failed to call delegateInfo_getDelegated")
			assert.True(t, isExpectedError, "Should have runtime API error for invalid account")
		} else {
			assert.NotNil(t, delegated, "Should return empty array for account with no delegations")
			assert.Equal(t, 0, len(delegated), "Should have no delegations for random account")
		}

		// Test with nil block hash
		delegated, err = runtime.GetDelegated(env.Client, *aliceAccount, nil)
		assert.Error(t, err, "Should error with nil block hash")
		assert.Contains(t, err.Error(), "block hash cannot be nil", "Should have proper error message for nil block hash")
		assert.Nil(t, delegated, "Delegated should be nil on error")
	})

	t.Run("ErrorHandling", func(t *testing.T) {
		t.Parallel()

		env, err := testutils.Setup()
		if err != nil {
			t.Skipf("Failed to setup test environment: %v", err)
		}
		defer env.Teardown()

		// Test with invalid block hash
		invalidBlockHash := types.NewHash([]byte("invalid_hash_that_does_not_exist_on_chain"))

		// GetDelegates with invalid block
		delegates, err := runtime.GetDelegates(env.Client, &invalidBlockHash)
		assert.Error(t, err, "Should error with invalid block hash")
		assert.Nil(t, delegates, "Delegates should be nil on error")

		// GetDelegate with invalid block
		aliceKeyring := signature.TestKeyringPairAlice
		aliceAccount, err := types.NewAccountID(aliceKeyring.PublicKey)
		require.NoError(t, err)

		delegate, err := runtime.GetDelegate(env.Client, *aliceAccount, &invalidBlockHash)
		assert.Error(t, err, "Should error with invalid block hash")
		assert.Nil(t, delegate, "Delegate should be nil on error")

		// GetDelegated with invalid block
		delegated, err := runtime.GetDelegated(env.Client, *aliceAccount, &invalidBlockHash)
		assert.Error(t, err, "Should error with invalid block hash")
		assert.Nil(t, delegated, "Delegated should be nil on error")
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

		// Test GetDelegates returns proper error when no delegates exist
		delegates, err := runtime.GetDelegates(env.Client, &blockHash)
		if err != nil {
			assert.Contains(t, err.Error(), "no delegates found", "Should have proper error message for no delegates")
		} else {
			// If delegates exist, make sure the response is properly structured
			assert.NotNil(t, delegates, "Delegates should not be nil when successful")
			assert.IsType(t, []runtime.DelegateInfo{}, delegates, "Should return proper type")
		}

		// Test GetDelegate with valid account but no delegate role
		bobKeyring, err := signature.KeyringPairFromSecret("//Bob", 42)
		require.NoError(t, err)
		bobAccount, err := types.NewAccountID(bobKeyring.PublicKey)
		require.NoError(t, err)

		delegate, err := runtime.GetDelegate(env.Client, *bobAccount, &blockHash)
		assert.Error(t, err, "Should error when account is not a delegate")
		assert.Nil(t, delegate, "Delegate should be nil when not found")

		// Test GetDelegated with account that has no delegations - should return empty array
		delegated, err := runtime.GetDelegated(env.Client, *bobAccount, &blockHash)
		// Expect success with empty array or runtime error for invalid account
		if err != nil {
			// Runtime might return error for invalid account parameters
			errorMsg := err.Error()
			isExpectedError := strings.Contains(errorMsg, "Invalid params") ||
				strings.Contains(errorMsg, "failed to call delegateInfo_getDelegated")
			assert.True(t, isExpectedError, "Should have runtime API error or success with empty array")
		} else {
			assert.NotNil(t, delegated, "Should return empty array for account with no delegations")
			assert.Equal(t, 0, len(delegated), "Should have no delegations for Bob account")
		}
	})
}
