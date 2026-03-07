# Custom Virtual Machine (VM)

This directory contains the implementation of the custom virtual machine (VM) for executing smart contracts on the blockchain.

## Goals
- Support Turing-complete smart contracts
- Enable contract deployment, execution, and upgradability
- Integrate with transaction and block processing
- Provide hooks for formal verification and gas/fee metering
- Allow future extensibility (privacy, oracles, governance, etc.)

## Initial Architecture
- `vm.go`: Core VM logic and instruction set
- `contract.go`: Contract interface, deployment, and lifecycle
- `storage.go`: Persistent contract storage
- `execution.go`: Execution context and gas metering

## Integration Plan
- Contracts will be deployed and invoked via special transactions
- The VM will be called during block processing to execute contract logic
- Results will be persisted in contract storage and reflected in state 