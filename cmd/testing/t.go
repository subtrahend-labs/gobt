// Random testing file
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/subtrahend-labs/gobt/client"
	"github.com/subtrahend-labs/gobt/runtime"
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

	res, err := runtime.GetNeurons(client, uint16(337), nil)
	if err != nil {
		log.Fatal(err.Error())
	}

	for i, v := range res {
		fmt.Printf("%d: %d\n", i, v.UID.Int64())
	}
}
