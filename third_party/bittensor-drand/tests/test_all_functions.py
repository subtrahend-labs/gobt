import time
import bittensor_drand as btcr


def test_get_latest_round():
    round_ = btcr.get_latest_round()
    assert isinstance(round_, int)
    assert round_ > 0


def test_encrypt_and_decrypt():
    data = b"hello, bittensor!"
    n_blocks = 1

    encrypted, reveal_round = btcr.encrypt(data, n_blocks)
    assert isinstance(encrypted, bytes)
    assert isinstance(reveal_round, int)

    print(f"Reveal round: {reveal_round}")
    current_round = btcr.get_latest_round()

    if current_round < reveal_round:
        print("Waiting for reveal round to arrive...")
        while btcr.get_latest_round() < reveal_round:
            time.sleep(3)

    decrypted = btcr.decrypt(encrypted)
    assert decrypted is not None
    assert decrypted == data


def test_get_encrypted_commitment():
    encrypted, round_ = btcr.get_encrypted_commitment("my_commitment", 1)
    assert isinstance(encrypted, bytes)
    assert isinstance(round_, int)


def test_get_encrypted_commit():
    uids = [0, 1]
    weights = [100, 200]
    version_key = 1
    tempo = 10
    current_block = 100
    netuid = 1
    subnet_reveal_period_epochs = 2

    encrypted, round_ = btcr.get_encrypted_commit(
        uids,
        weights,
        version_key,
        tempo,
        current_block,
        netuid,
        subnet_reveal_period_epochs,
    )
    assert isinstance(encrypted, bytes)
    assert isinstance(round_, int)
