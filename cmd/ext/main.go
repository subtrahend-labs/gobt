package main

import (
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types/extrinsic"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types/extrinsic/extensions"
	"github.com/joho/godotenv"
	"github.com/subtrahend-labs/gobt/client"
	gobtext "github.com/subtrahend-labs/gobt/extrinsics"
	"github.com/subtrahend-labs/gobt/storage"
	"github.com/vedhavyas/go-subkey/v2"
)

const (
	network uint16 = 42
)

var (
	endpoint  string
	sender    signature.KeyringPair
	recipient types.MultiAddress
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	seed := os.Getenv("SOURCE_SEED")
	endpoint = os.Getenv("ENDPOINT")
	destination := os.Getenv("DEST_ACCOUNT_ID")
	sender, err = signature.KeyringPairFromSecret(seed, network)
	if err != nil {
		log.Fatalf("Error creating sender: %s", err)
	}

	recipient, err = types.NewMultiAddressFromHexAccountID(destination)
	if err != nil {
		log.Fatalf("Error creating recipient: %s", err)
	}
}

func main() {
	senderSS58 := subkey.SS58Encode(sender.PublicKey, network)
	recipientSS58 := subkey.SS58Encode(recipient.AsID.ToBytes(), network)
	fmt.Println("sender ss58: ", senderSS58)
	fmt.Println("recipient ss58: ", recipientSS58)

	c, err := client.NewClient(endpoint, client.WithKeyring(&sender))
	if err != nil {
		log.Fatalf("Error creating client: %s", err)
	}

	senderInfo, err := storage.GetAccountInfo(c, sender.PublicKey)
	if err != nil {
		log.Fatalf("Error getting storage: %s", err)
	}

	recipientInfo, err := storage.GetAccountInfo(c, recipient.AsID.ToBytes())
	if err != nil {
		log.Fatalf("Error getting storage: %s", err)
	}

	fmt.Printf("Recipient Free balance before: %v\n", recipientInfo.Data.Free)

	bal, ok := new(big.Int).SetString("100000000", 10)
	if !ok {
		log.Fatalf("Error converting string to big.Int")
	}

	ext := gobtext.NewTransferKeepAlive(c, recipient, bal)

	genesisHash, err := c.Api.RPC.Chain.GetBlockHash(0)
	if err != nil {
		log.Fatalf("Error getting genesis hash: %s", err)
	}

	rv, err := c.Api.RPC.State.GetRuntimeVersionLatest()
	if err != nil {
		log.Fatalf("Error getting runtime version: %s", err)
	}

	extrinsic.PayloadMutatorFns[extensions.SignedExtensionName("SubtensorSignedExtension")] = func(payload *extrinsic.Payload) {}
	extrinsic.PayloadMutatorFns[extensions.SignedExtensionName("CommitmentsSignedExtension")] = func(payload *extrinsic.Payload) {}
	err = ext.Sign(
		sender,
		c.Meta,
		extrinsic.WithEra(types.ExtrinsicEra{IsImmortalEra: true}, genesisHash),
		extrinsic.WithNonce(types.NewUCompactFromUInt(uint64(senderInfo.Nonce))),
		extrinsic.WithTip(types.NewUCompactFromUInt(0)),
		extrinsic.WithSpecVersion(rv.SpecVersion),
		extrinsic.WithTransactionVersion(rv.TransactionVersion),
		extrinsic.WithGenesisHash(genesisHash),
		extrinsic.WithMetadataMode(extensions.CheckMetadataModeDisabled, extensions.CheckMetadataHash{Hash: types.NewEmptyOption[types.H256]()}),
	)
	if err != nil {
		log.Fatalf("Error signing: %s", err)
	}

	_, err = c.Api.RPC.Author.SubmitExtrinsic(*ext)
	if err != nil {
		log.Fatalf("Error submitting extrinsic: %s", err)
	}

	fmt.Println("Extrinsic submitted")
	time.Sleep(12 * time.Second)

	recipientInfo, err = storage.GetAccountInfo(c, recipient.AsID.ToBytes())
	fmt.Println("recipient balance after transfer: ", recipientInfo.Data.Free)
}
