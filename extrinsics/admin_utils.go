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
