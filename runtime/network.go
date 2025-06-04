package runtime

import (
	"fmt"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types/codec"
	"github.com/subtrahend-labs/gobt/client"
)

// GetNetworkLockCost retrieves the current network lock cost
func GetNetworkLockCost(c *client.Client, blockHash *types.Hash) (types.U64, error) {
	if blockHash == nil {
		return 0, fmt.Errorf("block hash cannot be nil")
	}

	var encodedResponse []byte
	err := c.Api.Client.Call(
		&encodedResponse,
		"subnetInfo_getNetworkLockCost",
		*blockHash,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to call subnetInfo_getNetworkLockCost: %v", err)
	}

	if len(encodedResponse) == 0 {
		return 0, fmt.Errorf("no network lock cost found")
	}

	var lockCost types.U64
	err = codec.Decode(encodedResponse, &lockCost)
	if err != nil {
		return 0, fmt.Errorf("failed to decode network lock cost: %v", err)
	}

	return lockCost, nil
}

// GetNetworkLastLock retrieves the last network lock value
func GetNetworkLastLock(c *client.Client, blockHash *types.Hash) (types.U64, error) {
	if blockHash == nil {
		return 0, fmt.Errorf("block hash cannot be nil")
	}

	var encodedResponse []byte
	err := c.Api.Client.Call(
		&encodedResponse,
		"subnetInfo_getNetworkLastLock",
		*blockHash,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to call subnetInfo_getNetworkLastLock: %v", err)
	}

	if len(encodedResponse) == 0 {
		return 0, fmt.Errorf("no network last lock found")
	}

	var lastLock types.U64
	err = codec.Decode(encodedResponse, &lastLock)
	if err != nil {
		return 0, fmt.Errorf("failed to decode network last lock: %v", err)
	}

	return lastLock, nil
}

// GetNetworkMinLock retrieves the minimum network lock value
func GetNetworkMinLock(c *client.Client, blockHash *types.Hash) (types.U64, error) {
	if blockHash == nil {
		return 0, fmt.Errorf("block hash cannot be nil")
	}

	var encodedResponse []byte
	err := c.Api.Client.Call(
		&encodedResponse,
		"subnetInfo_getNetworkMinLock",
		*blockHash,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to call subnetInfo_getNetworkMinLock: %v", err)
	}

	if len(encodedResponse) == 0 {
		return 0, fmt.Errorf("no network min lock found")
	}

	var minLock types.U64
	err = codec.Decode(encodedResponse, &minLock)
	if err != nil {
		return 0, fmt.Errorf("failed to decode network min lock: %v", err)
	}

	return minLock, nil
}

// GetNetworkLastLockBlock retrieves the block number of the last network lock
func GetNetworkLastLockBlock(c *client.Client, blockHash *types.Hash) (types.U64, error) {
	if blockHash == nil {
		return 0, fmt.Errorf("block hash cannot be nil")
	}

	var encodedResponse []byte
	err := c.Api.Client.Call(
		&encodedResponse,
		"subnetInfo_getNetworkLastLockBlock",
		*blockHash,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to call subnetInfo_getNetworkLastLockBlock: %v", err)
	}

	if len(encodedResponse) == 0 {
		return 0, fmt.Errorf("no network last lock block found")
	}

	var lastLockBlock types.U64
	err = codec.Decode(encodedResponse, &lastLockBlock)
	if err != nil {
		return 0, fmt.Errorf("failed to decode network last lock block: %v", err)
	}

	return lastLockBlock, nil
}

// GetLockReductionInterval retrieves the lock reduction interval
func GetLockReductionInterval(c *client.Client, blockHash *types.Hash) (types.U64, error) {
	if blockHash == nil {
		return 0, fmt.Errorf("block hash cannot be nil")
	}

	var encodedResponse []byte
	err := c.Api.Client.Call(
		&encodedResponse,
		"subnetInfo_getLockReductionInterval",
		*blockHash,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to call subnetInfo_getLockReductionInterval: %v", err)
	}

	if len(encodedResponse) == 0 {
		return 0, fmt.Errorf("no lock reduction interval found")
	}

	var interval types.U64
	err = codec.Decode(encodedResponse, &interval)
	if err != nil {
		return 0, fmt.Errorf("failed to decode lock reduction interval: %v", err)
	}

	return interval, nil
}

// GetCurrentBlock retrieves the current block number
func GetCurrentBlock(c *client.Client, blockHash *types.Hash) (types.U64, error) {
	if blockHash == nil {
		return 0, fmt.Errorf("block hash cannot be nil")
	}

	var encodedResponse []byte
	err := c.Api.Client.Call(
		&encodedResponse,
		"subnetInfo_getCurrentBlock",
		*blockHash,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to call subnetInfo_getCurrentBlock: %v", err)
	}

	if len(encodedResponse) == 0 {
		return 0, fmt.Errorf("no current block found")
	}

	var currentBlock types.U64
	err = codec.Decode(encodedResponse, &currentBlock)
	if err != nil {
		return 0, fmt.Errorf("failed to decode current block: %v", err)
	}

	return currentBlock, nil
}
