package main

import (
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v4"
	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/joho/godotenv"
	"github.com/vedhavyas/go-subkey/v2"
)

type AccountInfo struct {
	Nonce types.U32
	Data  struct {
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

	bal, ok := new(big.Int).SetString("500000", 10) // Just 0.0005 TAO
	if !ok {
		log.Fatalf("Error converting string to big.Int")
	}

	// Print all modules and their calls
	// fmt.Println("\nAVAILABLE MODULES AND CALLS:")
	// for _, module := range meta.AsMetadataV14.Pallets {
	// 	fmt.Printf("Module: %s (Index: %d)\n", module.Name, module.Index)

	// 	// Print calls if they exist
	// 	if module.HasCalls {
	// 		// Get the call type ID
	// 		callTypeID := module.Calls.Type.Int64()
	// 		fmt.Printf("  Call Type ID: %d\n", callTypeID)

	// 		// Find the type in the lookup
	// 		if callType, ok := meta.AsMetadataV14.EfficientLookup[callTypeID]; ok {
	// 			if callType.Def.IsVariant {
	// 				fmt.Println("  Available Calls:")
	// 				for _, variant := range callType.Def.Variant.Variants {
	// 					fmt.Printf("    %s (Index: %d)\n", variant.Name, variant.Index)
	// 				}
	// 			}
	// 		} else {
	// 			fmt.Printf("  Call type not found in lookup\n")
	// 		}
	// 	} else {
	// 		fmt.Printf("  No calls available\n")
	// 	}
	// 	fmt.Println()
	// }

	c, err := types.NewCall(meta, "Balances.transfer_keep_alive", recipient, types.NewUCompact(bal))
	if err != nil {
		log.Fatalf("Error creating call: %s", err)
	}

	ext := types.NewExtrinsic(c)

	genesisHash, err := api.RPC.Chain.GetBlockHash(0)
	if err != nil {
		log.Fatalf("Error getting genesis hash: %s", err)
	}

	rv, err := api.RPC.State.GetRuntimeVersionLatest()
	if err != nil {
		log.Fatalf("Error getting runtime version: %s", err)
	}

	currentHash, err := api.RPC.Chain.GetBlockHashLatest()
	if err != nil {
		log.Fatalf("Error getting current block: %s", err)
	}

	o := types.SignatureOptions{
		BlockHash:          currentHash,
		Era:                types.ExtrinsicEra{IsMortalEra: false},
		GenesisHash:        genesisHash,
		Nonce:              types.NewUCompactFromUInt(uint64(uint32(accountInfo.Nonce))),
		SpecVersion:        rv.SpecVersion,
		Tip:                types.NewUCompactFromUInt(0),
		TransactionVersion: rv.TransactionVersion,
	}

	err = ext.Sign(sender, o)
	if err != nil {
		log.Fatalf("Error signing: %s", err)
	}

	_, err = api.RPC.Author.SubmitExtrinsic(ext)
	if err != nil {
		log.Fatalf("Error submitting extrinsic: %s", err)
	}

	fmt.Println("Extrinsic submitted")

	// Wait for the extrinsic to be included in a block
	time.Sleep(12 * time.Second)

	// Get the account data
	ok, err = api.RPC.State.GetStorageLatest(senderKey, &accountInfo)
	if err != nil {
		log.Fatalf("Error getting storage: %s", err)
	}
	if !ok {
		log.Fatalf("No storage found")
	}
	fmt.Println("Sender balance after transfer: ", accountInfo.Data.Free)

}
