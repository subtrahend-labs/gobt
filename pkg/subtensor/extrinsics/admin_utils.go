package extrinsics

import (
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types/extrinsic"
	"github.com/subtrahend-labs/gobt/pkg/client"
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

func SudoSetCommitRevealWeightsEnabledCall(c *client.Client, netuid types.U16, enabled bool) (types.Call, error) {
	call, err := types.NewCall(
		c.Meta,
		"AdminUtils.sudo_set_commit_reveal_weights_enabled",
		netuid,
		enabled,
	)
	if err != nil {
		return types.Call{}, err
	}
	return call, nil
}

func SudoSetCommitRevealWeightsEnabledExt(c *client.Client, netuid types.U16, enabled bool) (*extrinsic.Extrinsic, error) {
	call, err := SudoSetCommitRevealWeightsEnabledCall(c, netuid, enabled)
	if err != nil {
		return nil, err
	}
	ext := extrinsic.NewExtrinsic(call)
	return &ext, nil
}

func SudoSetWeightsSetRateLimitCall(c *client.Client, netuid types.U16, weights_set_rate_limit types.U64) (types.Call, error) {
	call, err := types.NewCall(
		c.Meta,
		"AdminUtils.sudo_set_weights_set_rate_limit",
		netuid,
		weights_set_rate_limit,
	)
	if err != nil {
		return types.Call{}, err
	}
	return call, nil
}

func SudoSetWeightsSetRateLimitExt(c *client.Client, netuid types.U16, weights_set_rate_limit types.U64) (*extrinsic.Extrinsic, error) {
	call, err := SudoSetWeightsSetRateLimitCall(c, netuid, weights_set_rate_limit)
	if err != nil {
		return nil, err
	}
	ext := extrinsic.NewExtrinsic(call)
	return &ext, nil
}
