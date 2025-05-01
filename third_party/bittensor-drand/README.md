# Usage
Python package `bittensor_drand` has one function.

```python
from bittensor_drand import get_encrypted_commit
```

To test the function in your terminal:
1. Spin up a local subtensor branch which includes CR3
2. Create a subnet with netuid 1 (or replace the netuid with the one you create)
```bash
mkdir test
cd test
python3 -m venv venv
. venv/bin/activate
pip install maturin bittensor ipython
cd ..

maturin develop
ipython

```

then copy-past to ipython
```python
import numpy as np
import bittensor_drand as crv3
from bittensor.utils.weight_utils import convert_weights_and_uids_for_emit
import bittensor as bt

uids = [1, 3]
weights = [0.3, 0.7]
version_key = 843000
netuid = 1

subtensor = bt.Subtensor("local")

subnet_reveal_period_epochs = subtensor.get_subnet_reveal_period_epochs(
        netuid=netuid
    )
tempo = subtensor.get_subnet_hyperparameters(netuid).tempo
current_block = subtensor.get_current_block()

if isinstance(uids, list):
    uids = np.array(uids, dtype=np.int64)
if isinstance(weights, list):
    weights = np.array(weights, dtype=np.float32)

uids, weights = convert_weights_and_uids_for_emit(uids, weights)

print(crv3.get_encrypted_commit(uids, weights, version_key, tempo, current_block, netuid, subnet_reveal_period_epochs))
```
expected result
```python
(b'\xb9\x96\xe4\xd1\xfd\xabm\x8cc\xeb\xe3W\r\xc7J\xb4\xea\xa9\xd5u}OG~\xae\xcc\x9a@\xdf\xee\x16\xa9\x0c\x8d7\xd6\xea_c\xc2<\xcb\xa6\xbe^K\x97|\x16\xc6|;\xb5Z\x97\xc9\xb4\x8em\xf1hv\x16\xcf\xea\x1e7\xbe-Z\xe7e\x1f$\n\xf8\x08\xcb\x18.\x94V\xa3\xd7\xcd\xc9\x04F::\t)Z\xc6\xbey \x00\x00\x00\x00\x00\x00\x00\xaaN\xe8\xe97\x8f\x99\xbb"\xdf\xad\xf6\\#%\xca:\xc2\xce\xf9\x96\x9d\x8f\x9d\xa2\xad\xfd\xc73j\x16\xda \x00\x00\x00\x00\x00\x00\x00\x84*\xb0\rw\xad\xdc\x02o\xf7i)\xbb^\x99e\xe2\\\xee\x02NR+-Q\xcd \xf7\x02\x83\xffV>\x00\x00\x00\x00\x00\x00\x00"\x00\x00\x00\x00\x00\x00\x00*\x13wXb\x93\xc5"F\x17F\x05\xcd\x15\xb0=\xe2d\xfco3\x16\xfd\xe9\xc6\xbc\xd1\xb3Y\x97\xf9\xb9!\x01\x0c\x00\x00\x00\x00\x00\x00\x00X\xa2\x8c\x18Wkq\xe5\xe6\x1c2\x86\x08\x00\x00\x00\x00\x00\x00\x00AES_GCM_', 13300875)
```

To test this in a local subnet:
1. Spin up a local node based on the subtensor branch `spiigot/add-pallet-drand` using command `./scripts/localnet.sh False`
2. Create a subnet
3. Change the following hyperparameters:
    - `commit_reveal_weights_enabled` -> `True`
    - `tempo` -> 10 (keep in mind that you need to provide this as `tempo` argument to `get_encrypted_commit` function. Use polkadot website for this action.)
    - `weights_rate_limit` -> 0 (Reduces the limit when you can set weights.)
4. Register 1 or more wallets to the subnet
5. Create and activate python virtual environment (`python3 -m venv venv && . venv/bin/activate`)
6. Checkout bittensor `feat/roman/cr-v-3` branch.
7. Install bittensor `pip install -e .`
8. Cd to directory you cloned `https://github.com/opentensor/bittensor-commit-reveal/tree/staging` (FFI for CRv3).
9. Install the `maturin` python package and build/install `bittensor-commit-reveal` package to your env using the command `pip install maturin && maturin develop`
10. Run the following script within your python environment:
```python
import requests
import time

from bittensor import Subtensor, logging, Wallet

DRAND_API_BASE_URL_Q = "https://api.drand.sh/52db9ba70e0cc0f6eaf7803dd07447a1f5477735fd3f661792ba94600c84e971"

logging.set_info()


def get_drand_info(uri):
    """Fetch Drand network information."""
    url = f"{uri}/info"
    response = requests.get(url)
    response.raise_for_status()
    return response.json()


def get_current_round(info):
    """Calculate the current round based on genesis_time and period."""
    current_time = int(time.time())
    genesis_time = info["genesis_time"]
    period = info["period"]
    return (current_time - genesis_time) // period + 1


def main():
    sub = Subtensor("local")

    uids = [0]
    weights = [0.7]

    wallet = Wallet()  # corresponds the subnet owner wallet

    result, message = sub.set_weights(
        wallet=wallet,
        netuid=1,
        uids=uids,
        weights=weights,
        wait_for_inclusion=True,
        wait_for_finalization=True,
    )
    logging.info(f">>> Success, [blue]{result}[/blue], message: [magenta]{message}[/magenta]")

    reveal_round = int(message.split(":")[-1])
    # Fetch Drand network info
    for uri in [DRAND_API_BASE_URL_Q]:
        print(f"Fetching info from {uri}...")
        info = get_drand_info(uri)
        print("Info:", info)

        while True:
            time.sleep(info["period"])
            current_round = get_current_round(info)
            logging.console.info(f"Current round: [yellow]{current_round}[/yellow]")
            if current_round == reveal_round:
                logging.console.warning(f">>> It's time to reveal the target round: [blue]{reveal_round}[/blue]")

                break


if __name__ == "__main__":
    main()
```
11. Wait until your target_round comes.

12. Check your weights with the following code:

```python
import bittensor as bt

sub = bt.Subtensor(network="local")

netuid = 1  # your created subnet's netuid

print(sub.weights(netuid=netuid))
```
13. You can now see the same weights which you committed earlier