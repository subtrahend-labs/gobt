from typing import Union, Optional

from bittensor_drand.bittensor_drand import (
    get_encrypted_commit as _get_encrypted_commit,
    get_encrypted_commitment as _get_encrypted_commitment,
    encrypt as _encrypt,
    decrypt as _decrypt,
    get_latest_round as _get_latest_round,
)


def get_encrypted_commit(
    uids: list[int],
    weights: list[int],
    version_key: int,
    tempo: int,
    current_block: int,
    netuid: int,
    subnet_reveal_period_epochs: int,
    block_time: int = 12,
) -> tuple[bytes, int]:
    """Returns encrypted commit and target round for `commit_crv3_weights` extrinsic.

    Arguments:
        uids: The uids to commit.
        weights: The weights associated with the uids.
        version_key: The version key to use for committing and revealing. Default is `bittensor.core.settings.version_as_int`.
        tempo: Number of blocks in one epoch.
        current_block: The current block number in the network.
        netuid: The network unique identifier (NetUID) for the subnet.
        subnet_reveal_period_epochs: Number of epochs after which the reveal will be performed. Corresponds to the hyperparameter `commit_reveal_weights_interval` of the subnet. In epochs.
        block_time: Amount of time in seconds for one block. Defaults to 12 seconds.

    Returns:
        commit (bytes): Raw bytes of the encrypted, and compressed uids & weights values for setting weights.
        target_round (int): Drand round number when weights have to be revealed. Based on Drand Quicknet network.

    Raises:
        ValueError: If the input parameters are invalid or encryption fails.
    """
    return _get_encrypted_commit(
        uids,
        weights,
        version_key,
        tempo,
        current_block,
        netuid,
        subnet_reveal_period_epochs,
        block_time,
    )


def get_encrypted_commitment(
    data: str, blocks_until_reveal: int, block_time: Union[int, float] = 12.0
) -> tuple[bytes, int]:
    """Encrypts arbitrary string data with time-lock encryption.

    Arguments:
        data: The string data to encrypt.
        blocks_until_reveal: Number of blocks until the data should be revealed.
        block_time: Amount of time in seconds for one block. Defaults to 12 seconds.

    Returns:
        encrypted_data (bytes): Raw bytes of the encrypted data.
        target_round (int): Drand round number when data can be revealed.

    Raises:
        ValueError: If encryption fails.
    """
    return _get_encrypted_commitment(data, blocks_until_reveal, block_time)


def encrypt(
    data: bytes, n_blocks: int, block_time: Union[int, float] = 12.0
) -> tuple[bytes, int]:
    """Encrypts arbitrary binary data with time-lock encryption.

    Arguments:
        data: The binary data to encrypt.
        n_blocks: Number of blocks until the data should be revealed.
        block_time: Amount of time in seconds for one block. Defaults to 12 seconds.

    Returns:
        encrypted_data (bytes): Raw bytes of the encrypted data.
        target_round (int): Drand round number when data can be revealed.

    Raises:
        ValueError: If encryption fails.
    """
    return _encrypt(data, n_blocks, block_time)


def decrypt(encrypted_data: bytes, no_errors: bool = True) -> Optional[bytes]:
    """Decrypts previously encrypted data if the reveal time has been reached.

    Arguments:
        encrypted_data: The encrypted data to decrypt.
        no_errors: If True, returns None instead of raising exceptions when decryption fails.
                  If False, raises exceptions on decryption failures.

    Returns:
        decrypted_data (Optional[bytes]): The decrypted data if successful, None otherwise.

    Raises:
        ValueError: If decryption fails and no_errors is False.
    """
    return _decrypt(encrypted_data, no_errors)


def get_latest_round() -> int:
    """Gets the latest revealed Drand round number.

    Returns:
        round (int): The latest revealed Drand round number.

    Raises:
        ValueError: If fetching the latest round fails.
    """
    return _get_latest_round()
