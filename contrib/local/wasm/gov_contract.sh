#!/bin/bash
set -o errexit -o nounset -o pipefail -eu

DIR="/tmp/test-contracts"
mkdir -p $DIR

echo "-----------------------"
echo "## Add new CosmWasm contract via gov proposal"
wget "https://github.com/CosmWasm/wasmd/raw/14688c09855ee928a12bcb7cd102a53b78e3cbfb/x/wasm/keeper/testdata/hackatom.wasm" -q -O $DIR/hackatom.wasm
VAL1_KEY=$($CHAIN_BIN keys show -a val1)
RESP=$($CHAIN_BIN tx gov submit-proposal wasm-store "$DIR/hackatom.wasm" \
  --title "hackatom" \
  --description "hackatom test contact" \
  --deposit 10000uxprt \
  --run-as $VAL1_KEY \
  --instantiate-everybody "true" \
  --keyring-backend test \
  --from $VAL1_KEY --gas auto --fees 10000uxprt -y \
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
INIT="{\"verifier\":\"$($CHAIN_BIN keys show val1 -a --keyring-backend test)\", \"beneficiary\":\"$($CHAIN_BIN keys show test1 -a)\"}"
$CHAIN_BIN tx wasm instantiate "$CODE_ID" "$INIT" --admin="$($CHAIN_BIN keys show val1 -a --keyring-backend test)" \
  --from val1 --amount "10000uxprt" --label "local0.1.0" --gas-adjustment 1.5 --fees "10000uxprt" \
  --gas "auto" -y --chain-id $CHAIN_ID -b block -o json | jq

CONTRACT=$($CHAIN_BIN query wasm list-contract-by-code "$CODE_ID" -o json | jq -r '.contracts[-1]')
echo "* Contract address: $CONTRACT"

echo "### Query all"
RESP=$($CHAIN_BIN query wasm contract-state all "$CONTRACT" -o json)
echo "$RESP" | jq
echo "### Query smart"
$CHAIN_BIN query wasm contract-state smart "$CONTRACT" '{"verifier":{}}' -o json | jq
echo "### Query raw"
KEY=$(echo "$RESP" | jq -r ".models[0].key")
$CHAIN_BIN query wasm contract-state raw "$CONTRACT" "$KEY" -o json | jq

echo "-----------------------"
echo "## Execute contract $CONTRACT"
MSG='{"release":{}}'
$CHAIN_BIN tx wasm execute "$CONTRACT" "$MSG" \
  --from val1 --gas-adjustment 1.5 --fees "10000uxprt" \
  --gas "auto" -y --chain-id $CHAIN_ID -b block -o json | jq

echo "-----------------------"
echo "## Set new admin"
echo "### Query old admin: $($CHAIN_BIN q wasm contract "$CONTRACT" -o json | jq -r '.contract_info.admin')"
echo "### Update contract"
$CHAIN_BIN tx wasm set-contract-admin "$CONTRACT" "$($CHAIN_BIN keys show test1 -a)" \
  --from val1 --gas-adjustment 1.5 --gas "auto" --fees "10000uxprt" -y --chain-id $CHAIN_ID -b block -o json | jq
echo "### Query new admin: $($CHAIN_BIN q wasm contract "$CONTRACT" -o json | jq -r '.contract_info.admin')"
