package extrinsics

import (
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types/extrinsic"
	"github.com/subtrahend-labs/gobt/client"
	"github.com/subtrahend-labs/gobt/runtime"
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

// - [ ] sudo_toggle_evm_precompile (Index: 62)
func SudoToggleEvmPrecompileCall(c *client.Client, precompileId types.U8, enabled types.Bool) (types.Call, error) {
	call, err := types.NewCall(
		c.Meta,
		"AdminUtils.sudo_toggle_evm_precompile",
		precompileId,
		enabled,
	)
	if err != nil {
		return types.Call{}, err
	}
	return call, nil
}

func SudoToggleEvmPrecompileExt(c *client.Client, precompileId types.U8, enabled types.Bool) (*extrinsic.Extrinsic, error) {
	call, err := SudoToggleEvmPrecompileCall(c, precompileId, enabled)
	if err != nil {
		return nil, err
	}
	ext := extrinsic.NewExtrinsic(call)
	return &ext, nil
}

func SudoSetSubnetMovingAlphaCall(c *client.Client, alpha runtime.I96F32) (types.Call, error) {
	call, err := types.NewCall(
		c.Meta,
		"AdminUtils.sudo_set_subnet_moving_alpha",
		alpha,
	)
	if err != nil {
		return types.Call{}, err
	}
	return call, nil
}

func SudoSetSubnetMovingAlphaExt(c *client.Client, alpha runtime.I96F32) (*extrinsic.Extrinsic, error) {
	call, err := SudoSetSubnetMovingAlphaCall(c, alpha)
	if err != nil {
		return nil, err
	}
	ext := extrinsic.NewExtrinsic(call)
	return &ext, nil
}

func SudoSetSubnetOwnerHotkeyCall(c *client.Client, netuid types.U16, hotkey types.AccountID) (types.Call, error) {
	call, err := types.NewCall(
		c.Meta,
		"AdminUtils.sudo_set_subnet_owner_hotkey",
		netuid,
		hotkey,
	)
	if err != nil {
		return types.Call{}, err
	}
	return call, nil
}

func SudoSetSubnetOwnerHotkeyExt(c *client.Client, netuid types.U16, hotkey types.AccountID) (*extrinsic.Extrinsic, error) {
	call, err := SudoSetSubnetOwnerHotkeyCall(c, netuid, hotkey)
	if err != nil {
		return nil, err
	}
	ext := extrinsic.NewExtrinsic(call)
	return &ext, nil
}

func SudoSetEmaPriceHalvingPeriodCall(c *client.Client, netuid types.U16, emaHalving types.U64) (types.Call, error) {
	call, err := types.NewCall(
		c.Meta,
		"AdminUtils.sudo_set_ema_price_halving_period",
		netuid,
		emaHalving,
	)
	if err != nil {
		return types.Call{}, err
	}
	return call, nil
}

func SudoSetEmaPriceHalvingPeriodExt(c *client.Client, netuid types.U16, emaHalving types.U64) (*extrinsic.Extrinsic, error) {
	call, err := SudoSetEmaPriceHalvingPeriodCall(c, netuid, emaHalving)
	if err != nil {
		return nil, err
	}
	ext := extrinsic.NewExtrinsic(call)
	return &ext, nil
}
