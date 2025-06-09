package runtime

import (
	"fmt"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types/codec"
	"github.com/subtrahend-labs/gobt/client"
)

type DelegateInfo struct {
	AccountID        types.AccountID
	TakeRate         types.UCompact
	NominatorStake   types.UCompact
	ValidatorStake   types.UCompact
	TotalStake       types.UCompact
	Registrations    []types.UCompact
	VotingPower      types.UCompact
	LastEpochLength  types.UCompact
	LastEpochCost    types.UCompact
	ValidatorPermits []types.UCompact
	Return           types.UCompact
}

type DelegatedInfo struct {
	DelegateInfo DelegateInfo
	NetUID       types.UCompact
	Amount       types.UCompact
}

func GetDelegates(c *client.Client, blockHash *types.Hash) ([]DelegateInfo, error) {
	var targetBlockHash types.Hash
	if blockHash == nil {
		latestHash, err := c.Api.RPC.Chain.GetBlockHashLatest()
		if err != nil {
			return nil, fmt.Errorf("failed to get latest block hash: %v", err)
		}
		targetBlockHash = latestHash
	} else {
		targetBlockHash = *blockHash
	}

	var encodedResponse []byte
	err := c.Api.Client.Call(
		&encodedResponse,
		"delegateInfo_getDelegates",
		targetBlockHash,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to call delegateInfo_getDelegates: %v", err)
	}

	var delegates []DelegateInfo
	err = codec.Decode(encodedResponse, &delegates)
	if err != nil {
		return nil, fmt.Errorf("failed to decode delegates: %v", err)
	}

	if delegates == nil {
		delegates = []DelegateInfo{}
	}

	return delegates, nil
}

func GetDelegate(c *client.Client, delegateAccount types.AccountID, blockHash *types.Hash) (*types.Option[DelegateInfo], error) {
	var targetBlockHash types.Hash
	if blockHash == nil {
		latestHash, err := c.Api.RPC.Chain.GetBlockHashLatest()
		if err != nil {
			return nil, fmt.Errorf("failed to get latest block hash: %v", err)
		}
		targetBlockHash = latestHash
	} else {
		targetBlockHash = *blockHash
	}

	var encodedResponse []byte
	err := c.Api.Client.Call(
		&encodedResponse,
		"delegateInfo_getDelegate",
		delegateAccount,
		targetBlockHash,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to call delegateInfo_getDelegate: %v", err)
	}

	var delegate types.Option[DelegateInfo]
	err = codec.Decode(encodedResponse, &delegate)
	if err != nil {
		return nil, fmt.Errorf("failed to decode delegate: %v", err)
	}

	return &delegate, nil
}

func GetDelegated(c *client.Client, delegateeAccount types.AccountID, blockHash *types.Hash) ([]DelegatedInfo, error) {
	var targetBlockHash types.Hash
	if blockHash == nil {
		latestHash, err := c.Api.RPC.Chain.GetBlockHashLatest()
		if err != nil {
			return nil, fmt.Errorf("failed to get latest block hash: %v", err)
		}
		targetBlockHash = latestHash
	} else {
		targetBlockHash = *blockHash
	}

	var encodedResponse []byte
	err := c.Api.Client.Call(
		&encodedResponse,
		"delegateInfo_getDelegated",
		delegateeAccount,
		targetBlockHash,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to call delegateInfo_getDelegated: %v", err)
	}

	var delegated []DelegatedInfo
	err = codec.Decode(encodedResponse, &delegated)
	if err != nil {
		return nil, fmt.Errorf("failed to decode delegated info: %v", err)
	}

	if delegated == nil {
		delegated = []DelegatedInfo{}
	}

	return delegated, nil
}
