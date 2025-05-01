package runtime

import (
	"fmt"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types/codec"
	"github.com/subtrahend-labs/gobt/pkg/client"
)

type AxonInfo struct {
	Block        types.U64
	Version      types.U32
	IP           types.U128
	Port         types.U16
	IPType       types.U8
	Protocol     types.U8
	Placeholder1 types.U8
	Placeholder2 types.U8
}

type PrometheusInfo struct {
	Block   types.U64
	Version types.U32
	IP      types.U128
	Port    types.U16
	IPType  types.U8
}

type NeuronInfo struct {
	Hotkey         types.AccountID
	Coldkey        types.AccountID
	UID            types.UCompact
	NetUID         types.UCompact
	Active         types.Bool
	AxonInfo       AxonInfo
	PrometheusInfo PrometheusInfo
	Stake          []struct {
		Account types.AccountID
		Amount  types.UCompact
	}
	Rank            types.UCompact
	Emission        types.UCompact
	Incentive       types.UCompact
	Consensus       types.UCompact
	Trust           types.UCompact
	ValidatorTrust  types.UCompact
	Dividends       types.UCompact
	LastUpdate      types.UCompact
	ValidatorPermit types.Bool
	Weights         []struct {
		UID    types.UCompact
		Weight types.UCompact
	}
	Bonds []struct {
		UID  types.UCompact
		Bond types.UCompact
	}
	PruningScore types.UCompact
}

func GetNeurons(c *client.Client, netuid uint16, blockHash *types.Hash) ([]NeuronInfo, error) {
	var encodedResponse []byte
	err := c.Api.Client.Call(
		&encodedResponse,
		"neuronInfo_getNeurons",
		netuid,
		*blockHash,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to call neuronInfo_getNeurons: %v", err)
	}

	if len(encodedResponse) == 0 {
		return nil, fmt.Errorf("no neurons found for netuid %d", netuid)
	}

	var neurons []NeuronInfo
	err = codec.Decode(encodedResponse, &neurons)
	if err != nil {
		return nil, fmt.Errorf("failed to decode neurons: %v", err)
	}

	return neurons, nil
}
