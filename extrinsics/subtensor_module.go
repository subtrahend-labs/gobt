package extrinsics

import (
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types/extrinsic"
	"github.com/subtrahend-labs/gobt/client"
)

// #### Module: SubtensorModule (Index: 7)
//     - [x] set_weights (Index: 0)
//     - [x] serve_axon (Index: 4)
//     - [x] serve_axon_tls (Index: 40)
//     - [x] register_network (Index: 59)
//     - [x] root_register (Index: 62)
//     - [x] register (Index: 6)
//     - [x] burned_register (Index: 7)

//     - [o] batch_set_weights (Index: 80)
//     - [o] commit_weights (Index: 96)
//     - [o] batch_commit_weights (Index: 100)
//     - [o] reveal_weights (Index: 97)

//     - [o] commit_crv3_weights (Index: 99)

//     - [o] batch_reveal_weights (Index: 98)
//     - [o] set_tao_weights (Index: 8)
//     - [o] become_delegate (Index: 1)
//     - [o] decrease_take (Index: 65)
//     - [o] increase_take (Index: 66)
//     - [o] add_stake (Index: 2)
//     - [o] remove_stake (Index: 3)
//     - [o] serve_prometheus (Index: 5)

//     - [o] adjust_senate (Index: 63)
//     - [o] swap_hotkey (Index: 70)
//     - [o] swap_coldkey (Index: 71)
//     - [o] set_childkey_take (Index: 75)
//     - [o] sudo_set_tx_childkey_take_rate_limit (Index: 69)
//     - [o] sudo_set_min_childkey_take (Index: 76)
//     - [o] sudo_set_max_childkey_take (Index: 77)
//     - [o] sudo (Index: 51)
//     - [o] sudo_unchecked_weight (Index: 52)
//     - [o] vote (Index: 55)
//     - [o] faucet (Index: 60)
//     - [o] dissolve_network (Index: 61)
//     - [o] set_children (Index: 67)
//     - [o] schedule_swap_coldkey (Index: 73)
//     - [o] schedule_dissolve_network (Index: 74)
//     - [o] set_identity (Index: 68)
//     - [o] set_subnet_identity (Index: 78)
//     - [o] register_network_with_identity (Index: 79)
//     - [o] unstake_all (Index: 83)
//     - [o] unstake_all_alpha (Index: 84)
//     - [o] move_stake (Index: 85)
//     - [o] transfer_stake (Index: 86)
//     - [o] swap_stake (Index: 87)
//     - [x] add_stake_limit (Index: 88)
//     - [o] remove_stake_limit (Index: 89)
//     - [o] swap_stake_limit (Index: 90)
//     - [o] try_associate_hotkey (Index: 91)

type SubnetIdentityV2 struct {
	SubnetName    types.Bytes
	GithubRepo    types.Bytes
	SubnetContact types.Bytes
	SubnetURL     types.Bytes
	Discord       types.Bytes
	Description   types.Bytes
	Additional    types.Bytes
}

// writing code for add stake here (first time)
// do i need comments for future reference?
func AddStakeCall(c *client.Client, hotkey types.AccountID, netuid types.U16, amount_staked types.U64) (types.Call, error) {
	call, err := types.NewCall(c.Meta, "SubtensorModule.add_stake", hotkey, netuid, amount_staked)
	if err != nil {
		return types.Call{}, err
	}
	return call, nil
}

func AddStakeExt(c *client.Client, hotkey types.AccountID, netuid types.U16, amount_staked types.U64) (*extrinsic.Extrinsic, error) {
	call, err := AddStakeCall(c, hotkey, netuid, amount_staked)
	if err != nil {
		return nil, err
	}
	ext := extrinsic.NewExtrinsic(call)
	return &ext, nil
}

// increase take
func IncreaseTakeCall(c *client.Client, hotkey types.AccountID, take types.U16) (types.Call, error) {
	call, err := types.NewCall(c.Meta, "SubtensorModule.increase_take", hotkey, take)
	if err != nil {
		return types.Call{}, err
	}
	return call, nil
}

func IncreaseTakeExt(c *client.Client, hotkey types.AccountID, take types.U16) (*extrinsic.Extrinsic, error) {
	call, err := IncreaseTakeCall(c, hotkey, take)
	if err != nil {
		return nil, err
	}
	ext := extrinsic.NewExtrinsic(call)
	return &ext, nil
}

// decrease take
func DecreaseTakeCall(c *client.Client, hotkey types.AccountID, take types.U16) (types.Call, error) {
	call, err := types.NewCall(c.Meta, "SubtensorModule.decrease_take", hotkey, take)
	if err != nil {
		return types.Call{}, err
	}
	return call, nil
}

func DecreaseTakeExt(c *client.Client, hotkey types.AccountID, take types.U16) (*extrinsic.Extrinsic, error) {
	call, err := DecreaseTakeCall(c, hotkey, take)
	if err != nil {
		return nil, err
	}
	ext := extrinsic.NewExtrinsic(call)
	return &ext, nil
}

// remove stake
func RemoveStakeCall(c *client.Client, hotkey types.AccountID, netuid types.U16, amount_unstaked types.U64) (types.Call, error) {
	call, err := types.NewCall(c.Meta, "SubtensorModule.remove_stake", hotkey, netuid, amount_unstaked)
	if err != nil {
		return types.Call{}, err
	}
	return call, nil
}

func RemoveStakeExt(c *client.Client, hotkey types.AccountID, netuid types.U16, amount_unstaked types.U64) (*extrinsic.Extrinsic, error) {
	call, err := RemoveStakeCall(c, hotkey, netuid, amount_unstaked)
	if err != nil {
		return nil, err
	}
	ext := extrinsic.NewExtrinsic(call)
	return &ext, nil
}

// SetChildkeyTakeCall creates a call to set the childkey take value
func SetChildkeyTakeCall(c *client.Client, coldkey types.AccountID, hotkey types.AccountID, netuid types.U16, take types.U16) (types.Call, error) {
	call, err := types.NewCall(
		c.Meta,
		"SubtensorModule.set_childkey_take",
		coldkey,
		hotkey,
		netuid,
		take,
	)
	if err != nil {
		return types.Call{}, err
	}
	return call, nil
}

// SetChildkeyTakeExt creates an extrinsic to set the childkey take value
func SetChildkeyTakeExt(c *client.Client, coldkey types.AccountID, hotkey types.AccountID, netuid types.U16, take types.U16) (*extrinsic.Extrinsic, error) {
	call, err := SetChildkeyTakeCall(c, coldkey, hotkey, netuid, take)
	if err != nil {
		return nil, err
	}
	ext := extrinsic.NewExtrinsic(call)
	return &ext, nil
}

// move stake
// MoveStakeCall creates a call to move stake from one hotkey to another across subnets
func MoveStakeCall(c *client.Client,
	origin_hotkey types.AccountID,
	destination_hotkey types.AccountID,
	origin_netuid types.U16,
	destination_netuid types.U16,
	alpha_amount types.U64,
) (types.Call, error) {
	call, err := types.NewCall(
		c.Meta,
		"SubtensorModule.move_stake",
		origin_hotkey,
		destination_hotkey,
		origin_netuid,
		destination_netuid,
		alpha_amount,
	)
	if err != nil {
		return types.Call{}, err
	}
	return call, nil
}

// MoveStakeExt creates an extrinsic to move stake from one hotkey to another across subnets
func MoveStakeExt(
	c *client.Client,
	origin_hotkey types.AccountID,
	destination_hotkey types.AccountID,
	origin_netuid types.U16,
	destination_netuid types.U16,
	alpha_amount types.U64,
) (*extrinsic.Extrinsic, error) {
	call, err := MoveStakeCall(
		c,
		origin_hotkey,
		destination_hotkey,
		origin_netuid,
		destination_netuid,
		alpha_amount,
	)
	if err != nil {
		return nil, err
	}
	ext := extrinsic.NewExtrinsic(call)
	return &ext, nil
}

// swap stake
func SwapStakeCall(
	c *client.Client,
	hotkey types.AccountID,
	origin_netuid types.U16,
	destination_netuid types.U16,
	alpha_amount types.U64,
) (types.Call, error) {
	call, err := types.NewCall(
		c.Meta,
		"SubtensorModule.swap_stake",
		hotkey,
		origin_netuid,
		destination_netuid,
		alpha_amount,
	)
	if err != nil {
		return types.Call{}, err
	}
	return call, nil
}

func SwapStakeExt(
	c *client.Client,
	hotkey types.AccountID,
	origin_netuid types.U16,
	destination_netuid types.U16,
	alpha_amount types.U64,
) (*extrinsic.Extrinsic, error) {
	call, err := SwapStakeCall(
		c,
		hotkey,
		origin_netuid,
		destination_netuid,
		alpha_amount,
	)
	if err != nil {
		return nil, err
	}
	ext := extrinsic.NewExtrinsic(call)
	return &ext, nil
}

// transfer stake
// TransferStakeCall creates a call to transfer stake between coldkeys while keeping the same hotkey
func TransferStakeCall(
	c *client.Client,
	destination_coldkey types.AccountID,
	hotkey types.AccountID,
	origin_netuid types.U16,
	destination_netuid types.U16,
	alpha_amount types.U64,
) (types.Call, error) {
	call, err := types.NewCall(
		c.Meta,
		"SubtensorModule.transfer_stake",
		destination_coldkey,
		hotkey,
		origin_netuid,
		destination_netuid,
		alpha_amount,
	)
	if err != nil {
		return types.Call{}, err
	}
	return call, nil
}

// TransferStakeExt creates an extrinsic to transfer stake between coldkeys while keeping the same hotkey
func TransferStakeExt(
	c *client.Client,
	destination_coldkey types.AccountID,
	hotkey types.AccountID,
	origin_netuid types.U16,
	destination_netuid types.U16,
	alpha_amount types.U64,
) (*extrinsic.Extrinsic, error) {
	call, err := TransferStakeCall(
		c,
		destination_coldkey,
		hotkey,
		origin_netuid,
		destination_netuid,
		alpha_amount,
	)
	if err != nil {
		return nil, err
	}
	ext := extrinsic.NewExtrinsic(call)
	return &ext, nil
}

// swap stake limit
// SwapStakeLimitCall creates a call to swap stake between subnets with a limit price
func SwapStakeLimitCall(
	c *client.Client,
	hotkey types.AccountID,
	origin_netuid types.U16,
	destination_netuid types.U16,
	alpha_amount types.U64,
	limit_price types.U64,
	allow_partial types.Bool,
) (types.Call, error) {
	call, err := types.NewCall(
		c.Meta,
		"SubtensorModule.swap_stake_limit",
		hotkey,
		origin_netuid,
		destination_netuid,
		alpha_amount,
		limit_price,
		allow_partial,
	)
	if err != nil {
		return types.Call{}, err
	}
	return call, nil
}

// SwapStakeLimitExt creates an extrinsic to swap stake between subnets with a limit price
func SwapStakeLimitExt(
	c *client.Client,
	hotkey types.AccountID,
	origin_netuid types.U16,
	destination_netuid types.U16,
	alpha_amount types.U64,
	limit_price types.U64,
	allow_partial types.Bool,
) (*extrinsic.Extrinsic, error) {
	call, err := SwapStakeLimitCall(
		c,
		hotkey,
		origin_netuid,
		destination_netuid,
		alpha_amount,
		limit_price,
		allow_partial,
	)
	if err != nil {
		return nil, err
	}
	ext := extrinsic.NewExtrinsic(call)
	return &ext, nil
}

// unstake all
func UnstakeAllCall(c *client.Client, hotkey types.AccountID) (types.Call, error) {
	call, err := types.NewCall(
		c.Meta,
		"SubtensorModule.unstake_all",
		hotkey,
	)
	if err != nil {
		return types.Call{}, err
	}
	return call, nil
}

func UnstakeAllExt(c *client.Client, hotkey types.AccountID) (*extrinsic.Extrinsic, error) {
	call, err := UnstakeAllCall(
		c,
		hotkey,
	)
	if err != nil {
		return nil, err
	}
	ext := extrinsic.NewExtrinsic(call)
	return &ext, nil
}

// unstake all alpha
func UnstakeAllAlphaCall(c *client.Client, hotkey types.AccountID) (types.Call, error) {
	call, err := types.NewCall(
		c.Meta,
		"SubtensorModule.unstake_all_alpha",
		hotkey,
	)
	if err != nil {
		return types.Call{}, err
	}
	return call, nil
}

func UnstakeAllAlphaExt(c *client.Client, hotkey types.AccountID) (*extrinsic.Extrinsic, error) {
	call, err := UnstakeAllAlphaCall(
		c,
		hotkey,
	)
	if err != nil {
		return nil, err
	}
	ext := extrinsic.NewExtrinsic(call)
	return &ext, nil
}

// try associate hotkey
func TryAssociateHotkeyCall(c *client.Client, coldkey types.AccountID, hotkey types.AccountID) (types.Call, error) {
	call, err := types.NewCall(
		c.Meta,
		"SubtensorModule.try_associate_hotkey",
		coldkey,
		hotkey,
	)
	if err != nil {
		return types.Call{}, err
	}
	return call, nil
}

func TryAssociateHotkeyExt(c *client.Client, coldkey types.AccountID, hotkey types.AccountID) (*extrinsic.Extrinsic, error) {
	call, err := TryAssociateHotkeyCall(
		c,
		coldkey,
		hotkey,
	)
	if err != nil {
		return nil, err
	}
	ext := extrinsic.NewExtrinsic(call)
	return &ext, nil
}

// swap coldkey
func SwapColdkeyCall(c *client.Client, old_coldkey types.AccountID, new_coldkey types.AccountID, swap_cost types.U64) (types.Call, error) {
	call, err := types.NewCall(
		c.Meta,
		"SubtensorModule.swap_coldkey",
		old_coldkey,
		new_coldkey,
		swap_cost,
	)
	if err != nil {
		return types.Call{}, err
	}
	return call, nil
}

func SwapColdkeyExt(c *client.Client, old_coldkey types.AccountID, new_coldkey types.AccountID, swap_cost types.U64) (*extrinsic.Extrinsic, error) {
	call, err := SwapColdkeyCall(c, old_coldkey, new_coldkey, swap_cost)
	if err != nil {
		return nil, err
	}
	ext := extrinsic.NewExtrinsic(call)
	return &ext, nil
}

// swap hotkey
func SwapHotkeyCall(c *client.Client, old_hotkey types.AccountID, new_hotkey types.AccountID) (types.Call, error) {
	call, err := types.NewCall(
		c.Meta,
		"SubtensorModule.swap_coldkey",
		old_hotkey,
		new_hotkey,
	)
	if err != nil {
		return types.Call{}, err
	}
	return call, nil
}

func SwapHotkeyExt(c *client.Client, old_hotkey types.AccountID, new_hotkey types.AccountID) (*extrinsic.Extrinsic, error) {
	call, err := SwapHotkeyCall(c, old_hotkey, new_hotkey)
	if err != nil {
		return nil, err
	}
	ext := extrinsic.NewExtrinsic(call)
	return &ext, nil
}

// adjust senate
func AdjustSenateCall(c *client.Client, hotkey types.AccountID) (types.Call, error) {
	call, err := types.NewCall(
		c.Meta,
		"SubtensorModule.adjust_senate",
		hotkey,
	)
	if err != nil {
		return types.Call{}, err
	}
	return call, nil
}

func AdjustSenateExt(c *client.Client, hotkey types.AccountID) (*extrinsic.Extrinsic, error) {
	call, err := AdjustSenateCall(c, hotkey)
	if err != nil {
		return nil, err
	}
	ext := extrinsic.NewExtrinsic(call)
	return &ext, nil
}

// sudo_set_tx_childkey_take_rate_limit
func SudoSetTxChildkeyTakeRateLimitCall(c *client.Client, tx_rate_limit types.U64) (types.Call, error) {
	call, err := types.NewCall(
		c.Meta,
		"SubtensorModule.sudo_set_tx_childkey_take_rate_limit",
		tx_rate_limit,
	)
	if err != nil {
		return types.Call{}, err
	}
	return call, nil
}

func SudoSetTxChildkeyTakeRateLimitExt(c *client.Client, tx_rate_limit types.U64) (*extrinsic.Extrinsic, error) {
	call, err := SudoSetTxChildkeyTakeRateLimitCall(c, tx_rate_limit)
	if err != nil {
		return nil, err
	}
	ext := extrinsic.NewExtrinsic(call)
	return &ext, nil
}

// sudo_set_min_childkey_take (Index: 76)
func SudoSetMinChildkeyTakeCall(c *client.Client, take types.U16) (types.Call, error) {
	call, err := types.NewCall(
		c.Meta,
		"SubtensorModule.sudo_set_min_childkey_take",
		take,
	)
	if err != nil {
		return types.Call{}, err
	}
	return call, nil
}

func SudoSetMinChildkeyTakeExt(c *client.Client, take types.U16) (*extrinsic.Extrinsic, error) {
	call, err := SudoSetMinChildkeyTakeCall(c, take)
	if err != nil {
		return nil, err
	}
	ext := extrinsic.NewExtrinsic(call)
	return &ext, nil
}

// - [ ] sudo_set_max_childkey_take (Index: 77)
func SudoSetMaxChildkeyTakeCall(c *client.Client, take types.U16) (types.Call, error) {
	call, err := types.NewCall(
		c.Meta,
		"SubtensorModule.sudo_set_max_childkey_take",
		take,
	)
	if err != nil {
		return types.Call{}, err
	}
	return call, nil
}

func SudoSetMaxChildkeyTakeExt(c *client.Client, take types.U16) (*extrinsic.Extrinsic, error) {
	call, err := SudoSetMaxChildkeyTakeCall(c, take)
	if err != nil {
		return nil, err
	}
	ext := extrinsic.NewExtrinsic(call)
	return &ext, nil
}

//     - [ ] sudo (Index: 51)
// what does type box mean?

//     - [ ] sudo_unchecked_weight (Index: 52)

// - [ ] vote (Index: 55)
func VoteCall(c *client.Client, proposal types.Hash, index types.U32, approve types.Bool) (types.Call, error) {
	call, err := types.NewCall(
		c.Meta,
		"SubtensorModule.vote",
		proposal,
		index,
		approve,
	)
	if err != nil {
		return types.Call{}, err
	}
	return call, nil
}

func VoteExt(c *client.Client, proposal types.Hash, index types.U32, approve types.Bool) (*extrinsic.Extrinsic, error) {
	call, err := VoteCall(c, proposal, index, approve)
	if err != nil {
		return nil, err
	}
	ext := extrinsic.NewExtrinsic(call)
	return &ext, nil
}

// - [ ] faucet (Index: 60)
func FaucetCall(c *client.Client, blockNumber types.U64, nonce types.U64, work types.Bytes) (types.Call, error) {
	call, err := types.NewCall(
		c.Meta,
		"SubtensorModule.faucet",
		blockNumber,
		nonce,
		work,
	)
	if err != nil {
		return types.Call{}, err
	}
	return call, nil
}

func FaucetExt(c *client.Client, blockNumber types.U64, nonce types.U64, work types.Bytes) (*extrinsic.Extrinsic, error) { // not sure if this should be type bytes or smth else (check later)
	call, err := FaucetCall(c, blockNumber, nonce, work)
	if err != nil {
		return nil, err
	}
	ext := extrinsic.NewExtrinsic(call)
	return &ext, nil
}

// - [ ] dissolve_network (Index: 61)
func DissolveNetworkCall(c *client.Client, coldkey types.AccountID, netuid types.U16) (types.Call, error) {
	call, err := types.NewCall(
		c.Meta,
		"SubtensorModule.dissolve_network",
		coldkey,
		netuid,
	)
	if err != nil {
		return types.Call{}, err
	}
	return call, nil
}

func DissolveNetworkExt(c *client.Client, coldkey types.AccountID, netuid types.U16) (*extrinsic.Extrinsic, error) {
	call, err := DissolveNetworkCall(c, coldkey, netuid)
	if err != nil {
		return nil, err
	}
	ext := extrinsic.NewExtrinsic(call)
	return &ext, nil
}

//   - [ ] set_children (Index: 67)
//
// Define the tuple struct
type ChildTuple struct {
	Stake   types.U64
	Account types.AccountID
}

// Update the function signatures
func SetChildrenCall(c *client.Client, coldkey types.AccountID, hotkey types.AccountID, netuid types.U16, children []ChildTuple) (types.Call, error) { // is this waht you do instead of a         children: Vec<(u64, T::AccountId)>,
	call, err := types.NewCall(
		c.Meta,
		"SubtensorModule.set_children",
		coldkey,
		hotkey,
		netuid,
		children,
	)
	if err != nil {
		return types.Call{}, err
	}
	return call, nil
}

func SetChildrenExt(c *client.Client, coldkey types.AccountID, hotkey types.AccountID, netuid types.U16, children []ChildTuple) (*extrinsic.Extrinsic, error) {
	call, err := SetChildrenCall(c, coldkey, hotkey, netuid, children)
	if err != nil {
		return nil, err
	}
	ext := extrinsic.NewExtrinsic(call)
	return &ext, nil
}

// - [ ] schedule_swap_coldkey (Index: 73)
func ScheduleSwapColdkeyCall(c *client.Client, new_coldkey types.AccountID) (types.Call, error) {
	call, err := types.NewCall(
		c.Meta,
		"SubtensorModule.schedule_swap_coldkey",
		new_coldkey,
	)
	if err != nil {
		return types.Call{}, err
	}
	return call, nil
}

func ScheduleSwapColdkeyExt(c *client.Client, new_coldkey types.AccountID) (*extrinsic.Extrinsic, error) {
	call, err := ScheduleSwapColdkeyCall(c, new_coldkey)
	if err != nil {
		return nil, err
	}
	ext := extrinsic.NewExtrinsic(call)
	return &ext, nil
}

// - [ ] schedule_dissolve_network (Index: 74)
func ScheduleDissolveNetworkCall(c *client.Client, netuid types.U16) (types.Call, error) {
	call, err := types.NewCall(
		c.Meta,
		"SubtensorModule.schedule_dissolve_network",
		netuid,
	)
	if err != nil {
		return types.Call{}, err
	}
	return call, nil
}

func ScheduleDissolveNetworkExt(c *client.Client, netuid types.U16) (*extrinsic.Extrinsic, error) {
	call, err := ScheduleDissolveNetworkCall(c, netuid)
	if err != nil {
		return nil, err
	}
	ext := extrinsic.NewExtrinsic(call)
	return &ext, nil
}

// - [ ] set_identity (Index: 68)
func SetIdentityCall(
	c *client.Client,
	name types.Bytes,
	url types.Bytes,
	githubRepo types.Bytes,
	image types.Bytes,
	discord types.Bytes,
	description types.Bytes,
	additional types.Bytes,
) (types.Call, error) {
	call, err := types.NewCall(
		c.Meta,
		"SubtensorModule.set_identity",
		name,
		url,
		githubRepo,
		image,
		discord,
		description,
		additional,
	)
	if err != nil {
		return types.Call{}, err
	}
	return call, nil
}

func SetIdentityExt(
	c *client.Client,
	name types.Bytes,
	url types.Bytes,
	githubRepo types.Bytes,
	image types.Bytes,
	discord types.Bytes,
	description types.Bytes,
	additional types.Bytes,
) (*extrinsic.Extrinsic, error) {
	call, err := SetIdentityCall(
		c,
		name,
		url,
		githubRepo,
		image,
		discord,
		description,
		additional,
	)
	if err != nil {
		return nil, err
	}
	ext := extrinsic.NewExtrinsic(call)
	return &ext, nil
}

//   - [ ] set_subnet_identity (Index: 78)
//
// SetSubnetIdentityCall creates a call to set identity information for a subnet
func SetSubnetIdentityCall(
	c *client.Client,
	netuid types.U16,
	subnetName types.Bytes,
	githubRepo types.Bytes,
	subnetContact types.Bytes,
	subnetURL types.Bytes,
	discord types.Bytes,
	description types.Bytes,
	additional types.Bytes,
) (types.Call, error) {
	call, err := types.NewCall(
		c.Meta,
		"SubtensorModule.set_subnet_identity",
		netuid,
		subnetName,
		githubRepo,
		subnetContact,
		subnetURL,
		discord,
		description,
		additional,
	)
	if err != nil {
		return types.Call{}, err
	}
	return call, nil
}

// SetSubnetIdentityExt creates an extrinsic to set identity information for a subnet
func SetSubnetIdentityExt(
	c *client.Client,
	netuid types.U16,
	subnetName types.Bytes,
	githubRepo types.Bytes,
	subnetContact types.Bytes,
	subnetURL types.Bytes,
	discord types.Bytes,
	description types.Bytes,
	additional types.Bytes,
) (*extrinsic.Extrinsic, error) {
	call, err := SetSubnetIdentityCall(
		c,
		netuid,
		subnetName,
		githubRepo,
		subnetContact,
		subnetURL,
		discord,
		description,
		additional,
	)
	if err != nil {
		return nil, err
	}
	ext := extrinsic.NewExtrinsic(call)
	return &ext, nil
}

// - [ ] register_network_with_identity (Index: 79)
func RegisterNetworkWithIdentityCall(
	c *client.Client,
	hotkey types.AccountID,
	subnetName types.Bytes,
	githubRepo types.Bytes,
	subnetContact types.Bytes,
	subnetURL types.Bytes,
	discord types.Bytes,
	description types.Bytes,
	additional types.Bytes,
) (types.Call, error) {
	call, err := types.NewCall(
		c.Meta,
		"SubtensorModule.register_network_with_identity",
		hotkey,
		subnetName,
		githubRepo,
		subnetContact,
		subnetURL,
		discord,
		description,
		additional,
	)
	if err != nil {
		return types.Call{}, err
	}
	return call, nil
}

func RegisterNetworkWithIdentityExt(
	c *client.Client,
	hotkey types.AccountID,
	subnetName types.Bytes,
	githubRepo types.Bytes,
	subnetContact types.Bytes,
	subnetURL types.Bytes,
	discord types.Bytes,
	description types.Bytes,
	additional types.Bytes,
) (*extrinsic.Extrinsic, error) {
	call, err := RegisterNetworkWithIdentityCall(
		c,
		hotkey,
		subnetName,
		githubRepo,
		subnetContact,
		subnetURL,
		discord,
		description,
		additional,
	)
	if err != nil {
		return nil, err
	}
	ext := extrinsic.NewExtrinsic(call)
	return &ext, nil
}

//   - [ ] batch_set_weights (Index: 80)
//
// WeightTuple represents a (uid, weight) pair for batch weight setting
type WeightTuple struct {
	UID    types.U16
	Weight types.U16
}

// BatchSetWeightsCall creates a call to batch set weights for multiple networks
func BatchSetWeightsCall(
	c *client.Client,
	netuids []types.U16,
	weights [][]WeightTuple,
	versionKeys []types.U64,
) (types.Call, error) {
	call, err := types.NewCall(
		c.Meta,
		"SubtensorModule.batch_set_weights",
		netuids,
		weights,
		versionKeys,
	)
	if err != nil {
		return types.Call{}, err
	}
	return call, nil
}

// BatchSetWeightsExt creates an extrinsic to batch set weights for multiple networks
func BatchSetWeightsExt(
	c *client.Client,
	netuids []types.U16,
	weights [][]WeightTuple,
	versionKeys []types.U64,
) (*extrinsic.Extrinsic, error) {
	call, err := BatchSetWeightsCall(
		c,
		netuids,
		weights,
		versionKeys,
	)
	if err != nil {
		return nil, err
	}
	ext := extrinsic.NewExtrinsic(call)
	return &ext, nil
}

// - [ ] commit_weights (Index: 96)
func CommitWeightsCall(
	c *client.Client,
	netuid types.U16,
	commitHash types.Hash,
) (types.Call, error) {
	call, err := types.NewCall(
		c.Meta,
		"SubtensorModule.commit_weights",
		netuid,
		commitHash,
	)
	if err != nil {
		return types.Call{}, err
	}
	return call, nil
}

func CommitWeightsExt(
	c *client.Client,
	netuid types.U16,
	commitHash types.Hash,
) (*extrinsic.Extrinsic, error) {
	call, err := CommitWeightsCall(
		c,
		netuid,
		commitHash,
	)
	if err != nil {
		return nil, err
	}
	ext := extrinsic.NewExtrinsic(call)
	return &ext, nil
}

// - [ ] batch_commit_weights (Index: 100)
func BatchCommitWeightsCall(
	c *client.Client,
	netuids []types.U16,
	commitHashes []types.Hash,
) (types.Call, error) {
	call, err := types.NewCall(
		c.Meta,
		"SubtensorModule.batch_commit_weights",
		netuids,
		commitHashes,
	)
	if err != nil {
		return types.Call{}, err
	}
	return call, nil
}

func BatchCommitWeightsExt(
	c *client.Client,
	netuids []types.U16,
	commitHashes []types.Hash,
) (*extrinsic.Extrinsic, error) {
	call, err := BatchCommitWeightsCall(
		c,
		netuids,
		commitHashes,
	)
	if err != nil {
		return nil, err
	}
	ext := extrinsic.NewExtrinsic(call)
	return &ext, nil
}

// - [ ] reveal_weights (Index: 97)
func RevealWeightsCall(
	c *client.Client,
	netuid types.U16,
	uids []types.U16,
	values []types.U16,
	salt []types.U16,
	versionKey types.U64,
) (types.Call, error) {
	call, err := types.NewCall(
		c.Meta,
		"SubtensorModule.reveal_weights",
		netuid,
		uids,
		values,
		salt,
		versionKey,
	)
	if err != nil {
		return types.Call{}, err
	}
	return call, nil
}

func RevealWeightsExt(
	c *client.Client,
	netuid types.U16,
	uids []types.U16,
	values []types.U16,
	salt []types.U16,
	versionKey types.U64,
) (*extrinsic.Extrinsic, error) {
	call, err := RevealWeightsCall(
		c,
		netuid,
		uids,
		values,
		salt,
		versionKey,
	)
	if err != nil {
		return nil, err
	}
	ext := extrinsic.NewExtrinsic(call)
	return &ext, nil
}

// - [ ] batch_reveal_weights (Index: 98)
func BatchRevealWeightsCall(
	c *client.Client,
	netuid types.U16,
) (types.Call, error) {
	call, err := types.NewCall(
		c.Meta,
		"SubtensorModule.batch_reveal_weights",
		netuid,
	)
	if err != nil {
		return types.Call{}, err
	}
	return call, nil
}

func BatchRevealWeightsExt(
	c *client.Client,
	netuid types.U16,
) (*extrinsic.Extrinsic, error) {
	call, err := BatchRevealWeightsCall(
		c,
		netuid,
	)
	if err != nil {
		return nil, err
	}
	ext := extrinsic.NewExtrinsic(call)
	return &ext, nil
}

// - [ ] set_tao_weights (Index: 8)
func SetTaoWeightsCall(
	c *client.Client,
	netuid types.U16,
	hotkey types.AccountID,
	dests []types.U16,
	weights []types.U16,
	versionKey types.U64,
) (types.Call, error) {
	call, err := types.NewCall(
		c.Meta,
		"SubtensorModule.set_tao_weights",
		netuid,
		hotkey,
		dests,
		weights,
		versionKey,
	)
	if err != nil {
		return types.Call{}, err
	}
	return call, nil
}

func SetTaoWeightsExt(
	c *client.Client,
	netuid types.U16,
	hotkey types.AccountID,
	dests []types.U16,
	weights []types.U16,
	versionKey types.U64,
) (*extrinsic.Extrinsic, error) {
	call, err := SetTaoWeightsCall(
		c,
		netuid,
		hotkey,
		dests,
		weights,
		versionKey,
	)
	if err != nil {
		return nil, err
	}
	ext := extrinsic.NewExtrinsic(call)
	return &ext, nil
}

// - [ ] become_delegate (Index: 1)// BecomeDelegateCall creates a call to become a delegate (DEPRECATED)
func BecomeDelegateCall(c *client.Client, hotkey types.AccountID) (types.Call, error) {
	call, err := types.NewCall(
		c.Meta,
		"SubtensorModule.become_delegate",
		hotkey,
	)
	if err != nil {
		return types.Call{}, err
	}
	return call, nil
}

func BecomeDelegateExt(c *client.Client, hotkey types.AccountID) (*extrinsic.Extrinsic, error) {
	call, err := BecomeDelegateCall(c, hotkey)
	if err != nil {
		return nil, err
	}
	ext := extrinsic.NewExtrinsic(call)
	return &ext, nil
}

// - [ ] serve_prometheus (Index: 5)
func ServePrometheusCall(
	c *client.Client,
	netuid types.U16,
	version types.U32,
	ip types.U128,
	port types.U16,
	ipType types.U8,
) (types.Call, error) {
	call, err := types.NewCall(
		c.Meta,
		"SubtensorModule.serve_prometheus",
		netuid,
		version,
		ip,
		port,
		ipType,
	)
	if err != nil {
		return types.Call{}, err
	}
	return call, nil
}

func ServePrometheusExt(
	c *client.Client,
	netuid types.U16,
	version types.U32,
	ip types.U128,
	port types.U16,
	ipType types.U8,
) (*extrinsic.Extrinsic, error) {
	call, err := ServePrometheusCall(
		c,
		netuid,
		version,
		ip,
		port,
		ipType,
	)
	if err != nil {
		return nil, err
	}
	ext := extrinsic.NewExtrinsic(call)
	return &ext, nil
}

// original code here below
func AddStakeLimitCall(c *client.Client, hotkey types.AccountID, netuid types.U16, amount_staked types.U64, limit_price types.U64, allow_partial types.Bool) (types.Call, error) {
	call, err := types.NewCall(c.Meta, "SubtensorModule.add_stake_limit", hotkey, netuid, amount_staked, limit_price, allow_partial)
	if err != nil {
		return types.Call{}, err
	}
	return call, nil
}

func AddStakeLimitExt(c *client.Client, hotkey types.AccountID, netuid types.U16, amount_staked types.U64, limit_price types.U64, allow_partial types.Bool) (*extrinsic.Extrinsic, error) {
	call, err := AddStakeLimitCall(c, hotkey, netuid, amount_staked, limit_price, allow_partial)
	if err != nil {
		return nil, err
	}
	ext := extrinsic.NewExtrinsic(call)
	return &ext, nil
}

func RemoveStakeLimitCall(c *client.Client, hotkey types.AccountID, netuid types.U16, amount_unstaked types.U64, limit_price types.U64, allow_partial types.Bool) (types.Call, error) {
	call, err := types.NewCall(c.Meta, "SubtensorModule.remove_stake_limit", hotkey, netuid, amount_unstaked, limit_price, allow_partial)
	if err != nil {
		return types.Call{}, err
	}
	return call, nil
}

func RemoveStakeLimitExt(c *client.Client, hotkey types.AccountID, netuid types.U16, amount_unstaked types.U64, limit_price types.U64, allow_partial types.Bool) (*extrinsic.Extrinsic, error) {
	call, err := RemoveStakeLimitCall(c, hotkey, netuid, amount_unstaked, limit_price, allow_partial)
	if err != nil {
		return nil, err
	}
	ext := extrinsic.NewExtrinsic(call)
	return &ext, nil
}

func SetWeightsCall(c *client.Client, netuid types.U16, uids []types.U16, weights []types.U16, versionKey types.U64) (types.Call, error) {
	call, err := types.NewCall(c.Meta, "SubtensorModule.set_weights", netuid, uids, weights, versionKey)
	if err != nil {
		return types.Call{}, err
	}
	return call, nil
}

func SetWeightsExt(c *client.Client, netuid types.U16, uids []types.U16, weights []types.U16, versionKey types.U64) (*extrinsic.Extrinsic, error) {
	call, err := SetWeightsCall(c, netuid, uids, weights, versionKey)
	if err != nil {
		return nil, err
	}
	ext := extrinsic.NewExtrinsic(call)
	return &ext, nil
}

func RootRegisterCall(c *client.Client, hotkey types.AccountID) (types.Call, error) {
	call, err := types.NewCall(c.Meta, "SubtensorModule.root_register", hotkey)
	if err != nil {
		return types.Call{}, err
	}
	return call, nil
}

func RootRegisterExt(c *client.Client, hotkey types.AccountID) (*extrinsic.Extrinsic, error) {
	call, err := RootRegisterCall(c, hotkey)
	if err != nil {
		return nil, err
	}
	ext := extrinsic.NewExtrinsic(call)
	return &ext, nil
}

func RegisterCall(c *client.Client, netuid types.U16, blockNumber types.U64, nonce types.U64, work types.Bytes, hotkey types.AccountID, coldkey types.AccountID) (types.Call, error) {
	call, err := types.NewCall(
		c.Meta,
		"SubtensorModule.register",
		netuid,
		blockNumber,
		nonce,
		work,
		hotkey,
		coldkey,
	)
	if err != nil {
		return types.Call{}, err
	}
	return call, nil
}

func RegisterExt(c *client.Client, netuid types.U16, blockNumber types.U64, nonce types.U64, work types.Bytes, hotkey types.AccountID, coldkey types.AccountID) (*extrinsic.Extrinsic, error) {
	call, err := RegisterCall(c, netuid, blockNumber, nonce, work, hotkey, coldkey)
	if err != nil {
		return nil, err
	}
	ext := extrinsic.NewExtrinsic(call)
	return &ext, nil
}

func BurnedRegisterCall(c *client.Client, hotkey types.AccountID, netuid types.U16) (types.Call, error) {
	call, err := types.NewCall(
		c.Meta,
		"SubtensorModule.burned_register",
		netuid,
		hotkey,
	)
	if err != nil {
		return types.Call{}, err
	}
	return call, nil
}

func BurnedRegisterExt(c *client.Client, hotkey types.AccountID, netuid types.U16) (*extrinsic.Extrinsic, error) {
	call, err := BurnedRegisterCall(c, hotkey, netuid)
	if err != nil {
		return nil, err
	}
	ext := extrinsic.NewExtrinsic(call)
	return &ext, nil
}

func RegisterNetworkCall(c *client.Client, hotkey types.AccountID) (types.Call, error) {
	call, err := types.NewCall(
		c.Meta,
		"SubtensorModule.register_network",
		hotkey,
	)
	if err != nil {
		return types.Call{}, err
	}
	return call, nil
}

func RegisterNetworkExt(c *client.Client, hotkey types.AccountID) (*extrinsic.Extrinsic, error) {
	call, err := RegisterNetworkCall(c, hotkey)
	if err != nil {
		return nil, err
	}
	ext := extrinsic.NewExtrinsic(call)
	return &ext, nil
}

func ServeAxonCall(c *client.Client, netuid types.U16, version types.U32, ip types.U128,
	port types.U16, ipType types.U8, protocol types.U8, placeholder1 types.U8,
	placeholder2 types.U8) (types.Call, error) {

	call, err := types.NewCall(
		c.Meta,
		"SubtensorModule.serve_axon",
		netuid,
		version,
		ip,
		port,
		ipType,
		protocol,
		placeholder1,
		placeholder2,
	)

	if err != nil {
		return types.Call{}, err
	}

	return call, nil
}

func ServeAxonExt(c *client.Client, netuid types.U16, version types.U32, ip types.U128,
	port types.U16, ipType types.U8, protocol types.U8, placeholder1 types.U8,
	placeholder2 types.U8) (*extrinsic.Extrinsic, error) {

	call, err := ServeAxonCall(c, netuid, version, ip, port, ipType, protocol,
		placeholder1, placeholder2)

	if err != nil {
		return nil, err
	}

	ext := extrinsic.NewExtrinsic(call)
	return &ext, nil
}

func ServeAxonTLSCall(c *client.Client, netuid types.U16, version types.U32, ip types.U128,
	port types.U16, ipType types.U8, protocol types.U8, placeholder1 types.U8,
	placeholder2 types.U8, certificate types.Bytes) (types.Call, error) {

	call, err := types.NewCall(
		c.Meta,
		"SubtensorModule.serve_axon_tls",
		netuid,
		version,
		ip,
		port,
		ipType,
		protocol,
		placeholder1,
		placeholder2,
		certificate,
	)

	if err != nil {
		return types.Call{}, err
	}

	return call, nil
}

func ServeAxonTLSExt(c *client.Client, netuid types.U16, version types.U32, ip types.U128,
	port types.U16, ipType types.U8, protocol types.U8, placeholder1 types.U8,
	placeholder2 types.U8, certificate types.Bytes) (*extrinsic.Extrinsic, error) {

	call, err := ServeAxonTLSCall(c, netuid, version, ip, port, ipType, protocol,
		placeholder1, placeholder2, certificate)

	if err != nil {
		return nil, err
	}

	ext := extrinsic.NewExtrinsic(call)
	return &ext, nil
}

func CommitCRV3WeightsCall(c *client.Client, netuid types.U16, commit types.Bytes, revealRound types.U64) (types.Call, error) {
	call, err := types.NewCall(
		c.Meta,
		"SubtensorModule.commit_crv3_weights",
		netuid,
		commit,
		revealRound,
	)

	if err != nil {
		return types.Call{}, err
	}

	return call, nil
}

func CommitCRV3WeightsExt(c *client.Client, netuid types.U16, commit types.Bytes, revealRound types.U64) (*extrinsic.Extrinsic, error) {
	call, err := CommitCRV3WeightsCall(c, netuid, commit, revealRound)
	if err != nil {
		return nil, err
	}
	ext := extrinsic.NewExtrinsic(call)
	return &ext, nil
}
