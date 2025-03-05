package main

import (
	"fmt"

	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v4"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
)

type SubtensorCRV3WeightsCommitted struct {
	Phase  types.Phase
	Hotkey []types.Hash
}

type CustomEventRecords struct {
	types.EventRecords
	SubtensorModule_CRV3WeightsCommitted []SubtensorCRV3WeightsCommitted
}

func main() {
	testURL := "ws://localhost:9944"
	api, err := gsrpc.NewSubstrateAPI(testURL)
	if err != nil {
		panic(err)
	}
	hash, err := api.RPC.Chain.GetBlockHashLatest()
	if err != nil {
		panic(err)
	}
	fmt.Println(hash.Hex())

	meta, err := api.RPC.State.GetMetadataLatest()
	if err != nil {
		panic(err)
	}

	for _, pallet := range meta.AsMetadataV14.Pallets {
		fmt.Println("====================================")
		fmt.Printf("Name: %v\n", pallet.Name)
		fmt.Println("====================================")
	}

}
