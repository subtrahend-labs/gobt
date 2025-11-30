package extrinsics

/*
#cgo CFLAGS: -I${SRCDIR}/../third_party/bittensor-drand
#cgo LDFLAGS: -L${SRCDIR}/../third_party/bittensor-drand/target/release -lbittensor_drand -ldl
#include "bindings.h"
*/
import "C"

import (
	"fmt"
	"unsafe"
)

// GenerateCommit computes the encrypted CRv4 commit and returns the bytes + reveal round.
func GenerateCommit(
	uids []uint16,
	vals []uint16,
	versionKey uint64,
	tempo uint64,
	currentBlock uint64,
	netuid uint16,
	revealEpochs uint64,
	blockTime float64,
	hotkey []byte,
) ([]byte, uint64, error) {
	if len(uids) != len(vals) {
		return nil, 0, fmt.Errorf("uids and vals must have same length")
	}

	// Prepare C pointers
	uPtr := (*C.uint16_t)(unsafe.Pointer(&uids[0]))
	vPtr := (*C.uint16_t)(unsafe.Pointer(&vals[0]))
	length := C.size_t(len(uids))
	
	var hotkeyPtr *C.uint8_t
	var hotkeyLen C.size_t
	if len(hotkey) > 0 {
		hotkeyPtr = (*C.uint8_t)(unsafe.Pointer(&hotkey[0]))
		hotkeyLen = C.size_t(len(hotkey))
	}

	var cRound C.uint64_t
	var cErr *C.char

	// Call into the Rust library
	buf := C.cr_generate_commit(
		uPtr, length,
		vPtr, length,
		C.uint64_t(versionKey),
		C.uint64_t(tempo),
		C.uint64_t(currentBlock),
		C.uint16_t(netuid),
		C.uint64_t(revealEpochs),
		C.double(blockTime),
		hotkeyPtr, hotkeyLen,
		&cRound,
		&cErr,
	)

	if cErr != nil {
		err := fmt.Errorf("error from Rust: %s", C.GoString(cErr))
		C.cr_free_str(cErr)
		return nil, 0, err
	}

	if buf.ptr == nil {
		return nil, 0, fmt.Errorf("received nil buffer from Rust")
	}

	// Convert C buffer to Go slice
	goBytes := C.GoBytes(unsafe.Pointer(buf.ptr), C.int(buf.len))
	
	// Free the Rust-allocated buffer
	C.cr_free(buf)

	revealRound := uint64(cRound)

	return goBytes, revealRound, nil
}