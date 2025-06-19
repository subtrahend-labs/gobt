package extrinsics

import (
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types/extrinsic"
	"github.com/subtrahend-labs/gobt/client"
)

func NewSudoCall(c *client.Client, ext *types.Call) (types.Call, error) {
	call, err := types.NewCall(c.Meta, "Sudo.sudo", ext)
	if err != nil {
		return types.Call{}, err
	}
	return call, nil
}

func NewSudoExt(c *client.Client, ext *types.Call) (*extrinsic.Extrinsic, error) {
	call, err := types.NewCall(c.Meta, "Sudo.sudo", ext)
	if err != nil {
		return nil, err
	}

	sudoExt := extrinsic.NewExtrinsic(call)
	return &sudoExt, nil
}

func NewSudoUncheckedWeightCall(c *client.Client, call *types.Call, weight types.Weight) (types.Call, error) {
	sudoCall, err := types.NewCall(c.Meta, "Sudo.sudo_unchecked_weight", call, weight)
	if err != nil {
		return types.Call{}, err
	}
	return sudoCall, nil
}

func NewSudoUncheckedWeightExt(c *client.Client, call *types.Call, weight types.Weight) (*extrinsic.Extrinsic, error) {
	sudoCall, err := types.NewCall(c.Meta, "Sudo.sudo_unchecked_weight", call, weight)
	if err != nil {
		return nil, err
	}

	sudoExt := extrinsic.NewExtrinsic(sudoCall)
	return &sudoExt, nil
}
