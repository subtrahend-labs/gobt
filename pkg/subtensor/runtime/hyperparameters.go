package runtime

import (
	"fmt"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types/codec"
	"github.com/subtrahend-labs/gobt/pkg/client"
)

type Hyperparameters struct {
	Rho                        types.UCompact // Compact<u16>
	Kappa                      types.UCompact // Compact<u16>
	ImmunityPeriod             types.UCompact // Compact<u16>
	MinAllowedWeights          types.UCompact // Compact<u16>
	MaxWeightsLimit            types.UCompact // Compact<u16>
	Tempo                      types.UCompact // Compact<u16>
	MinDifficulty              types.UCompact // Compact<u64>
	MaxDifficulty              types.UCompact // Compact<u64>
	WeightsVersion             types.UCompact // Compact<u64>
	WeightsRateLimit           types.UCompact // Compact<u64>
	AdjustmentInterval         types.UCompact // Compact<u16>
	ActivityCutoff             types.UCompact // Compact<u16>
	RegistrationAllowed        types.Bool     // bool
	TargetRegsPerInterval      types.UCompact // Compact<u16>
	MinBurn                    types.UCompact // Compact<u64>
	MaxBurn                    types.UCompact // Compact<u64>
	BondsMovingAvg             types.UCompact // Compact<u64>
	MaxRegsPerBlock            types.UCompact // Compact<u16>
	ServingRateLimit           types.UCompact // Compact<u64>
	MaxValidators              types.UCompact // Compact<u16>
	AdjustmentAlpha            types.UCompact // Compact<u64>
	Difficulty                 types.UCompact // Compact<u64>
	CommitRevealPeriod         types.UCompact // Compact<u64>
	CommitRevealWeightsEnabled types.Bool     // bool
	AlphaHigh                  types.UCompact // Compact<u16>
	AlphaLow                   types.UCompact // Compact<u16>
	LiquidAlphaEnabled         types.Bool     // bool
}

// GetMetagraph retrieves the metagraph for a specific subnet
func GetHyperparameters(c *client.Client, netuid uint16, blockHash *types.Hash) (*Hyperparameters, error) {
	// First, try to see what's being returned from the API call
	var encodedResponse []byte
	err := c.Api.Client.Call(
		&encodedResponse,
		"subnetInfo_getSubnetHyperparams",
		netuid,
		blockHash,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to call subnetInfo_getMetagraph: %v", err)
	}

	if len(encodedResponse) == 0 {
		return nil, fmt.Errorf("no metagraph found for netuid %d", netuid)
	}

	// Try to decode as a Vec<u8> first
	var h Hyperparameters
	if err := codec.Decode(encodedResponse, &h); err != nil {
		return nil, fmt.Errorf("failed to decode metagraph: %v", err)
	}
	return &h, nil
}
