package main

import (
	"log"

	"github.com/subtrahend-labs/gobt/client"
	"github.com/subtrahend-labs/gobt/internal/info"
)

func main() {
	c, err := client.NewClient("wss://test.finney.opentensor.ai:443")
	if err != nil {
		log.Fatalf("Error creating client: %s", err)
	}

	info.PrintModulesAndCalls(c.Meta, nil)
}
