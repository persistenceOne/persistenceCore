#!/bin/bash
set -o nounset -o errtrace -o pipefail

DIR="$HOME/test-contracts"
mkdir -p $DIR

echo "-----------------------"
echo "## Add dummy Oracle contract"

# copy oracle.wasm in ./tmp/trash/
wget "https://raw.github.com/Tikaryan/persistence_contract/tikaryan/update-stargate-queries/artifacts/oracle.wasm" -q -O $DIR/oracle.wasm

RESP=$($CHAIN_BIN tx wasm store "$DIR/oracle.wasm" --keyring-backend test \
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
INIT="{\"owner\":\"$($CHAIN_BIN keys show val1 -a --keyring-backend test)\", \"symbol\":\"$ASSET\"}"
$CHAIN_BIN tx wasm instantiate "$CODE_ID" "$INIT" --admin="$($CHAIN_BIN keys show val1 -a --keyring-backend test)" \
  --from val1 --amount "10000uxprt" --label "local0.1.0" --gas-adjustment 1.5 --fees "10000uxprt" \
  --gas "auto" -y --chain-id $CHAIN_ID -b block -o json | jq

CONTRACT=$($CHAIN_BIN query wasm list-contract-by-code "$CODE_ID" -o json | jq -r '.contracts[-1]')
echo "* Contract address: $CONTRACT"

echo "### Query all"
RESP=$($CHAIN_BIN query wasm contract-state all "$CONTRACT" -o json)
echo "$RESP" | jq

# Read the response and wait for the exchange rate to be updated by the oracle feeder script.
# while retry count is not zero
RETRYCOUNT=20
while [ $RETRYCOUNT -gt 0 ]; do
  RESP=$($CHAIN_BIN query wasm contract-state smart "$CONTRACT" '{"get_exchange_rate":{"symbol":"DUMMY"}}' 2>&1)

  if [[ $RESP == *"invalid request"* ]]; then
    echo "Invalid request. Feeder has not updated the exchange rate yet. Waiting for 5 seconds" && sleep 5

    RETRYCOUNT=$((RETRYCOUNT-1))

    # if retry count is zero, exit with error.
    if [ $RETRYCOUNT -eq 0 ]; then
      echo "Exchange rate not updated by feeder. Exiting."
      exit 1
    fi
  elif [[ $RESP == *"exchange_rate"* ]]; then
    EXCHANGE_RATE=$(echo $RESP | jq -r '.data.exchange_rate')
    echo "Exchange rate of asset $ASSET in USD: $EXCHANGE_RATE"
    break
  fi
done
