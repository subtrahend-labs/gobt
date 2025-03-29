package extrinsics

import (
	"log"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types/extrinsic"
	"github.com/subtrahend-labs/gobt/client"
)

func NewSudo(c *client.Client, ext *extrinsic.Extrinsic) *extrinsic.Extrinsic {
	call, err := types.NewCall(c.Meta, "Sudo.sudo", *ext)
	if err != nil {
		log.Fatalf("Error creating call: %s", err)
	}
	sudoExt := extrinsic.NewExtrinsic(call)
	return &sudoExt
}
