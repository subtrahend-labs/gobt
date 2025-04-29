package storage

import (
	"fmt"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/subtrahend-labs/gobt/client"
)

type AccountInfo struct {
	Nonce       types.U32
	Consumers   types.U32
	Providers   types.U32
	Sufficients types.U32
	Data        struct {
		Free     types.U64
		Reserved types.U64
		Frozen   types.U64
		Flags    types.U128
	}
}

func GetAccountInfo(c *client.Client, accountID []byte, block *types.Hash) (*AccountInfo, error) {
	meta, err := c.Api.RPC.State.GetMetadataLatest()
	if err != nil {
		return nil, fmt.Errorf("failed to get metadata: %v", err)
	}

	storageKey, err := types.CreateStorageKey(meta, "System", "Account", accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to create storage key: %v", err)
	}

	var res AccountInfo
	err = getStorageOptionalBlock(c, storageKey, &res, block)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
