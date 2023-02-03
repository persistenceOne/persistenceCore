#!/bin/bash
set -o errexit -o nounset -o pipefail -eu

DIRNAME="$(dirname $(realpath ${BASH_SOURCE[0]}))"

echo "Replacing local persistenceCore binary with v7 binary"
cp $DIRNAME/../../build/persistenceCoreV7 ~/go/bin/persistenceCore

PERSISTENCE_KEY="persistenceCore keys -a --keyring-backend=test show"
PERSISTENCE_BALANCE="persistenceCore q bank balances"
GAIA_KEY="gaiad keys -a --keyring-backend=test show"
GAIA_BALANCE="gaiad q bank balances"

echo "=> Getting IBC channel info"
GAIA_CHANNEL_INFO=$(gaiad q ibc channel channels -o json | jq 'first(.channels[] | select(.state == "STATE_OPEN") | select(.port_id == "transfer"))')
P_CHANNEL_INFO=$(persistenceCore q ibc channel channels -o json | jq 'first(.channels[] | select(.state == "STATE_OPEN") | select(.port_id == "transfer"))')

echo "=> Channel info:"
echo $GAIA_CHANNEL_INFO | jq
echo $P_CHANNEL_INFO | jq

if [[ -z $GAIA_CHANNEL_INFO ]]; then
    echo "No open transfer port and connection... exiting";
    exit 1;
fi

if [[ -z $P_CHANNEL_INFO ]]; then
    echo "No open transfer port and connection... exiting";
    exit 1;
fi

GAIA_PORT="$(echo $GAIA_CHANNEL_INFO | jq -r '.port_id')"
GAIA_CHANNEL="$(echo $GAIA_CHANNEL_INFO | jq -r '.channel_id')"
PERSISTENCE_PORT="$(echo $P_CHANNEL_INFO | jq -r '.port_id')"
PERSISTENCE_CHANNEL="$(echo $P_CHANNEL_INFO | jq -r '.channel_id')"

check_balance() {
  echo "persistence val1: $($PERSISTENCE_BALANCE $($PERSISTENCE_KEY val1) | jq -r '.balances')"
  echo "persistence val2: $($PERSISTENCE_BALANCE $($PERSISTENCE_KEY val2) | jq -r '.balances')"
  echo "gaia val1: $($GAIA_BALANCE $($GAIA_KEY val1) | jq -r '.balances')"
  echo "gaia val2: $($GAIA_BALANCE $($GAIA_KEY val2) | jq -r '.balances')"
  echo "gaia val3: $($GAIA_BALANCE $($GAIA_KEY val3) | jq -r '.balances')"
}

persistence_transfer() {
  let "tokens=$1 * 1000000"
  denom="$2"
  echo "=> Transfer $tokens $denom from persistence:$3 to gaia:$4"
  persistenceCore tx ibc-transfer transfer "$GAIA_PORT" "$GAIA_CHANNEL" \
    $($GAIA_KEY $4) --fees 5000uxprt "$tokens$denom" \
    --from "$3" --gas auto --gas-adjustment 1.2 -y --keyring-backend test \
    -b block -o json | jq -r '{height, txhash, code, raw_log}'
}

gaia_transfer() {
  let "tokens=$1 * 1000000"
  denom="$2"
  echo "=> Transfer $tokens $denom from gaia:$3 to persistence:$4"
  gaiad tx ibc-transfer transfer "$PERSISTENCE_PORT" "$PERSISTENCE_CHANNEL" \
    $($PERSISTENCE_KEY $4) --fees 5000uatom "$tokens$denom" \
    --from "$3" --gas auto --gas-adjustment 1.2 -y --keyring-backend test \
    -b block -o json | jq -r '{height, txhash, code, raw_log}'
}

echo "=> Check balance after upgrade"
check_balance

echo "=> IBC Transfer token from source chain"
persistence_transfer 10 uxprt val1 val2
gaia_transfer 10 uatom val2 val1

echo "=> Waiting for a bit to let ibc-transfer happen"
sleep 4

echo "=> Balances after transfer from source chain"
check_balance

IBC_DENOM_PERSISTENCE="ibc/$(persistenceCore q ibc-transfer denom-hash $PERSISTENCE_PORT/$PERSISTENCE_CHANNEL/uatom | jq -r '.hash')"
IBC_DENOM_GAIA="ibc/$(gaiad q ibc-transfer denom-hash $GAIA_PORT/$GAIA_CHANNEL/uxprt | jq -r '.hash')"

echo "=> IBC Transfer ibc token back to source chain"
gaia_transfer 5 $IBC_DENOM_GAIA val2 val1
persistence_transfer 11 $IBC_DENOM_PERSISTENCE val1 val2
persistence_transfer 165 $IBC_DENOM_PERSISTENCE val2 val3

echo "=> Waiting for a bit to let ibc-transfer happen"
sleep 4

echo "=> Balances after transfer ibc tokens back to source chain"
check_balance

echo "=> Execute existing wasm contract"
bash -e $DIRNAME/execute-existing-contract.sh

echo "=> Testing wasm contract upload/interact/migrate"
UPLOAD_AGAIN=false bash -e $DIRNAME/contract.sh
