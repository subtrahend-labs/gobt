package main

import (
	"fmt"
	"log"
	// "math/big"
	"os"

	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v4"
	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/joho/godotenv"
	"github.com/vedhavyas/go-subkey/v2"
)

type AccountInfo struct {
	Data struct {
		Free     types.U64 // Balance as u64
		Reserved types.U64
		Frozen   types.U64 // Now Frozen instead of MiscFrozen
		Flags    types.U64 // ExtraFlags field
	}
}

const (
	network uint16 = 42
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	seed := os.Getenv("SOURCE_SEED")
	endpoint := os.Getenv("ENDPOINT")
	destination := os.Getenv("DEST_ACCOUNT_ID")
	sender, err := signature.KeyringPairFromSecret(seed, network)
	if err != nil {
		log.Fatalf("Error creating sender: %s", err)
	}

	recipient, err := types.NewMultiAddressFromHexAccountID(destination)
	if err != nil {
		log.Fatalf("Error creating recipient: %s", err)
	}

	senderSS58 := subkey.SS58Encode(sender.PublicKey, network)
	recipientSS58 := subkey.SS58Encode(recipient.AsID.ToBytes(), network)
	fmt.Println("sender ss58: ", senderSS58)
	fmt.Println("recipient ss58: ", recipientSS58)

	api, err := gsrpc.NewSubstrateAPI(endpoint)
	if err != nil {
		log.Fatalf("Error creating API instance: %s", err)
	}

	meta, err := api.RPC.State.GetMetadataLatest()
	if err != nil {
		log.Fatalf("Error getting metadata: %s", err)
	}

	_ = meta

	senderKey, err := types.CreateStorageKey(meta, "System", "Account", sender.PublicKey)
	if err != nil {
		log.Fatalf("Error creating storage senderKey: %s", err)
	}

	// Get the account data
	var accountInfo AccountInfo
	ok, err := api.RPC.State.GetStorageLatest(senderKey, &accountInfo)
	if err != nil {
		log.Fatalf("Error getting storage: %s", err)
	}
	if !ok {
		log.Fatalf("No storage found")
	}

	// Print the balance information
	fmt.Printf("Free balance: %v\n", accountInfo.Data.Free)
	fmt.Printf("Reserved balance: %v\n", accountInfo.Data.Reserved)
	fmt.Printf("Frozen balance: %v\n", accountInfo.Data.Frozen)
	fmt.Printf("Flags: %v\n", accountInfo.Data.Flags)

	recipientKey, err := types.CreateStorageKey(meta, "System", "Account", recipient.AsID.ToBytes())
	if err != nil {
		log.Fatalf("Error creating storage recipientKey: %s", err)
	}

	// Get the account data
	var recipientInfo AccountInfo
	ok, err = api.RPC.State.GetStorageLatest(recipientKey, &recipientInfo)
	if err != nil {
		log.Fatalf("Error getting storage: %s", err)
	}
	if !ok {
		log.Fatalf("No storage found")
	}

	// Print the balance information
	fmt.Printf("Free balance: %v\n", recipientInfo.Data.Free)
	fmt.Printf("Reserved balance: %v\n", recipientInfo.Data.Reserved)
	fmt.Printf("Frozen balance: %v\n", recipientInfo.Data.Frozen)
	fmt.Printf("Flags: %v\n", recipientInfo.Data.Flags)

	//bal, ok := new(big.Int).SetString("100000000000000", 10)
	//if !ok {
	//	panic(fmt.Errorf("failed to convert balance"))
	//}

	//
	//	c, err := types.NewCall(meta, "Balances.transfer", user, types.NewUCompact(bal))
	//	if err != nil {
	//		panic(err)
	//	}
	//
	//	ext := types.NewExtrinsic(c)
	//
	//	genesisHash, err := api.RPC.Chain.GetBlockHash(0)
	//	if err != nil {
	//		panic(err)
	//	}

	//	rv, err := api.RPC.State.GetRuntimeVersion()
	//	if err != nil {
	//		panic(err)
	//	}
	//
	//	senderKey, err := types.CreateStorageKey(meta, "System", "Account", user.AsID[:])

}
