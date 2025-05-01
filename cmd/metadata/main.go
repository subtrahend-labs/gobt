package main

import (
	"log"

	"github.com/subtrahend-labs/gobt/internal/info"
	"github.com/subtrahend-labs/gobt/pkg/client"
)

func main() {
	c, err := client.NewClient("wss://test.finney.opentensor.ai:443")
	if err != nil {
		log.Fatalf("Error creating client: %s", err)
	}

	info.LookupExtrinsicArgs(c.Meta, "SubtensorModule", "commit_crv3_weights")
	// info.PrintModulesAndCalls(c.Meta, nil)
	// info.PrintExtensions(c.Meta, nil)
	// info.PrintExtensionDetails(c.Meta, nil, "SubtensorSignedExtension")
}
