// Random testing file
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/joho/godotenv"
	"github.com/subtrahend-labs/gobt/client"
	"github.com/subtrahend-labs/gobt/storage"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	endpoint := os.Getenv("ENDPOINT")

	client, err := client.NewClient(endpoint)
	if err != nil {
		log.Fatalf("Error creating client: %s", err)
	}

	res, err := storage.GetSubnetTaoInEmission(client, types.NewU16(4), nil)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Printf("%v\n", *res)
}
