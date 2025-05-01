package runtime

import (
	"fmt"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types/codec"
	"github.com/subtrahend-labs/gobt/client"
)

// #[freeze_struct("7b506df55bd44646")]
// #[derive(Decode, Encode, PartialEq, Eq, Clone, Debug, TypeInfo)]
// pub struct SubnetHyperparams {
//     rho: Compact<u16>,
//     kappa: Compact<u16>,
//     immunity_period: Compact<u16>,
//     min_allowed_weights: Compact<u16>,
//     max_weights_limit: Compact<u16>,
//     tempo: Compact<u16>,
//     min_difficulty: Compact<u64>,
//     max_difficulty: Compact<u64>,
//     weights_version: Compact<u64>,
//     weights_rate_limit: Compact<u64>,
//     adjustment_interval: Compact<u16>,
//     activity_cutoff: Compact<u16>,
//     pub registration_allowed: bool,
//     target_regs_per_interval: Compact<u16>,
//     min_burn: Compact<u64>,
//     max_burn: Compact<u64>,
//     bonds_moving_avg: Compact<u64>,
//     max_regs_per_block: Compact<u16>,
//     serving_rate_limit: Compact<u64>,
//     max_validators: Compact<u16>,
//     adjustment_alpha: Compact<u64>,
//     difficulty: Compact<u64>,
//     commit_reveal_period: Compact<u64>,
//     commit_reveal_weights_enabled: bool,
//     alpha_high: Compact<u16>,
//     alpha_low: Compact<u16>,
//     liquid_alpha_enabled: bool,
// }

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
