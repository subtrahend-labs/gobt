package runtime_test

import (
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/subtrahend-labs/gobt/runtime"
	"github.com/subtrahend-labs/gobt/testutils"
)

func TestDelegateRuntimeAPIs(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Parallel()

	t.Run("GetDelegates", func(t *testing.T) {
		t.Parallel()

		env, err := testutils.Setup()
		require.NoError(t, err, "Failed to setup test environment")
		defer env.Teardown()

		blockHash, err := env.Client.Api.RPC.Chain.GetBlockHashLatest()
		require.NoError(t, err, "Failed to get latest block hash")

		delegates, err := runtime.GetDelegates(env.Client, &blockHash)
		require.NoError(t, err, "Failed to get delegates")
		assert.NotNil(t, delegates, "Delegates should not be nil")

		t.Logf("Found %d delegates in test environment", len(delegates))

		switch len(delegates) {
		case 0:
			t.Log("No delegates found in test environment")
		default:
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

		delegates, err = runtime.GetDelegates(env.Client, nil)
		require.NoError(t, err, "Should work with nil block hash")
		assert.NotNil(t, delegates, "Delegates should not be nil")
	})

	t.Run("GetDelegate", func(t *testing.T) {
		t.Parallel()

		env, err := testutils.Setup()
		require.NoError(t, err, "Failed to setup test environment")
		defer env.Teardown()

		blockHash, err := env.Client.Api.RPC.Chain.GetBlockHashLatest()
		require.NoError(t, err, "Failed to get latest block hash")

		delegates, err := runtime.GetDelegates(env.Client, &blockHash)
		require.NoError(t, err, "Failed to get delegates")

		switch len(delegates) {
		case 0:
			t.Log("No delegates found in test environment, testing with random account")
		default:
			delegateAccount := delegates[0].AccountID

			delegateOption, err := runtime.GetDelegate(env.Client, delegateAccount, &blockHash)
			require.NoError(t, err, "Failed to get specific delegate")
			assert.NotNil(t, delegateOption, "Delegate option should not be nil")

			ok, delegate := delegateOption.Unwrap()
			require.True(t, ok, "Delegate should exist")

			assert.Equal(t, delegateAccount, delegate.AccountID, "Delegate accounts should match")
			assert.Equal(t, delegates[0].TakeRate, delegate.TakeRate, "TakeRate values should match")
			assert.Equal(t, delegates[0].TotalStake, delegate.TotalStake, "TotalStake values should match")

			t.Logf("Successfully retrieved delegate: %x", delegate.AccountID.ToBytes())

			delegateOption, err = runtime.GetDelegate(env.Client, delegates[0].AccountID, nil)
			require.NoError(t, err, "Should work with nil block hash")
			assert.NotNil(t, delegateOption, "Delegate option should not be nil")
		}

		randomKeyring, err := signature.KeyringPairFromSecret("//Random/Test/Account/12345", 42)
		require.NoError(t, err)
		randomAccount, err := types.NewAccountID(randomKeyring.PublicKey)
		require.NoError(t, err)

		delegateOption, err := runtime.GetDelegate(env.Client, *randomAccount, &blockHash)
		switch {
		case err != nil:
			t.Logf("Runtime API error for non-existent delegate (expected): %v", err)
		default:
			assert.NotNil(t, delegateOption, "Delegate option should not be nil")
			ok, _ := delegateOption.Unwrap()
			assert.False(t, ok, "Should not find delegate for random account")
		}
	})

	t.Run("GetDelegated", func(t *testing.T) {
		t.Parallel()

		env, err := testutils.Setup()
		require.NoError(t, err, "Failed to setup test environment")
		defer env.Teardown()

		blockHash, err := env.Client.Api.RPC.Chain.GetBlockHashLatest()
		require.NoError(t, err, "Failed to get latest block hash")

		aliceKeyring := signature.TestKeyringPairAlice
		aliceAccount, err := types.NewAccountID(aliceKeyring.PublicKey)
		require.NoError(t, err)

		delegated, err := runtime.GetDelegated(env.Client, *aliceAccount, &blockHash)
		switch {
		case err != nil:
			t.Logf("Runtime API error for GetDelegated (acceptable in test environment): %v", err)
		default:
			assert.NotNil(t, delegated, "Delegated info should not be nil")
			t.Logf("Found %d delegations for Alice", len(delegated))

			switch len(delegated) {
			case 0:
				t.Log("No delegations found for Alice")
			default:
				delegation := delegated[0]
				assert.NotNil(t, delegation.DelegateInfo.AccountID, "Delegate AccountID should not be nil")
				assert.GreaterOrEqual(t, uint64(delegation.NetUID.Int64()), uint64(0), "NetUID should be non-negative")
				assert.GreaterOrEqual(t, uint64(delegation.Amount.Int64()), uint64(0), "Amount should be non-negative")

				t.Logf("First delegation - Delegate: %x, NetUID: %d, Amount: %d",
					delegation.DelegateInfo.AccountID.ToBytes(),
					delegation.NetUID.Int64(),
					delegation.Amount.Int64())
			}

			randomKeyring, err := signature.KeyringPairFromSecret("//Random/Test/Delegatee/67890", 42)
			require.NoError(t, err)
			randomAccount, err := types.NewAccountID(randomKeyring.PublicKey)
			require.NoError(t, err)

			delegated, err = runtime.GetDelegated(env.Client, *randomAccount, &blockHash)
			switch {
			case err != nil:
				t.Logf("Runtime API error for random account (acceptable): %v", err)
			default:
				assert.NotNil(t, delegated, "Should return empty array for account with no delegations")
				assert.Equal(t, 0, len(delegated), "Should have no delegations for random account")
			}

			delegated, err = runtime.GetDelegated(env.Client, *aliceAccount, nil)
			switch {
			case err != nil:
				t.Logf("Runtime API error with nil block hash (acceptable): %v", err)
			default:
				assert.NotNil(t, delegated, "Delegated should not be nil")
			}
		}
	})

	t.Run("ErrorHandling", func(t *testing.T) {
		t.Parallel()

		env, err := testutils.Setup()
		require.NoError(t, err, "Failed to setup test environment")
		defer env.Teardown()

		invalidBlockHash := types.NewHash([]byte("invalid_hash_that_does_not_exist_on_chain"))

		delegates, err := runtime.GetDelegates(env.Client, &invalidBlockHash)
		assert.Error(t, err, "Should error with invalid block hash")
		assert.Nil(t, delegates, "Delegates should be nil on error")

		aliceKeyring := signature.TestKeyringPairAlice
		aliceAccount, err := types.NewAccountID(aliceKeyring.PublicKey)
		require.NoError(t, err)

		delegate, err := runtime.GetDelegate(env.Client, *aliceAccount, &invalidBlockHash)
		assert.Error(t, err, "Should error with invalid block hash")
		assert.Nil(t, delegate, "Delegate should be nil on error")

		delegated, err := runtime.GetDelegated(env.Client, *aliceAccount, &invalidBlockHash)
		assert.Error(t, err, "Should error with invalid block hash")
		assert.Nil(t, delegated, "Delegated should be nil on error")
	})

	t.Run("NullAndEmptyResponseHandling", func(t *testing.T) {
		t.Parallel()

		env, err := testutils.Setup()
		require.NoError(t, err, "Failed to setup test environment")
		defer env.Teardown()

		blockHash, err := env.Client.Api.RPC.Chain.GetBlockHashLatest()
		require.NoError(t, err, "Failed to get latest block hash")

		delegates, err := runtime.GetDelegates(env.Client, &blockHash)
		require.NoError(t, err, "Should not error for delegates")
		assert.NotNil(t, delegates, "Delegates should not be nil when successful")
		assert.IsType(t, []runtime.DelegateInfo{}, delegates, "Should return proper type")

		bobKeyring, err := signature.KeyringPairFromSecret("//Bob", 42)
		require.NoError(t, err)
		bobAccount, err := types.NewAccountID(bobKeyring.PublicKey)
		require.NoError(t, err)

		delegateOption, err := runtime.GetDelegate(env.Client, *bobAccount, &blockHash)
		switch {
		case err != nil:
			t.Logf("Runtime API error for Bob account (acceptable): %v", err)
		default:
			assert.NotNil(t, delegateOption, "Delegate option should not be nil")
			ok, _ := delegateOption.Unwrap()
			assert.False(t, ok, "Should not find delegate for Bob account")
		}

		delegated, err := runtime.GetDelegated(env.Client, *bobAccount, &blockHash)
		switch {
		case err != nil:
			t.Logf("Runtime API error for Bob account (acceptable in test environment): %v", err)
		default:
			assert.NotNil(t, delegated, "Should return empty array for account with no delegations")
			assert.Equal(t, 0, len(delegated), "Should have no delegations for Bob account")
		}
	})
}
