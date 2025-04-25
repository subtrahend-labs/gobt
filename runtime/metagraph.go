package runtime

import (
	"fmt"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types/codec"
	"github.com/subtrahend-labs/gobt/client"
)

// SubnetIdentityV2 represents identity information for a subnet
type SubnetIdentityV2 struct {
	SubnetName     types.Bytes
	GithubRepo     types.Bytes
	SubnetContact  types.Bytes
	SubnetURL      types.Bytes
	Discord        types.Bytes
	Description    types.Bytes
	AdditionalInfo types.Bytes
}

// ChainIdentityOfV2 represents identity information for a chain account
type ChainIdentityOfV2 struct {
	Name        types.Bytes
	URL         types.Bytes
	GithubRepo  types.Bytes
	Image       types.Bytes
	Discord     types.Bytes
	Description types.Bytes
	Additional  types.Bytes
}

// I96F32 represents a fixed-point number with 96 bits integer part and 32 bits fractional part
type I96F32 struct {
	Bits types.U128
}

// AccountAmountPair represents a tuple of account and amount
type AccountAmountPair struct {
	Account types.AccountID
	Amount  types.UCompact
}

type Metagraph struct {
	// Subnet index
	Netuid types.UCompact // Compact<u16>

	// Name and symbol
	Name   []types.UCompact // Vec<Compact<u8>>
	Symbol []types.UCompact // Vec<Compact<u8>>

	// Identity
	Identity            types.Option[SubnetIdentityV2] // Option<SubnetIdentityV2>
	NetworkRegisteredAt types.UCompact                 // Compact<u64>
	OwnerHotkey         types.AccountID                // AccountId
	OwnerColdkey        types.AccountID                // AccountId

	// Tempo terms
	Block               types.UCompact // Compact<u64>
	Tempo               types.UCompact // Compact<u16>
	LastStep            types.UCompact // Compact<u64>
	BlocksSinceLastStep types.UCompact // Compact<u64>

	// Subnet emission terms
	SubnetEmission       types.UCompact // Compact<u64>
	AlphaIn              types.UCompact // Compact<u64>
	AlphaOut             types.UCompact // Compact<u64>
	TaoIn                types.UCompact // Compact<u64>
	AlphaOutEmission     types.UCompact // Compact<u64>
	AlphaInEmission      types.UCompact // Compact<u64>
	TaoInEmission        types.UCompact // Compact<u64>
	PendingAlphaEmission types.UCompact // Compact<u64>
	PendingRootEmission  types.UCompact // Compact<u64>
	SubnetVolume         types.UCompact // Compact<u128>
	MovingPrice          I96F32         // fixed-point

	// Hparams for epoch
	Rho   types.UCompact // Compact<u16>
	Kappa types.UCompact // Compact<u16>

	// Validator params
	MinAllowedWeights types.UCompact // Compact<u16>
	MaxWeightsLimit   types.UCompact // Compact<u16>
	WeightsVersion    types.UCompact // Compact<u64>
	WeightsRateLimit  types.UCompact // Compact<u64>
	ActivityCutoff    types.UCompact // Compact<u16>
	MaxValidators     types.UCompact // Compact<u16>

	// Registration
	NumUids                types.UCompact // Compact<u16>
	MaxUids                types.UCompact // Compact<u16>
	Burn                   types.UCompact // Compact<u64>
	Difficulty             types.UCompact // Compact<u64>
	RegistrationAllowed    types.Bool     // bool
	PowRegistrationAllowed types.Bool     // bool
	ImmunityPeriod         types.UCompact // Compact<u16>
	MinDifficulty          types.UCompact // Compact<u64>
	MaxDifficulty          types.UCompact // Compact<u64>
	MinBurn                types.UCompact // Compact<u64>
	MaxBurn                types.UCompact // Compact<u64>
	AdjustmentAlpha        types.UCompact // Compact<u64>
	AdjustmentInterval     types.UCompact // Compact<u16>
	TargetRegsPerInterval  types.UCompact // Compact<u16>
	MaxRegsPerBlock        types.UCompact // Compact<u16>
	ServingRateLimit       types.UCompact // Compact<u64>

	// Commit‐reveal
	CommitRevealWeightsEnabled types.Bool     // bool
	CommitRevealPeriod         types.UCompact // Compact<u64>

	// Bonds
	LiquidAlphaEnabled types.Bool     // bool
	AlphaHigh          types.UCompact // Compact<u16>
	AlphaLow           types.UCompact // Compact<u16>
	BondsMovingAvg     types.UCompact // Compact<u64>

	// Metagraph info
	Hotkeys         []types.AccountID // Vec<AccountId>
	Coldkeys        []types.AccountID
	Identities      []types.Option[ChainIdentityOfV2]
	Axons           []AxonInfo
	Active          []types.Bool
	ValidatorPermit []types.Bool

	PruningScore []types.UCompact // Vec<Compact<u16>>
	LastUpdate   []types.UCompact // Vec<Compact<u64>>
	Emission     []types.UCompact // Vec<Compact<u64>>
	Dividends    []types.UCompact // Vec<Compact<u16>>
	Incentives   []types.UCompact // Vec<Compact<u16>>
	Consensus    []types.UCompact // Vec<Compact<u16>>
	Trust        []types.UCompact // Vec<Compact<u16>>
	Rank         []types.UCompact // Vec<Compact<u16>>

	BlockAtRegistration []types.UCompact // Vec<Compact<u64>>
	AlphaStake          []types.UCompact // Vec<Compact<u64>>
	TaoStake            []types.UCompact // Vec<Compact<u64>>
	TotalStake          []types.UCompact // Vec<Compact<u64>>

	// Dividend breakdown
	TaoDividendsPerHotkey   []AccountAmountPair // Vec<(AccountId, Compact<u64>)>
	AlphaDividendsPerHotkey []AccountAmountPair
}

// GetMetagraph retrieves the metagraph for a specific subnet
func GetMetagraph(c *client.Client, netuid uint16, blockHash *types.Hash) (*Metagraph, error) {
	// First, try to see what's being returned from the API call
	var encodedResponse []byte
	err := c.Api.Client.Call(
		&encodedResponse,
		"subnetInfo_getMetagraph",
		netuid,
		blockHash,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to call subnetInfo_getMetagraph: %v", err)
	}

	if len(encodedResponse) == 0 {
		return nil, fmt.Errorf("no metagraph found for netuid %d", netuid)
	}

	// Try to decode as a Vec<u8> first
	var m Metagraph
	if err := codec.Decode(encodedResponse, &m); err != nil {
		return nil, fmt.Errorf("failed to decode metagraph: %v", err)
	}
	return &m, nil
}
