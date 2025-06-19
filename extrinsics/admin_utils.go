package extrinsics

import (
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types/extrinsic"
	"github.com/subtrahend-labs/gobt/client"
)

func SudoSetNetworkRateLimitCall(c *client.Client, rateLimit types.U64) (types.Call, error) {
	call, err := types.NewCall(
		c.Meta,
		"AdminUtils.sudo_set_network_rate_limit",
		rateLimit,
	)
	if err != nil {
		return types.Call{}, err
	}
	return call, nil
}

func SudoSetNetworkRateLimitExt(c *client.Client, rateLimit types.U64) (*extrinsic.Extrinsic, error) {
	call, err := SudoSetNetworkRateLimitCall(c, rateLimit)
	if err != nil {
		return nil, err
	}
	ext := extrinsic.NewExtrinsic(call)
	return &ext, nil
}

func SudoSetMaxDifficultyCall(c *client.Client, netuid uint16, max_difficulty types.U64) (types.Call, error) {
	call, err := types.NewCall(
		c.Meta,
		"AdminUtils.sudo_set_max_difficulty",
		netuid,
		max_difficulty,
	)
	if err != nil {
		return types.Call{}, err
	}
	return call, nil
}

func SudoSetMaxDifficultyExt(c *client.Client, netuid uint16, max_difficulty types.U64) (*extrinsic.Extrinsic, error) {
	call, err := SudoSetMaxDifficultyCall(c, netuid, max_difficulty)
	if err != nil {
		return nil, err
	}
	ext := extrinsic.NewExtrinsic(call)
	return &ext, nil
}

func SudoSetMinDifficultyCall(c *client.Client, netuid uint16, min_difficulty types.U64) (types.Call, error) {
	call, err := types.NewCall(
		c.Meta,
		"AdminUtils.sudo_set_min_difficulty",
		netuid,
		min_difficulty,
	)
	if err != nil {
		return types.Call{}, err
	}
	return call, nil
}

func SudoSetMinDifficultyExt(c *client.Client, netuid uint16, min_difficulty types.U64) (*extrinsic.Extrinsic, error) {
	call, err := SudoSetMinDifficultyCall(c, netuid, min_difficulty)
	if err != nil {
		return nil, err
	}
	ext := extrinsic.NewExtrinsic(call)
	return &ext, nil
}

func SudoSetDifficultyCall(c *client.Client, netuid uint16, default_difficulty types.U64) (types.Call, error) {
	call, err := types.NewCall(
		c.Meta,
		"AdminUtils.sudo_set_difficulty",
		netuid,
		default_difficulty,
	)
	if err != nil {
		return types.Call{}, err
	}
	return call, nil
}

func SudoSetDifficultyExt(c *client.Client, netuid uint16, default_difficulty types.U64) (*extrinsic.Extrinsic, error) {
	call, err := SudoSetDifficultyCall(c, netuid, default_difficulty)
	if err != nil {
		return nil, err
	}
	ext := extrinsic.NewExtrinsic(call)
	return &ext, nil
}

func SudoSetWeightsVersionKeyCall(c *client.Client, netuid uint16, weights_version_key types.U64) (types.Call, error) {
	call, err := types.NewCall(
		c.Meta,
		"AdminUtils.sudo_set_weights_version_key",
		netuid,
		weights_version_key,
	)
	if err != nil {
		return types.Call{}, err
	}
	return call, nil
}

func SudoSetWeightsVersionKeyExt(c *client.Client, netuid uint16, weights_version_key types.U64) (*extrinsic.Extrinsic, error) {
	call, err := SudoSetWeightsVersionKeyCall(c, netuid, weights_version_key)
	if err != nil {
		return nil, err
	}
	ext := extrinsic.NewExtrinsic(call)
	return &ext, nil
}

func SudoSetTempoCall(c *client.Client, netuid uint16, tempo uint16) (types.Call, error) {
	call, err := types.NewCall(
		c.Meta,
		"AdminUtils.sudo_set_tempo",
		netuid,
		tempo,
	)
	if err != nil {
		return types.Call{}, err
	}
	return call, nil
}

func SudoSetTempoExt(c *client.Client, netuid uint16, tempo uint16) (*extrinsic.Extrinsic, error) {
	call, err := SudoSetTempoCall(c, netuid, tempo)
	if err != nil {
		return nil, err
	}
	ext := extrinsic.NewExtrinsic(call)
	return &ext, nil
}

func SudoSetDefaultTakeCall(c *client.Client, default_take types.U16) (types.Call, error) {
	call, err := types.NewCall(
		c.Meta,
		"AdminUtils.sudo_set_default_take",
		default_take,
	)
	if err != nil {
		return types.Call{}, err
	}
	return call, nil
}

func SudoSetDefaultTakeExt(c *client.Client, default_take types.U16) (*extrinsic.Extrinsic, error) {
	call, err := SudoSetDefaultTakeCall(c, default_take)
	if err != nil {
		return nil, err
	}
	ext := extrinsic.NewExtrinsic(call)
	return &ext, nil
}

func SudoSetTxRateLimitCall(c *client.Client, tx_rate_limit types.U64) (types.Call, error) {
	call, err := types.NewCall(
		c.Meta,
		"AdminUtils.sudo_set_tx_rate_limit",
		tx_rate_limit,
	)
	if err != nil {
		return types.Call{}, err
	}
	return call, nil
}

func SudoSetTxRateLimitExt(c *client.Client, tx_rate_limit types.U64) (*extrinsic.Extrinsic, error) {
	call, err := SudoSetTxRateLimitCall(c, tx_rate_limit)
	if err != nil {
		return nil, err
	}
	ext := extrinsic.NewExtrinsic(call)
	return &ext, nil
}

func SudoSetServingRateLimitCall(c *client.Client, netuid types.U16, serving_rate_limit types.U64) (types.Call, error) {
	call, err := types.NewCall(
		c.Meta,
		"AdminUtils.sudo_set_serving_rate_limit",
		netuid,
		serving_rate_limit,
	)
	if err != nil {
		return types.Call{}, err
	}
	return call, nil
}

func SudoSetServingRateLimitExt(c *client.Client, netuid types.U16, serving_rate_limit types.U64) (*extrinsic.Extrinsic, error) {
	call, err := SudoSetServingRateLimitCall(c, netuid, serving_rate_limit)
	if err != nil {
		return nil, err
	}
	ext := extrinsic.NewExtrinsic(call)
	return &ext, nil
}

func SudoSetAdjustmentIntervalCall(c *client.Client, netuid types.U16, adjustment_interval types.U16) (types.Call, error) {
	call, err := types.NewCall(
		c.Meta,
		"AdminUtils.sudo_set_adjustment_interval",
		netuid,
		adjustment_interval,
	)
	if err != nil {
		return types.Call{}, err
	}
	return call, nil
}

func SudoSetAdjustmentIntervalExt(c *client.Client, netuid types.U16, adjustment_interval types.U16) (*extrinsic.Extrinsic, error) {
	call, err := SudoSetAdjustmentIntervalCall(c, netuid, adjustment_interval)
	if err != nil {
		return nil, err
	}
	ext := extrinsic.NewExtrinsic(call)
	return &ext, nil
}

func SudoSetAdjustmentAlphaCall(c *client.Client, netuid types.U16, adjustment_alpha types.U64) (types.Call, error) {
	call, err := types.NewCall(
		c.Meta,
		"AdminUtils.sudo_set_adjustment_alpha",
		netuid,
		adjustment_alpha,
	)
	if err != nil {
		return types.Call{}, err
	}
	return call, nil
}

func SudoSetAdjustmentAlphaExt(c *client.Client, netuid types.U16, adjustment_alpha types.U64) (*extrinsic.Extrinsic, error) {
	call, err := SudoSetAdjustmentAlphaCall(c, netuid, adjustment_alpha)
	if err != nil {
		return nil, err
	}
	ext := extrinsic.NewExtrinsic(call)
	return &ext, nil
}

func SudoSetMaxWeightLimitCall(c *client.Client, netuid types.U16, max_weight_limit types.U16) (types.Call, error) {
	call, err := types.NewCall(
		c.Meta,
		"AdminUtils.sudo_set_max_weight_limit",
		netuid,
		max_weight_limit,
	)
	if err != nil {
		return types.Call{}, err
	}
	return call, nil
}

func SudoSetMaxWeightLimitExt(c *client.Client, netuid types.U16, max_weight_limit types.U16) (*extrinsic.Extrinsic, error) {
	call, err := SudoSetMaxWeightLimitCall(c, netuid, max_weight_limit)
	if err != nil {
		return nil, err
	}
	ext := extrinsic.NewExtrinsic(call)
	return &ext, nil
}

func SudoSetImmunityPeriodCall(c *client.Client, netuid types.U16, immunity_period types.U16) (types.Call, error) {
	call, err := types.NewCall(
		c.Meta,
		"AdminUtils.sudo_set_immunity_period",
		netuid,
		immunity_period,
	)
	if err != nil {
		return types.Call{}, err
	}
	return call, nil
}

func SudoSetImmunityPeriodExt(c *client.Client, netuid types.U16, immunity_period types.U16) (*extrinsic.Extrinsic, error) {
	call, err := SudoSetImmunityPeriodCall(c, netuid, immunity_period)
	if err != nil {
		return nil, err
	}
	ext := extrinsic.NewExtrinsic(call)
	return &ext, nil
}

func SudoSetMinAllowedWeightsCall(c *client.Client, netuid types.U16, min_allowed_weights types.U16) (types.Call, error) {
	call, err := types.NewCall(
		c.Meta,
		"AdminUtils.sudo_set_min_allowed_weights",
		netuid,
		min_allowed_weights,
	)
	if err != nil {
		return types.Call{}, err
	}
	return call, nil
}

func SudoSetMinAllowedWeightsExt(c *client.Client, netuid types.U16, min_allowed_weights types.U16) (*extrinsic.Extrinsic, error) {
	call, err := SudoSetMinAllowedWeightsCall(c, netuid, min_allowed_weights)
	if err != nil {
		return nil, err
	}
	ext := extrinsic.NewExtrinsic(call)
	return &ext, nil
}

func SudoSetMaxAllowedUidsCall(c *client.Client, netuid types.U16, max_allowed_uids types.U16) (types.Call, error) {
	call, err := types.NewCall(
		c.Meta,
		"AdminUtils.sudo_set_max_allowed_uids",
		netuid,
		max_allowed_uids,
	)
	if err != nil {
		return types.Call{}, err
	}
	return call, nil
}

func SudoSetMaxAllowedUidsExt(c *client.Client, netuid types.U16, max_allowed_uids types.U16) (*extrinsic.Extrinsic, error) {
	call, err := SudoSetMaxAllowedUidsCall(c, netuid, max_allowed_uids)
	if err != nil {
		return nil, err
	}
	ext := extrinsic.NewExtrinsic(call)
	return &ext, nil
}

func SudoSetKappaCall(c *client.Client, netuid types.U16, kappa types.U16) (types.Call, error) {
	call, err := types.NewCall(
		c.Meta,
		"AdminUtils.sudo_set_kappa",
		netuid,
		kappa,
	)
	if err != nil {
		return types.Call{}, err
	}
	return call, nil
}

func SudoSetKappaExt(c *client.Client, netuid types.U16, kappa types.U16) (*extrinsic.Extrinsic, error) {
	call, err := SudoSetKappaCall(c, netuid, kappa)
	if err != nil {
		return nil, err
	}
	ext := extrinsic.NewExtrinsic(call)
	return &ext, nil
}

func SudoSetRhoCall(c *client.Client, netuid types.U16, rho types.U16) (types.Call, error) {
	call, err := types.NewCall(
		c.Meta,
		"AdminUtils.sudo_set_rho",
		netuid,
		rho,
	)
	if err != nil {
		return types.Call{}, err
	}
	return call, nil
}

func SudoSetRhoExt(c *client.Client, netuid types.U16, rho types.U16) (*extrinsic.Extrinsic, error) {
	call, err := SudoSetRhoCall(c, netuid, rho)
	if err != nil {
		return nil, err
	}
	ext := extrinsic.NewExtrinsic(call)
	return &ext, nil
}

func SudoSetActivityCutoffCall(c *client.Client, netuid types.U16, activity_cutoff types.U16) (types.Call, error) {
	call, err := types.NewCall(
		c.Meta,
		"AdminUtils.sudo_set_activity_cutoff",
		netuid,
		activity_cutoff,
	)
	if err != nil {
		return types.Call{}, err
	}
	return call, nil
}

func SudoSetActivityCutoffExt(c *client.Client, netuid types.U16, activity_cutoff types.U16) (*extrinsic.Extrinsic, error) {
	call, err := SudoSetActivityCutoffCall(c, netuid, activity_cutoff)
	if err != nil {
		return nil, err
	}
	ext := extrinsic.NewExtrinsic(call)
	return &ext, nil
}

func SudoSetNetworkRegistrationAllowedCall(c *client.Client, netuid types.U16, registration_allowed bool) (types.Call, error) {
	call, err := types.NewCall(
		c.Meta,
		"AdminUtils.sudo_set_network_registration_allowed",
		netuid,
		registration_allowed,
	)
	if err != nil {
		return types.Call{}, err
	}
	return call, nil
}

func SudoSetNetworkRegistrationAllowedExt(c *client.Client, netuid types.U16, registration_allowed bool) (*extrinsic.Extrinsic, error) {
	call, err := SudoSetNetworkRegistrationAllowedCall(c, netuid, registration_allowed)
	if err != nil {
		return nil, err
	}
	ext := extrinsic.NewExtrinsic(call)
	return &ext, nil
}

func SudoSetNetworkPowRegistrationAllowedCall(c *client.Client, netuid types.U16, registration_allowed bool) (types.Call, error) {
	call, err := types.NewCall(
		c.Meta,
		"AdminUtils.sudo_set_network_pow_registration_allowed",
		netuid,
		registration_allowed,
	)
	if err != nil {
		return types.Call{}, err
	}
	return call, nil
}

func SudoSetNetworkPowRegistrationAllowedExt(c *client.Client, netuid types.U16, registration_allowed bool) (*extrinsic.Extrinsic, error) {
	call, err := SudoSetNetworkPowRegistrationAllowedCall(c, netuid, registration_allowed)
	if err != nil {
		return nil, err
	}
	ext := extrinsic.NewExtrinsic(call)
	return &ext, nil
}

func SudoSetTargetRegistrationsPerIntervalCall(c *client.Client, netuid types.U16, target_registrations_per_interval types.U16) (types.Call, error) {
	call, err := types.NewCall(
		c.Meta,
		"AdminUtils.sudo_set_target_registrations_per_interval",
		netuid,
		target_registrations_per_interval,
	)
	if err != nil {
		return types.Call{}, err
	}
	return call, nil
}

func SudoSetTargetRegistrationsPerIntervalExt(c *client.Client, netuid types.U16, target_registrations_per_interval types.U16) (*extrinsic.Extrinsic, error) {
	call, err := SudoSetTargetRegistrationsPerIntervalCall(c, netuid, target_registrations_per_interval)
	if err != nil {
		return nil, err
	}
	ext := extrinsic.NewExtrinsic(call)
	return &ext, nil
}

func SudoSetMinBurnCall(c *client.Client, netuid types.U16, min_burn types.U64) (types.Call, error) {
	call, err := types.NewCall(
		c.Meta,
		"AdminUtils.sudo_set_min_burn",
		netuid,
		min_burn,
	)
	if err != nil {
		return types.Call{}, err
	}
	return call, nil
}

func SudoSetMinBurnExt(c *client.Client, netuid types.U16, min_burn types.U64) (*extrinsic.Extrinsic, error) {
	call, err := SudoSetMinBurnCall(c, netuid, min_burn)
	if err != nil {
		return nil, err
	}
	ext := extrinsic.NewExtrinsic(call)
	return &ext, nil
}

func SudoSetMaxBurnCall(c *client.Client, netuid types.U16, max_burn types.U64) (types.Call, error) {
	call, err := types.NewCall(
		c.Meta,
		"AdminUtils.sudo_set_max_burn",
		netuid,
		max_burn,
	)
	if err != nil {
		return types.Call{}, err
	}
	return call, nil
}

func SudoSetMaxBurnExt(c *client.Client, netuid types.U16, max_burn types.U64) (*extrinsic.Extrinsic, error) {
	call, err := SudoSetMaxBurnCall(c, netuid, max_burn)
	if err != nil {
		return nil, err
	}
	ext := extrinsic.NewExtrinsic(call)
	return &ext, nil
}


func SudoSetMaxAllowedValidatorsCall(c *client.Client, netuid types.U16, max_allowed_validators types.U16) (types.Call, error) {
	call, err := types.NewCall(
		c.Meta,
		"AdminUtils.sudo_set_max_allowed_validators",
		netuid,
		max_allowed_validators,
	)
	if err != nil {
		return types.Call{}, err
	}
	return call, nil
}

func SudoSetMaxAllowedValidatorsExt(c *client.Client, netuid types.U16, max_allowed_validators types.U16) (*extrinsic.Extrinsic, error) {
	call, err := SudoSetMaxAllowedValidatorsCall(c, netuid, max_allowed_validators)
	if err != nil {
		return nil, err
	}
	ext := extrinsic.NewExtrinsic(call)
	return &ext, nil
}

func SudoSetBondsMovingAverageCall(c *client.Client, netuid types.U16, bonds_moving_average types.U64) (types.Call, error) {
	call, err := types.NewCall(
		c.Meta,
		"AdminUtils.sudo_set_bonds_moving_average",
		netuid,
		bonds_moving_average,
	)
	if err != nil {
		return types.Call{}, err
	}
	return call, nil
}

func SudoSetBondsMovingAverageExt(c *client.Client, netuid types.U16, bonds_moving_average types.U64) (*extrinsic.Extrinsic, error) {
	call, err := SudoSetBondsMovingAverageCall(c, netuid, bonds_moving_average)
	if err != nil {
		return nil, err
	}
	ext := extrinsic.NewExtrinsic(call)
	return &ext, nil
}

func SudoSetBondsPenaltyCall(c *client.Client, netuid types.U16, bonds_penalty types.U16) (types.Call, error) {
	call, err := types.NewCall(
		c.Meta,
		"AdminUtils.sudo_set_bonds_penalty",
		netuid,
		bonds_penalty,
	)
	if err != nil {
		return types.Call{}, err
	}
	return call, nil
}

func SudoSetBondsPenaltyExt(c *client.Client, netuid types.U16, bonds_penalty types.U16) (*extrinsic.Extrinsic, error) {
	call, err := SudoSetBondsPenaltyCall(c, netuid, bonds_penalty)
	if err != nil {
		return nil, err
	}
	ext := extrinsic.NewExtrinsic(call)
	return &ext, nil
}

func SudoSetMaxRegistrationsPerBlockCall(c *client.Client, netuid types.U16, max_registrations_per_block types.U16) (types.Call, error) {
	call, err := types.NewCall(
		c.Meta,
		"AdminUtils.sudo_set_max_registrations_per_block",
		netuid,
		max_registrations_per_block,
	)
	if err != nil {
		return types.Call{}, err
	}
	return call, nil
}

func SudoSetMaxRegistrationsPerBlockExt(c *client.Client, netuid types.U16, max_registrations_per_block types.U16) (*extrinsic.Extrinsic, error) {
	call, err := SudoSetMaxRegistrationsPerBlockCall(c, netuid, max_registrations_per_block)
	if err != nil {
		return nil, err
	}
	ext := extrinsic.NewExtrinsic(call)
	return &ext, nil
}

func SudoSetSubnetOwnerCutCall(c *client.Client, subnet_owner_cut types.U16) (types.Call, error) {
	call, err := types.NewCall(
		c.Meta,
		"AdminUtils.sudo_set_subnet_owner_cut",
		subnet_owner_cut,
	)
	if err != nil {
		return types.Call{}, err
	}
	return call, nil
}

func SudoSetSubnetOwnerCutExt(c *client.Client, subnet_owner_cut types.U16) (*extrinsic.Extrinsic, error) {
	call, err := SudoSetSubnetOwnerCutCall(c, subnet_owner_cut)
	if err != nil {
		return nil, err
	}
	ext := extrinsic.NewExtrinsic(call)
	return &ext, nil
}
