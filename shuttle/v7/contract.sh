#!/bin/bash
set -o errexit -o nounset -o pipefail -eu

CHAIN_BIN="persistenceCore"
CHAIN_ID="core-1"
SHOW_KEY="persistenceCore keys show --keyring-backend test -a"

DIRNAME="$(dirname $(realpath ${BASH_SOURCE[0]}))"
CONTRACT_DIR="$DIRNAME/test-contracts"
CONTRACT_FILE1="$CONTRACT_DIR/hackatom.wasm"
CONTRACT_FILE2="$CONTRACT_DIR/burner.wasm"

UPLOAD_AGAIN="${UPLOAD_AGAIN:=true}"

echo "------------TESTING-WASM-CONTRACT-----------"

mkdir -p $CONTRACT_DIR

if [ ! -f "$CONTRACT_FILE1" ]; then
    echo "=> Downloading CosmWasm contract: hackatom.wasm (contract1)"
    wget "https://github.com/CosmWasm/wasmd/raw/14688c09855ee928a12bcb7cd102a53b78e3cbfb/x/wasm/keeper/testdata/hackatom.wasm" -q -O $CONTRACT_FILE1
else
    echo "=> Already downloaded: hackatom.wasm (contract1)"
fi

if [ ! -f "$CONTRACT_FILE2" ]; then
    echo "=> Downloading CosmWasm contract: burner.wasm (contract2)"
    wget "https://github.com/CosmWasm/wasmd/raw/14688c09855ee928a12bcb7cd102a53b78e3cbfb/x/wasm/keeper/testdata/burner.wasm" -q -O $CONTRACT_FILE2
else
    echo "=> Already downloaded: burner.wasm (contract2)"
fi

echo "--------------------------------------------"
echo "=> Uploading contract1"
RESP=$($CHAIN_BIN tx wasm store "$CONTRACT_FILE1" --keyring-backend test \
  --from val1 --gas auto --fees 60000uxprt -y --chain-id $CHAIN_ID -b block -o json --gas-adjustment 1.5)
echo "$RESP" | jq  -r '{height, txhash, code, raw_log}'

CODE_ID1=$(echo "$RESP" | jq -r '.logs[0].events[] | select(.type == "store_code") | .attributes[] | select(.key == "code_id") | .value')
echo "=> Code ID: $CODE_ID1"

echo "=> Download wasm contract1"
TMPDIR=$(mktemp -t wasmdXXXXXX)
$CHAIN_BIN q wasm code $CODE_ID1 $TMPDIR
rm -rf $TMPDIR

echo "--------------------------------------------"
echo "=> Instantiate wasm contract1"
INIT="{\"verifier\":\"$($SHOW_KEY val1)\", \"beneficiary\":\"$($SHOW_KEY test1)\"}"
$CHAIN_BIN tx wasm instantiate "$CODE_ID1" "$INIT" --admin="$($SHOW_KEY val1)" \
  --from val1 --amount "10000uxprt" --label "local0.1.0" --gas-adjustment 1.5 --fees "10000uxprt" \
  --gas "auto" -y --chain-id $CHAIN_ID -b block \
  -o json | jq -r '{height, txhash, code, raw_log}'

CONTRACT_ADDR=$($CHAIN_BIN query wasm list-contract-by-code "$CODE_ID1" -o json | jq -r '.contracts[-1]')
echo "-> Contract address: $CONTRACT_ADDR"

echo "--------------------------------------------"
echo "=> Query contract-state"
echo "-> query all"
STATE_ALL=$($CHAIN_BIN query wasm contract-state all "$CONTRACT_ADDR" -o json)
echo "$STATE_ALL" | jq

echo "-> query smart"
$CHAIN_BIN query wasm contract-state smart "$CONTRACT_ADDR" '{"verifier":{}}' -o json | jq

echo "-> query raw"
KEY=$(echo "$STATE_ALL" | jq -r ".models[0].key")
$CHAIN_BIN query wasm contract-state raw "$CONTRACT_ADDR" "$KEY" -o json | jq

echo "--------------------------------------------"
echo "=> Execute wasm contract: $CONTRACT_ADDR"
MSG='{"release":{}}'
$CHAIN_BIN tx wasm execute "$CONTRACT_ADDR" "$MSG" \
  --from val1 --keyring-backend test --gas-adjustment 1.5 \
  --fees "10000uxprt" --gas "auto" -y --chain-id $CHAIN_ID \
  -b block -o json | jq -r '{height, txhash, code, raw_log}'

echo "--------------------------------------------"
echo "=> Set new admin"
echo "-> admin before set: $($CHAIN_BIN q wasm contract "$CONTRACT_ADDR" -o json | jq -r '.contract_info.admin')"
echo "-> Run tx wasm set-contract-admin"
$CHAIN_BIN tx wasm set-contract-admin "$CONTRACT_ADDR" "$($SHOW_KEY test1)" \
  --from val1 --gas-adjustment 1.5 --gas "auto" --fees "10000uxprt" -y --chain-id $CHAIN_ID \
  -b block -o json | jq -r '{height, txhash, code, raw_log}'
echo "-> admin after set: $($CHAIN_BIN q wasm contract "$CONTRACT_ADDR" -o json | jq -r '.contract_info.admin')"

echo "--------------------------------------------"
echo "=> Migrate wasm contract"
echo "-> Uploading contract2"
RESP=$($CHAIN_BIN tx wasm store "$CONTRACT_FILE2" --keyring-backend test \
  --from val1 --gas auto --fees 60000uxprt -y --chain-id $CHAIN_ID -b block -o json --gas-adjustment 1.5)
echo $RESP | jq -r '{height, txhash, code, raw_log}'
CODE_ID2=$(echo "$RESP" | jq -r '.logs[0].events[] | select(.type == "store_code") | .attributes[] | select(.key == "code_id") | .value')

echo "-> Migrating to code id: $CODE_ID2"
DEST_ACCOUNT=$($SHOW_KEY test1)
$CHAIN_BIN tx wasm migrate "$CONTRACT_ADDR" "$CODE_ID2" "{\"payout\": \"$DEST_ACCOUNT\"}" --from test1 \
  --keyring-backend test --chain-id $CHAIN_ID --gas "auto" --gas-adjustment 1.5 --fees "10000uxprt" -b block \
  -y -o json | jq -r '{height, txhash, code, raw_log}'

echo "-> Query destination account balance"
$CHAIN_BIN q bank balances "$DEST_ACCOUNT" -o json | jq
echo "-> Query contract meta data: $CONTRACT_ADDR"
$CHAIN_BIN q wasm contract "$CONTRACT_ADDR" -o json | jq
echo "-> Query contract meta history: $CONTRACT_ADDR"
$CHAIN_BIN q wasm contract-history "$CONTRACT_ADDR" -o json | jq

echo "--------------------------------------------"
echo "=> Clear contract admin"
echo "-> admin before clear: $($CHAIN_BIN q wasm contract "$CONTRACT_ADDR" -o json | jq -r '.contract_info.admin')"
echo "-> Run tx wasm clear-contract-admin"
$CHAIN_BIN tx wasm clear-contract-admin "$CONTRACT_ADDR" \
  --from test1 -y --chain-id $CHAIN_ID -b block -o json --keyring-backend test \
  --gas "auto" --gas-adjustment 1.5 --fees "10000uxprt" | jq -r '{height, txhash, code, raw_log}'
echo "-> admin after clear: $($CHAIN_BIN q wasm contract "$CONTRACT_ADDR" -o json | jq -r '.contract_info.admin')"

if [[ "$UPLOAD_AGAIN" == "true" ]]; then
    echo "--------------------------------------------"
    echo "=> Uploading contract1 again (to be executed again in post-upgrade script)"
    RESP=$($CHAIN_BIN tx wasm store "$CONTRACT_FILE1" --keyring-backend test \
    --from val1 --gas auto --fees 60000uxprt -y --chain-id $CHAIN_ID -b block -o json --gas-adjustment 1.5)
    echo "$RESP" | jq  -r '{height, txhash, code, raw_log}'

    CODE_ID3=$(echo "$RESP" | jq -r '.logs[0].events[] | select(.type == "store_code") | .attributes[] | select(.key == "code_id") | .value')
    echo "=> Code ID: $CODE_ID3"

    echo "--------------------------------------------"
    echo "=> Instantiate wasm contract1"
    INIT="{\"verifier\":\"$($SHOW_KEY val1)\", \"beneficiary\":\"$($SHOW_KEY test1)\"}"
    $CHAIN_BIN tx wasm instantiate "$CODE_ID3" "$INIT" --admin="$($SHOW_KEY val1)" \
    --from val1 --amount "10000uxprt" --label "local0.1.0" --gas-adjustment 1.5 --fees "10000uxprt" \
    --gas "auto" -y --chain-id $CHAIN_ID -b block \
    -o json | jq -r '{height, txhash, code, raw_log}'

    CONTRACT_ADDR=$($CHAIN_BIN query wasm list-contract-by-code "$CODE_ID3" -o json | jq -r '.contracts[-1]')
    echo "-> Contract address: $CONTRACT_ADDR"

    echo "--------------------------------------------"
    echo "=> Execute wasm contract: $CONTRACT_ADDR"
    MSG='{"release":{}}'
    $CHAIN_BIN tx wasm execute "$CONTRACT_ADDR" "$MSG" \
    --from val1 --keyring-backend test --gas-adjustment 1.5 \
    --fees "10000uxprt" --gas "auto" -y --chain-id $CHAIN_ID \
    -b block -o json | jq -r '{height, txhash, code, raw_log}'
fi

echo "-------------------DONE---------------------"