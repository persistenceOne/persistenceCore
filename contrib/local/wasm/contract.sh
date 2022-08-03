#!/bin/bash
set -o errexit -o nounset -o pipefail -eu

DIR="$HOME/test-contracts"
mkdir -p $DIR

echo "-----------------------"
echo "## Add new CosmWasm contract"
wget "https://github.com/CosmWasm/wasmd/raw/14688c09855ee928a12bcb7cd102a53b78e3cbfb/x/wasm/keeper/testdata/hackatom.wasm" -q -O $DIR/hackatom.wasm
RESP=$($CHAIN_BIN tx wasm store "$DIR/hackatom.wasm" --keyring-backend test \
  --from val1 --gas auto --fees 10000uxprt -y --chain-id $CHAIN_ID -b block -o json --gas-adjustment 1.5)
echo "$RESP"
CODE_ID=$(echo "$RESP" | jq -r '.logs[0].events[1].attributes[-1].value')
echo "* Code id: $CODE_ID"
echo "* Download code"

TMPDIR=$(mktemp -t wasmdXXXXXX)
$CHAIN_BIN q wasm code "$CODE_ID" "$TMPDIR"
rm -f "$TMPDIR"
echo "-----------------------"
echo "## List code"
$CHAIN_BIN query wasm list-code --chain-id $CHAIN_ID -o json | jq

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

echo "-----------------------"
echo "## Migrate contract"
echo "### Upload new code"
wget "https://github.com/CosmWasm/wasmd/raw/14688c09855ee928a12bcb7cd102a53b78e3cbfb/x/wasm/keeper/testdata/burner.wasm" -q -O $DIR/burner.wasm
RESP=$($CHAIN_BIN tx wasm store "$DIR/burner.wasm" --gas-adjustment 1.5 --fees "10000uxprt" \
  --from val1 --gas "auto" -y --chain-id $CHAIN_ID -b block -o json)

BURNER_CODE_ID=$(echo "$RESP" | jq -r '.logs[0].events[1].attributes[-1].value')
echo "### Migrate to code id: $BURNER_CODE_ID"

DEST_ACCOUNT=$($CHAIN_BIN keys show test1 -a)
$CHAIN_BIN tx wasm migrate "$CONTRACT" "$BURNER_CODE_ID" "{\"payout\": \"$DEST_ACCOUNT\"}" --from test1 \
  --chain-id $CHAIN_ID --gas "auto" --gas-adjustment 1.5 --fees "10000uxprt" -b block -y -o json | jq

echo "### Query destination account: $BURNER_CODE_ID"
$CHAIN_BIN q bank balances "$DEST_ACCOUNT" -o json | jq
echo "### Query contract meta data: $CONTRACT"
$CHAIN_BIN q wasm contract "$CONTRACT" -o json | jq

echo "### Query contract meta history: $CONTRACT"
$CHAIN_BIN q wasm contract-history "$CONTRACT" -o json | jq

echo "-----------------------"
echo "## Clear contract admin"
echo "### Query old admin: $($CHAIN_BIN q wasm contract "$CONTRACT" -o json | jq -r '.contract_info.admin')"
echo "### Update contract"
$CHAIN_BIN tx wasm clear-contract-admin "$CONTRACT" \
  --from test1 -y --chain-id $CHAIN_ID -b block -o json \
  --gas "auto" --gas-adjustment 1.5 --fees "10000uxprt" | jq
echo "### Query new admin: $($CHAIN_BIN q wasm contract "$CONTRACT" -o json | jq -r '.contract_info.admin')"
