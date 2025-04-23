package testutils

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/centrifuge/go-substrate-rpc-client/v4/registry/retriever"
	"github.com/centrifuge/go-substrate-rpc-client/v4/registry/state"
	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types/extrinsic"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types/extrinsic/extensions"
	"github.com/docker/go-connections/nat"
	"github.com/stretchr/testify/require"
	"github.com/subtrahend-labs/gobt/client"
	"github.com/subtrahend-labs/gobt/sigtools"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type TestEnv struct {
	Container testcontainers.Container
	Client    *client.Client
}

func Setup() (*TestEnv, error) {
	ctx := context.Background()

	// Define container request
	nodePort := "9944/tcp"
	req := testcontainers.ContainerRequest{
		Image:        "subtensor-local:latest",
		ExposedPorts: []string{nodePort},
		Cmd: []string{
			"/bin/bash",
			"-c",
			"node-subtensor --dev --rpc-external --rpc-cors all --rpc-methods=unsafe --offchain-worker never",
		},
		WaitingFor: wait.ForAll(
			wait.ForListeningPort(nat.Port(nodePort)),
			wait.ForLog("Running JSON-RPC server").WithStartupTimeout(30*time.Second),
		),
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	// Start the container
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start container: %v", err)
	}

	// Get mapped port and host
	mappedPort, err := container.MappedPort(ctx, nat.Port(nodePort))
	if err != nil {
		return nil, fmt.Errorf("failed to get mapped port: %v", err)
	}
	host, err := container.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get host: %v", err)
	}

	wsURL := fmt.Sprintf("ws://%s:%s", host, mappedPort.Port())

	time.Sleep(3 * time.Second)

	keyringAlice := signature.TestKeyringPairAlice
	cl, err := client.NewClient(wsURL, client.WithKeyring(&keyringAlice))
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %v", err)
	}

	return &TestEnv{
		Container: container,
		Client:    cl,
	}, nil
}

func (env *TestEnv) Teardown() {
	ctx := context.Background()
	env.Container.Terminate(ctx)
}

func SignAndSubmit(t *testing.T, cl *client.Client, ext *extrinsic.Extrinsic, signer signature.KeyringPair, nonce uint32) types.Hash {
	t.Helper()

	// Register custom extension mutators
	extrinsic.PayloadMutatorFns[extensions.SignedExtensionName("SubtensorSignedExtension")] = func(payload *extrinsic.Payload) {}
	extrinsic.PayloadMutatorFns[extensions.SignedExtensionName("CommitmentsSignedExtension")] = func(payload *extrinsic.Payload) {}

	tip := types.NewUCompactFromUInt(0)
	n := types.NewUCompactFromUInt(uint64(nonce))
	sc := sigtools.NewSigningContext(&tip, &n)
	ops, err := sigtools.CreateSigningOptions(cl, signer, sc)

	// Sign the extrinsic
	err = ext.Sign(
		signer,
		cl.Meta,
		ops...,
	)
	require.NoError(t, err, "Failed to sign extrinsic")

	txnSub, err := cl.Api.RPC.Author.SubmitAndWatchExtrinsic(*ext)
	require.NoError(t, err, "Failed to submit extrinsic")

	var blockHash types.Hash
	for {
		status := <-txnSub.Chan()
		t.Logf("Transaction status: %v", status)
		if status.IsInBlock {
			blockHash = status.AsInBlock
			t.Logf("Transaction included in block: %v", blockHash)
			break
		}
		if status.IsFinalized {
			t.Fatalf("Transaction shouldn't be finalized yet: %v", status.AsFinalized)
		}
		if status.IsDropped || status.IsInvalid || status.IsRetracted {
			t.Fatalf("Transaction failed: dropped=%v, invalid=%v, retracted=%v", status.IsDropped, status.IsInvalid, status.IsRetracted)
		}
	}
	txnSub.Unsubscribe()
	evtr, err := retriever.NewDefaultEventRetriever(
		state.NewEventProvider(cl.Api.RPC.State),
		cl.Api.RPC.State,
	)
	require.NoError(t, err)

	events, err := evtr.GetEvents(blockHash)
	require.NoError(t, err)

	if derr := ExtractDispatchError(*cl.Meta, events); derr != nil {
		t.Fatalf("extrinsic dispatch failed: %v", derr)
	}

	return blockHash

}
