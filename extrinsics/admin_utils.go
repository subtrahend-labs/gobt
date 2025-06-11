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
