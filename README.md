# persistenceCore

[![LoC](https://tokei.rs/b1/github/persistenceOne/persistenceCore)](https://github.com/persistenceOne/persistenceCore)

This project implements an application for the Persistence Core chain that all the other chains in the ecosystem connect to as a raised and open moderator for interoperability, shared security, and as a gateway to other ecosystems and chains.

## Talk to us!
*   [Telegram](https://t.me/PersistenceOneChat)
*   [Twitter](https://twitter.com/PersistenceOne)

## SetUp:

### Prerequisite

Install goLang 1.15+

### To connect to a chain:

1. `persistenceCore init [node_name]`
2. Replace `${HOME}/.persistenceCore/config/genesis.json` with the genesis file of the chain.
3. Add `persistent_peers` or `seeds` in `${HOME}/.persistenceCore/config/config.toml`
4. Start the chain: `persistenceCore start`

### To add keys

`persistenceCore keys add [key_name]`

or

`persistenceCore keys add [key_name] --recover` (to give your own mnemonics)

### To start a new chain
1. Initialize: `persistenceCore init [node_name] --chain-id [chain_name]`
2. Add key for genesis account `persistenceCore keys add [genesis_key_name]`
3. Add genesis account `persistenceCore add-genesis-account [genesis_key_name] 10000000000000000000stake`
4. Create a validator at genesis `persistenceCore gentx [genesis_key_name] 10000000stake --chain-id [chain_name]`
5. Collect gentxs `persistenceCore collect-gentxs`
6. Start the chain `persistenceCore start`
7. To start rest server set `enable=true` in `config/app.toml` under `[api]` and restart the chain
