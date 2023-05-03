#!/bin/bash
set -o errexit -o nounset -o pipefail -eu

CHAIN_BIN="persistenceCore"
CHAIN_ID="test-core-1"
SHOW_KEY="persistenceCore keys show --keyring-backend test -a"

CODE_ID=$($CHAIN_BIN q wasm list-code --chain-id $CHAIN_ID -o json | jq -r '.code_infos[-1].code_id')

if [[ $CODE_ID == "null" ]]; then
    echo "No contract found";
    exit 1;
fi

CONTRACT_ADDR=$($CHAIN_BIN q wasm list-contract-by-code "$CODE_ID" -o json | jq -r '.contracts[-1]')

echo "--------------------------------------------"
echo "=> Query contract-state"
echo "-> query all"
STATE_ALL=$($CHAIN_BIN q wasm contract-state all "$CONTRACT_ADDR" -o json)
echo "$STATE_ALL" | jq

echo "-> query smart"
$CHAIN_BIN q wasm contract-state smart "$CONTRACT_ADDR" '{"verifier":{}}' -o json | jq

echo "-> query raw"
KEY=$(echo "$STATE_ALL" | jq -r ".models[0].key")
$CHAIN_BIN q wasm contract-state raw "$CONTRACT_ADDR" "$KEY" -o json | jq

echo "--------------------------------------------"
echo "=> Execute wasm contract: $CONTRACT_ADDR"
MSG='{"release":{}}'
$CHAIN_BIN tx wasm execute "$CONTRACT_ADDR" "$MSG" \
  --from val1 --keyring-backend test --gas-adjustment 1.5 \
  --fees "10000uxprt" --gas "auto" -y --chain-id $CHAIN_ID \
  -b block -o json | jq -r '{height, txhash, code, raw_log}'

echo "--------------------------------------------"
echo "=> Create one more instance"
INIT="{\"verifier\":\"$($SHOW_KEY val1)\", \"beneficiary\":\"$($SHOW_KEY test1)\"}"
$CHAIN_BIN tx wasm instantiate "$CODE_ID" "$INIT" --admin="$($SHOW_KEY val1)" \
  --from val1 --amount "10000uxprt" --label "local0.1.0" --gas-adjustment 1.5 --fees "10000uxprt" \
  --gas "auto" -y --chain-id $CHAIN_ID -b block \
  -o json | jq -r '{height, txhash, code, raw_log}'

CONTRACT_ADDR=$($CHAIN_BIN query wasm list-contract-by-code "$CODE_ID" -o json | jq -r '.contracts[-1]')
echo "-> New instantiated contract address: $CONTRACT_ADDR"

echo "--------------------------------------------"
echo "=> Execute wasm contract: $CONTRACT_ADDR"
MSG='{"release":{}}'
$CHAIN_BIN tx wasm execute "$CONTRACT_ADDR" "$MSG" \
  --from val1 --keyring-backend test --gas-adjustment 1.5 \
  --fees "10000uxprt" --gas "auto" -y --chain-id $CHAIN_ID \
  -b block -o json | jq -r '{height, txhash, code, raw_log}'

echo "-------------------DONE---------------------"