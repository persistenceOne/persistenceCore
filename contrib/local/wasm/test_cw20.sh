#!/bin/bash

set -o errexit -o nounset -o pipefail -eu

DIR="/tmp/test-contracts"
mkdir -p $DIR

VAL1_KEY=$($CHAIN_BIN keys show -a val1 --keyring-backend test)
TEST1_KEY=$($CHAIN_BIN keys show -a test1 --keyring-backend test)
TEST2_KEY=$($CHAIN_BIN keys show -a test2 --keyring-backend test)

echo "-----------------------"
echo "## Add cw CosmWasm contract via gov proposal"
wget "https://github.com/CosmWasm/cw-plus/releases/download/v0.13.4/cw20_base.wasm" -q -O $DIR/cw20_base.wasm

RESP=$($CHAIN_BIN tx gov submit-proposal wasm-store "$DIR/cw20_base.wasm" \
  --title "Add cw20_base" \
  --description "cw20_base contact" \
  --deposit 10000uxprt \
  --run-as $VAL1_KEY \
  --instantiate-everybody "true" \
  --keyring-backend test \
  --from val1 --gas auto --fees 10000uxprt -y \
  --chain-id $CHAIN_ID \
  -b block -o json --gas-adjustment 1.5)
echo "$RESP"
PROPOSAL_ID=$(echo "$RESP" | jq -r '.logs[0].events[] | select(.type == "submit_proposal") | .attributes[] | select(.key == "proposal_id") | .value')

echo "### Query proposal prevote"
$CHAIN_BIN q gov proposal $PROPOSAL_ID -o json | jq

echo "### Vote proposal"
$CHAIN_BIN tx gov vote $PROPOSAL_ID yes --from val1 --yes --chain-id $CHAIN_ID \
    --fees 500uxprt --gas auto --gas-adjustment 1.5 -b block -o json | jq
$CHAIN_BIN tx gov vote $PROPOSAL_ID yes --from test1 --yes --chain-id $CHAIN_ID \
    --fees 500uxprt --gas auto --gas-adjustment 1.5 -b block -o json | jq
$CHAIN_BIN tx gov vote $PROPOSAL_ID yes --from test2 --yes --chain-id $CHAIN_ID \
    --fees 500uxprt --gas auto --gas-adjustment 1.5 -b block -o json | jq

echo "### Query proposal postvote"
$CHAIN_BIN q gov proposal $PROPOSAL_ID -o json | jq

echo "### Waiting for voting period..."
sleep 40
$CHAIN_BIN q wasm list-code -o json | jq

CODE_ID=$($CHAIN_BIN q wasm list-code -o json | jq -r ".code_infos[-1].code_id")

echo "-----------------------"
echo "## Create new contract instance"
INIT=$(cat <<EOF
{
  "name": "My first token",
  "symbol": "FRST",
  "decimals": 6,
  "initial_balances": [{
    "address": "$TEST1_KEY",
    "amount": "123456789000"
  }]
}
EOF
)
$CHAIN_BIN tx wasm instantiate "$CODE_ID" "$INIT" --admin="$TEST1_KEY" \
  --from $TEST1_KEY --amount "10000uxprt" --label "First Coin" --gas-adjustment 1.5 --fees "10000uxprt" \
  --gas "auto" -y --chain-id $CHAIN_ID -b block -o json | jq

CONTRACT=$($CHAIN_BIN query wasm list-contract-by-code "$CODE_ID" -o json | jq -r '.contracts[-1]')
echo "* Contract address: $CONTRACT"

echo "### Query test balance: expected balance: "
$CHAIN_BIN query wasm contract-state smart $CONTRACT "{\"balance\":{\"address\":\"$TEST1_KEY\"}}"
echo "### Query non balance: expected balance: "
$CHAIN_BIN query wasm contract-state smart $CONTRACT "{\"balance\":{\"address\":\"$TEST2_KEY\"}}"

echo "-----------------------"
echo "## Execute contract $CONTRACT"
TRANSFER=$(cat <<EOF
{
  "transfer": {
    "recipient": "$TEST2_KEY",
    "amount": "987654321"
  }
}
EOF
)

echo $TRANSFER | jq
$CHAIN_BIN tx wasm execute $CONTRACT "$TRANSFER" \
  --from $TEST1_KEY --gas-adjustment 1.5 --fees "10000uxprt" \
  --gas "auto" -y --chain-id $CHAIN_ID -b block -o json | jq

echo "### Query balance after transfer: expected balance: "
$CHAIN_BIN query wasm contract-state smart $CONTRACT "{\"balance\":{\"address\":\"$TEST1_KEY\"}}"
echo "### Query non balance after transfer: expected balance: "
$CHAIN_BIN query wasm contract-state smart $CONTRACT "{\"balance\":{\"address\":\"$TEST2_KEY\"}}"
