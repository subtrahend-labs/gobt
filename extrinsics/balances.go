package extrinsics

import (
	"log"
	"math/big"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types/extrinsic"
	"github.com/subtrahend-labs/gobt/client"
)

func NewTransferAllowDeath(c *client.Client, recipient types.MultiAddress, amount *big.Int) *extrinsic.Extrinsic {
	call, err := types.NewCall(c.Meta, "Balances.transfer_allow_death", recipient, types.NewUCompact(amount))
	if err != nil {
		log.Fatalf("Error creating call: %s", err)
	}
	ext := extrinsic.NewExtrinsic(call)
	return &ext
}

func NewForceTransfer(c *client.Client, source types.AccountID, recipient types.MultiAddress, amount *big.Int) *extrinsic.Extrinsic {
	call, err := types.NewCall(c.Meta, "Balances.force_transfer", source, recipient, types.NewUCompact(amount))
	if err != nil {
		log.Fatalf("Error creating call: %s", err)
	}
	ext := extrinsic.NewExtrinsic(call)
	return &ext
}

func NewTransferKeepAlive(c *client.Client, recipient types.MultiAddress, amount *big.Int) *extrinsic.Extrinsic {
	call, err := types.NewCall(c.Meta, "Balances.transfer_keep_alive", recipient, types.NewUCompact(amount))
	if err != nil {
		log.Fatalf("Error creating call: %s", err)
	}
	ext := extrinsic.NewExtrinsic(call)
	return &ext
}

func NewTransferAll(c *client.Client, recipient types.MultiAddress, keepAlive bool) *extrinsic.Extrinsic {
	call, err := types.NewCall(c.Meta, "Balances.transfer_all", recipient, keepAlive)
	if err != nil {
		log.Fatalf("Error creating call: %s", err)
	}
	ext := extrinsic.NewExtrinsic(call)
	return &ext
}

func NewForceUnreserve(c *client.Client, who types.AccountID, currencyId types.UCompact) *extrinsic.Extrinsic {
	call, err := types.NewCall(c.Meta, "Balances.force_unreserve", who, currencyId)
	if err != nil {
		log.Fatalf("Error creating call: %s", err)
	}
	ext := extrinsic.NewExtrinsic(call)
	return &ext
}

func NewUpgradeAccounts(c *client.Client, newAccount types.AccountID, numSlashingSpans types.U32) *extrinsic.Extrinsic {
	call, err := types.NewCall(c.Meta, "Balances.upgrade_accounts", newAccount, numSlashingSpans)
	if err != nil {
		log.Fatalf("Error creating call: %s", err)
	}
	ext := extrinsic.NewExtrinsic(call)
	return &ext
}

func NewForceSetBalance(c *client.Client, who types.AccountID, newFree types.U128, newReserved types.U128) *extrinsic.Extrinsic {
	call, err := types.NewCall(c.Meta, "Balances.force_set_balance", who, newFree, newReserved)
	if err != nil {
		log.Fatalf("Error creating call: %s", err)
	}
	ext := extrinsic.NewExtrinsic(call)
	return &ext
}

func NewForceAdjustTotalIssuance(c *client.Client, newTotal types.U128) *extrinsic.Extrinsic {
	call, err := types.NewCall(c.Meta, "Balances.force_adjust_total_issuance", newTotal)
	if err != nil {
		log.Fatalf("Error creating call: %s", err)
	}
	ext := extrinsic.NewExtrinsic(call)
	return &ext
}

func NewBurn(c *client.Client, currencyId types.UCompact, amount types.U128) *extrinsic.Extrinsic {
	call, err := types.NewCall(c.Meta, "Balances.burn", currencyId, amount)
	if err != nil {
		log.Fatalf("Error creating call: %s", err)
	}
	ext := extrinsic.NewExtrinsic(call)
	return &ext
}
