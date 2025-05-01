use crate::drand;
use codec::{Decode, Encode};
use pyo3::exceptions::PyValueError;
use pyo3::prelude::*;
use pyo3::types::PyBytes;
use pyo3::{pyfunction, pymodule, wrap_pyfunction, Bound, PyResult, Python};
use std::time::{SystemTime, UNIX_EPOCH};

/// Returns a timelock-encrypted commitment and its corresponding Drand reveal round.
///
/// This function is used to generate an encrypted commitment based on UID weights and
/// subnet reveal parameters. The commitment will be decryptable only after a calculated
/// Drand round based on `tempo`, `current_block`, and `subnet_reveal_period_epochs`.
///
/// Args:
///     uids (List[int]): List of UID integers.
///     weights (List[int]): Corresponding list of weight values (same length as `uids`).
///     version_key (int): A version identifier for this commitment.
///     tempo (int): Block interval for the subnet (tempo parameter).
///     current_block (int): The current block number in the chain.
///     netuid (int): Subnet identifier.
///     subnet_reveal_period_epochs (int): Number of epochs to wait before decryption.
///     block_time (float, optional): Block time in seconds (default = 12.0).
///
/// Returns:
///     Tuple[bytes, int]: A tuple containing:
///         - the encrypted commitment (as bytes)
///         - the reveal round number (int) when it can be decrypted.
#[pyfunction]
#[pyo3(signature = (uids, weights, version_key, tempo, current_block, netuid, subnet_reveal_period_epochs, block_time=12.0))]
fn get_encrypted_commit(
    py: Python,
    uids: Vec<u16>,
    weights: Vec<u16>,
    version_key: u64,
    tempo: u64,
    current_block: u64,
    netuid: u16,
    subnet_reveal_period_epochs: u64,
    block_time: f64,
) -> PyResult<(Py<PyBytes>, u64)> {
    // create runtime to make async call
    let result = drand::generate_commit(
        uids,
        weights,
        version_key,
        tempo,
        current_block,
        netuid,
        subnet_reveal_period_epochs,
        block_time,
    );
    // matching the result
    match result {
        Ok((ciphertext, target_round)) => {
            let py_bytes = PyBytes::new_bound(py, &ciphertext).into();
            Ok((py_bytes, target_round))
        }
        Err(e) => Err(PyValueError::new_err(format!("{:?}", e))),
    }
}

/// Encrypts a string commitment with a timelock for a future Drand round.
///
/// This function encrypts arbitrary string data, ensuring it will be decryptable
/// only after a number of blocks pass. Useful for general-purpose timelock encryption.
///
/// Args:
///     data (str): The string to encrypt.
///     blocks_until_reveal (int): Number of blocks to wait before the data is decryptable.
///     block_time (float, optional): Block time in seconds (default = 12.0).
///
/// Returns:
///     Tuple[bytes, int]: A tuple containing:
///         - the encrypted bytes
///         - the Drand reveal round (int)
#[pyfunction]
#[pyo3(signature = (data, blocks_until_reveal, block_time=12.0))]
fn get_encrypted_commitment(
    py: Python,
    data: &str,
    blocks_until_reveal: u64,
    block_time: f64,
) -> PyResult<(Py<PyBytes>, u64)> {
    let result = drand::encrypt_commitment(data, blocks_until_reveal, block_time);
    // matching the result
    match result {
        Ok((ciphertext, target_round)) => {
            let py_bytes = PyBytes::new_bound(py, &ciphertext).into();
            Ok((py_bytes, target_round))
        }
        Err(e) => Err(PyValueError::new_err(format!("{:?}", e))),
    }
}

#[pyfunction(name = "get_latest_round")]
fn get_latest_round_py() -> PyResult<u64> {
    let response = drand::get_round_info(None)
        .map_err(|e| PyValueError::new_err(format!("Drand fetch error: {}", e)))?;

    Ok(response.round)
}

/// Encrypts binary data for a future Drand round based on block delay.
///
/// This method timelock-encrypts the provided binary `data`, such that it
/// becomes decryptable only after `n_blocks` have passed. Internally, it maps
/// this to a specific Drand round number.
///
/// Args:
///     data (bytes): Data to encrypt.
///     n_blocks (int): Number of blocks to wait before decryption is possible.
///     block_time (float, optional): Block time in seconds (default = 12.0).
///
/// Returns:
///     Tuple[bytes, int]: A tuple containing:
///         - the encrypted payload
///         - the Drand reveal round number
#[pyfunction]
#[pyo3(signature = (data, n_blocks, block_time=12.0))]
fn encrypt(
    py: Python,
    data: &[u8],
    n_blocks: u64,
    block_time: f64,
) -> PyResult<(Py<PyBytes>, u64)> {
    let now = SystemTime::now()
        .duration_since(UNIX_EPOCH)
        .map_err(|e| PyValueError::new_err(format!("SystemTime error: {:?}", e)))?
        .as_secs_f64();

    let reveal_timestamp = (n_blocks as f64 * block_time + now).ceil() as u64 - drand::GENESIS_TIME;
    let reveal_round = reveal_timestamp / drand::DRAND_PERIOD;

    let encrypted_data = drand::encrypt_and_compress(data, reveal_round)
        .map_err(|e| PyValueError::new_err(format!("Encryption failed: {:?}", e)))?;

    let encrypted_with_reveal_round = drand::UserData {
        encrypted_data,
        reveal_round,
    }
    .encode();

    Ok((
        PyBytes::new_bound(py, &encrypted_with_reveal_round).into(),
        reveal_round,
    ))
}

/// Attempts to decrypt data previously encrypted with Drand timelock encryption.
///
/// This function automatically extracts the reveal round from the encrypted message,
/// fetches the corresponding Drand signature (if available), and decrypts the message.
///
/// Args:
///     encrypted_data (bytes): Data previously returned from `encrypt` or `get_encrypted_commit`.
///     no_errors (bool, optional): If True, suppresses errors and returns None instead (default = True).
///
/// Returns:
///     Optional[bytes]: Decrypted data if successful, otherwise None or raises an error.
#[pyfunction]
#[pyo3(signature = (encrypted_data, no_errors=true))]
fn decrypt(py: Python, encrypted_data: &[u8], no_errors: bool) -> PyResult<Option<Py<PyBytes>>> {
    let user_data = match drand::UserData::decode(&mut &encrypted_data[..]) {
        Ok(data) => data,
        Err(e) => {
            return if no_errors {
                Ok(None)
            } else {
                Err(PyValueError::new_err(format!(
                    "Error deserializing data: {:?}",
                    e
                )))
            }
        }
    };

    let signature_opt = drand::get_reveal_round_signature(Some(user_data.reveal_round), no_errors)
        .map_err(|e| PyValueError::new_err(e))?;

    let signature_str = match signature_opt {
        Some(s) => s,
        None => {
            return if no_errors {
                Ok(None)
            } else {
                Err(PyValueError::new_err("Signature not available"))
            }
        }
    };

    let signature_bytes = hex::decode(signature_str)
        .map_err(|e| PyValueError::new_err(format!("Invalid hex in signature: {:?}", e)))?;

    let decoded_data = drand::decrypt_and_decompress(&user_data.encrypted_data, &signature_bytes)
        .map_err(|e| PyValueError::new_err(e))?;

    Ok(Some(PyBytes::new_bound(py, &decoded_data).into()))
}

#[pymodule]
fn bittensor_drand(m: &Bound<'_, PyModule>) -> PyResult<()> {
    m.add_function(wrap_pyfunction!(get_encrypted_commit, m)?)?;
    m.add_function(wrap_pyfunction!(get_encrypted_commitment, m)?)?;
    m.add_function(wrap_pyfunction!(encrypt, m)?)?;
    m.add_function(wrap_pyfunction!(decrypt, m)?)?;
    m.add_function(wrap_pyfunction!(get_latest_round_py, m)?)?;
    Ok(())
}
