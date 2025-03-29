package info

import (
	"fmt"
	"io"
	"os"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
)

func LookupExtrinsicArgs(meta *types.Metadata, palletName, callName string) {
	for _, pallet := range meta.AsMetadataV14.Pallets {
		if string(pallet.Name) != palletName || !pallet.HasCalls {
			continue
		}

		callTypeID := pallet.Calls.Type.Int64()
		callType, ok := meta.AsMetadataV14.EfficientLookup[callTypeID]
		if !ok {
			fmt.Printf("Call type not found\n")
			return
		}

		if !callType.Def.IsVariant {
			fmt.Printf("Call type is not a variant\n")
			return
		}

		for _, variant := range callType.Def.Variant.Variants {
			if string(variant.Name) != callName {
				continue
			}

			fmt.Printf("Call: %s::%s\n", palletName, callName)
			for i, field := range variant.Fields {
				typeID := field.Type.Int64()
				fieldType, _ := meta.AsMetadataV14.EfficientLookup[typeID]

				name := field.Name
				if name == "" {
					name = types.NewText(fmt.Sprintf("arg%d", i))
				}

				fmt.Printf("  Arg %d: %s (Type: ", i, name)
				if fieldType != nil && len(fieldType.Path) > 0 {
					for j, part := range fieldType.Path {
						if j > 0 {
							fmt.Printf("::")
						}
						fmt.Printf("%s", part)
					}
				} else {
					fmt.Printf("unknown")
				}
				fmt.Printf(")\n")
			}
			return
		}

		fmt.Printf("Call '%s' not found in module '%s'\n", callName, palletName)
		return
	}

	fmt.Printf("Module '%s' not found\n", palletName)
}

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

func PrintStorageItemInfo(meta *types.Metadata, out io.Writer, palletName, storageName string) {
	if out == nil {
		out = os.Stdout
	}

	found := false
	for _, module := range meta.AsMetadataV14.Pallets {
		if string(module.Name) == palletName {
			if !module.HasStorage {
				fmt.Fprintf(out, "Pallet '%s' has no storage items\n", palletName)
				return
			}

			for _, storageEntry := range module.Storage.Items {
				if string(storageEntry.Name) == storageName {
					found = true

					fmt.Fprintf(out, "Storage Item: %s::%s\n", palletName, storageName)
					fmt.Fprintf(out, "  Documentation: %s\n", storageEntry.Documentation)
					fmt.Fprintf(out, "  Type: %v\n", storageEntry.Type)

					if storageEntry.Type.IsMap {
						fmt.Fprintf(out, "  Storage Type: Map\n")
						fmt.Fprintf(out, "  Key Type ID: %d\n", storageEntry.Type.AsMap.Key.Int64())
						fmt.Fprintf(out, "  Value Type ID: %d\n", storageEntry.Type.AsMap.Value.Int64())

						// Get key type info
						keyTypeId := storageEntry.Type.AsMap.Key.Int64()
						if keyType, found := meta.AsMetadataV14.EfficientLookup[keyTypeId]; found {
							fmt.Fprintf(out, "  Key Type Definition: %+v\n", keyType)
						}

						// Get value type info
						valueTypeId := storageEntry.Type.AsMap.Value.Int64()
						if valueType, found := meta.AsMetadataV14.EfficientLookup[valueTypeId]; found {
							fmt.Fprintf(out, "  Value Type Definition: %+v\n", valueType)
						}
					} else if storageEntry.Type.IsPlainType {
						fmt.Fprintf(out, "  Storage Type: Plain\n")
						fmt.Fprintf(out, "  Type ID: %d\n", storageEntry.Type.AsPlainType.Int64())

						typeId := storageEntry.Type.AsPlainType.Int64()
						if typeInfo, found := meta.AsMetadataV14.EfficientLookup[typeId]; found {
							fmt.Fprintf(out, "  Type Definition: %+v\n", typeInfo)
						}
					} else {
						fmt.Fprintf(out, "  Storage Type: Other/Unknown\n")
					}

					break
				}
			}

			if !found {
				fmt.Fprintf(out, "Storage item '%s' not found in pallet '%s'\n", storageName, palletName)
			}

			return
		}
	}

	fmt.Fprintf(out, "Pallet '%s' not found\n", palletName)
}
