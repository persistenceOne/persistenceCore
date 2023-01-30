#!/bin/bash

OFFSET_HEIGHT=40
UPGRADE_NAME=v7

set -o errexit -o nounset -o pipefail -eu

$CHAIN_BIN status 2>&1 | jq

CURRENT_HEIGHT=$($CHAIN_BIN status 2>&1 | jq -r ".SyncInfo.latest_block_height")
UPGRADE_HEIGHT=`expr $CURRENT_HEIGHT + $OFFSET_HEIGHT`
echo "Starting software upgrade"

echo "### Submit proposal from val1"
RESP=$($CHAIN_BIN tx gov submit-proposal software-upgrade $UPGRADE_NAME --yes --title "$UPGRADE_NAME" --description "$UPGRADE_NAME" \
    --upgrade-height $UPGRADE_HEIGHT --from val1 --chain-id $CHAIN_ID --deposit 100uxprt \
    --fees 2000uxprt --gas auto --gas-adjustment 1.5 -b block -o json)
PROPOSAL_ID=$(echo "$RESP" | jq -r '.logs[0].events[] | select(.type == "submit_proposal") | .attributes[] | select(.key == "proposal_id") | .value')

echo "### Query proposal prevote"
$CHAIN_BIN q gov proposal $PROPOSAL_ID -o json | jq

echo "### Vote proposal"
$CHAIN_BIN tx gov vote $PROPOSAL_ID yes --from val1 --yes --chain-id $CHAIN_ID \
    --fees 200uxprt --gas auto --gas-adjustment 1.5 -b block -o json | jq
$CHAIN_BIN tx gov vote $PROPOSAL_ID yes --from test1 --yes --chain-id $CHAIN_ID \
    --fees 200uxprt --gas auto --gas-adjustment 1.5 -b block -o json | jq
$CHAIN_BIN tx gov vote $PROPOSAL_ID yes --from test2 --yes --chain-id $CHAIN_ID \
    --fees 200uxprt --gas auto --gas-adjustment 1.5 -b block -o json | jq

echo "###Proposal voting period"
sleep 40
echo "### Query proposal postvote"
$CHAIN_BIN q gov proposal $PROPOSAL_ID -o json | jq

echo "Upgrade Height: $UPGRADE_HEIGHT"