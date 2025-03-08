package main

import (
	"fmt"
	"log"
	"math/big"
	"os"

	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v4"
	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/joho/godotenv"
	"github.com/vedhavyas/go-subkey/v2"
)

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

	fmt.Println("Sender SS58:", subkey.SS58Encode(sender.PublicKey, network))
	fmt.Println("Recipient SS58:", subkey.SS58Encode(recipient.AsID.ToBytes(), network))

	api, err := gsrpc.NewSubstrateAPI(endpoint)
	if err != nil {
		log.Fatalf("Error creating API: %s", err)
	}

	meta, err := api.RPC.State.GetMetadataLatest()
	if err != nil {
		log.Fatalf("Error getting metadata: %s", err)
	}

	rv, err := api.RPC.State.GetRuntimeVersionLatest()
	if err != nil {
		log.Fatalf("Error getting runtime version: %s", err)
	}

	// Fetch sender nonce
	key, err := types.CreateStorageKey(meta, "System", "Account", sender.PublicKey)
	if err != nil {
		log.Fatalf("Error creating storage key: %s", err)
	}
	var accountInfo types.AccountInfo // Use standard type
	ok, err := api.RPC.State.GetStorageLatest(key, &accountInfo)
	if err != nil || !ok {
		log.Fatalf("Error getting account info: %s", err)
	}

	// Set transfer amount (0.1 TAO = 10,000,000 units)
	bal, _ := new(big.Int).SetString("10000000", 10)
	c, err := types.NewCall(meta, "Balances.transfer_allow_death", recipient, types.NewUCompact(bal))
	if err != nil {
		log.Fatalf("Error creating call: %s", err)
	}
	ext := types.NewExtrinsic(c)

	genesisHash, err := api.RPC.Chain.GetBlockHash(0)
	if err != nil {
		log.Fatalf("Error getting genesis hash: %s", err)
	}

	// Signature options (mimic example)
	o := types.SignatureOptions{
		BlockHash:          genesisHash,                            // Use genesis hash like example
		Era:                types.ExtrinsicEra{IsMortalEra: false}, // Immortal
		GenesisHash:        genesisHash,
		Nonce:              types.NewUCompactFromUInt(uint64(accountInfo.Nonce)),
		SpecVersion:        rv.SpecVersion,
		Tip:                types.NewUCompactFromUInt(100000), // 0.001 TAO
		TransactionVersion: rv.TransactionVersion,
	}

	// Sign and log extrinsic
	err = ext.Sign(sender, o)
	if err != nil {
		log.Fatalf("Error signing: %s", err)
	}

	// Submit
	hash, err := api.RPC.Author.SubmitExtrinsic(ext)
	if err != nil {
		fmt.Printf("Error details: %+v\n", err)
		log.Fatalf("Error submitting extrinsic: %s", err)
	}
	fmt.Printf("Extrinsic submitted with hash: %v\n", hash)
}
