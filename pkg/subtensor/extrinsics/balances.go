package extrinsics

import (
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types/extrinsic"
	"github.com/subtrahend-labs/gobt/pkg/client"
)

func TransferAllowDeathCall(c *client.Client, recipient types.MultiAddress, amount types.UCompact) (types.Call, error) {
	call, err := types.NewCall(c.Meta, "Balances.transfer_allow_death", recipient, amount)
	if err != nil {
		return types.Call{}, err
	}
	return call, nil
}

func TransferAllowDeathExt(c *client.Client, recipient types.MultiAddress, amount types.UCompact) (*extrinsic.Extrinsic, error) {
	call, err := TransferAllowDeathCall(c, recipient, amount)
	if err != nil {
		return nil, err
	}
	ext := extrinsic.NewExtrinsic(call)
	return &ext, nil
}

func TransferKeepAliveCall(c *client.Client, recipient types.MultiAddress, amount types.UCompact) (types.Call, error) {
	call, err := types.NewCall(c.Meta, "Balances.transfer_keep_alive", recipient, amount)
	if err != nil {
		return types.Call{}, err
	}
	return call, nil
}

func TransferKeepAliveExt(c *client.Client, recipient types.MultiAddress, amount types.UCompact) (*extrinsic.Extrinsic, error) {
	call, err := TransferKeepAliveCall(c, recipient, amount)
	if err != nil {
		return nil, err
	}
	ext := extrinsic.NewExtrinsic(call)
	return &ext, nil
}

// TODO: should include or only through sudo
func ForceTransferCall(c *client.Client, source types.MultiAddress, recipient types.MultiAddress, amount types.UCompact) (types.Call, error) {
	call, err := types.NewCall(c.Meta, "Balances.force_transfer", source, recipient, amount)
	if err != nil {
		return types.Call{}, err
	}
	return call, nil
}

// TODO: should include or only through sudo
func ForceTransferExt(c *client.Client, source types.MultiAddress, recipient types.MultiAddress, amount types.UCompact) (*extrinsic.Extrinsic, error) {
	call, err := ForceTransferCall(c, source, recipient, amount)
	if err != nil {
		return nil, err
	}
	ext := extrinsic.NewExtrinsic(call)
	return &ext, nil
}

func TransferAllCall(c *client.Client, recipient types.MultiAddress, keepAlive bool) (types.Call, error) {
	call, err := types.NewCall(c.Meta, "Balances.transfer_all", recipient, types.NewBool(keepAlive))
	if err != nil {
		return types.Call{}, err
	}
	return call, nil
}

func TransferAllExt(c *client.Client, recipient types.MultiAddress, keepAlive bool) (*extrinsic.Extrinsic, error) {
	call, err := TransferAllCall(c, recipient, keepAlive)
	if err != nil {
		return nil, err
	}
	ext := extrinsic.NewExtrinsic(call)
	return &ext, nil
}

func ForceUnreserveCall(c *client.Client, who types.AccountID, currencyId types.UCompact) (types.Call, error) {
	call, err := types.NewCall(c.Meta, "Balances.force_unreserve", who, currencyId)
	if err != nil {
		return types.Call{}, err
	}
	return call, nil
}

func ForceUnreserveExt(c *client.Client, who types.AccountID, currencyId types.UCompact) (*extrinsic.Extrinsic, error) {
	call, err := ForceUnreserveCall(c, who, currencyId)
	if err != nil {
		return nil, err
	}
	ext := extrinsic.NewExtrinsic(call)
	return &ext, nil
}

func UpgradeAccountsCall(c *client.Client, newAccount types.AccountID, numSlashingSpans types.U32) (types.Call, error) {
	call, err := types.NewCall(c.Meta, "Balances.upgrade_accounts", newAccount, numSlashingSpans)
	if err != nil {
		return types.Call{}, err
	}
	return call, nil
}

func UpgradeAccountsExt(c *client.Client, newAccount types.AccountID, numSlashingSpans types.U32) (*extrinsic.Extrinsic, error) {
	call, err := UpgradeAccountsCall(c, newAccount, numSlashingSpans)
	if err != nil {
		return nil, err
	}
	ext := extrinsic.NewExtrinsic(call)
	return &ext, nil
}

func ForceSetBalanceCall(c *client.Client, who types.MultiAddress, newFree types.UCompact, newReserved types.UCompact) (types.Call, error) {
	call, err := types.NewCall(c.Meta, "Balances.force_set_balance", who, newFree, newReserved)
	if err != nil {
		return types.Call{}, err
	}
	return call, nil
}

func ForceSetBalanceExt(c *client.Client, who types.MultiAddress, newFree types.UCompact, newReserved types.UCompact) (*extrinsic.Extrinsic, error) {
	call, err := ForceSetBalanceCall(c, who, newFree, newReserved)
	if err != nil {
		return nil, err
	}
	ext := extrinsic.NewExtrinsic(call)
	return &ext, nil
}

func ForceSetTotalIssuanceCall(c *client.Client, newTotal types.U128) (types.Call, error) {
	call, err := types.NewCall(c.Meta, "Balances.force_set_total_issuance", newTotal)
	if err != nil {
		return types.Call{}, err
	}
	return call, nil
}

func ForceSetTotalIssuanceExt(c *client.Client, newTotal types.U128) (*extrinsic.Extrinsic, error) {
	call, err := ForceSetTotalIssuanceCall(c, newTotal)
	if err != nil {
		return nil, err
	}
	ext := extrinsic.NewExtrinsic(call)
	return &ext, nil
}

func BurnCall(c *client.Client, currencyId types.UCompact, amount types.U128) (types.Call, error) {
	call, err := types.NewCall(c.Meta, "Balances.burn", currencyId, amount)
	if err != nil {
		return types.Call{}, err
	}
	return call, nil
}

func BurnExt(c *client.Client, currencyId types.UCompact, amount types.U128) (*extrinsic.Extrinsic, error) {
	call, err := BurnCall(c, currencyId, amount)
	if err != nil {
		return nil, err
	}
	ext := extrinsic.NewExtrinsic(call)
	return &ext, nil
}
