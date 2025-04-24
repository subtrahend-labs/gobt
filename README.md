# gobt

[![Go Reference](https://pkg.go.dev/badge/github.com/subtrahend-labs/gobt.svg)](https://pkg.go.dev/github.com/subtrahend-labs/gobt) [![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](./LICENSE)

> **Go bindings** for Subtensor: wrap extrinsics, RPC calls, and storage so you can interact with Bittensor in Golang

## Table of Contents

- [Overview](#overview)  
- [Legend](#legend)  
- [Installation](#installation)  
- [Features](#features)
- [Contributing](#contributing)  
- [License](#license)  
- [Disclaimer](#disclaimer)

---

## Overview

**gobt** wraps Substrate extrinsics, RPC endpoints, and storage queries in pure Go.  
Whether you need to submit transactions, fetch on-chain data structures (like Metagraphs and Neurons), or drive CLIs, gobt structures everything into:

- Strongly-typed Go structs  
- Ergonomic client API  
- Reusable, testable modules  

---

## Legend

This library is incomplete and under active development. Included are 
markdown files that list the extrinsics and rpc calls that remain to be 
implemented.  We use the following legend to indicate the status of each item: 

- **[-]** `won't implement`  
- **[x]** `implemented and tested`  
- **[o]** `not tested`

---

## Installation

**gobt** is a pure Go library, so you can install it with `go get`:

```bash
go get github.com/subtrahend-labs/gobt
```

---

## Features

- 🚀 **Extrinsics:** high-level Go functions for all your pallet calls  
- 🔗 **RPC wrappers:** fetch `Metagraph`, `Neurons`, and all other on-chain data
- 🏗️ **Runtime types:** auto-aligned SCALE-decodeable structs  
- 📦 **Modular:** pick & choose extrinsics, RPC call, and storage helpers  
- 🧪 **Testable:** built-in unit & integration tests for core functionality  

--- 

## Contributing

1. Fork the repo
2. Create a feature branch
3. Add tests & documentation
4. Open a pull request

Please check existing issues before filing new ones.

---

## License 

This project is licensed under the MIT License. 

--- 

<a name="disclaimer"></a>
## ⚠️ Disclaimer 

**USE AT YOUR OWN RISK.** This library is provided “as-is,” with **no warranties** of any kind, express or implied, including but not limited to merchantability, fitness for a particular purpose, or non-infringement. The authors, contributors, and maintainers of this project **expressly disclaim all liability** for any direct, indirect, incidental, special, consequential, or punitive damages arising out of your use of, or inability to use, this software—even if advised of the possibility of such damages.  

By using this library to interact with any blockchain, you acknowledge and agree that you bear full responsibility for your actions, transactions, and any network fees or losses incurred. You should review, test, and audit this code thoroughly before using it in any production or value-sensitive environment.


