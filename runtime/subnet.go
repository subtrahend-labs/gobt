package runtime

import (
	"errors"
	"fmt"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types/codec"
	"github.com/subtrahend-labs/gobt/client"
)

type SubnetInfo struct {
	NetUID                   types.U16
	Rho                      types.U16
	Kappa                    types.U16
	Difficulty               types.U64
	ImmunityPeriod           types.U16
	ValidatorBatchSize       types.U16
	ValidatorSequenceLength  types.U16
	ValidatorEpochLength     types.U16
	ValidatorEpochsPerReset  types.U16
	ValidatorExcludeQuantile types.U16
	ValidatorPruneLen        types.U64
	ValidatorLogitsWindow    types.U16
	TrimUnderweight          types.Bool
	MaxRegsPerBlock          types.U16
	MaxValidators            types.U16
	BondsMovingAvg           types.U64
	MaxBurn                  types.U64
	NetworkModality          types.U16
	NetworkConnect           [][]types.U16
	EmissionValues           types.U64
	Burn                     types.U64
	Owner                    types.AccountID
}

type SubnetInfov2 struct {
	NetUID                         types.U16
	Rho                            types.U16
	Kappa                          types.U16
	Difficulty                     types.U64
	ImmunityPeriod                 types.U16
	ValidatorBatchSize             types.U16
	ValidatorSequenceLength        types.U16
	ValidatorEpochLength           types.U16
	ValidatorEpochsPerReset        types.U16
	ValidatorExcludeQuantile       types.U16
	ValidatorPruneLen              types.U64
	ValidatorLogitsWindow          types.U16
	TrimUnderweight                types.Bool
	MaxRegsPerBlock                types.U16
	MaxValidators                  types.U16
	BondsMovingAvg                 types.U64
	MaxBurn                        types.U64
	NetworkModality                types.U16
	NetworkConnect                 [][]types.U16
	EmissionValues                 types.U64
	Burn                           types.U64
	Owner                          types.AccountID
	FoundationAccount              types.AccountID
	FoundationDistribution         types.U64
	IncentivePruningDenominator    types.U64
	StakePruningDenominator        types.U64
	StakePruningMin                types.U64
	SynergyScalingLawPower         types.U16
	SubnetworkN                    types.U16
	MinValidatorPermitStake        types.U64
	MaxAllowedModules              types.U16
	BondsMovingAverage             types.U64
	MaxRegistrationsPerBlock       types.U16
	TargetRegistrationsPerInterval types.U16
	TargetRegistrationsInterval    types.U16
	MinBurn                        types.U64
	MaxBurnAsPercent               types.U16
	TxRateLimit                    types.U64
	TxDelegateTakeRateLimit        types.U64
	TxChildKeyTakeRateLimit        types.U64
	ChildKeyTakeFrom               []types.AccountID
	MaxChildRegistrations          types.U16
}

type SubnetHyperparams struct {
	Rho                            types.U16
	Kappa                          types.U16
	ImmunityPeriod                 types.U16
	ValidatorBatchSize             types.U16
	ValidatorSequenceLength        types.U16
	ValidatorEpochLength           types.U16
	ValidatorEpochsPerReset        types.U16
	ValidatorExcludeQuantile       types.U16
	ValidatorPruneLen              types.U64
	ValidatorLogitsWindow          types.U16
	TrimUnderweight                types.Bool
	MaxRegsPerBlock                types.U16
	MaxValidators                  types.U16
	BondsMovingAvg                 types.U64
	MaxBurn                        types.U64
	NetworkModality                types.U16
	NetworkConnect                 [][]types.U16
	IncentivePruningDenominator    types.U64
	StakePruningDenominator        types.U64
	StakePruningMin                types.U64
	SynergyScalingLawPower         types.U16
	SubnetworkN                    types.U16
	MinValidatorPermitStake        types.U64
	MaxAllowedModules              types.U16
	BondsMovingAverage             types.U64
	MaxRegistrationsPerBlock       types.U16
	TargetRegistrationsPerInterval types.U16
	TargetRegistrationsInterval    types.U16
	MinBurn                        types.U64
	MaxBurnAsPercent               types.U16
	TxRateLimit                    types.U64
	TxDelegateTakeRateLimit        types.U64
	TxChildKeyTakeRateLimit        types.U64
	ChildKeyTakeFrom               []types.AccountID
	MaxChildRegistrations          types.U16
}

type DynamicInfo struct {
	NetUID                         types.U16
	Difficulty                     types.U64
	Burn                           types.U64
	MinBurn                        types.U64
	MaxBurn                        types.U64
	MaxRegistrationsPerBlock       types.U16
	TargetRegistrationsPerInterval types.U16
	TargetRegistrationsInterval    types.U16
	Owner                          types.AccountID
	MaxAllowedModules              types.U16
}

type SubnetState struct {
	NetUID             types.U16
	N                  types.U16
	StakeThreshold     types.U64
	Founder            types.AccountID
	FounderShare       types.U16
	IncentiveMechanism types.U8
	MaxAllowedModules  types.U16
	MaxAllowedWeights  types.U16
	MinAllowedWeights  types.U16
	TxRateLimit        types.U64
}

type SelectiveMetagraph struct {
	// Subnet index
	Netuid types.UCompact // Compact<u16>

	// Selected fields based on metagraph_indexes
	Hotkeys         []types.AccountID // Vec<AccountId>
	Coldkeys        []types.AccountID
	Active          []types.Bool
	ValidatorPermit []types.Bool

	PruningScore []types.UCompact // Vec<Compact<u16>>
	LastUpdate   []types.UCompact // Vec<Compact<u64>>
	Emission     []types.UCompact // Vec<Compact<u64>>
	Dividends    []types.UCompact // Vec<Compact<u16>>
	Incentives   []types.UCompact // Vec<Compact<u16>>
	Consensus    []types.UCompact // Vec<Compact<u16>>
	Trust        []types.UCompact // Vec<Compact<u16>>
	Rank         []types.UCompact // Vec<Compact<u16>>

	BlockAtRegistration []types.UCompact // Vec<Compact<u64>>
	AlphaStake          []types.UCompact // Vec<Compact<u64>>
	TaoStake            []types.UCompact // Vec<Compact<u64>>
	TotalStake          []types.UCompact // Vec<Compact<u64>>
}

func GetSubnetInfo(c *client.Client, netuid uint16, blockHash *types.Hash) (*SubnetInfo, error) {
	if blockHash == nil {
		return nil, fmt.Errorf("block hash cannot be nil")
	}

	var encodedResponse []byte
	err := c.Api.Client.Call(
		&encodedResponse,
		"subnetInfo_getSubnetInfo",
		netuid,
		*blockHash,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to call subnetInfo_getSubnetInfo: %v", err)
	}

	if len(encodedResponse) == 0 {
		return nil, fmt.Errorf("no subnet info found for netuid %d", netuid)
	}

	var subnet types.Option[SubnetInfo]
	err = codec.Decode(encodedResponse, &subnet)
	if err != nil {
		return nil, fmt.Errorf("failed to decode subnet info: %v", err)
	}

	ok, s := subnet.Unwrap()
	if ok {
		return &s, nil
	}
	return nil, errors.New("no subnet info found")
}

func GetSubnetsInfo(c *client.Client, blockHash *types.Hash) ([]types.Option[SubnetInfo], error) {
	if blockHash == nil {
		return nil, fmt.Errorf("block hash cannot be nil")
	}

	var encodedResponse []byte
	err := c.Api.Client.Call(
		&encodedResponse,
		"subnetInfo_getSubnetsInfo",
		*blockHash,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to call subnetInfo_getSubnetsInfo: %v", err)
	}

	if len(encodedResponse) == 0 {
		return nil, fmt.Errorf("no subnets info found")
	}

	var subnets []types.Option[SubnetInfo]
	err = codec.Decode(encodedResponse, &subnets)
	if err != nil {
		return nil, fmt.Errorf("failed to decode subnets info: %v", err)
	}

	return subnets, nil
}

func GetSubnetInfoV2(c *client.Client, netuid uint16, blockHash *types.Hash) (*SubnetInfov2, error) {
	if blockHash == nil {
		return nil, fmt.Errorf("block hash cannot be nil")
	}

	var encodedResponse []byte
	err := c.Api.Client.Call(
		&encodedResponse,
		"subnetInfo_getSubnetInfoV2",
		netuid,
		*blockHash,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to call subnetInfo_getSubnetInfoV2: %v", err)
	}

	if len(encodedResponse) == 0 {
		return nil, fmt.Errorf("no subnet info v2 found for netuid %d", netuid)
	}

	var subnet types.Option[SubnetInfov2]
	err = codec.Decode(encodedResponse, &subnet)
	if err != nil {
		return nil, fmt.Errorf("failed to decode subnet info v2: %v", err)
	}

	ok, s := subnet.Unwrap()
	if ok {
		return &s, nil
	}
	return nil, errors.New("no subnet info v2 found")
}

func GetSubnetsInfoV2(c *client.Client, blockHash *types.Hash) ([]types.Option[SubnetInfov2], error) {
	if blockHash == nil {
		return nil, fmt.Errorf("block hash cannot be nil")
	}

	var encodedResponse []byte
	err := c.Api.Client.Call(
		&encodedResponse,
		"subnetInfo_getSubnetsInfoV2",
		*blockHash,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to call subnetInfo_getSubnetsInfoV2: %v", err)
	}

	if len(encodedResponse) == 0 {
		return nil, fmt.Errorf("no subnets info v2 found")
	}

	var subnets []types.Option[SubnetInfov2]
	err = codec.Decode(encodedResponse, &subnets)
	if err != nil {
		return nil, fmt.Errorf("failed to decode subnets info v2: %v", err)
	}

	return subnets, nil
}

func GetSubnetHyperparams(c *client.Client, netuid uint16, blockHash *types.Hash) (*SubnetHyperparams, error) {
	if blockHash == nil {
		return nil, fmt.Errorf("block hash cannot be nil")
	}

	var encodedResponse []byte
	err := c.Api.Client.Call(
		&encodedResponse,
		"subnetInfo_getSubnetHyperparams",
		netuid,
		*blockHash,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to call subnetInfo_getSubnetHyperparams: %v", err)
	}

	if len(encodedResponse) == 0 {
		return nil, fmt.Errorf("no subnet hyperparams found for netuid %d", netuid)
	}

	var hyperparams types.Option[SubnetHyperparams]
	err = codec.Decode(encodedResponse, &hyperparams)
	if err != nil {
		return nil, fmt.Errorf("failed to decode subnet hyperparams: %v", err)
	}

	ok, h := hyperparams.Unwrap()
	if ok {
		return &h, nil
	}
	return nil, errors.New("no subnet hyperparams found")
}

func GetAllDynamicInfo(c *client.Client, blockHash *types.Hash) ([]types.Option[DynamicInfo], error) {
	if blockHash == nil {
		return nil, fmt.Errorf("block hash cannot be nil")
	}

	var encodedResponse []byte
	err := c.Api.Client.Call(
		&encodedResponse,
		"subnetInfo_getAllDynamicInfo",
		*blockHash,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to call subnetInfo_getAllDynamicInfo: %v", err)
	}

	if len(encodedResponse) == 0 {
		return nil, fmt.Errorf("no dynamic info found")
	}

	var dynamicInfo []types.Option[DynamicInfo]
	err = codec.Decode(encodedResponse, &dynamicInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to decode dynamic info: %v", err)
	}

	return dynamicInfo, nil
}

func GetDynamicInfo(c *client.Client, netuid uint16, blockHash *types.Hash) (*DynamicInfo, error) {
	if blockHash == nil {
		return nil, fmt.Errorf("block hash cannot be nil")
	}

	var encodedResponse []byte
	err := c.Api.Client.Call(
		&encodedResponse,
		"subnetInfo_getDynamicInfo",
		netuid,
		*blockHash,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to call subnetInfo_getDynamicInfo: %v", err)
	}

	if len(encodedResponse) == 0 {
		return nil, fmt.Errorf("no dynamic info found for netuid %d", netuid)
	}

	var dynamicInfo types.Option[DynamicInfo]
	err = codec.Decode(encodedResponse, &dynamicInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to decode dynamic info: %v", err)
	}

	ok, d := dynamicInfo.Unwrap()
	if ok {
		return &d, nil
	}
	return nil, errors.New("no dynamic info found")
}

func GetSubnetState(c *client.Client, netuid uint16, blockHash *types.Hash) (*SubnetState, error) {
	if blockHash == nil {
		return nil, fmt.Errorf("block hash cannot be nil")
	}

	var encodedResponse []byte
	err := c.Api.Client.Call(
		&encodedResponse,
		"subnetInfo_getSubnetState",
		netuid,
		*blockHash,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to call subnetInfo_getSubnetState: %v", err)
	}

	if len(encodedResponse) == 0 {
		return nil, fmt.Errorf("no subnet state found for netuid %d", netuid)
	}

	var subnetState types.Option[SubnetState]
	err = codec.Decode(encodedResponse, &subnetState)
	if err != nil {
		return nil, fmt.Errorf("failed to decode subnet state: %v", err)
	}

	ok, s := subnetState.Unwrap()
	if ok {
		return &s, nil
	}
	return nil, errors.New("no subnet state found")
}

func GetAllMetagraphs(c *client.Client, blockHash *types.Hash) ([]types.Option[Metagraph], error) {
	if blockHash == nil {
		return nil, fmt.Errorf("block hash cannot be nil")
	}

	var encodedResponse []byte
	err := c.Api.Client.Call(
		&encodedResponse,
		"subnetInfo_getAllMetagraphs",
		*blockHash,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to call subnetInfo_getAllMetagraphs: %v", err)
	}

	if len(encodedResponse) == 0 {
		return nil, fmt.Errorf("no metagraphs found")
	}

	var metagraphs []types.Option[Metagraph]
	err = codec.Decode(encodedResponse, &metagraphs)
	if err != nil {
		return nil, fmt.Errorf("failed to decode metagraphs: %v", err)
	}

	return metagraphs, nil
}

func GetMetagraphSubnet(c *client.Client, netuid uint16, blockHash *types.Hash) (*Metagraph, error) {
	if blockHash == nil {
		return nil, fmt.Errorf("block hash cannot be nil")
	}

	// Use the existing GetMetagraph function from metagraph.go to avoid duplication
	return GetMetagraph(c, netuid, blockHash)
}

func GetSelectiveMetagraph(c *client.Client, netuid uint16, metagraphIndexes []uint16, blockHash *types.Hash) (*SelectiveMetagraph, error) {
	if blockHash == nil {
		return nil, fmt.Errorf("block hash cannot be nil")
	}

	// Convert []uint16 to []types.U16 for the API call
	indexes := make([]types.U16, len(metagraphIndexes))
	for i, idx := range metagraphIndexes {
		indexes[i] = types.NewU16(idx)
	}

	var encodedResponse []byte
	err := c.Api.Client.Call(
		&encodedResponse,
		"subnetInfo_getSelectiveMetagraph",
		netuid,
		indexes,
		*blockHash,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to call subnetInfo_getSelectiveMetagraph: %v", err)
	}

	if len(encodedResponse) == 0 {
		return nil, fmt.Errorf("no selective metagraph found for netuid %d", netuid)
	}

	var selectiveMetagraph types.Option[SelectiveMetagraph]
	err = codec.Decode(encodedResponse, &selectiveMetagraph)
	if err != nil {
		return nil, fmt.Errorf("failed to decode selective metagraph: %v", err)
	}

	ok, sm := selectiveMetagraph.Unwrap()
	if ok {
		return &sm, nil
	}
	return nil, errors.New("no selective metagraph found")
}
