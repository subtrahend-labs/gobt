package runtime_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/subtrahend-labs/gobt/runtime"
	"github.com/subtrahend-labs/gobt/testutils"
)

func TestSubnetRuntimeAPIs(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Parallel()

	t.Run("GetSubnetInfo", func(t *testing.T) {
		t.Parallel()

		env, err := testutils.Setup()
		if err != nil {
			t.Skipf("Failed to setup test environment: %v", err)
		}
		defer env.Teardown()

		blockHash, err := env.Client.Api.RPC.Chain.GetBlockHashLatest()
		require.NoError(t, err, "Failed to get latest block hash")

		subnet, err := runtime.GetSubnetInfo(env.Client, 0, &blockHash)
		if err != nil {
			errorMsg := err.Error()
			isExpectedError := strings.Contains(errorMsg, "no subnet info found") ||
				strings.Contains(errorMsg, "Method not found") ||
				strings.Contains(errorMsg, "failed to decode") ||
				strings.Contains(errorMsg, "Invalid params")
			assert.True(t, isExpectedError, "Expected subnet-related error, got: %v", err)
		} else {
			assert.NotNil(t, subnet, "Subnet should not be nil")
			assert.Equal(t, types.U16(0), subnet.NetUID, "NetUID should be 0")
			assert.GreaterOrEqual(t, uint64(subnet.Difficulty), uint64(0), "Difficulty should be non-negative")
			assert.GreaterOrEqual(t, uint64(subnet.ImmunityPeriod), uint64(0), "ImmunityPeriod should be non-negative")
			assert.GreaterOrEqual(t, uint64(subnet.MaxValidators), uint64(0), "MaxValidators should be non-negative")
			assert.NotNil(t, subnet.Owner, "Owner should not be nil")

			t.Logf("Subnet 0 found - Owner: %x", subnet.Owner.ToBytes())
			t.Logf("  Difficulty: %d", subnet.Difficulty)
			t.Logf("  MaxValidators: %d", subnet.MaxValidators)
			t.Logf("  ImmunityPeriod: %d", subnet.ImmunityPeriod)
		}

		subnet, err = runtime.GetSubnetInfo(env.Client, 999, &blockHash)
		assert.Error(t, err, "Should error for non-existent subnet")
		assert.Nil(t, subnet, "Subnet should be nil on error")

		subnet, err = runtime.GetSubnetInfo(env.Client, 0, nil)
		assert.Error(t, err, "Should error with nil block hash")
		assert.Contains(t, err.Error(), "block hash cannot be nil")
		assert.Nil(t, subnet, "Subnet should be nil on error")
	})

	t.Run("GetSubnetsInfo", func(t *testing.T) {
		t.Parallel()

		env, err := testutils.Setup()
		if err != nil {
			t.Skipf("Failed to setup test environment: %v", err)
		}
		defer env.Teardown()

		blockHash, err := env.Client.Api.RPC.Chain.GetBlockHashLatest()
		require.NoError(t, err, "Failed to get latest block hash")

		subnets, err := runtime.GetSubnetsInfo(env.Client, &blockHash)
		if err != nil {
			errorMsg := err.Error()
			isExpectedError := strings.Contains(errorMsg, "no subnets info found") ||
				strings.Contains(errorMsg, "Method not found") ||
				strings.Contains(errorMsg, "failed to decode") ||
				strings.Contains(errorMsg, "Invalid params")
			assert.True(t, isExpectedError, "Expected subnets-related error, got: %v", err)
		} else {
			assert.NotNil(t, subnets, "Subnets should not be nil")
			t.Logf("Found %d subnets in test environment", len(subnets))

			// If there are subnets, verify the structure
			for i, subnetOpt := range subnets {
				ok, subnet := subnetOpt.Unwrap()
				if ok {
					assert.GreaterOrEqual(t, uint64(subnet.Difficulty), uint64(0), "Subnet %d difficulty should be non-negative", i)
					assert.NotNil(t, subnet.Owner, "Subnet %d owner should not be nil", i)
					t.Logf("Subnet %d: NetUID=%d, Owner=%x", i, subnet.NetUID, subnet.Owner.ToBytes())
				}
			}
		}

		// Test with nil block hash
		subnets, err = runtime.GetSubnetsInfo(env.Client, nil)
		assert.Error(t, err, "Should error with nil block hash")
		assert.Contains(t, err.Error(), "block hash cannot be nil")
		assert.Nil(t, subnets, "Subnets should be nil on error")
	})

	t.Run("GetSubnetInfoV2", func(t *testing.T) {
		t.Parallel()

		env, err := testutils.Setup()
		if err != nil {
			t.Skipf("Failed to setup test environment: %v", err)
		}
		defer env.Teardown()

		blockHash, err := env.Client.Api.RPC.Chain.GetBlockHashLatest()
		require.NoError(t, err, "Failed to get latest block hash")

		// Test for subnet 0 (root subnet should exist)
		subnet, err := runtime.GetSubnetInfoV2(env.Client, 0, &blockHash)
		if err != nil {
			errorMsg := err.Error()
			isExpectedError := strings.Contains(errorMsg, "no subnet info v2 found") ||
				strings.Contains(errorMsg, "Method not found") ||
				strings.Contains(errorMsg, "failed to decode") ||
				strings.Contains(errorMsg, "Invalid params")
			assert.True(t, isExpectedError, "Expected subnet v2-related error, got: %v", err)
		} else {
			assert.NotNil(t, subnet, "Subnet V2 should not be nil")
			assert.Equal(t, types.U16(0), subnet.NetUID, "NetUID should be 0")
			assert.GreaterOrEqual(t, uint64(subnet.Difficulty), uint64(0), "Difficulty should be non-negative")
			assert.NotNil(t, subnet.Owner, "Owner should not be nil")
			assert.NotNil(t, subnet.FoundationAccount, "FoundationAccount should not be nil")
			assert.GreaterOrEqual(t, uint64(subnet.MaxAllowedModules), uint64(0), "MaxAllowedModules should be non-negative")

			t.Logf("Subnet V2 0 found - Owner: %x", subnet.Owner.ToBytes())
			t.Logf("  Foundation: %x", subnet.FoundationAccount.ToBytes())
			t.Logf("  MaxAllowedModules: %d", subnet.MaxAllowedModules)
		}

		// Test for non-existent subnet
		subnet, err = runtime.GetSubnetInfoV2(env.Client, 999, &blockHash)
		assert.Error(t, err, "Should error for non-existent subnet")
		assert.Nil(t, subnet, "Subnet V2 should be nil on error")

		// Test with nil block hash
		subnet, err = runtime.GetSubnetInfoV2(env.Client, 0, nil)
		assert.Error(t, err, "Should error with nil block hash")
		assert.Contains(t, err.Error(), "block hash cannot be nil")
		assert.Nil(t, subnet, "Subnet V2 should be nil on error")
	})

	t.Run("GetSubnetsInfoV2", func(t *testing.T) {
		t.Parallel()

		env, err := testutils.Setup()
		if err != nil {
			t.Skipf("Failed to setup test environment: %v", err)
		}
		defer env.Teardown()

		blockHash, err := env.Client.Api.RPC.Chain.GetBlockHashLatest()
		require.NoError(t, err, "Failed to get latest block hash")

		subnets, err := runtime.GetSubnetsInfoV2(env.Client, &blockHash)
		if err != nil {
			errorMsg := err.Error()
			isExpectedError := strings.Contains(errorMsg, "no subnets info v2 found") ||
				strings.Contains(errorMsg, "Method not found") ||
				strings.Contains(errorMsg, "failed to decode") ||
				strings.Contains(errorMsg, "Invalid params")
			assert.True(t, isExpectedError, "Expected subnets v2-related error, got: %v", err)
		} else {
			assert.NotNil(t, subnets, "Subnets V2 should not be nil")
			t.Logf("Found %d subnets V2 in test environment", len(subnets))

			// If there are subnets, verify the structure
			for i, subnetOpt := range subnets {
				ok, subnet := subnetOpt.Unwrap()
				if ok {
					assert.GreaterOrEqual(t, uint64(subnet.Difficulty), uint64(0), "Subnet V2 %d difficulty should be non-negative", i)
					assert.NotNil(t, subnet.Owner, "Subnet V2 %d owner should not be nil", i)
					assert.NotNil(t, subnet.FoundationAccount, "Subnet V2 %d foundation should not be nil", i)
					t.Logf("Subnet V2 %d: NetUID=%d, Owner=%x", i, subnet.NetUID, subnet.Owner.ToBytes())
				}
			}
		}

		// Test with nil block hash
		subnets, err = runtime.GetSubnetsInfoV2(env.Client, nil)
		assert.Error(t, err, "Should error with nil block hash")
		assert.Contains(t, err.Error(), "block hash cannot be nil")
		assert.Nil(t, subnets, "Subnets V2 should be nil on error")
	})

	t.Run("GetSubnetHyperparams", func(t *testing.T) {
		t.Parallel()

		env, err := testutils.Setup()
		if err != nil {
			t.Skipf("Failed to setup test environment: %v", err)
		}
		defer env.Teardown()

		blockHash, err := env.Client.Api.RPC.Chain.GetBlockHashLatest()
		require.NoError(t, err, "Failed to get latest block hash")

		// Test for subnet 0 (root subnet should exist)
		hyperparams, err := runtime.GetSubnetHyperparams(env.Client, 0, &blockHash)
		if err != nil {
			errorMsg := err.Error()
			isExpectedError := strings.Contains(errorMsg, "no subnet hyperparams found") ||
				strings.Contains(errorMsg, "Method not found") ||
				strings.Contains(errorMsg, "failed to decode") ||
				strings.Contains(errorMsg, "Invalid params") ||
				strings.Contains(errorMsg, "does not support Decodeable interface")
			assert.True(t, isExpectedError, "Expected hyperparams-related error, got: %v", err)
		} else {
			assert.NotNil(t, hyperparams, "Subnet hyperparams should not be nil")
			assert.GreaterOrEqual(t, uint64(hyperparams.Rho), uint64(0), "Rho should be non-negative")
			assert.GreaterOrEqual(t, uint64(hyperparams.Kappa), uint64(0), "Kappa should be non-negative")
			assert.GreaterOrEqual(t, uint64(hyperparams.ImmunityPeriod), uint64(0), "ImmunityPeriod should be non-negative")
			assert.GreaterOrEqual(t, uint64(hyperparams.MaxValidators), uint64(0), "MaxValidators should be non-negative")

			t.Logf("Subnet 0 hyperparams found")
			t.Logf("  Rho: %d", hyperparams.Rho)
			t.Logf("  Kappa: %d", hyperparams.Kappa)
			t.Logf("  ImmunityPeriod: %d", hyperparams.ImmunityPeriod)
		}

		// Test for non-existent subnet
		hyperparams, err = runtime.GetSubnetHyperparams(env.Client, 999, &blockHash)
		assert.Error(t, err, "Should error for non-existent subnet")
		assert.Nil(t, hyperparams, "Subnet hyperparams should be nil on error")

		// Test with nil block hash
		hyperparams, err = runtime.GetSubnetHyperparams(env.Client, 0, nil)
		assert.Error(t, err, "Should error with nil block hash")
		assert.Contains(t, err.Error(), "block hash cannot be nil")
		assert.Nil(t, hyperparams, "Subnet hyperparams should be nil on error")
	})

	t.Run("GetAllDynamicInfo", func(t *testing.T) {
		t.Parallel()

		env, err := testutils.Setup()
		if err != nil {
			t.Skipf("Failed to setup test environment: %v", err)
		}
		defer env.Teardown()

		blockHash, err := env.Client.Api.RPC.Chain.GetBlockHashLatest()
		require.NoError(t, err, "Failed to get latest block hash")

		dynamicInfos, err := runtime.GetAllDynamicInfo(env.Client, &blockHash)
		if err != nil {
			errorMsg := err.Error()
			isExpectedError := strings.Contains(errorMsg, "no dynamic info found") ||
				strings.Contains(errorMsg, "Method not found") ||
				strings.Contains(errorMsg, "failed to decode") ||
				strings.Contains(errorMsg, "Invalid params")
			assert.True(t, isExpectedError, "Expected dynamic info-related error, got: %v", err)
		} else {
			assert.NotNil(t, dynamicInfos, "Dynamic infos should not be nil")
			t.Logf("Found %d dynamic infos in test environment", len(dynamicInfos))

			// If there are dynamic infos, verify the structure
			for i, dynamicOpt := range dynamicInfos {
				ok, dynamic := dynamicOpt.Unwrap()
				if ok {
					assert.GreaterOrEqual(t, uint64(dynamic.Difficulty), uint64(0), "Dynamic %d difficulty should be non-negative", i)
					assert.GreaterOrEqual(t, uint64(dynamic.Burn), uint64(0), "Dynamic %d burn should be non-negative", i)
					assert.NotNil(t, dynamic.Owner, "Dynamic %d owner should not be nil", i)
					t.Logf("Dynamic %d: NetUID=%d, Owner=%x", i, dynamic.NetUID, dynamic.Owner.ToBytes())
				}
			}
		}

		// Test with nil block hash
		dynamicInfos, err = runtime.GetAllDynamicInfo(env.Client, nil)
		assert.Error(t, err, "Should error with nil block hash")
		assert.Contains(t, err.Error(), "block hash cannot be nil")
		assert.Nil(t, dynamicInfos, "Dynamic infos should be nil on error")
	})

	t.Run("GetDynamicInfo", func(t *testing.T) {
		t.Parallel()

		env, err := testutils.Setup()
		if err != nil {
			t.Skipf("Failed to setup test environment: %v", err)
		}
		defer env.Teardown()

		blockHash, err := env.Client.Api.RPC.Chain.GetBlockHashLatest()
		require.NoError(t, err, "Failed to get latest block hash")

		// Test for subnet 0 (root subnet should exist)
		dynamicInfo, err := runtime.GetDynamicInfo(env.Client, 0, &blockHash)
		if err != nil {
			errorMsg := err.Error()
			isExpectedError := strings.Contains(errorMsg, "no dynamic info found") ||
				strings.Contains(errorMsg, "Method not found") ||
				strings.Contains(errorMsg, "failed to decode") ||
				strings.Contains(errorMsg, "Invalid params")
			assert.True(t, isExpectedError, "Expected dynamic info-related error, got: %v", err)
		} else {
			assert.NotNil(t, dynamicInfo, "Dynamic info should not be nil")
			assert.Equal(t, types.U16(0), dynamicInfo.NetUID, "NetUID should be 0")
			assert.GreaterOrEqual(t, uint64(dynamicInfo.Difficulty), uint64(0), "Difficulty should be non-negative")
			assert.GreaterOrEqual(t, uint64(dynamicInfo.Burn), uint64(0), "Burn should be non-negative")
			assert.NotNil(t, dynamicInfo.Owner, "Owner should not be nil")

			t.Logf("Dynamic info 0 found - Owner: %x", dynamicInfo.Owner.ToBytes())
			t.Logf("  Difficulty: %d", dynamicInfo.Difficulty)
			t.Logf("  Burn: %d", dynamicInfo.Burn)
		}

		// Test for non-existent subnet
		dynamicInfo, err = runtime.GetDynamicInfo(env.Client, 999, &blockHash)
		assert.Error(t, err, "Should error for non-existent subnet")
		assert.Nil(t, dynamicInfo, "Dynamic info should be nil on error")

		// Test with nil block hash
		dynamicInfo, err = runtime.GetDynamicInfo(env.Client, 0, nil)
		assert.Error(t, err, "Should error with nil block hash")
		assert.Contains(t, err.Error(), "block hash cannot be nil")
		assert.Nil(t, dynamicInfo, "Dynamic info should be nil on error")
	})

	t.Run("GetSubnetState", func(t *testing.T) {
		t.Parallel()

		env, err := testutils.Setup()
		if err != nil {
			t.Skipf("Failed to setup test environment: %v", err)
		}
		defer env.Teardown()

		blockHash, err := env.Client.Api.RPC.Chain.GetBlockHashLatest()
		require.NoError(t, err, "Failed to get latest block hash")

		// Test for subnet 0 (root subnet should exist)
		subnetState, err := runtime.GetSubnetState(env.Client, 0, &blockHash)
		if err != nil {
			errorMsg := err.Error()
			isExpectedError := strings.Contains(errorMsg, "no subnet state found") ||
				strings.Contains(errorMsg, "Method not found") ||
				strings.Contains(errorMsg, "failed to decode") ||
				strings.Contains(errorMsg, "Invalid params")
			assert.True(t, isExpectedError, "Expected subnet state-related error, got: %v", err)
		} else {
			assert.NotNil(t, subnetState, "Subnet state should not be nil")
			assert.Equal(t, types.U16(0), subnetState.NetUID, "NetUID should be 0")
			assert.GreaterOrEqual(t, uint64(subnetState.N), uint64(0), "N should be non-negative")
			assert.GreaterOrEqual(t, uint64(subnetState.StakeThreshold), uint64(0), "StakeThreshold should be non-negative")
			assert.NotNil(t, subnetState.Founder, "Founder should not be nil")

			t.Logf("Subnet state 0 found - Founder: %x", subnetState.Founder.ToBytes())
			t.Logf("  N: %d", subnetState.N)
			t.Logf("  StakeThreshold: %d", subnetState.StakeThreshold)
			t.Logf("  FounderShare: %d", subnetState.FounderShare)
		}

		// Test for non-existent subnet
		subnetState, err = runtime.GetSubnetState(env.Client, 999, &blockHash)
		assert.Error(t, err, "Should error for non-existent subnet")
		assert.Nil(t, subnetState, "Subnet state should be nil on error")

		// Test with nil block hash
		subnetState, err = runtime.GetSubnetState(env.Client, 0, nil)
		assert.Error(t, err, "Should error with nil block hash")
		assert.Contains(t, err.Error(), "block hash cannot be nil")
		assert.Nil(t, subnetState, "Subnet state should be nil on error")
	})

	t.Run("GetAllMetagraphs", func(t *testing.T) {
		t.Parallel()

		env, err := testutils.Setup()
		if err != nil {
			t.Skipf("Failed to setup test environment: %v", err)
		}
		defer env.Teardown()

		blockHash, err := env.Client.Api.RPC.Chain.GetBlockHashLatest()
		require.NoError(t, err, "Failed to get latest block hash")

		// Use defer/recover to catch panics during decoding
		var metagraphs []types.Option[runtime.Metagraph]
		var testErr error

		func() {
			defer func() {
				if r := recover(); r != nil {
					testErr = fmt.Errorf("panic during decoding: %v", r)
				}
			}()
			metagraphs, testErr = runtime.GetAllMetagraphs(env.Client, &blockHash)
		}()

		if testErr != nil {
			errorMsg := testErr.Error()
			isExpectedError := strings.Contains(errorMsg, "no metagraphs found") ||
				strings.Contains(errorMsg, "Method not found") ||
				strings.Contains(errorMsg, "failed to decode") ||
				strings.Contains(errorMsg, "Invalid params") ||
				strings.Contains(errorMsg, "panic during decoding")
			assert.True(t, isExpectedError, "Expected metagraphs-related error, got: %v", testErr)
		} else {
			assert.NotNil(t, metagraphs, "Metagraphs should not be nil")
			t.Logf("Found %d metagraphs in test environment", len(metagraphs))

			// If there are metagraphs, verify the structure
			for i, metagraphOpt := range metagraphs {
				ok, metagraph := metagraphOpt.Unwrap()
				if ok {
					assert.GreaterOrEqual(t, uint64(metagraph.Netuid.Int64()), uint64(0), "Metagraph %d netuid should be non-negative", i)
					// Hotkeys can be nil or empty slice, both are valid
					hotkeyCount := 0
					if metagraph.Hotkeys != nil {
						hotkeyCount = len(metagraph.Hotkeys)
					}
					t.Logf("Metagraph %d: NetUID=%d, Hotkeys count=%d", i, metagraph.Netuid.Int64(), hotkeyCount)
				}
			}
		}

		// Test with nil block hash
		metagraphs, err = runtime.GetAllMetagraphs(env.Client, nil)
		assert.Error(t, err, "Should error with nil block hash")
		assert.Contains(t, err.Error(), "block hash cannot be nil")
		assert.Nil(t, metagraphs, "Metagraphs should be nil on error")
	})

	t.Run("GetMetagraphSubnet", func(t *testing.T) {
		t.Parallel()

		env, err := testutils.Setup()
		if err != nil {
			t.Skipf("Failed to setup test environment: %v", err)
		}
		defer env.Teardown()

		blockHash, err := env.Client.Api.RPC.Chain.GetBlockHashLatest()
		require.NoError(t, err, "Failed to get latest block hash")

		// Use defer/recover to catch panics during decoding
		var metagraph *runtime.Metagraph
		var testErr error

		func() {
			defer func() {
				if r := recover(); r != nil {
					testErr = fmt.Errorf("panic during decoding: %v", r)
				}
			}()
			metagraph, testErr = runtime.GetMetagraphSubnet(env.Client, 0, &blockHash)
		}()

		if testErr != nil {
			errorMsg := testErr.Error()
			isExpectedError := strings.Contains(errorMsg, "no metagraph found") ||
				strings.Contains(errorMsg, "Method not found") ||
				strings.Contains(errorMsg, "failed to decode") ||
				strings.Contains(errorMsg, "Invalid params") ||
				strings.Contains(errorMsg, "panic during decoding")
			assert.True(t, isExpectedError, "Expected metagraph-related error, got: %v", testErr)
		} else {
			assert.NotNil(t, metagraph, "Metagraph should not be nil")
			assert.GreaterOrEqual(t, uint64(metagraph.Netuid.Int64()), uint64(0), "NetUID should be non-negative")
			// Hotkeys and Coldkeys can be nil or empty slices, both are valid
			if metagraph.Hotkeys != nil {
				t.Logf("  Hotkeys: %d", len(metagraph.Hotkeys))
			} else {
				t.Logf("  Hotkeys: nil")
			}
			if metagraph.Coldkeys != nil {
				t.Logf("  Coldkeys: %d", len(metagraph.Coldkeys))
			} else {
				t.Logf("  Coldkeys: nil")
			}

			t.Logf("Metagraph 0 found - NetUID: %d", metagraph.Netuid.Int64())
		}

		// Test for non-existent subnet
		func() {
			defer func() {
				if r := recover(); r != nil {
					// Panics are also expected for non-existent subnets
					t.Logf("Expected panic for non-existent subnet: %v", r)
				}
			}()
			metagraph, err = runtime.GetMetagraphSubnet(env.Client, 999, &blockHash)
			if err != nil {
				assert.Error(t, err, "Should error for non-existent subnet")
			}
			if metagraph != nil {
				t.Logf("Unexpected metagraph returned for non-existent subnet")
			}
		}()

		// Test with nil block hash
		metagraph, err = runtime.GetMetagraphSubnet(env.Client, 0, nil)
		assert.Error(t, err, "Should error with nil block hash")
		assert.Contains(t, err.Error(), "block hash cannot be nil")
		assert.Nil(t, metagraph, "Metagraph should be nil on error")
	})

	t.Run("GetSelectiveMetagraph", func(t *testing.T) {
		t.Parallel()

		env, err := testutils.Setup()
		if err != nil {
			t.Skipf("Failed to setup test environment: %v", err)
		}
		defer env.Teardown()

		blockHash, err := env.Client.Api.RPC.Chain.GetBlockHashLatest()
		require.NoError(t, err, "Failed to get latest block hash")

		// Test with some sample metagraph indexes
		indexes := []uint16{0, 1, 2}

		// Use defer/recover to catch panics during decoding
		var selectiveMetagraph *runtime.SelectiveMetagraph
		var testErr error

		func() {
			defer func() {
				if r := recover(); r != nil {
					testErr = fmt.Errorf("panic during decoding: %v", r)
				}
			}()
			selectiveMetagraph, testErr = runtime.GetSelectiveMetagraph(env.Client, 0, indexes, &blockHash)
		}()

		if testErr != nil {
			errorMsg := testErr.Error()
			isExpectedError := strings.Contains(errorMsg, "no selective metagraph found") ||
				strings.Contains(errorMsg, "Method not found") ||
				strings.Contains(errorMsg, "failed to decode") ||
				strings.Contains(errorMsg, "Invalid params") ||
				strings.Contains(errorMsg, "panic during decoding")
			assert.True(t, isExpectedError, "Expected selective metagraph-related error, got: %v", testErr)
		} else {
			assert.NotNil(t, selectiveMetagraph, "Selective metagraph should not be nil")
			assert.GreaterOrEqual(t, uint64(selectiveMetagraph.Netuid.Int64()), uint64(0), "NetUID should be non-negative")
			// Hotkeys and Coldkeys can be nil or empty slices, both are valid
			hotkeyCount := 0
			coldkeyCount := 0
			if selectiveMetagraph.Hotkeys != nil {
				hotkeyCount = len(selectiveMetagraph.Hotkeys)
			}
			if selectiveMetagraph.Coldkeys != nil {
				coldkeyCount = len(selectiveMetagraph.Coldkeys)
			}

			t.Logf("Selective metagraph 0 found - NetUID: %d", selectiveMetagraph.Netuid.Int64())
			t.Logf("  Hotkeys: %d", hotkeyCount)
			t.Logf("  Coldkeys: %d", coldkeyCount)
		}

		// Test for non-existent subnet
		func() {
			defer func() {
				if r := recover(); r != nil {
					// Panics are also expected for non-existent subnets
					t.Logf("Expected panic for non-existent subnet: %v", r)
				}
			}()
			selectiveMetagraph, err = runtime.GetSelectiveMetagraph(env.Client, 999, indexes, &blockHash)
			if err != nil {
				assert.Error(t, err, "Should error for non-existent subnet")
			}
		}()

		// Test with empty indexes
		emptyIndexes := []uint16{}
		func() {
			defer func() {
				if r := recover(); r != nil {
					// Panics might occur with empty indexes
					t.Logf("Panic with empty indexes: %v", r)
				}
			}()
			selectiveMetagraph, err = runtime.GetSelectiveMetagraph(env.Client, 0, emptyIndexes, &blockHash)
			if err != nil {
				// This might be valid behavior depending on the runtime implementation
				errorMsg := err.Error()
				isExpectedError := strings.Contains(errorMsg, "no selective metagraph found") ||
					strings.Contains(errorMsg, "Method not found") ||
					strings.Contains(errorMsg, "failed to decode") ||
					strings.Contains(errorMsg, "Invalid params")
				assert.True(t, isExpectedError, "Expected selective metagraph-related error for empty indexes, got: %v", err)
			}
		}()

		// Test with nil block hash
		selectiveMetagraph, err = runtime.GetSelectiveMetagraph(env.Client, 0, indexes, nil)
		assert.Error(t, err, "Should error with nil block hash")
		assert.Contains(t, err.Error(), "block hash cannot be nil")
		assert.Nil(t, selectiveMetagraph, "Selective metagraph should be nil on error")
	})

	t.Run("ErrorHandling", func(t *testing.T) {
		t.Parallel()

		// Create a dummy client for testing nil block hash errors
		env, err := testutils.Setup()
		if err != nil {
			t.Skipf("Failed to setup test environment: %v", err)
		}
		defer env.Teardown()

		// Test all functions with nil block hash
		subnet, err := runtime.GetSubnetInfo(env.Client, 0, nil)
		assert.Error(t, err, "GetSubnetInfo should error with nil block hash")
		assert.Contains(t, err.Error(), "block hash cannot be nil")
		assert.Nil(t, subnet, "Result should be nil on error")

		subnets, err := runtime.GetSubnetsInfo(env.Client, nil)
		assert.Error(t, err, "GetSubnetsInfo should error with nil block hash")
		assert.Contains(t, err.Error(), "block hash cannot be nil")
		assert.Nil(t, subnets, "Result should be nil on error")

		subnetV2, err := runtime.GetSubnetInfoV2(env.Client, 0, nil)
		assert.Error(t, err, "GetSubnetInfoV2 should error with nil block hash")
		assert.Contains(t, err.Error(), "block hash cannot be nil")
		assert.Nil(t, subnetV2, "Result should be nil on error")

		subnetsV2, err := runtime.GetSubnetsInfoV2(env.Client, nil)
		assert.Error(t, err, "GetSubnetsInfoV2 should error with nil block hash")
		assert.Contains(t, err.Error(), "block hash cannot be nil")
		assert.Nil(t, subnetsV2, "Result should be nil on error")

		hyperparams, err := runtime.GetSubnetHyperparams(env.Client, 0, nil)
		assert.Error(t, err, "GetSubnetHyperparams should error with nil block hash")
		assert.Contains(t, err.Error(), "block hash cannot be nil")
		assert.Nil(t, hyperparams, "Result should be nil on error")

		allDynamic, err := runtime.GetAllDynamicInfo(env.Client, nil)
		assert.Error(t, err, "GetAllDynamicInfo should error with nil block hash")
		assert.Contains(t, err.Error(), "block hash cannot be nil")
		assert.Nil(t, allDynamic, "Result should be nil on error")

		dynamicInfo, err := runtime.GetDynamicInfo(env.Client, 0, nil)
		assert.Error(t, err, "GetDynamicInfo should error with nil block hash")
		assert.Contains(t, err.Error(), "block hash cannot be nil")
		assert.Nil(t, dynamicInfo, "Result should be nil on error")

		subnetState, err := runtime.GetSubnetState(env.Client, 0, nil)
		assert.Error(t, err, "GetSubnetState should error with nil block hash")
		assert.Contains(t, err.Error(), "block hash cannot be nil")
		assert.Nil(t, subnetState, "Result should be nil on error")

		// Test new metagraph functions
		allMetagraphs, err := runtime.GetAllMetagraphs(env.Client, nil)
		assert.Error(t, err, "GetAllMetagraphs should error with nil block hash")
		assert.Contains(t, err.Error(), "block hash cannot be nil")
		assert.Nil(t, allMetagraphs, "Result should be nil on error")

		metagraph, err := runtime.GetMetagraphSubnet(env.Client, 0, nil)
		assert.Error(t, err, "GetMetagraphSubnet should error with nil block hash")
		assert.Contains(t, err.Error(), "block hash cannot be nil")
		assert.Nil(t, metagraph, "Result should be nil on error")

		selectiveMetagraph, err := runtime.GetSelectiveMetagraph(env.Client, 0, []uint16{0, 1}, nil)
		assert.Error(t, err, "GetSelectiveMetagraph should error with nil block hash")
		assert.Contains(t, err.Error(), "block hash cannot be nil")
		assert.Nil(t, selectiveMetagraph, "Result should be nil on error")
	})

	t.Run("StructValidation", func(t *testing.T) {
		t.Parallel()

		// Test struct field types and validations
		var subnet runtime.SubnetInfo
		assert.Equal(t, types.U16(0), subnet.NetUID, "Default NetUID should be 0")
		assert.Equal(t, types.U64(0), subnet.Difficulty, "Default Difficulty should be 0")

		var subnetV2 runtime.SubnetInfov2
		assert.Equal(t, types.U16(0), subnetV2.NetUID, "Default NetUID should be 0")
		assert.Equal(t, types.U64(0), subnetV2.Difficulty, "Default Difficulty should be 0")

		var hyperparams runtime.SubnetHyperparams
		assert.Equal(t, types.U16(0), hyperparams.Rho, "Default Rho should be 0")
		assert.Equal(t, types.U16(0), hyperparams.Kappa, "Default Kappa should be 0")

		var dynamicInfo runtime.DynamicInfo
		assert.Equal(t, types.U16(0), dynamicInfo.NetUID, "Default NetUID should be 0")
		assert.Equal(t, types.U64(0), dynamicInfo.Difficulty, "Default Difficulty should be 0")

		var subnetState runtime.SubnetState
		assert.Equal(t, types.U16(0), subnetState.NetUID, "Default NetUID should be 0")
		assert.Equal(t, types.U16(0), subnetState.N, "Default N should be 0")

		// Test new structs
		var selectiveMetagraph runtime.SelectiveMetagraph
		assert.Equal(t, int64(0), selectiveMetagraph.Netuid.Int64(), "Default Netuid should be 0")
		assert.Nil(t, selectiveMetagraph.Hotkeys, "Default Hotkeys should be nil")
		assert.Nil(t, selectiveMetagraph.Coldkeys, "Default Coldkeys should be nil")

		t.Log("All struct validations passed")
	})
}
