package runtime

import (
	"errors"
	"fmt"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types/codec"
	"github.com/subtrahend-labs/gobt/client"
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

type NeuronInfoLite struct {
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

func GetNeuron(c *client.Client, netuid uint16, uid uint16, blockHash *types.Hash) (*NeuronInfo, error) {
	var encodedResponse []byte
	err := c.Api.Client.Call(
		&encodedResponse,
		"neuronInfo_getNeuron",
		netuid,
		uid,
		*blockHash,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to call neuronInfo_getNeurons: %v", err)
	}

	if len(encodedResponse) == 0 {
		return nil, fmt.Errorf("no neurons found for netuid %d", netuid)
	}

	var neuron types.Option[NeuronInfo]
	err = codec.Decode(encodedResponse, &neuron)
	if err != nil {
		return nil, fmt.Errorf("failed to decode neurons: %v", err)
	}
	ok, n := neuron.Unwrap()
	if ok {
		return &n, nil
	}
	return nil, errors.New("no neuron found")
}

func GetNeuronsLite(c *client.Client, netuid uint16, blockHash *types.Hash) ([]NeuronInfoLite, error) {
	var encodedResponse []byte
	err := c.Api.Client.Call(
		&encodedResponse,
		"neuronInfo_getNeuronsLite",
		netuid,
		*blockHash,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to call neuronInfo_getNeuronsLite: %v", err)
	}

	if len(encodedResponse) == 0 {
		return nil, fmt.Errorf("no neurons lite found for netuid %d", netuid)
	}

	var neurons []NeuronInfoLite
	err = codec.Decode(encodedResponse, &neurons)
	if err != nil {
		return nil, fmt.Errorf("failed to decode neurons lite: %v", err)
	}

	return neurons, nil
}

func GetNeuronLite(c *client.Client, netuid uint16, uid uint16, blockHash *types.Hash) (*NeuronInfoLite, error) {
	var encodedResponse []byte
	err := c.Api.Client.Call(
		&encodedResponse,
		"neuronInfo_getNeuronLite",
		netuid,
		uid,
		*blockHash,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to call neuronInfo_getNeuronLite: %v", err)
	}

	if len(encodedResponse) == 0 {
		return nil, fmt.Errorf("no neuron lite found for netuid %d, uid %d", netuid, uid)
	}

	var neuron types.Option[NeuronInfoLite]
	err = codec.Decode(encodedResponse, &neuron)
	if err != nil {
		return nil, fmt.Errorf("failed to decode neuron lite: %v", err)
	}
	ok, n := neuron.Unwrap()
	if ok {
		return &n, nil
	}
	return nil, errors.New("no neuron lite found")
}
