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
