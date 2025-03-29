package testutils

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types/extrinsic"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types/extrinsic/extensions"
	"github.com/docker/go-connections/nat"
	"github.com/stretchr/testify/require"
	"github.com/subtrahend-labs/gobt/client"
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

	genesisHash, err := cl.Api.RPC.Chain.GetBlockHash(0)
	require.NoError(t, err, "Failed to get genesis hash")

	rv, err := cl.Api.RPC.State.GetRuntimeVersionLatest()
	require.NoError(t, err, "Failed to get runtime version")

	// Register custom extension mutators
	extrinsic.PayloadMutatorFns[extensions.SignedExtensionName("SubtensorSignedExtension")] = func(payload *extrinsic.Payload) {}
	extrinsic.PayloadMutatorFns[extensions.SignedExtensionName("CommitmentsSignedExtension")] = func(payload *extrinsic.Payload) {}

	// Sign the extrinsic
	err = ext.Sign(
		signer,
		cl.Meta,
		extrinsic.WithEra(types.ExtrinsicEra{IsImmortalEra: true}, genesisHash),
		extrinsic.WithNonce(types.NewUCompactFromUInt(uint64(nonce))),
		extrinsic.WithTip(types.NewUCompactFromUInt(0)),
		extrinsic.WithSpecVersion(rv.SpecVersion),
		extrinsic.WithTransactionVersion(rv.TransactionVersion),
		extrinsic.WithGenesisHash(genesisHash),
		extrinsic.WithMetadataMode(extensions.CheckMetadataModeDisabled, extensions.CheckMetadataHash{Hash: types.NewEmptyOption[types.H256]()}),
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
		if status.IsDropped || status.IsInvalid {
			t.Fatalf("Transaction failed: dropped=%v, invalid=%v", status.IsDropped, status.IsInvalid)
		}
	}
	txnSub.Unsubscribe()

	return blockHash
}
