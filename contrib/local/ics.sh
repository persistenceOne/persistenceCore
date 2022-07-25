#!/bin/bash
set -o errexit -o nounset -o pipefail -eu
CHAIN_BIN=persistenceCore
CHAIN_DIR_1=~/.persistenceCore
CHAIN_ID=test-1
DIR="/tmp/test-contracts"
mkdir -p $DIR

echo "-----------------------"
echo "## Add new CosmWasm contract via gov proposal"
wget "https://github.com/CosmWasm/wasmd/raw/14688c09855ee928a12bcb7cd102a53b78e3cbfb/x/wasm/keeper/testdata/hackatom.wasm" -q -O $DIR/hackatom.wasm
VAL1_KEY=$($CHAIN_BIN keys show -a val1)
RESP=$(persistenceCore tx gov submit-proposal wasm-store "cw20_ics20.wasm" \
  --title "ics " \
  --description "ics20 contract " \
  --deposit 10000stake \
  --run-as $VAL1_KEY \
  --instantiate-everybody "true" \
  --keyring-backend test \
  --from val1 --gas auto --fees 10000stake -y \
  --chain-id test-1 \
  -b block -o json --gas-adjustment 1.5)
echo "$RESP"
PROPOSAL_ID=$(echo "$RESP" | jq -r '.logs[0].events[] | select(.type == "submit_proposal") | .attributes[] | select(.key == "proposal_id") | .value')

echo "### Query proposal prevote"
$CHAIN_BIN q gov proposal $PROPOSAL_ID -o json | jq

echo "### Vote proposal"
$CHAIN_BIN tx gov vote $PROPOSAL_ID yes --from val1 --yes --chain-id $CHAIN_ID \
    --fees 500stake --gas auto --gas-adjustment 1.5 -b block --keyring-backend test -o json | jq
$CHAIN_BIN tx gov vote $PROPOSAL_ID yes --from demo1 --yes --chain-id $CHAIN_ID \
    --fees 500stake --gas auto --gas-adjustment 1.5 -b block --keyring-backend test -o json | jq
$CHAIN_BIN tx gov vote $PROPOSAL_ID yes --from demo2 --yes --chain-id $CHAIN_ID \
    --fees 500stake --gas auto --gas-adjustment 1.5 -b block --keyring-backend test -o json | jq

echo "### Query proposal postvote"
$CHAIN_BIN q gov proposal $PROPOSAL_ID -o json | jq

echo "### Waiting for voting period..."
sleep 40
$CHAIN_BIN q wasm list-code -o json | jq

CODE_ID=$($CHAIN_BIN q wasm list-code -o json | jq -r ".code_infos[-1].code_id")

echo "-----------------------"
echo "## Create new contract instance"
INIT="{\"default_timeout\": 30 ,\"default_gas_limit\": 40000,\"gov_contract\":\"$($CHAIN_BIN keys show val1 -a --keyring-backend test)\", \"allowlist\": []}"
$CHAIN_BIN tx wasm instantiate "$CODE_ID" "$INIT" --admin="$($CHAIN_BIN keys show val1 -a --keyring-backend test)" \
  --from val1 --amount "10000stake" --label "local0.1.0" --gas-adjustment 1.5 --fees "10000stake" \
  --gas "auto" -y --chain-id $CHAIN_ID -b block --keyring-backend test -o json | jq

CONTRACT=$($CHAIN_BIN query wasm list-contract-by-code "$CODE_ID" -o json | jq -r '.contracts[-1]')
echo "* Contract address: $CONTRACT"


