#!/bin/bash

DENOM="${DENOM:=uxprt}"
CHAIN_DATA_DIR="${CHAIN_DATA_DIR:=.persistenceCore}"

VALIDATOR_CONFIG="configs/validators.json"
KEYS_CONFIG="configs/keys.json"

# Variables
COINS="100000000000000000$DENOM"

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
$CHAIN_BIN add-genesis-account $($CHAIN_BIN keys show -a $(jq -r .genesis[0].name $VALIDATOR_CONFIG) --keyring-backend="test") $COINS --keyring-backend="test"

# Give Validator addresses initial coins
for ((i=0; i<$(jq -r '.validators | length' $VALIDATOR_CONFIG); i++))
do
  $CHAIN_BIN add-genesis-account $($CHAIN_BIN keys show -a $(jq -r .validators[$i].name $VALIDATOR_CONFIG) --keyring-backend="test") $COINS --keyring-backend="test"
done

for ((i=0; i<$(jq -r '.keys | length' $KEYS_CONFIG); i++))
do
  $CHAIN_BIN add-genesis-account $($CHAIN_BIN keys show -a $(jq -r .keys[$i].name $KEYS_CONFIG) --keyring-backend="test") $COINS --keyring-backend="test"
done

$CHAIN_BIN gentx $(jq -r ".genesis[0].name" $VALIDATOR_CONFIG) 5000000000$DENOM --keyring-backend="test" --chain-id $CHAIN_ID

echo "Output of gentx"
cat $HOME/$CHAIN_DATA_DIR/config/gentx/*.json | jq

echo "Running collect-gentxs"
$CHAIN_BIN collect-gentxs

ls $HOME/$CHAIN_DATA_DIR/config

echo "Update app.toml file"
sed -i -e 's#keyring-backend = "os"#keyring-backend = "test"#g' $HOME/$CHAIN_DATA_DIR/config/client.toml
sed -i -e 's#output = "text"#output = "json"#g' $HOME/$CHAIN_DATA_DIR/config/client.toml
sed -i -e 's#broadcast-mode = "sync"#broadcast-mode = "block"#g' $HOME/$CHAIN_DATA_DIR/config/client.toml
sed -i -e "s#chain-id = \"\"#chain-id = \"$CHAIN_ID\"#g" $HOME/$CHAIN_DATA_DIR/config/client.toml

echo "Update config.toml file"
sed -i -e 's#"tcp://127.0.0.1:26657"#"tcp://0.0.0.0:26657"#g' $HOME/$CHAIN_DATA_DIR/config/config.toml
sed -i -e 's/timeout_commit = "5s"/timeout_commit = "1s"/g' $HOME/$CHAIN_DATA_DIR/config/config.toml
sed -i -e 's/timeout_propose = "3s"/timeout_propose = "1s"/g' $HOME/$CHAIN_DATA_DIR/config/config.toml
sed -i -e 's/index_all_keys = false/index_all_keys = true/g' $HOME/$CHAIN_DATA_DIR/config/config.toml

echo "Update genesis.json file with updated local params"
sed -i -e "s/pstake/1PROTECTED1/g; s/restake/2PROTECTED2/g; s/unstake/3PROTECTED3/g; s/stake/$DENOM/g; s/1PROTECTED1/pstake/g; s/2PROTECTED2/restake/g; s/3PROTECTED3/unstake/g" $HOME/$CHAIN_DATA_DIR/config/genesis.json

jq -r '.app_state.staking.params.unbonding_time |= "30s"' $HOME/$CHAIN_DATA_DIR/config/genesis.json > /tmp/genesis.json; mv /tmp/genesis.json $HOME/$CHAIN_DATA_DIR/config/genesis.json
jq -r '.app_state.slashing.params.downtime_jail_duration |= "6s"' $HOME/$CHAIN_DATA_DIR/config/genesis.json > /tmp/genesis.json; mv /tmp/genesis.json $HOME/$CHAIN_DATA_DIR/config/genesis.json
jq -r '.app_state.gov.deposit_params.max_deposit_period |= "30s"' $HOME/$CHAIN_DATA_DIR/config/genesis.json > /tmp/genesis.json; mv /tmp/genesis.json $HOME/$CHAIN_DATA_DIR/config/genesis.json
jq -r '.app_state.gov.deposit_params.min_deposit[0].amount |= "10"' $HOME/$CHAIN_DATA_DIR/config/genesis.json > /tmp/genesis.json; mv /tmp/genesis.json $HOME/$CHAIN_DATA_DIR/config/genesis.json
jq -r '.app_state.gov.voting_params.voting_period |= "30s"' $HOME/$CHAIN_DATA_DIR/config/genesis.json > /tmp/genesis.json; mv /tmp/genesis.json $HOME/$CHAIN_DATA_DIR/config/genesis.json
jq -r '.app_state.gov.tally_params.quorum |= "0.000000000000000000"' $HOME/$CHAIN_DATA_DIR/config/genesis.json > /tmp/genesis.json; mv /tmp/genesis.json $HOME/$CHAIN_DATA_DIR/config/genesis.json
jq -r '.app_state.gov.tally_params.threshold |= "0.000000000000000000"' $HOME/$CHAIN_DATA_DIR/config/genesis.json > /tmp/genesis.json; mv /tmp/genesis.json $HOME/$CHAIN_DATA_DIR/config/genesis.json
jq -r '.app_state.gov.tally_params.veto_threshold |= "0.000000000000000000"' $HOME/$CHAIN_DATA_DIR/config/genesis.json > /tmp/genesis.json; mv /tmp/genesis.json $HOME/$CHAIN_DATA_DIR/config/genesis.json

# Set wasm as permissioned or permissionless based on environment variable
wasm_permission="Nobody"
if [ $WASM_PERMISSIONLESS == "true" ]
then
  wasm_permission="Everybody"
fi

jq -r ".app_state.wasm.params.code_upload_access.permission |= \"${wasm_permission}\"" $HOME/$CHAIN_DATA_DIR/config/genesis.json > /tmp/genesis.json; mv /tmp/genesis.json $HOME/$CHAIN_DATA_DIR/config/genesis.json
jq -r ".app_state.wasm.params.instantiate_default_permission |= \"${wasm_permission}\"" $HOME/$CHAIN_DATA_DIR/config/genesis.json > /tmp/genesis.json; mv /tmp/genesis.json $HOME/$CHAIN_DATA_DIR/config/genesis.json

$CHAIN_BIN tendermint show-node-id
