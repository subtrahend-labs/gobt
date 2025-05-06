package sigtools

import (
	"math/big"

	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types/extrinsic"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types/extrinsic/extensions"
	"github.com/subtrahend-labs/gobt/client"
	"github.com/subtrahend-labs/gobt/storage"
)

type SigningContext struct {
	Tip   *types.UCompact
	Nonce *types.UCompact
}

func NewSigningContext(t *types.UCompact, n *types.UCompact) *SigningContext {
	return &SigningContext{Tip: t, Nonce: n}
}

func CreateSigningOptions(c *client.Client, keypair signature.KeyringPair, sc *SigningContext) ([]extrinsic.SigningOption, error) {
	extrinsic.PayloadMutatorFns[extensions.SignedExtensionName("SubtensorSignedExtension")] = func(payload *extrinsic.Payload) {}
	extrinsic.PayloadMutatorFns[extensions.SignedExtensionName("CommitmentsSignedExtension")] = func(payload *extrinsic.Payload) {}
	var options []extrinsic.SigningOption

	// tip
	tip := types.NewUCompactFromUInt(0)
	if sc != nil && sc.Tip != nil {
		tip = *sc.Tip
	}
	options = append(options,
		extrinsic.WithTip(tip),
	)

	// Nonce
	nonce := types.NewUCompact(big.NewInt(0))
	if sc != nil && sc.Nonce != nil {
		nonce = *sc.Nonce
	} else if s, err := storage.GetAccountInfo(c, keypair.PublicKey, nil); err == nil {
		nonce = types.NewUCompactFromUInt(uint64(s.Nonce))
	}
	options = append(options,
		extrinsic.WithNonce(nonce),
	)

	// Spec & transaction Version
	rv, err := c.Api.RPC.State.GetRuntimeVersionLatest()
	if err != nil {
		return nil, err
	}
	options = append(options,
		extrinsic.WithSpecVersion(rv.SpecVersion),
	)
	options = append(options,
		extrinsic.WithTransactionVersion(rv.TransactionVersion),
	)

	// Metadat Mode
	options = append(options,
		extrinsic.WithMetadataMode(extensions.CheckMetadataModeDisabled, extensions.CheckMetadataHash{Hash: types.NewEmptyOption[types.H256]()}),
	)

	// Genesis and era
	genesisHash, err := c.Api.RPC.Chain.GetBlockHash(0)
	if err != nil {
		return nil, err
	}

	options = append(options,
		extrinsic.WithGenesisHash(genesisHash),
	)
	options = append(options,
		extrinsic.WithEra(types.ExtrinsicEra{IsImmortalEra: true}, genesisHash),
	)

	return options, nil
}
