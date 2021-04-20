# persistenceCore

[![LoC](https://tokei.rs/b1/github/persistenceOne/persistenceCore)](https://github.com/persistenceOne/persistenceCore)

This project implements an application for the Persistence Core chain that all the other chains in the ecosystem connect to as a raised and open moderator for interoperability, shared security, and as a gateway to other ecosystems and chains.

## Talk to us!
*   [Twitter](https://twitter.com/PersistenceOne)
*   [Telegram](https://t.me/PersistenceOneChat)
*   [Discord](https://discord.com/channels/796174129077813248)

## Hardware Requirements 
* **Minimal**
    * 1 GB RAM
    * 25 GB HDD
    * 1.4 GHz CPU
* **Recommended**
    * 2 GB RAM
    * 100 GB HDD
    * 2.0 GHz x2 CPU

> NOTE: SSDs have limited TBW before non-catastrophic data errors. Running a full node requires a TB+ writes per day, causing rapid deterioration of SSDs over HDDs of comparable quality.

## Operating System
* Linux/Windows/MacOS(x86)
* **Recommended**
    * Linux(x86_64)

## Installation Steps
>Prerequisite: go1.15+ required. [ref](https://golang.org/doc/install)

>Prerequisite: git. [ref](https://github.com/git/git)

>Optional requirement: GNU make. [ref](https://www.gnu.org/software/make/manual/html_node/index.html)


* Clone git repository
```shell
git clone https://github.com/persistenceOne/persistenceCore.git
```
* Checkout release tag
```shell
git fetch --tags
git checkout [vX.X.X]
```
* Install
```shell
cd persistenceCore
make all
```

### Generate keys

`persistenceCore keys add [key_name]`

or

`persistenceCore keys add [key_name] --recover` to regenerate keys with your [BIP39](https://github.com/bitcoin/bips/tree/master/bip-0039) mnemonic

### Connect to a chain and start node
* [Install](#installation-steps) persistenceCore application
* Initialize node
```shell
persistenceCore init [NODE_NAME]
```
* Replace `${HOME}/.persistenceCore/config/genesis.json` with the genesis file of the chain.
* Add `persistent_peers` or `seeds` in `${HOME}/.persistenceCore/config/config.toml`
* Start node
```shell
persistenceCore start
```

### Initialize a new chain and start node 
* Initialize: `persistenceCore init [node_name] --chain-id [chain_name]`
* Add key for genesis account `persistenceCore keys add [genesis_key_name]`
* Add genesis account `persistenceCore add-genesis-account [genesis_key_name] 10000000000000000000stake`
* Create a validator at genesis `persistenceCore gentx [genesis_key_name] 10000000stake --chain-id [chain_name]`
* Collect genesis transactions `persistenceCore collect-gentxs`
* Start node `persistenceCore start`
* To start rest server set `enable=true` in `config/app.toml` under `[api]` and restart the chain

### Ledger Support( Experimental)

> NOTE: This is an experimental feature. Persistence uses coin-type 750; generating keys through this method below will create keys with coin-type 118(cosmos) and will only be supported by CLI and not by current or future wallets. Ledger support for the Persistence application is coming soon.  

* Install the Cosmos application on the Ledger device. [ref](https://hub.cosmos.network/main/resources/ledger.html#install-the-cosmos-ledger-application)
* Connect the Ledger device to a system with persistenceCore binary and open the Cosmos application on it.
* Add key
```shell
persistenceCore keys add [key_name] --ledger
```
* Sign transaction
```shell
persistenceCore tx [transaction parameters] --ledger
```

### Reset chain
```shell
rm -rf ~/.persistenceCore
```

### Shutdown node
```shell
killall persistenceCore
```

### Check version
```shell
persistenceCore version
```

## Test-nets
* [test-core-1](https://github.com/persistenceOne/genesisTransactions/tree/master/test-core-1)

## Main-net
* [core-1](https://github.com/persistenceOne/genesisTransactions/tree/master/core-1)
