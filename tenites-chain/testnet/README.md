# Tenites Testnet — tenites-testnet-1

## Chain Info
| Parameter | Value |
|-----------|-------|
| Chain ID | tenites-testnet-1 |
| Binary | tenitesd |
| Cosmos SDK | v0.50.10 |
| CometBFT | v0.38.12 |
| IBC | ibc-go v8.8.0 — core module and transfer module wired |
| Go version | 1.21 |
| Native denom | utsc |
| Bech32 prefix | tenites |

## Network Endpoints (once deployed)
| Endpoint | URL |
|----------|-----|
| RPC | https://rpc.tenites-testnet.io |
| REST/LCD | https://lcd.tenites-testnet.io |
| gRPC | grpc.tenites-testnet.io:9090 |
| Explorer | https://explorer.tenites-testnet.io |

## Validator Setup
See `setup-validator.sh` for automated setup.
Manual steps:
1. Initialize node with correct chain ID
2. Copy genesis.json from this directory
3. Configure persistent peers
4. Start node via systemd

## Custom Modules
| Module | Purpose |
|--------|---------|
| x/tsc | TSC Stablecoin — mint/burn/transfer |
| x/wallet | Multi-currency wallet management |
| x/identity | TenitesID identity registry |
| x/compliance | On-chain compliance decisions |
| x/settlement | Settlement finality records |
| x/obligation | Payment obligation lifecycle |
| x/trust | Entity trust scoring |
| x/dispute | Settlement dispute resolution |
| x/audit | Immutable audit trail |
| x/governance | Extended governance proposals |
| x/policy | Policy-as-code enforcement |
| x/regulator | Regulatory view and reporting |
| x/intent | Intent-based routing and fills |
| x/valset | Permissioned validator set |
| x/aa | Account abstraction |

## IBC Configuration
| Parameter | Value |
|-----------|-------|
| IBC version | ibc-go v8.8.0 |
| Light client | 07-Tendermint |
| Transfer | ICS-20 enabled |
| Channels | None yet (needs relayer) |
| Osmosis channel | Pending (Phase 3) |

See `docs/ibc/hermes-config.toml` for Hermes relayer configuration.

## TSC Operations

### Mint
```bash
tenitesd tx tsc mint [authority-id] [wallet-id] \
  [amount] [purpose] [jurisdiction] \
  --from [key] --chain-id tenites-testnet-1
```

### Query Supply
```bash
tenitesd query tsc total-supply \
  --node https://rpc.tenites-testnet.io
```

## Genesis Checksum
Verify genesis file integrity before starting a validator:
```bash
sha256sum genesis.json
# Expected: db86316666393943d222ec1d2b1d18da22d9da614b24a8a678d5b1a905105bec
```
