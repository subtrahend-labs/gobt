CR_DRAND := third_party/bittensor-drand

.PHONY: crv3
crv3:
	cd $(CR_DRAND) && \
	cargo build --release --no-default-features

.PHONY: clean-crv3
clean-crv3:
	cd $(CR_DRAND) && cargo clean

