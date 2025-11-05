// Package typetools
package typetools

import (
	"math/big"

	"github.com/centrifuge/go-substrate-rpc-client/v4/scale"
)

type U96F32 big.Int

// NewU96F32 creates a U96F32 value from an integer and a fractional part (scaled by 2^32).
func NewU96F32(integer uint64, fractional uint32) *U96F32 {
	// Shift the integer part left by 32 bits and add the fractional part.
	val := new(big.Int).SetUint64(integer)
	val.Lsh(val, 32)
	val.Or(val, new(big.Int).SetUint64(uint64(fractional)))
	return (*U96F32)(val)
}

// Float64 converts the U96F32 value to a float64 for display or interoperation.
func (u *U96F32) Float64() float64 {
	val := (*big.Int)(u)
	// Convert to big.Float for precise floating-point conversion
	f := new(big.Float).SetInt(val)
	// Divide by 2^32
	scale := new(big.Float).SetUint64(1 << 32)
	f.Quo(f, scale)
	res, _ := f.Float64()
	return res
}

func (u *U96F32) Decode(decoder scale.Decoder) error {
	bs := make([]byte, 16)
	err := decoder.Read(bs)
	if err != nil {
		return err
	}
	scale.Reverse(bs)
	b := big.NewInt(0).SetBytes(bs)
	*u = *(*U96F32)(b)
	return nil
}
