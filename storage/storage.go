package storage

import (
	"fmt"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/subtrahend-labs/gobt/client"
)

func getStorageOptionalBlock(c *client.Client, key types.StorageKey, res any, block *types.Hash) error {
	var err error
	var ok bool
	if block == nil {
		ok, err = c.Api.RPC.State.GetStorageLatest(key, res)
	} else {
		ok, err = c.Api.RPC.State.GetStorage(key, res, *block)
	}
	if err != nil {
		return fmt.Errorf("failed to get storage: %v", err)
	}

	if !ok {
		return fmt.Errorf("storage not found")
	}
	return nil
}
