package client

import (
	"fmt"

	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v4"
	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
)

type Client struct {
	api     *gsrpc.SubstrateAPI
	meta    *types.Metadata
	keyring *signature.KeyringPair
	network uint16
}

func NewClient(endpoint string) (*Client, error) {
	client := &Client{
		network: 42,
	}

	api, err := gsrpc.NewSubstrateAPI(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to subtensor node: %w", err)
	}

	client.api = api

	meta, err := api.RPC.State.GetMetadataLatest()
	if err != nil {
		return nil, fmt.Errorf("failed to get metadata from subtensor node: %w", err)
	}

	client.meta = meta

	return client, nil
}
