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

	var emission types.U64
	var ok bool
	if block == nil {
		ok, err = c.Api.RPC.State.GetStorageLatest(storageKey, &emission)
	} else {
		ok, err = c.Api.RPC.State.GetStorage(storageKey, &emission, *block)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get storage: %v", err)
	}

	if !ok {
		return nil, fmt.Errorf("storage not found")
	}

	return &emission, nil
}
