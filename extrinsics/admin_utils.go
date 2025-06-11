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
