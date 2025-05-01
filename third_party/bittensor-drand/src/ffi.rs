//! Thin C-ABI wrappers around crate::drand helpers.
//! *All* heap ownership is transferred to the caller – remember to call `cr_free`!
//!
//! This module provides FFI bindings for the Drand-based timelock encryption functionality,
//! allowing C/C++ applications to use the Rust implementation.

use crate::drand;
use codec::{Decode, Encode};
use std::ffi::CString;
use std::os::raw::c_char;
use std::ptr;

/// A buffer structure for transferring byte arrays across the FFI boundary.
///
/// This structure is used to safely transfer ownership of heap-allocated memory
/// from Rust to C. The caller is responsible for freeing the memory using `cr_free`.
#[derive(Clone, Copy)]
#[repr(C)]
pub struct CRByteBuffer {
    /// Pointer to the buffer data
    pub ptr: *mut u8,
    /// Length of the buffer in bytes
    pub len: usize,
    /// Capacity of the allocated buffer
    pub cap: usize, // frees via Vec::from_raw_parts
}

impl CRByteBuffer {
    /// Creates a new CRByteBuffer from a Vec<u8>, transferring ownership to the caller.
    ///
    /// This function takes ownership of the vector and forgets it in Rust,
    /// allowing the memory to be managed by the C caller.
    fn from_vec(mut v: Vec<u8>) -> Self {
        let buf = CRByteBuffer {
            ptr: v.as_mut_ptr(),
            len: v.len(),
            cap: v.capacity(),
        };
        std::mem::forget(v); // give ownership to caller
        buf
    }
}

/// Frees a CRByteBuffer that was previously returned by a function in this module.
///
/// # Safety
///
/// The buffer must have been created by one of the functions in this module.
/// Calling this function with an invalid buffer or calling it twice on the same buffer
/// will result in undefined behavior.
#[no_mangle]
pub extern "C" fn cr_free(buf: CRByteBuffer) {
    if buf.ptr.is_null() {
        return;
    }
    unsafe {
        let _ = Vec::from_raw_parts(buf.ptr, buf.len, buf.cap);
    }
}

/// Frees a C string that was previously returned by a function in this module.
///
/// # Safety
///
/// The string must have been created by one of the functions in this module.
/// Calling this function with an invalid pointer or calling it twice on the same pointer
/// will result in undefined behavior.
#[no_mangle]
pub extern "C" fn cr_free_str(s: *mut c_char) {
    if !s.is_null() {
        unsafe {
            let _ = CString::from_raw(s); // Take ownership and drop the CString
        }
    }
}

// ---------- helpers for error propagation ----------
/// Converts any error type that implements ToString into a C-compatible string.
///
/// The returned pointer must be freed with `cr_free_str`.
fn err_to_cstring<E: ToString>(e: E) -> *mut c_char {
    CString::new(e.to_string()).unwrap().into_raw()
}

/// Encrypts binary data using Drand timelock encryption.
///
/// # Parameters
///
/// * `data_ptr` - Pointer to the data to encrypt
/// * `data_len` - Length of the data in bytes
/// * `n_blocks` - Number of blocks to wait before the data can be decrypted
/// * `block_time` - Duration of a single block in seconds (e.g., 12.0 for mainnet)
/// * `round_out` - Output parameter that will be set to the reveal round number
/// * `err_out` - Output parameter that will be set to an error message on failure
///
/// # Returns
///
/// A `CRByteBuffer` containing the encrypted data, or an empty buffer on error.
/// The caller is responsible for freeing the buffer with `cr_free`.
///
/// # Safety
///
/// The caller must ensure that:
/// - `data_ptr` points to valid memory of at least `data_len` bytes
/// - `round_out` and `err_out` point to valid memory locations
#[no_mangle]
pub extern "C" fn cr_encrypt(
    data_ptr: *const u8,
    data_len: usize,
    n_blocks: u64,
    block_time: f64,
    round_out: *mut u64,
    err_out: *mut *mut c_char,
) -> CRByteBuffer {
    unsafe { *err_out = ptr::null_mut() }
    if data_ptr.is_null() {
        unsafe { *err_out = err_to_cstring("data ptr is null") };
        return CRByteBuffer {
            ptr: ptr::null_mut(),
            len: 0,
            cap: 0,
        };
    }
    let data = unsafe { std::slice::from_raw_parts(data_ptr, data_len) };
    let now = std::time::SystemTime::now()
        .duration_since(std::time::UNIX_EPOCH)
        .unwrap()
        .as_secs_f64();
    let reveal_timestamp = (n_blocks as f64 * block_time + now).ceil() as u64 - drand::GENESIS_TIME;
    let reveal_round = reveal_timestamp / drand::DRAND_PERIOD;

    match drand::encrypt_and_compress(data, reveal_round) {
        Ok(ct) => {
            // pack ciphertext + round into SCALE-encoded UserData
            let user_data = drand::UserData {
                encrypted_data: ct,
                reveal_round,
            };
            let encoded = user_data.encode();

            unsafe { *round_out = reveal_round };
            CRByteBuffer::from_vec(encoded)
        }
        Err(e) => {
            unsafe { *err_out = err_to_cstring(format!("{:?}", e)) };
            CRByteBuffer {
                ptr: ptr::null_mut(),
                len: 0,
                cap: 0,
            }
        }
    }
}

/// Decrypts data that was previously encrypted with timelock encryption.
///
/// # Parameters
///
/// * `enc_ptr` - Pointer to the encrypted data
/// * `enc_len` - Length of the encrypted data in bytes
/// * `no_errors` - If true, suppresses error messages and returns null on failure
/// * `out_len` - Output parameter that will be set to the length of the decrypted data
/// * `err_out` - Output parameter that will be set to an error message on failure
///
/// # Returns
///
/// A pointer to the decrypted data, or null on error.
/// The caller is responsible for freeing the memory (using standard free/delete).
///
/// # Safety
///
/// The caller must ensure that:
/// - `enc_ptr` points to valid memory of at least `enc_len` bytes
/// - `out_len` and `err_out` point to valid memory locations
#[no_mangle]
pub extern "C" fn cr_decrypt(
    enc_ptr: *const u8,
    enc_len: usize,
    no_errors: bool,
    out_len: *mut usize,
    err_out: *mut *mut c_char,
) -> *mut u8 {
    unsafe {
        *err_out = ptr::null_mut();
        *out_len = 0;
    }
    if enc_ptr.is_null() {
        unsafe {
            *err_out = err_to_cstring("enc ptr is null");
        }
        return ptr::null_mut();
    }
    let enc_slice = unsafe { std::slice::from_raw_parts(enc_ptr, enc_len) };

    let user_data = match drand::UserData::decode(&mut &enc_slice[..]) {
        Ok(d) => d,
        Err(e) => {
            if no_errors {
                return ptr::null_mut();
            }
            unsafe {
                *err_out = err_to_cstring(format!("{:?}", e));
            }
            return ptr::null_mut();
        }
    };

    let sig_opt = match drand::get_reveal_round_signature(Some(user_data.reveal_round), no_errors) {
        Ok(s) => s,
        Err(e) => {
            if no_errors {
                return ptr::null_mut();
            }
            unsafe {
                *err_out = err_to_cstring(e);
            }
            return ptr::null_mut();
        }
    };

    let sig_hex = match sig_opt {
        Some(s) => s,
        None => return ptr::null_mut(),
    };

    let sig_bytes = match hex::decode(&sig_hex) {
        Ok(b) => b,
        Err(e) => {
            if no_errors {
                return ptr::null_mut();
            }
            unsafe {
                *err_out = err_to_cstring(format!("{:?}", e));
            }
            return ptr::null_mut();
        }
    };

    match drand::decrypt_and_decompress(&user_data.encrypted_data, &sig_bytes) {
        Ok(plain) => {
            unsafe {
                *out_len = plain.len();
            }
            let mut v = plain; // take ownership
            let ptr = v.as_mut_ptr();
            std::mem::forget(v);
            ptr
        }
        Err(e) => {
            if !no_errors {
                unsafe {
                    *err_out = err_to_cstring(e);
                }
            }
            ptr::null_mut()
        }
    }
}

/// Retrieves the latest round number from the Drand network.
///
/// # Parameters
///
/// * `err_out` - Output parameter that will be set to an error message on failure
///
/// # Returns
///
/// The latest round number, or 0 on error.
///
/// # Safety
///
/// The caller must ensure that `err_out` points to a valid memory location.
#[no_mangle]
pub extern "C" fn cr_get_latest_round(err_out: *mut *mut c_char) -> u64 {
    unsafe {
        *err_out = ptr::null_mut();
    }
    match drand::get_round_info(None) {
        Ok(resp) => resp.round,
        Err(e) => {
            unsafe {
                *err_out = err_to_cstring(e);
            }
            0
        }
    }
}

/// Encrypts a string-based commitment using Drand timelock encryption.
///
/// This function calculates the reveal round based on the current time plus
/// the specified number of blocks, then encrypts the data to be revealed at that time.
///
/// # Parameters
///
/// * `data_ptr` - Pointer to the UTF-8 encoded string data to encrypt
/// * `data_len` - Length of the data in bytes
/// * `blocks_until_reveal` - Number of blocks to wait before the data can be decrypted
/// * `block_time` - Duration of a single block in seconds (e.g., 12.0 for mainnet)
/// * `round_out` - Output parameter that will be set to the reveal round number
/// * `err_out` - Output parameter that will be set to an error message on failure
///
/// # Returns
///
/// A `CRByteBuffer` containing the encrypted data, or an empty buffer on error.
/// The caller is responsible for freeing the buffer with `cr_free`.
///
/// # Safety
///
/// The caller must ensure that:
/// - `data_ptr` points to valid UTF-8 encoded memory of at least `data_len` bytes
/// - `round_out` and `err_out` point to valid memory locations
#[no_mangle]
pub extern "C" fn cr_encrypt_commitment(
    data_ptr: *const u8,
    data_len: usize,
    blocks_until_reveal: u64,
    block_time: f64,
    round_out: *mut u64,
    err_out: *mut *mut c_char,
) -> CRByteBuffer {
    unsafe { *err_out = ptr::null_mut() }
    if data_ptr.is_null() {
        unsafe { *err_out = err_to_cstring("data ptr is null") };
        return CRByteBuffer {
            ptr: ptr::null_mut(),
            len: 0,
            cap: 0,
        };
    }

    let bytes = unsafe { std::slice::from_raw_parts(data_ptr, data_len) };
    let data = match std::str::from_utf8(bytes) {
        Ok(s) => s,
        Err(e) => {
            unsafe { *err_out = err_to_cstring(format!("utf8 error: {e}")) };
            return CRByteBuffer {
                ptr: ptr::null_mut(),
                len: 0,
                cap: 0,
            };
        }
    };

    match drand::encrypt_commitment(data, blocks_until_reveal, block_time) {
        Ok((ct, rr)) => {
            unsafe { *round_out = rr }
            CRByteBuffer::from_vec(ct)
        }
        Err((ioe, msg)) => {
            unsafe { *err_out = err_to_cstring(format!("{msg}: {ioe}")) };
            CRByteBuffer {
                ptr: ptr::null_mut(),
                len: 0,
                cap: 0,
            }
        }
    }
}

/// Generates a commitment for a set of UIDs and values using Drand timelock encryption.
///
/// This function creates a commitment for voting/scoring in the Bittensor network,
/// encrypting the UIDs and their corresponding values to be revealed at a future time.
///
/// # Parameters
///
/// * `uids_ptr` - Pointer to an array of UIDs (u16 values)
/// * `uids_len` - Length of the UIDs array
/// * `vals_ptr` - Pointer to an array of values (u16 values)
/// * `vals_len` - Length of the values array
/// * `version_key` - Version key for the commitment
/// * `tempo` - Tempo value for the commitment
/// * `current_block` - Current block number
/// * `netuid` - Network UID
/// * `subnet_reveal_epochs` - Number of epochs to wait before revealing
/// * `block_time` - Duration of a single block in seconds
/// * `round_out` - Output parameter that will be set to the reveal round number
/// * `err_out` - Output parameter that will be set to an error message on failure
///
/// # Returns
///
/// A `CRByteBuffer` containing the encrypted commitment, or an empty buffer on error.
/// The caller is responsible for freeing the buffer with `cr_free`.
///
/// # Safety
///
/// The caller must ensure that:
/// - `uids_ptr` points to valid memory of at least `uids_len` u16 elements
/// - `vals_ptr` points to valid memory of at least `vals_len` u16 elements
/// - `round_out` and `err_out` point to valid memory locations
#[no_mangle]
pub extern "C" fn cr_generate_commit(
    uids_ptr: *const u16,
    uids_len: usize,
    vals_ptr: *const u16,
    vals_len: usize,
    version_key: u64,
    tempo: u64,
    current_block: u64,
    netuid: u16,
    subnet_reveal_epochs: u64,
    block_time: f64,
    round_out: *mut u64,
    err_out: *mut *mut c_char,
) -> CRByteBuffer {
    unsafe { *err_out = ptr::null_mut() }

    if (uids_ptr.is_null() && uids_len > 0) || (vals_ptr.is_null() && vals_len > 0) {
        unsafe { *err_out = err_to_cstring("uids/values ptr is null") };
        return CRByteBuffer {
            ptr: ptr::null_mut(),
            len: 0,
            cap: 0,
        };
    }

    if uids_len != vals_len {
        unsafe { *err_out = err_to_cstring("uids_len != vals_len") };
        return CRByteBuffer {
            ptr: ptr::null_mut(),
            len: 0,
            cap: 0,
        };
    }

    let uids = unsafe { std::slice::from_raw_parts(uids_ptr, uids_len) }.to_vec();
    let values = unsafe { std::slice::from_raw_parts(vals_ptr, vals_len) }.to_vec();

    match drand::generate_commit(
        uids,
        values,
        version_key,
        tempo,
        current_block,
        netuid,
        subnet_reveal_epochs,
        block_time,
    ) {
        Ok((ct, rr)) => {
            unsafe { *round_out = rr }
            CRByteBuffer::from_vec(ct)
        }
        Err((ioe, msg)) => {
            unsafe { *err_out = err_to_cstring(format!("{msg}: {ioe}")) };
            CRByteBuffer {
                ptr: ptr::null_mut(),
                len: 0,
                cap: 0,
            }
        }
    }
}

// TODO: add valgrind leak detection CI step.
// ============================================================================

#[cfg(test)]
mod tests {
    use super::*;
    use libc::free;
    use std::ffi::CStr;
    use std::ptr;

    // helper to free *mut c_char returned from FFI
    unsafe fn drop_cstring(ptr: *mut c_char) {
        if !ptr.is_null() {
            let _ = CString::from_raw(ptr);
        }
    }

    #[test]
    fn test_encrypt_success() {
        let msg = b"hello";
        let mut round: u64 = 0;
        let mut err_ptr: *mut c_char = ptr::null_mut();

        // SAFETY: we pass a valid pointer/length pair.
        let buf = unsafe {
            cr_encrypt(
                msg.as_ptr(),
                msg.len(),
                0,   // n_blocks (0 = immediate)
                1.0, // block_time
                &mut round,
                &mut err_ptr,
            )
        };

        assert!(round > 0, "round should be set");
        assert!(err_ptr.is_null(), "err_out must be NULL on success");
        assert!(!buf.ptr.is_null(), "buffer pointer must be non-NULL");
        assert!(buf.len > 0, "ciphertext must be non-empty");

        // Free resources
        unsafe {
            cr_free(buf);
        }
    }

    #[test]
    fn test_encrypt_null_ptr() {
        let mut round: u64 = 0;
        let mut err_ptr: *mut c_char = ptr::null_mut();

        let buf = unsafe {
            cr_encrypt(
                ptr::null(), // <- bad pointer
                5,
                0,
                1.0,
                &mut round,
                &mut err_ptr,
            )
        };

        assert!(buf.ptr.is_null(), "buffer should be NULL on error");
        assert!(round == 0, "round must stay 0 on failure");
        assert!(!err_ptr.is_null(), "err_out should be set");

        // grab error message for debugging (optional)
        unsafe {
            let msg = CStr::from_ptr(err_ptr).to_string_lossy();
            assert!(
                msg.contains("null"),
                "expected null-ptr message, got '{}'",
                msg
            );
            drop_cstring(err_ptr);
        }
    }

    #[test]
    fn test_encrypt_zero_len() {
        // Non-null ptr but len == 0
        let dummy = 0u8;
        let mut round: u64 = 0;
        let mut err_ptr: *mut c_char = ptr::null_mut();

        let buf = unsafe {
            cr_encrypt(
                &dummy as *const u8,
                0, // <- zero length
                0,
                1.0,
                &mut round,
                &mut err_ptr,
            )
        };

        // We don't mandate success vs. error – just require determinism:
        if err_ptr.is_null() {
            // Succeeded: buf.ptr must be non-NULL, free it.
            assert!(!buf.ptr.is_null(), "ptr may not be NULL on success");
            unsafe { cr_free(buf) };
        } else {
            // Failed: buf.ptr must be NULL, free the error string.
            assert!(buf.ptr.is_null(), "ptr must be NULL on failure");
            unsafe { drop_cstring(err_ptr) };
        }
    }

    // ---------------------------------------------------------------
    // 4. encrypt → decrypt round-trip
    // ---------------------------------------------------------------
    #[test]
    fn test_decrypt_roundtrip() {
        let msg = b"round-trip bytes";
        let mut round: u64 = 0;
        let mut err_ptr: *mut c_char = ptr::null_mut();

        // encrypt
        let ct_buf = unsafe {
            cr_encrypt(
                msg.as_ptr(),
                msg.len(),
                0,   // unlock immediately
                1.0, // block_time (doesn't matter)
                &mut round,
                &mut err_ptr,
            )
        };
        assert!(err_ptr.is_null());
        assert!(!ct_buf.ptr.is_null());

        // copy ciphertext to safe Vec (because we'll free ct_buf soon)
        let ciphertext = unsafe { std::slice::from_raw_parts(ct_buf.ptr, ct_buf.len) }.to_vec();

        // free Rust allocation backing the ciphertext buffer
        unsafe { cr_free(ct_buf) };

        // decrypt
        let mut out_len: usize = 0;
        let mut dec_err: *mut c_char = ptr::null_mut();
        let plain_ptr = unsafe {
            cr_decrypt(
                ciphertext.as_ptr(),
                ciphertext.len(),
                false, // no_errors
                &mut out_len,
                &mut dec_err,
            )
        };

        assert!(dec_err.is_null(), "decrypt error: {:?}", unsafe {
            if dec_err.is_null() {
                None
            } else {
                Some(CStr::from_ptr(dec_err))
            }
        });
        assert!(!plain_ptr.is_null());
        assert_eq!(out_len, msg.len());

        // compare plaintext
        let plain = unsafe { std::slice::from_raw_parts(plain_ptr, out_len) };
        assert_eq!(plain, msg);

        // free memory returned by decrypt (uses libc malloc)
        unsafe { free(plain_ptr as *mut _) };
    }

    // ---------------------------------------------------------------
    // 5. invalid ciphertext -> error
    // ---------------------------------------------------------------
    #[test]
    fn test_decrypt_invalid_ciphertext() {
        let garbage = [0u8; 32]; // any junk bytes are fine; doesn't need RNG
        let mut out_len = 0usize;
        let mut err_ptr: *mut c_char = ptr::null_mut();
        let plain_ptr = unsafe {
            cr_decrypt(
                garbage.as_ptr(),
                garbage.len(),
                false, // capture errors
                &mut out_len,
                &mut err_ptr,
            )
        };

        assert!(plain_ptr.is_null());
        assert_eq!(out_len, 0);
        assert!(!err_ptr.is_null(), "error message expected");

        unsafe {
            // optional assertion on message contents
            let msg = CStr::from_ptr(err_ptr).to_string_lossy();
            assert!(
                msg.contains("Error") || msg.contains("error"),
                "unexpected msg: {msg}"
            );
            let _ = CString::from_raw(err_ptr); // drop
        }
    }

    // ---------------------------------------------------------------
    // 6. invalid ciphertext with no_errors=true
    // ---------------------------------------------------------------
    #[test]
    fn test_decrypt_no_errors() {
        let garbage = [0u8; 16];
        let mut out_len = 0usize;
        let mut err_ptr: *mut c_char = ptr::null_mut();
        let plain_ptr = unsafe {
            cr_decrypt(
                garbage.as_ptr(),
                garbage.len(),
                true, // suppress errors
                &mut out_len,
                &mut err_ptr,
            )
        };

        assert!(plain_ptr.is_null());
        assert_eq!(out_len, 0);
        assert!(
            err_ptr.is_null(),
            "err_out must stay NULL when no_errors=true"
        );
    }

    // ---------------------------------------------------------------
    // 7. decrypt: NULL pointer input
    // ---------------------------------------------------------------
    #[test]
    fn test_decrypt_null_ptr() {
        let mut out_len: usize = 0;
        let mut err_ptr: *mut c_char = ptr::null_mut();

        let plain_ptr = unsafe {
            cr_decrypt(
                ptr::null(), // bad pointer
                10,
                false,
                &mut out_len,
                &mut err_ptr,
            )
        };

        assert!(plain_ptr.is_null());
        assert_eq!(out_len, 0);
        assert!(!err_ptr.is_null());

        unsafe { drop_cstring(err_ptr) };
    }

    // ---------------------------------------------------------------
    // 8. latest round should normally succeed (network-dependent)
    // ---------------------------------------------------------------
    #[test]
    fn test_get_latest_round_success() {
        let mut err_ptr: *mut c_char = ptr::null_mut();
        let round = unsafe { cr_get_latest_round(&mut err_ptr) };

        if err_ptr.is_null() {
            // success path
            assert!(round > 0, "round must be >0 on success");
        } else {
            // network failure path: we only require proper error propagation
            assert_eq!(round, 0);
            unsafe { drop_cstring(err_ptr) };
        }
    }

    // ---------------------------------------------------------------
    // 9. encrypt_commitment happy path
    // ---------------------------------------------------------------
    #[test]
    fn test_encrypt_commitment_success() {
        let msg = "Bittensor FTW";
        let mut round: u64 = 0;
        let mut err_ptr: *mut c_char = ptr::null_mut();

        let buf = unsafe {
            cr_encrypt_commitment(
                msg.as_ptr(),
                msg.len(),
                5,    // blocks_until_reveal
                12.0, // block_time
                &mut round,
                &mut err_ptr,
            )
        };

        assert!(err_ptr.is_null());
        assert!(round > 0);
        assert!(!buf.ptr.is_null());
        assert!(buf.len > 0);

        unsafe { cr_free(buf) };
    }

    // ---------------------------------------------------------------
    // 10. encrypt_commitment with non-UTF8 bytes
    // ---------------------------------------------------------------
    #[test]
    fn test_encrypt_commitment_utf8_error() {
        let bad_bytes = [0xffu8, 0xfe, 0xfd];
        let mut round: u64 = 0;
        let mut err_ptr: *mut c_char = ptr::null_mut();

        let buf = unsafe {
            cr_encrypt_commitment(
                bad_bytes.as_ptr(),
                bad_bytes.len(),
                5,
                12.0,
                &mut round,
                &mut err_ptr,
            )
        };

        assert!(buf.ptr.is_null());
        assert_eq!(buf.len, 0);
        assert_eq!(round, 0);
        assert!(!err_ptr.is_null());

        unsafe { drop_cstring(err_ptr) };
    }

    // ---------------------------------------------------------------
    // 11. generate_commit success
    // ---------------------------------------------------------------
    #[test]
    fn test_generate_commit_success() {
        let uids: [u16; 3] = [1, 2, 3];
        let vals: [u16; 3] = [10, 20, 30];

        let mut round: u64 = 0;
        let mut err_ptr: *mut c_char = ptr::null_mut();

        let buf = unsafe {
            cr_generate_commit(
                uids.as_ptr(),
                uids.len(),
                vals.as_ptr(),
                vals.len(),
                42,     // version_key
                20,     // tempo
                10_000, // current_block
                1,      // netuid
                2,      // subnet_reveal_epochs
                12.0,   // block_time
                &mut round,
                &mut err_ptr,
            )
        };

        assert!(err_ptr.is_null(), "err_out should be NULL on success");
        assert!(round > 0);
        assert!(!buf.ptr.is_null() && buf.len > 0);

        unsafe { cr_free(buf) };
    }

    // ---------------------------------------------------------------
    // 12. generate_commit with NULL uids pointer
    // ---------------------------------------------------------------
    #[test]
    fn test_generate_commit_null_uids() {
        let vals: [u16; 2] = [1, 2];
        let mut round: u64 = 0;
        let mut err_ptr: *mut c_char = ptr::null_mut();

        let buf = unsafe {
            cr_generate_commit(
                ptr::null(), // NULL uids
                2,           // non-zero len
                vals.as_ptr(),
                vals.len(),
                0,
                0,
                0,
                0,
                0,
                12.0,
                &mut round,
                &mut err_ptr,
            )
        };

        assert!(buf.ptr.is_null());
        assert!(!err_ptr.is_null());
        unsafe { drop_cstring(err_ptr) };
    }

    // ---------------------------------------------------------------
    // 13. generate_commit mismatched lengths
    // ---------------------------------------------------------------
    #[test]
    fn test_generate_commit_mismatched_lengths() {
        let uids: [u16; 2] = [1, 2];
        let vals: [u16; 3] = [10, 20, 30];
        let mut round: u64 = 0;
        let mut err_ptr: *mut c_char = ptr::null_mut();

        let buf = unsafe {
            cr_generate_commit(
                uids.as_ptr(),
                uids.len(),
                vals.as_ptr(),
                vals.len(),
                0,
                0,
                0,
                0,
                0,
                12.0,
                &mut round,
                &mut err_ptr,
            )
        };

        // implementation may error either inside Rust or FFI guard
        if err_ptr.is_null() {
            // unexpected success – treat as failure but avoid leak
            panic!("expected error on mismatched lengths");
        } else {
            assert!(buf.ptr.is_null());
            unsafe { drop_cstring(err_ptr) };
        }
    }

    // ---------------------------------------------------------------
    // 14. double free on NULL buffer (should NO-OP)
    // ---------------------------------------------------------------
    #[test]
    fn test_double_free_no_crash() {
        let null_buf = CRByteBuffer {
            ptr: ptr::null_mut(),
            len: 0,
            cap: 0,
        };
        unsafe {
            cr_free(null_buf);
            cr_free(null_buf); // second call must not crash
        }
    }

    // ---------------------------------------------------------------
    // 15. large input (10 MB) round-trip
    // ---------------------------------------------------------------
    #[test]
    fn test_large_input() {
        let mut data = vec![0u8; 10 * 1024 * 1024]; // 10 MB zeros
        for (i, b) in data.iter_mut().enumerate() {
            *b = (i % 256) as u8; // non-trivial pattern
        }

        let mut round: u64 = 0;
        let mut err: *mut c_char = ptr::null_mut();
        let buf = unsafe { cr_encrypt(data.as_ptr(), data.len(), 0, 1.0, &mut round, &mut err) };
        assert!(err.is_null());
        let ct = unsafe { std::slice::from_raw_parts(buf.ptr, buf.len) }.to_vec();
        unsafe { cr_free(buf) };

        let mut out_len = 0usize;
        let mut derr: *mut c_char = ptr::null_mut();
        let plain_ptr =
            unsafe { cr_decrypt(ct.as_ptr(), ct.len(), false, &mut out_len, &mut derr) };
        assert!(derr.is_null());
        assert_eq!(out_len, data.len());
        let plain = unsafe { std::slice::from_raw_parts(plain_ptr, out_len) };
        assert_eq!(plain, &data[..]);
        unsafe { free(plain_ptr as *mut _) };
    }

    // ---------------------------------------------------------------
    // 16. threaded usage stress
    // ---------------------------------------------------------------
    #[test]
    fn test_threaded_usage() {
        const THREADS: usize = 8;
        let handles: Vec<_> = (0..THREADS)
            .map(|i| {
                std::thread::spawn(move || {
                    let msg = format!("thread-{i}-payload");
                    let mut round = 0u64;
                    let mut err: *mut c_char = ptr::null_mut();
                    let buf = unsafe {
                        cr_encrypt(msg.as_ptr(), msg.len(), 0, 1.0, &mut round, &mut err)
                    };
                    assert!(err.is_null());
                    let ct = unsafe { std::slice::from_raw_parts(buf.ptr, buf.len) }.to_vec();
                    unsafe { cr_free(buf) };

                    let mut out_len = 0usize;
                    let mut derr: *mut c_char = ptr::null_mut();
                    let plain_ptr = unsafe {
                        cr_decrypt(ct.as_ptr(), ct.len(), false, &mut out_len, &mut derr)
                    };
                    assert!(derr.is_null());
                    let plain = unsafe { std::slice::from_raw_parts(plain_ptr, out_len) };
                    assert_eq!(plain, msg.as_bytes());
                    unsafe { free(plain_ptr as *mut _) };
                })
            })
            .collect();

        for h in handles {
            h.join().expect("thread panicked");
        }
    }

    // ---------------------------------------------------------------
    // 17. decrypting ciphertext whose reveal-round is far in the future
    // ---------------------------------------------------------------
    #[test]
    fn test_decrypt_future_round() {
        let msg = b"future round payload";
        let mut round: u64 = 0;
        let mut err_ptr: *mut c_char = ptr::null_mut();

        // Encrypt for a reveal round ≈14 days in the future (100 000 × 12 s).
        let buf = unsafe {
            cr_encrypt(
                msg.as_ptr(),
                msg.len(),
                100_000, // n_blocks
                12.0,    // block_time
                &mut round,
                &mut err_ptr,
            )
        };
        assert!(err_ptr.is_null(), "encryption must succeed");
        assert!(!buf.ptr.is_null());

        // Copy ciphertext, then release original buffer.
        let ct = unsafe { std::slice::from_raw_parts(buf.ptr, buf.len) }.to_vec();
        unsafe { cr_free(buf) };

        // Attempt to decrypt immediately – should fail (signature not available yet).
        let mut out_len: usize = 0;
        let mut dec_err: *mut c_char = ptr::null_mut();
        let plain_ptr = unsafe {
            cr_decrypt(
                ct.as_ptr(),
                ct.len(),
                false, // capture errors
                &mut out_len,
                &mut dec_err,
            )
        };

        assert!(plain_ptr.is_null(), "decryption must fail for future round");
        assert_eq!(out_len, 0);
        assert!(!dec_err.is_null(), "err_out should contain message");

        unsafe { drop_cstring(dec_err) };
    }
}
