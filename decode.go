package main

import (
	"fmt"
	"log"
	"os"

	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v4"
	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/joho/godotenv"
	"github.com/vedhavyas/go-subkey/v2"
)

type AccountInfo struct {
	// Nonce       types.UCompact
	// Consumers   types.UCompact
	// Providers   types.UCompact
	// Sufficients types.UCompact
	Nonce       types.U32
	Consumers   types.U32
	Providers   types.U32
	Sufficients types.U32
	Data        struct {
		Free     types.U64
		Reserved types.U64
		Frozen   types.U64
		Flags    types.U128 // Assuming ExtraFlags is u128
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

	for _, module := range meta.AsMetadataV14.Pallets {
		if string(module.Name) == "System" {
			// Check if Storage exists
			if module.HasStorage {
				// Access entries through Value
				for _, storageEntry := range module.Storage.Items {
					if string(storageEntry.Name) == "Account" {
						fmt.Printf("Account type: %v\n", storageEntry.Type)

						// Get the type ID if it's a map
						if storageEntry.Type.IsMap {
							valueTypeId := storageEntry.Type.AsMap.Value.Int64()
							fmt.Printf("Account value type ID: %d\n", valueTypeId)

							// Look up this type ID in the type registry
							if accountType, found := meta.AsMetadataV14.EfficientLookup[valueTypeId]; found {
								fmt.Printf("Account type definition: %+v\n", accountType)
							}
						}
						break
					}
				}
				break
			}
		}
	}

}
