package main

import (
	"fmt"
	"log"
	"os"

	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v4"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types/codec"
	"github.com/joho/godotenv"
)

const (
	network uint16 = 42
)

type AxonInfo struct {
	IP      types.U128
	Port    types.U32
	IPType  types.U8
	Version types.U32
}

type PrometheusInfo struct {
	IP      types.U128
	Port    types.U32
	Version types.U32
}

type NeuronInfo struct {
	Hotkey         types.AccountID
	Coldkey        types.AccountID
	UID            types.UCompact
	NetUID         types.UCompact
	Active         types.Bool
	AxonInfo       AxonInfo
	PrometheusInfo PrometheusInfo
	Stake          []struct {
		Account types.AccountID
		Amount  types.UCompact
	}
	Rank            types.UCompact
	Emission        types.UCompact
	Incentive       types.UCompact
	Consensus       types.UCompact
	Trust           types.UCompact
	ValidatorTrust  types.UCompact
	Dividends       types.UCompact
	LastUpdate      types.UCompact
	ValidatorPermit types.Bool
	Weights         []struct {
		UID    types.UCompact
		Weight types.UCompact
	}
	Bonds []struct {
		UID  types.UCompact
		Bond types.UCompact
	}
	PruningScore types.UCompact
}

func getNeurons(api *gsrpc.SubstrateAPI, netuid uint16, blockHash *types.Hash) ([]NeuronInfo, error) {
	// Expect a hex-encoded Vec<u8>
	var encodedResponse string
	err := api.Client.Call(
		&encodedResponse,
		"neuronInfo_getNeurons",
		netuid,
		*blockHash,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to call neuronInfo_getNeurons: %v", err)
	}

	// Print the raw hex response
	fmt.Printf("Raw response for netuid %d: %s\n", netuid, encodedResponse)

	// Handle empty response
	if encodedResponse == "" || encodedResponse == "0x00" {
		fmt.Printf("No neurons found for netuid %d\n", netuid)
		return []NeuronInfo{}, nil
	}

	// Decode hex string to bytes
	bytes, err := codec.HexDecodeString(encodedResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to decode hex response: %v", err)
	}

	// Decode bytes into []NeuronInfo
	var neurons []NeuronInfo
	err = codec.Decode(bytes, &neurons)
	if err != nil {
		return nil, fmt.Errorf("failed to decode neurons: %v", err)
	}

	return neurons, nil
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	endpoint := os.Getenv("ENDPOINT")

	api, err := gsrpc.NewSubstrateAPI(endpoint)
	if err != nil {
		log.Fatalf("Error creating API instance: %s", err)
	}

	// Get the latest block hash
	blockHash, err := api.RPC.Chain.GetBlockHashLatest()
	if err != nil {
		log.Fatalf("Error getting latest block hash: %s", err)
	}

	netuid := 4
	fmt.Printf("\nTesting netuid %d:\n", netuid)
	neurons, err := getNeurons(api, uint16(netuid), &blockHash)
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
