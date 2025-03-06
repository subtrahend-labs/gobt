package main

import (
	"fmt"
	"log"
	"os"

	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v4"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/joho/godotenv"
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
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	endpoint := os.Getenv("ENDPOINT")
	accountID := os.Getenv("ACCOUNT_ID")

	api, err := gsrpc.NewSubstrateAPI(endpoint)
	if err != nil {
		log.Fatalf("Error creating API instance: %s", err)
	}

	meta, err := api.RPC.State.GetMetadataLatest()
	if err != nil {
		log.Fatalf("Error getting metadata: %s", err)
	}

	fmt.Println(meta)

	user, err := types.NewMultiAddressFromHexAccountID(accountID)
	if err != nil {
		log.Fatalf("Error creating user address: %s", err)
	}

	fmt.Println("User details %v", user)

}
