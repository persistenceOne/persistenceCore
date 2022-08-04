#!/bin/bash

RELAYER_CONFIG="${RELAYER_CONFIG:=configs/relayer.toml}"
KEYS_CONFIG="${KEYS_CONFIG:=configs/relayer-keys.json}"

# Add hermes configuration properly
mkdir -p $HOME/.hermes && touch $HOME/.hermes/config.toml
cp $RELAYER_CONFIG $HOME/.hermes/config.toml

# Add chain keys
for ((i=0; i<$(jq -r ".chains | length" $KEYS_CONFIG); i++))
do
  jq -r ".chains[$i].keys[0]" $KEYS_CONFIG > /tmp/key$i.json
  hermes keys add \
    --chain $(jq -r ".chains[$i].id" $KEYS_CONFIG) \
    --key-file /tmp/key$i.json \
    --hd-path $(jq -r ".chains[$i].hdpath" $KEYS_CONFIG)
done
