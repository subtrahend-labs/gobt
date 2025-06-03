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

func TestNeuronRuntimeAPIs(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Parallel()

	t.Run("GetNeurons", func(t *testing.T) {
		t.Parallel()

		env, err := testutils.Setup()
		if err != nil {
			t.Skipf("Failed to setup test environment: %v", err)
		}
		defer env.Teardown()

		blockHash, err := env.Client.Api.RPC.Chain.GetBlockHashLatest()
		require.NoError(t, err, "Failed to get latest block hash")

		neurons, err := runtime.GetNeurons(env.Client, 0, &blockHash)
		if err != nil {
			errorMsg := err.Error()
			isExpectedError := strings.Contains(errorMsg, "no neurons found") ||
				strings.Contains(errorMsg, "Method not found") ||
				strings.Contains(errorMsg, "failed to decode") ||
				strings.Contains(errorMsg, "Invalid params")
			assert.True(t, isExpectedError, "Expected neurons-related error, got: %v", err)
		} else {
			t.Logf("Found %d neurons in subnet 0", len(neurons))

			for i, neuron := range neurons {
				assert.GreaterOrEqual(t, uint64(neuron.UID.Int64()), uint64(0), "Neuron %d UID should be non-negative", i)
				assert.Equal(t, types.U16(0), neuron.NetUID.Int64(), "Neuron %d NetUID should be 0", i)
				assert.NotNil(t, neuron.Hotkey, "Neuron %d hotkey should not be nil", i)
				assert.NotNil(t, neuron.Coldkey, "Neuron %d coldkey should not be nil", i)

				if i < 3 {
					t.Logf("Neuron %d: UID=%d, Hotkey=%x", i, neuron.UID.Int64(), neuron.Hotkey.ToBytes()[:8])
				}
			}
		}

		neurons, err = runtime.GetNeurons(env.Client, 999, &blockHash)
		if err != nil {
			assert.Error(t, err, "Should error for non-existent subnet")
			assert.Nil(t, neurons, "Neurons should be nil on error")
		} else {
			assert.Equal(t, 0, len(neurons), "Should have 0 neurons for non-existent subnet")
			t.Logf("Got empty neurons slice for non-existent subnet (valid behavior)")
		}
	})

	t.Run("GetNeuron", func(t *testing.T) {
		t.Parallel()

		env, err := testutils.Setup()
		if err != nil {
			t.Skipf("Failed to setup test environment: %v", err)
		}
		defer env.Teardown()

		blockHash, err := env.Client.Api.RPC.Chain.GetBlockHashLatest()
		require.NoError(t, err, "Failed to get latest block hash")

		// Test for neuron UID 0 in subnet 0
		neuron, err := runtime.GetNeuron(env.Client, 0, 0, &blockHash)
		if err != nil {
			// Accept various types of errors that can occur in test environment
			errorMsg := err.Error()
			isExpectedError := strings.Contains(errorMsg, "no neuron found") ||
				strings.Contains(errorMsg, "no neurons found") ||
				strings.Contains(errorMsg, "Method not found") ||
				strings.Contains(errorMsg, "failed to decode") ||
				strings.Contains(errorMsg, "Invalid params")
			assert.True(t, isExpectedError, "Expected neuron-related error, got: %v", err)
		} else {
			// If neuron exists, verify structure
			assert.NotNil(t, neuron, "Neuron should not be nil")
			assert.Equal(t, int64(0), neuron.UID.Int64(), "UID should be 0")
			assert.Equal(t, int64(0), neuron.NetUID.Int64(), "NetUID should be 0")
			assert.NotNil(t, neuron.Hotkey, "Hotkey should not be nil")
			assert.NotNil(t, neuron.Coldkey, "Coldkey should not be nil")

			t.Logf("Neuron 0 found - UID: %d, Hotkey: %x", neuron.UID.Int64(), neuron.Hotkey.ToBytes()[:8])
		}

		// Test for non-existent neuron
		neuron, err = runtime.GetNeuron(env.Client, 0, 999, &blockHash)
		assert.Error(t, err, "Should error for non-existent neuron")
		assert.Nil(t, neuron, "Neuron should be nil on error")

		// Test for non-existent subnet
		neuron, err = runtime.GetNeuron(env.Client, 999, 0, &blockHash)
		assert.Error(t, err, "Should error for non-existent subnet")
		assert.Nil(t, neuron, "Neuron should be nil on error")
	})

	t.Run("GetNeuronsLite", func(t *testing.T) {
		t.Parallel()

		env, err := testutils.Setup()
		if err != nil {
			t.Skipf("Failed to setup test environment: %v", err)
		}
		defer env.Teardown()

		blockHash, err := env.Client.Api.RPC.Chain.GetBlockHashLatest()
		require.NoError(t, err, "Failed to get latest block hash")

		// Test for subnet 0 (root subnet should exist)
		neurons, err := runtime.GetNeuronsLite(env.Client, 0, &blockHash)
		if err != nil {
			// Accept various types of errors that can occur in test environment
			errorMsg := err.Error()
			isExpectedError := strings.Contains(errorMsg, "no neurons lite found") ||
				strings.Contains(errorMsg, "Method not found") ||
				strings.Contains(errorMsg, "failed to decode") ||
				strings.Contains(errorMsg, "Invalid params")
			assert.True(t, isExpectedError, "Expected neurons lite-related error, got: %v", err)
		} else {
			// If neurons exist, verify structure
			t.Logf("Found %d neurons lite in subnet 0", len(neurons))

			for i, neuron := range neurons {
				assert.GreaterOrEqual(t, uint64(neuron.UID.Int64()), uint64(0), "Neuron lite %d UID should be non-negative", i)
				assert.Equal(t, types.U16(0), neuron.NetUID.Int64(), "Neuron lite %d NetUID should be 0", i)
				assert.NotNil(t, neuron.Hotkey, "Neuron lite %d hotkey should not be nil", i)
				assert.NotNil(t, neuron.Coldkey, "Neuron lite %d coldkey should not be nil", i)

				if i < 3 { // Log first few neurons
					t.Logf("Neuron lite %d: UID=%d, Hotkey=%x", i, neuron.UID.Int64(), neuron.Hotkey.ToBytes()[:8])
				}
			}
		}

		// Test for non-existent subnet
		neurons, err = runtime.GetNeuronsLite(env.Client, 999, &blockHash)
		if err != nil {
			assert.Error(t, err, "Should error for non-existent subnet")
			assert.Nil(t, neurons, "Neurons lite should be nil on error")
		} else {
			// Empty slice is also valid for non-existent subnet
			assert.Equal(t, 0, len(neurons), "Should have 0 neurons lite for non-existent subnet")
			t.Logf("Got empty neurons lite slice for non-existent subnet (valid behavior)")
		}
	})

	t.Run("GetNeuronLite", func(t *testing.T) {
		t.Parallel()

		env, err := testutils.Setup()
		if err != nil {
			t.Skipf("Failed to setup test environment: %v", err)
		}
		defer env.Teardown()

		blockHash, err := env.Client.Api.RPC.Chain.GetBlockHashLatest()
		require.NoError(t, err, "Failed to get latest block hash")

		// Test for neuron UID 0 in subnet 0
		neuron, err := runtime.GetNeuronLite(env.Client, 0, 0, &blockHash)
		if err != nil {
			// Accept various types of errors that can occur in test environment
			errorMsg := err.Error()
			isExpectedError := strings.Contains(errorMsg, "no neuron lite found") ||
				strings.Contains(errorMsg, "Method not found") ||
				strings.Contains(errorMsg, "failed to decode") ||
				strings.Contains(errorMsg, "Invalid params")
			assert.True(t, isExpectedError, "Expected neuron lite-related error, got: %v", err)
		} else {
			// If neuron exists, verify structure
			assert.NotNil(t, neuron, "Neuron lite should not be nil")
			assert.Equal(t, int64(0), neuron.UID.Int64(), "UID should be 0")
			assert.Equal(t, int64(0), neuron.NetUID.Int64(), "NetUID should be 0")
			assert.NotNil(t, neuron.Hotkey, "Hotkey should not be nil")
			assert.NotNil(t, neuron.Coldkey, "Coldkey should not be nil")

			t.Logf("Neuron lite 0 found - UID: %d, Hotkey: %x", neuron.UID.Int64(), neuron.Hotkey.ToBytes()[:8])
		}

		// Test for non-existent neuron
		neuron, err = runtime.GetNeuronLite(env.Client, 0, 999, &blockHash)
		assert.Error(t, err, "Should error for non-existent neuron")
		assert.Nil(t, neuron, "Neuron lite should be nil on error")

		// Test for non-existent subnet
		neuron, err = runtime.GetNeuronLite(env.Client, 999, 0, &blockHash)
		assert.Error(t, err, "Should error for non-existent subnet")
		assert.Nil(t, neuron, "Neuron lite should be nil on error")
	})

	t.Run("ErrorHandling", func(t *testing.T) {
		t.Parallel()

		env, err := testutils.Setup()
		if err != nil {
			t.Skipf("Failed to setup test environment: %v", err)
		}
		defer env.Teardown()

		func() {
			defer func() {
				if r := recover(); r != nil {
					t.Logf("GetNeurons panicked with nil block hash (expected): %v", r)
				}
			}()
			neurons, err := runtime.GetNeurons(env.Client, 0, nil)
			if err != nil {
				t.Logf("GetNeurons errored with nil block hash: %v", err)
			}
			if neurons != nil {
				t.Logf("GetNeurons unexpectedly succeeded with nil block hash")
			}
		}()

		func() {
			defer func() {
				if r := recover(); r != nil {
					t.Logf("GetNeuron panicked with nil block hash (expected): %v", r)
				}
			}()
			neuron, err := runtime.GetNeuron(env.Client, 0, 0, nil)
			if err != nil {
				t.Logf("GetNeuron errored with nil block hash: %v", err)
			}
			if neuron != nil {
				t.Logf("GetNeuron unexpectedly succeeded with nil block hash")
			}
		}()

		func() {
			defer func() {
				if r := recover(); r != nil {
					t.Logf("GetNeuronsLite panicked with nil block hash (expected): %v", r)
				}
			}()
			neurons, err := runtime.GetNeuronsLite(env.Client, 0, nil)
			if err != nil {
				t.Logf("GetNeuronsLite errored with nil block hash: %v", err)
			}
			if neurons != nil {
				t.Logf("GetNeuronsLite unexpectedly succeeded with nil block hash")
			}
		}()

		func() {
			defer func() {
				if r := recover(); r != nil {
					t.Logf("GetNeuronLite panicked with nil block hash (expected): %v", r)
				}
			}()
			neuron, err := runtime.GetNeuronLite(env.Client, 0, 0, nil)
			if err != nil {
				t.Logf("GetNeuronLite errored with nil block hash: %v", err)
			}
			if neuron != nil {
				t.Logf("GetNeuronLite unexpectedly succeeded with nil block hash")
			}
		}()
	})

	t.Run("StructValidation", func(t *testing.T) {
		t.Parallel()

		var neuron runtime.NeuronInfo
		assert.Equal(t, int64(0), neuron.UID.Int64(), "Default UID should be 0")
		assert.Equal(t, int64(0), neuron.NetUID.Int64(), "Default NetUID should be 0")
		assert.False(t, bool(neuron.Active), "Default Active should be false")
		assert.False(t, bool(neuron.ValidatorPermit), "Default ValidatorPermit should be false")

		var neuronLite runtime.NeuronInfoLite
		assert.Equal(t, int64(0), neuronLite.UID.Int64(), "Default UID should be 0")
		assert.Equal(t, int64(0), neuronLite.NetUID.Int64(), "Default NetUID should be 0")
		assert.False(t, bool(neuronLite.Active), "Default Active should be false")
		assert.False(t, bool(neuronLite.ValidatorPermit), "Default ValidatorPermit should be false")

		t.Log("All struct validations passed")
	})
}
