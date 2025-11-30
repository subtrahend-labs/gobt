# Makefile for gobt project

CR_DRAND := third_party/bittensor-drand

# Build CRv4 bindings without Python dependencies
.PHONY: crv4
crv4:
	cd $(CR_DRAND) && \
	cargo build --release --no-default-features

# Clean CRv4 build artifacts
.PHONY: clean-crv4
clean-crv4:
	cd $(CR_DRAND) && cargo clean

# Generate C bindings header
.PHONY: bindings
bindings:
	cd $(CR_DRAND) && \
	cbindgen --crate bittensor_drand --output bindings.h

# Build everything
.PHONY: all
all: crv4

# Clean everything
.PHONY: clean
clean: clean-crv4