package extrinsics

import (
	"log"
	"math/big"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types/extrinsic"
	"github.com/subtrahend-labs/gobt/client"
)

func NewTransferKeepAlive(c *client.Client, recipient types.MultiAddress, amount *big.Int) *extrinsic.Extrinsic {
	call, err := types.NewCall(c.Meta, "Balances.transfer_keep_alive", recipient, types.NewUCompact(amount))
	if err != nil {
		log.Fatalf("Error creating call: %s", err)
	}
	ext := extrinsic.NewExtrinsic(call)
	return &ext
}
