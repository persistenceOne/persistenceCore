# CHANGELOG

## [Unreleased]

### Overview

- Upgrades cosmos-sdk from `v0.45.14` to `v0.47.3` LSM fork created by persistence, including [barberry](https://forum.cosmos.network/t/cosmos-sdk-security-advisory-barberry/10825) security fix
- Migrates from tendermint `v0.34.26` to cometbft `v0.37.2`
- Upgrades ibc-go from `v4.3.1` to `v7.1.0` LSM fork including [huckleberry](https://forum.cosmos.network/t/ibc-security-advisory-huckleberry/10731) security fix
- Upgrades wasmd from `v0.30.0` to `v0.40.2` LSM fork & wasmvm to `1.2.4` including [cherry](https://github.com/CosmWasm/advisories/blob/main/CWAs/CWA-2023-002.md) bugfix
- Some SDK 47 things to keep in mind:
  - The SDK version includes some key store migration for the CLI. Make sure you backup your private keys before testing this! You can not switch back to v45 keys
  - CLI: `add-genesis-account`, `gentx`, `add-genesis-account`, `collect-gentxs` and others are now under `genesis` command as parent
  - CLI: `--broadcast-mode block` was removed. You need to query the result for a TX with `persistenceCore q tx <hash>` instead
  - ...add more?
- Upgrades persistence-sdk from `v2.0.1` to `vx.x.x`
- Upgrades pstake-native from `v2.0.0` to `vx.x.x`
- Adds wasm-bindings
- Modules integrated from `persistence-sdk`
  - IBC hooks
  - PFM
  - Oracle
- Modules integrated from `pstake-native`
  - liquidstakeibc - this deprecates `lscosmos` module
- Integrated POB module for MEV - disabled aution txs for now.

#### MinCommissionRate

- `MinCommissionRate` is set to 5%, which was proposed [here](https://www.mintscan.io/persistence/proposals/18)

    > **Note**  
    > During upgrade,  
    > Validator's `CommissionRate` will be set to 5%, if it is lower than the `MinCommissionRate` (i.e. 5%),  
    > and Validator's `MaxCommissionRate` will be set to 10% (if lower than 10%) to give validator some margin to work with.

### Changes

- ([#211](https://github.com/persistenceOne/persistenceCore/pull/211)) Enfoce minimum limit for `CommissionRate` & `MaxCommissionRate`
- ([#207](https://github.com/persistenceOne/persistenceCore/pull/207)) adds POB module for skip-mev
- ([#205](https://github.com/persistenceOne/persistenceCore/pull/205)) bump cosmos-sdk to `v0.47.3-lsm` and deps (includes new modules: IBC hooks, PFM, liquidstakeibc)
- ([#198](https://github.com/persistenceOne/persistenceCore/pull/198), [#206](https://github.com/persistenceOne/persistenceCore/pull/206)) starship e2e upgrade tests
- ([#184](https://github.com/persistenceOne/persistenceCore/pull/184)) removal of unused exposer
- ([#182](https://github.com/persistenceOne/persistenceCore/pull/182)) app restructure
- ([#179](https://github.com/persistenceOne/persistenceCore/pull/179), [#194](https://github.com/persistenceOne/persistenceCore/pull/194)) add wasm-bindings and integrate oracle module
- ([#170](https://github.com/persistenceOne/persistenceCore/pull/170)) fix: cleanup action release
