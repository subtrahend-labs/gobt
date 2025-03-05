package main

import (
	"fmt"

	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v4"
)

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
