package storage

import (
	"fmt"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/subtrahend-labs/gobt/client"
	"github.com/subtrahend-labs/gobt/typetools"
)

// Returns sn tao emission percentage * 1e7
// block is optional, but more performant if already subscribed to chain
func GetSubnetTaoInEmission(c *client.Client, netuid types.U16, block *types.Hash) (*types.U64, error) {
	meta, err := c.Api.RPC.State.GetMetadataLatest()
	if err != nil {
		return nil, fmt.Errorf("failed to get metadata: %v", err)
	}

	storageKey, err := types.CreateStorageKey(meta, "SubtensorModule", "SubnetTaoInEmission", typetools.Uint16ToBytes(uint16(netuid)))
	if err != nil {
		return nil, fmt.Errorf("failed to create storage key: %v", err)
	}

	var res types.U64
	err = getStorageOptionalBlock(c, storageKey, &res, block)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

// Validator permits per uid
func GetValidatorPermits(c *client.Client, netuid types.U16, block *types.Hash) (*[]types.Bool, error) {
	meta, err := c.Api.RPC.State.GetMetadataLatest()
	if err != nil {
		return nil, fmt.Errorf("failed to get metadata: %v", err)
	}

	storageKey, err := types.CreateStorageKey(meta, "SubtensorModule", "ValidatorPermit", typetools.Uint16ToBytes(uint16(netuid)))
	if err != nil {
		return nil, fmt.Errorf("failed to create storage key: %v", err)
	}

	var res []types.Bool
	err = getStorageOptionalBlock(c, storageKey, &res, block)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
