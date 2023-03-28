#!/bin/bash
set -o nounset -o pipefail -eu

DIR="/tmp/test-oracle"
mkdir -p $DIR
cp ./oracle/proposal.json $DIR/proposal.json

VAL1_KEY=$($CHAIN_BIN keys show -a val1 --keyring-backend test)
TEST1_KEY=$($CHAIN_BIN keys show -a test1 --keyring-backend test)
TEST2_KEY=$($CHAIN_BIN keys show -a test2 --keyring-backend test)

echo "-----------------------"
echo "## Change oracle params(AcceptList) via gov proposal"

echo "### Sleep so that few blocks are mined before submitting proposal"
sleep 5

## Submit tx to change oracle params via gov proposal using proposal.sh
RESP=$($CHAIN_BIN tx gov submit-proposal param-change  "$DIR/proposal.json" \
  --keyring-backend test \
  --from $VAL1_KEY --gas auto --fees 10000uxprt -y \
  --chain-id $CHAIN_ID \
  -b block -o json --gas-adjustment 1.5)

echo "$RESP" | jq
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
$CHAIN_BIN q oracle params -o json | jq

echo "### Query params"
$CHAIN_BIN q oracle params -o json | jq