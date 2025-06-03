package runtime

import (
	"errors"
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
	if blockHash == nil {
		return nil, fmt.Errorf("block hash cannot be nil")
	}

	var encodedResponse []byte
	err := c.Api.Client.Call(
		&encodedResponse,
		"delegateInfo_getDelegates",
		*blockHash,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to call delegateInfo_getDelegates: %v", err)
	}

	var delegates []DelegateInfo
	err = codec.Decode(encodedResponse, &delegates)
	if err != nil {
		return nil, fmt.Errorf("failed to decode delegates: %v", err)
	}

	// Return empty array instead of nil when no delegates exist
	if delegates == nil {
		delegates = []DelegateInfo{}
	}

	return delegates, nil
}

func GetDelegate(c *client.Client, delegateAccount types.AccountID, blockHash *types.Hash) (*DelegateInfo, error) {
	if blockHash == nil {
		return nil, fmt.Errorf("block hash cannot be nil")
	}

	var encodedResponse []byte
	err := c.Api.Client.Call(
		&encodedResponse,
		"delegateInfo_getDelegate",
		delegateAccount,
		*blockHash,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to call delegateInfo_getDelegate: %v", err)
	}

	if len(encodedResponse) == 0 {
		return nil, fmt.Errorf("no delegate found for account %v", delegateAccount)
	}

	var delegate types.Option[DelegateInfo]
	err = codec.Decode(encodedResponse, &delegate)
	if err != nil {
		return nil, fmt.Errorf("failed to decode delegate: %v", err)
	}

	ok, d := delegate.Unwrap()
	if ok {
		return &d, nil
	}
	return nil, errors.New("no delegate found")
}

func GetDelegated(c *client.Client, delegateeAccount types.AccountID, blockHash *types.Hash) ([]DelegatedInfo, error) {
	if blockHash == nil {
		return nil, fmt.Errorf("block hash cannot be nil")
	}

	var encodedResponse []byte
	err := c.Api.Client.Call(
		&encodedResponse,
		"delegateInfo_getDelegated",
		delegateeAccount,
		*blockHash,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to call delegateInfo_getDelegated: %v", err)
	}

	var delegated []DelegatedInfo
	err = codec.Decode(encodedResponse, &delegated)
	if err != nil {
		return nil, fmt.Errorf("failed to decode delegated info: %v", err)
	}

	// Return empty array instead of nil when no delegations exist
	if delegated == nil {
		delegated = []DelegatedInfo{}
	}

	return delegated, nil
}

// Runtime API functions for delegate information
// Based on Rust runtime API:
// - get_delegates() -> Vec<DelegateInfo>
// - get_delegate(delegate_account: AccountId32) -> Option<DelegateInfo>
// - get_delegated(delegatee_account: AccountId32) -> Vec<(DelegateInfo, (Compact<u64>, Compact<u16>))>
