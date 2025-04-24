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

//     - [ ] batch_set_weights (Index: 80)
//     - [ ] commit_weights (Index: 96)
//     - [ ] batch_commit_weights (Index: 100)
//     - [ ] reveal_weights (Index: 97)
//     - [ ] commit_crv3_weights (Index: 99)
//     - [ ] batch_reveal_weights (Index: 98)
//     - [ ] set_tao_weights (Index: 8)
//     - [ ] become_delegate (Index: 1)
//     - [ ] decrease_take (Index: 65)
//     - [ ] increase_take (Index: 66)
//     - [ ] add_stake (Index: 2)
//     - [ ] remove_stake (Index: 3)
//     - [ ] serve_prometheus (Index: 5)

//     - [ ] adjust_senate (Index: 63)
//     - [ ] swap_hotkey (Index: 70)
//     - [ ] swap_coldkey (Index: 71)
//     - [ ] set_childkey_take (Index: 75)
//     - [ ] sudo_set_tx_childkey_take_rate_limit (Index: 69)
//     - [ ] sudo_set_min_childkey_take (Index: 76)
//     - [ ] sudo_set_max_childkey_take (Index: 77)
//     - [ ] sudo (Index: 51)
//     - [ ] sudo_unchecked_weight (Index: 52)
//     - [ ] vote (Index: 55)
//     - [ ] faucet (Index: 60)
//     - [ ] dissolve_network (Index: 61)
//     - [ ] set_children (Index: 67)
//     - [ ] schedule_swap_coldkey (Index: 73)
//     - [ ] schedule_dissolve_network (Index: 74)
//     - [ ] set_identity (Index: 68)
//     - [ ] set_subnet_identity (Index: 78)
//     - [ ] register_network_with_identity (Index: 79)
//     - [ ] unstake_all (Index: 83)
//     - [ ] unstake_all_alpha (Index: 84)
//     - [ ] move_stake (Index: 85)
//     - [ ] transfer_stake (Index: 86)
//     - [ ] swap_stake (Index: 87)
//     - [x] add_stake_limit (Index: 88)
//     - [x] remove_stake_limit (Index: 89)
//     - [ ] swap_stake_limit (Index: 90)
//     - [ ] try_associate_hotkey (Index: 91)

type SubnetIdentityV2 struct {
	SubnetName    types.Bytes
	GithubRepo    types.Bytes
	SubnetContact types.Bytes
	SubnetURL     types.Bytes
	Discord       types.Bytes
	Description   types.Bytes
	Additional    types.Bytes
}

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
	placeholder2 types.U8, certificate []byte) (types.Call, error) {

	// build a generic Option[Bytes]
	//	var certOption types.Option[types.Bytes]
	//	if len(certificate) > 0 {
	//		// wrap raw []byte in the Bytes type, then in an Option
	//		certOption = types.NewOption(types.NewBytes(certificate))
	//	} else {
	//		certOption = types.NewEmptyOption[types.Bytes]()
	//	}
	//
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
	placeholder2 types.U8, certificate []byte) (*extrinsic.Extrinsic, error) {

	call, err := ServeAxonCall(c, netuid, version, ip, port, ipType, protocol,
		placeholder1, placeholder2, certificate)

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
