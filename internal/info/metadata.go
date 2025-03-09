package info

import (
	"fmt"
	"io"
	"os"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
)

func PrintModulesAndCalls(meta *types.Metadata, out io.Writer) {
	if out == nil {
		out = os.Stdout
	}

	fmt.Fprintln(out, "\nAVAILABLE MODULES AND CALLS:")
	for _, module := range meta.AsMetadataV14.Pallets {
		fmt.Fprintf(out, "Module: %s (Index: %d)\n", module.Name, module.Index)

		// Print calls if they exist
		if module.HasCalls {
			// Get the call type ID
			callTypeID := module.Calls.Type.Int64()
			fmt.Fprintf(out, "  Call Type ID: %d\n", callTypeID)

			// Find the type in the lookup
			if callType, ok := meta.AsMetadataV14.EfficientLookup[callTypeID]; ok {
				if callType.Def.IsVariant {
					fmt.Fprintln(out, "  Available Calls:")
					for _, variant := range callType.Def.Variant.Variants {
						fmt.Fprintf(out, "    %s (Index: %d)\n", variant.Name, variant.Index)
					}
				}
			} else {
				fmt.Fprintf(out, "  Call type not found in lookup\n")
			}
		} else {
			fmt.Fprintf(out, "  No calls available\n")
		}
		fmt.Fprintln(out)
	}
}

func PrintExtensions(meta *types.Metadata, out io.Writer) {
	if out == nil {
		out = os.Stdout
	}

	fmt.Fprintln(out, "Signed Extensions:")
	for _, ext := range meta.AsMetadataV14.Extrinsic.SignedExtensions {
		fmt.Fprintf(out, "- %s (Identifier: %s, Type: %v, AdditionalSigned: %v)\n",
			ext.Identifier, ext.Identifier, ext.Type, ext.AdditionalSigned)
	}
}

func PrintExtensionDetails(meta *types.Metadata, out io.Writer, extensionNames ...string) {
	if out == nil {
		out = os.Stdout
	}

	fmt.Fprintln(out, "Extension Details:")
	for _, ext := range meta.AsMetadataV14.Extrinsic.SignedExtensions {
		if len(extensionNames) > 0 {
			found := false
			for _, name := range extensionNames {
				if string(ext.Identifier) == name {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		typeID := ext.Type.Int64()
		if def, ok := meta.AsMetadataV14.EfficientLookup[typeID]; ok {
			fmt.Fprintf(out, "Extension %s Type Definition: %+v\n", ext.Identifier, def)
		} else {
			fmt.Fprintf(out, "Extension %s Type ID %d not found in lookup\n", ext.Identifier, typeID)
		}
	}
}
