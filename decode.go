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

	blockHash, err := client.Api.RPC.Chain.GetBlockHashLatest()
	if err != nil {
		log.Fatalf("Error getting latest block hash: %s", err)
	}

	netuid := 4
	fmt.Printf("\nTesting netuid %d:\n", netuid)
	neurons, err := runtime.GetNeurons(client, uint16(netuid), &blockHash)
	if err != nil {
		log.Printf("Error fetching neurons for netuid %d: %s", netuid, err)
	}

	if len(neurons) == 0 {
		fmt.Printf("No neurons found for netuid %d\n", netuid)
	}

	for i, neuron := range neurons {
		fmt.Printf("Neuron %d (netuid %d):\n", i, netuid)
		fmt.Printf("  Hotkey: %v\n", neuron.Hotkey)
		fmt.Printf("  Coldkey: %v\n", neuron.Coldkey)
		fmt.Printf("  UID: %v\n", neuron.NetUID)
		fmt.Printf("  NetUID: %v\n", neuron.NetUID)
		fmt.Printf("  Active: %v\n", neuron.Active)
	}
}
