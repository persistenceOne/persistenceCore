#!/bin/bash

VALIDATOR_CONFIG="configs/validators.json"
KEYS_CONFIG="configs/keys.json"
# Set home to chain dir for easy setup

# Variables
COINS="100000000000000000stake"

set -eu

jq -r ".genesis[0].mnemonic" $VALIDATOR_CONFIG | $CHAIN_BIN init $CHAIN_ID --chain-id $CHAIN_ID --recover
jq -r ".genesis[0].mnemonic" $VALIDATOR_CONFIG | $CHAIN_BIN keys add $(jq -r ".genesis[0].name" $VALIDATOR_CONFIG) --recover --keyring-backend="test"

# Add keys to keyringg
for ((i=0; i<$(jq -r '.validators | length' $VALIDATOR_CONFIG); i++))
do
  jq -r ".validators[$i].mnemonic" $VALIDATOR_CONFIG | $CHAIN_BIN keys add $(jq -r ".validators[$i].name" $VALIDATOR_CONFIG) --recover --keyring-backend="test"
done

for ((i=0; i<$(jq -r '.keys | length' $KEYS_CONFIG); i++))
do
  jq -r ".keys[$i].mnemonic" $KEYS_CONFIG | $CHAIN_BIN keys add $(jq -r ".keys[$i].name" $KEYS_CONFIG) --recover --keyring-backend="test"
done

# Provide genesis validator self deligations
$CHAIN_BIN add-genesis-account $(jq -r .genesis[0].address $VALIDATOR_CONFIG) $COINS --keyring-backend="test"

# Give Validator addresses initial coins
for ((i=0; i<$(jq -r '.validators | length' $VALIDATOR_CONFIG); i++))
do
  $CHAIN_BIN add-genesis-account $(jq -r .validators[$i].address $VALIDATOR_CONFIG) $COINS --keyring-backend="test"
done

for ((i=0; i<$(jq -r '.keys | length' $KEYS_CONFIG); i++))
do
  $CHAIN_BIN add-genesis-account $(jq -r .keys[$i].address $KEYS_CONFIG) $COINS --keyring-backend="test"
done

$CHAIN_BIN gentx $(jq -r ".genesis[0].name" $VALIDATOR_CONFIG) 5000000000stake --keyring-backend="test" --chain-id $CHAIN_ID
echo "Output of gentx"
cat $HOME/.persistenceCore/config/gentx/*.json | jq

echo "Running collect-gentxs"
$CHAIN_BIN collect-gentxs

echo "Update app.toml file"
sed -i 's#keyring-backend = "os"#keyring-backend = "test"#g' $HOME/.persistenceCore/config/client.toml

echo "Update config.toml file"
sed -i 's#"tcp://127.0.0.1:26657"#"tcp://0.0.0.0:26657"#g' $HOME/.persistenceCore/config/config.toml
sed -i 's/timeout_commit = "5s"/timeout_commit = "1s"/g' $HOME/.persistenceCore/config/config.toml
sed -i 's/timeout_propose = "3s"/timeout_propose = "1s"/g' $HOME/.persistenceCore/config/config.toml
sed -i 's/index_all_keys = false/index_all_keys = true/g' $HOME/.persistenceCore/config/config.toml

echo "Update genesis.json file with updated local params"
jq -r '.app_state.staking.params.unbonding_time |= "30s"' $HOME/.persistenceCore/config/genesis.json > /tmp/trash/genesis.json; mv /tmp/trash/genesis.json $HOME/.persistenceCore/config/genesis.json
jq -r '.app_state.slashing.params.downtime_jail_duration |= "6s"' $HOME/.persistenceCore/config/genesis.json > /tmp/trash/genesis.json; mv /tmp/trash/genesis.json $HOME/.persistenceCore/config/genesis.json
jq -r '.app_state.gov.deposit_params.max_deposit_period |= "30s"' $HOME/.persistenceCore/config/genesis.json > /tmp/trash/genesis.json; mv /tmp/trash/genesis.json $HOME/.persistenceCore/config/genesis.json
jq -r '.app_state.gov.deposit_params.min_deposit[0].amount |= "10"' $HOME/.persistenceCore/config/genesis.json > /tmp/trash/genesis.json; mv /tmp/trash/genesis.json $HOME/.persistenceCore/config/genesis.json
jq -r '.app_state.gov.voting_params.voting_period |= "30s"' $HOME/.persistenceCore/config/genesis.json > /tmp/trash/genesis.json; mv /tmp/trash/genesis.json $HOME/.persistenceCore/config/genesis.json
jq -r '.app_state.gov.tally_params.quorum |= "0.000000000000000000"' $HOME/.persistenceCore/config/genesis.json > /tmp/trash/genesis.json; mv /tmp/trash/genesis.json $HOME/.persistenceCore/config/genesis.json
jq -r '.app_state.gov.tally_params.threshold |= "0.000000000000000000"' $HOME/.persistenceCore/config/genesis.json > /tmp/trash/genesis.json; mv /tmp/trash/genesis.json $HOME/.persistenceCore/config/genesis.json
jq -r '.app_state.gov.tally_params.veto_threshold |= "0.000000000000000000"' $HOME/.persistenceCore/config/genesis.json > /tmp/trash/genesis.json; mv /tmp/trash/genesis.json $HOME/.persistenceCore/config/genesis.json

$CHAIN_BIN tendermint show-node-id
