package client

import (
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	// without options
	client, err := NewClient("wss://test.finney.opentensor.ai:443")
	assert.Nil(t, err, "Should not fail on valid endpoint")
	assert.NotNil(t, client, "Should return a client object")
	assert.Equal(t, uint16(42), client.Network, "Should have default network")
	assert.NotNil(t, client.Api, "Should have an API object from valid subtensor")
	assert.Nil(t, client.keyring, "Should not have a keyring")

	// with options
	client, err = NewClient(
		"wss://test.finney.opentensor.ai:443",
		WithNetwork(1),
	)

	assert.Nil(t, err, "Should not fail on valid endpoint")
	assert.NotNil(t, client, "Should return a client object")
	assert.Equal(t, uint16(1), client.Network, "Should have network 1")
	assert.NotNil(t, client.Api, "Should have an API object from valid subtensor")
	assert.Nil(t, client.keyring, "Should not have a keyring")

	testKeyring := signature.TestKeyringPairAlice
	client, err = NewClient(
		"wss://test.finney.opentensor.ai:443",
		WithNetwork(1),
		WithKeyring(&testKeyring),
	)

	assert.Nil(t, err, "Should not fail on valid endpoint")
	assert.NotNil(t, client, "Should return a client object")
	assert.Equal(t, uint16(1), client.Network, "Should have network 1")
	assert.NotNil(t, client.Api, "Should have an API object from valid subtensor")
	assert.NotNil(t, client.keyring, "Should have a keyring")
	assert.Equal(t, testKeyring, *client.keyring, "Should have the same keyring")
}
