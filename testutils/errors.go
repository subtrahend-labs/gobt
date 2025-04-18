package testutils

import (
	"encoding/binary"
	"fmt"

	"github.com/centrifuge/go-substrate-rpc-client/v4/registry"
	"github.com/centrifuge/go-substrate-rpc-client/v4/registry/parser"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
)

// moduleErrorCode decodes the 4‑byte little‑endian ModuleError.Error into a uint32.
func moduleErrorCode(b [4]types.U8) uint32 {
	return binary.LittleEndian.Uint32([]byte{byte(b[0]), byte(b[1]), byte(b[2]), byte(b[3])})
}

// ExtractDispatchError inspects parser.Events and returns a Go error describing exactly why
// System.ExtrinsicFailed fired.  If no failure event is found, returns nil.
func ExtractDispatchError(
	meta types.Metadata,
	events []*parser.Event,
) error {
	for _, ev := range events {
		if ev.Name != "System.ExtrinsicFailed" {
			continue
		}
		if len(ev.Fields) == 0 {
			return fmt.Errorf("ExtrinsicFailed: no DispatchError field")
		}

		val := ev.Fields[0].Value

		// 1) Native DispatchError
		if de, ok := val.(types.DispatchError); ok {
			if de.IsModule {
				me := de.ModuleError
				idx := moduleErrorCode(me.Error)
				metaErr, err := meta.FindError(me.Index, me.Error)
				if err != nil {
					return fmt.Errorf("module error %d/%d (lookup failed: %v)",
						me.Index, idx, err)
				}
				return fmt.Errorf("%s: %s", metaErr.Name, metaErr.Value)
			}
			switch {
			case de.IsBadOrigin:
				return fmt.Errorf("BadOrigin")
			case de.IsCannotLookup:
				return fmt.Errorf("CannotLookup")
			default:
				return fmt.Errorf("DispatchError: unknown variant")
			}
		}

		// 2) Registry‑based failure
		if df, ok := val.(registry.DecodedFields); ok {
			// The outer df typically has a single field whose Value is itself DecodedFields
			for _, f := range df {
				if nested, ok := f.Value.(registry.DecodedFields); ok {
					// nested should have exactly two entries: "index" and "error"
					var palletIndex uint8
					var errorBytes [4]types.U8

					for _, nf := range nested {
						switch nf.Name {
						case "index":
							// nf.Value is usually a uint64 or types.U8
							switch v := nf.Value.(type) {
							case uint8:
								palletIndex = v
							case uint64:
								palletIndex = uint8(v)
							case types.U8:
								palletIndex = uint8(v)
							}
						case "error":
							// nf.Value is []interface{}{u8, u8, u8, u8}
							if arr, ok := nf.Value.([]interface{}); ok {
								for i, elt := range arr {
									switch b := elt.(type) {
									case uint8:
										errorBytes[i] = types.U8(b)
									case types.U8:
										errorBytes[i] = b
									}
								}
							}
						}
					}

					// Now look it up in the metadata
					metaErr, err := meta.FindError(types.U8(palletIndex), errorBytes)
					if err != nil {
						return fmt.Errorf("module error %d/%v (lookup failed: %v)",
							palletIndex, errorBytes, err)
					}
					return fmt.Errorf("%s: %s", metaErr.Name, metaErr.Value)
				}
			}
		}

		// 3) Something completely unexpected
		return fmt.Errorf("dispatch failed; unexpected field type %T: %#v", val, val)
	}
	return nil
}
