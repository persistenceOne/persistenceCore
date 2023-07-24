#!/bin/bash

CHAIN_BIN="${CHAIN_BIN:=persistenceCore}"
CHAIN_DATA_DIR="${CHAIN_DATA_DIR:=.persistenceCore}"
CHAIN_ID="${CHAIN_ID:=test-core-2}"
NODE_HOST="${NODE_HOST:=localhost}"
NODE_PORT="${NODE_PORT:=26657}"
EXPOSER_PORT="${EXPOSER_PORT:=8081}"

set -o errexit -o nounset -o pipefail -eu

echo "Download mnemonic from exposer"
MNEMONIC_CONFIG="/tmp/mnemonic.json"
curl -o $MNEMONIC_CONFIG http://$NODE_HOST:$EXPOSER_PORT/keys

echo "Starting to add keys to the keyring"
jq -r ".genesis[0].mnemonic" $MNEMONIC_CONFIG | $CHAIN_BIN keys add $(jq -r ".genesis[0].name" $MNEMONIC_CONFIG) --recover --keyring-backend="test"

# Add keys to keyring
for ((i=0; i<$(jq -r '.validators | length' $MNEMONIC_CONFIG); i++))
do
  jq -r ".validators[$i].mnemonic" $MNEMONIC_CONFIG | $CHAIN_BIN keys add $(jq -r ".validators[$i].name" $MNEMONIC_CONFIG) --recover --keyring-backend="test"
done

for ((i=0; i<$(jq -r '.keys | length' $MNEMONIC_CONFIG); i++))
do
  jq -r ".keys[$i].mnemonic" $MNEMONIC_CONFIG | $CHAIN_BIN keys add $(jq -r ".keys[$i].name" $MNEMONIC_CONFIG) --recover --keyring-backend="test"
done

echo "Update client.toml file"
sed -i -e 's#keyring-backend = ".*"#keyring-backend = "test"#g' $HOME/$CHAIN_DATA_DIR/config/client.toml
sed -i -e 's#output = ".*"#output = "json"#g' $HOME/$CHAIN_DATA_DIR/config/client.toml
sed -i -e 's#broadcast-mode = ".*"#broadcast-mode = "block"#g' $HOME/$CHAIN_DATA_DIR/config/client.toml
sed -i -e "s#chain-id = \".*\"#chain-id = \"$CHAIN_ID\"#g" $HOME/$CHAIN_DATA_DIR/config/client.toml
sed -i -e "s#node = \".*\"#node = \"tcp://$NODE_HOST:$NODE_PORT\"#g" $HOME/$CHAIN_DATA_DIR/config/client.toml

$CHAIN_BIN status 2>&1
