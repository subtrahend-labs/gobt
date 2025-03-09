package client

import (
	"fmt"

	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v4"
	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
)

type Client struct {
	Api     *gsrpc.SubstrateAPI
	Meta    *types.Metadata
	keyring *signature.KeyringPair
	Network uint16
}

func NewClient(endpoint string, opts ...Option) (*Client, error) {
	client := &Client{
		Network: 42,
	}

	for _, opt := range opts {
		opt(client)
	}

	api, err := gsrpc.NewSubstrateAPI(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to subtensor node: %w", err)
	}

	client.Api = api

	meta, err := api.RPC.State.GetMetadataLatest()
	if err != nil {
		return nil, fmt.Errorf("failed to get metadata from subtensor node: %w", err)
	}

	client.Meta = meta

	return client, nil
}
