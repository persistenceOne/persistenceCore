# persistenceCore

[![Go Report Card](https://goreportcard.com/badge/github.com/persistenceOne/persistenceCore)](https://goreportcard.com/report/github.com/persistenceOne/persistenceCore)
[![license](https://img.shields.io/github/license/persistenceOne/persistenceCore.svg)](https://github.com/persistenceOne/persistenceCore/blob/master/LICENSE)
[![LoC](https://tokei.rs/b1/github/persistenceOne/persistenceCore)](https://github.com/persistenceOne/persistenceCore)


Application implementing the minimum clique of PersistenceSDK modules enabling interNFT definition, issuance, ownership transfer and decentralized exchange.

## Talk to us!
*   [Telegram](https://t.me/PersistenceOneChat)
*   [Twitter](https://twitter.com/PersistenceOne)

## SetUp:

### Prerequisite

Install goLang 1.15+

### To connect to a chain:

1. `persistenceNode init [node_name]`
2. Replace `${HOME}/.persistenceCore/config/genesis.json` with the genesis file of the chain.
3. Add `persistent_peers` or `seeds` in `${HOME}/.persistenceCore/config/config.toml`
4. Start the chain: `persistenceNode start`

### To add keys

`persistenceClient keys add [key_name]`

or

`persistenceClient keys add [key_name] --recover` (to give your own mnemonics)

### To start a new chain
1. Initialize: `persistenceNode init [node_name] --chain-id [chain_name]`
2. Add key for genesis account `persistenceClient keys add [genesis_key_name]`
3. Add genesis account `persistenceNode add-genesis-account [genesis_key_name] 10000000000000000000stake`
4. Create a validator at genesis `persistenceNode gentx --name [genesis_key_name]`
5. Collect gentxs `persistenceNode collect-gentxs`
6. Start the chain `persistenceNode start`
7. To start rest server `persistenceClient rest-server --chain-id=test --node=tcp://0.0.0.0:26657 --laddr=tcp://0.0.0.0:1317`