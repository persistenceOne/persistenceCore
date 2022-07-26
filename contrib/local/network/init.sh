#!/bin/bash

BINARY=persistenceCore
BINARY_2=gaiad
CHAIN_DIR=./data
CHAINID_1=test-1
CHAINID_2=test-2
VAL_MNEMONIC_1="together chief must vocal account off apart dinosaur move canvas spring whisper improve cruise idea earn reflect flash goat illegal mistake blood earn ridge"
VAL_MNEMONIC_2="angry twist harsh drastic left brass behave host shove marriage fall update business leg direct reward object ugly security warm tuna model broccoli choice"
DEMO_MNEMONIC_1="marble allow december print trial know resource cry next segment twice nose because steel omit confirm hair extend shrimp seminar one minor phone deputy"
DEMO_MNEMONIC_2="veteran try aware erosion drink dance decade comic dawn museum release episode original list ability owner size tuition surface ceiling depth seminar capable only"
RLY_MNEMONIC_1="axis decline final suggest denial erupt satisfy weekend utility fortune dry glory recall real other evil spatial speed seek rubber struggle wolf tortoise large"
RLY_MNEMONIC_2="record gift you once hip style during joke field prize dust unique length more pencil transfer quit train device arrive energy sort steak upset"
P2PPORT_1=16656
P2PPORT_2=26656
RPCPORT_1=16657
RPCPORT_2=26657
RESTPORT_1=1316
RESTPORT_2=1317
ROSETTA_1=8080
ROSETTA_2=8081

# Stop if it is already running
if pgrep -x "$BINARY" >/dev/null; then
    echo "Terminating $BINARY..."
    killall $BINARY
fi

echo "Removing previous data..."
rm -rf $CHAIN_DIR/$CHAINID_1 &> /dev/null
rm -rf $CHAIN_DIR/$CHAINID_2 &> /dev/null

# Add directories for both chains, exit if an error occurs
if ! mkdir -p $CHAIN_DIR/$CHAINID_1 2>/dev/null; then
    echo "Failed to create chain folder. Aborting..."
    exit 1
fi

if ! mkdir -p $CHAIN_DIR/$CHAINID_2 2>/dev/null; then
    echo "Failed to create chain folder. Aborting..."
    exit 1
fi

echo "Initializing $CHAINID_1..."
echo "Initializing $CHAINID_2..."
$BINARY init test --home $CHAIN_DIR/$CHAINID_1 --chain-id=$CHAINID_1
$BINARY_2 init test --home $CHAIN_DIR/$CHAINID_2 --chain-id=$CHAINID_2

echo "Adding genesis accounts..."
echo $VAL_MNEMONIC_1 | $BINARY keys add val1 --home $CHAIN_DIR/$CHAINID_1 --recover --keyring-backend test
echo $VAL_MNEMONIC_2 | $BINARY_2 keys add val2 --home $CHAIN_DIR/$CHAINID_2 --recover --keyring-backend test
echo $DEMO_MNEMONIC_1 | $BINARY keys add demowallet1 --home $CHAIN_DIR/$CHAINID_1 --recover --keyring-backend test
echo $DEMO_MNEMONIC_2 | $BINARY_2 keys add demowallet2 --home $CHAIN_DIR/$CHAINID_2 --recover --keyring-backend test
echo $RLY_MNEMONIC_1 | $BINARY keys add rly1 --home $CHAIN_DIR/$CHAINID_1 --recover --keyring-backend test
echo $RLY_MNEMONIC_2 | $BINARY_2 keys add rly2 --home $CHAIN_DIR/$CHAINID_2 --recover --keyring-backend test

$BINARY add-genesis-account $($BINARY --home $CHAIN_DIR/$CHAINID_1 keys show val1 --keyring-backend test -a) 100000000000stake --home $CHAIN_DIR/$CHAINID_1
$BINARY_2 add-genesis-account $($BINARY_2 --home $CHAIN_DIR/$CHAINID_2 keys show val2 --keyring-backend test -a) 100000000000stake  --home $CHAIN_DIR/$CHAINID_2
$BINARY add-genesis-account $($BINARY --home $CHAIN_DIR/$CHAINID_1 keys show demowallet1 --keyring-backend test -a) 100000000000stake --home $CHAIN_DIR/$CHAINID_1
$BINARY_2 add-genesis-account $($BINARY_2 --home $CHAIN_DIR/$CHAINID_2 keys show demowallet2 --keyring-backend test -a) 100000000000stake  --home $CHAIN_DIR/$CHAINID_2
$BINARY add-genesis-account $($BINARY --home $CHAIN_DIR/$CHAINID_1 keys show rly1 --keyring-backend test -a) 100000000000stake --home $CHAIN_DIR/$CHAINID_1
$BINARY_2 add-genesis-account $($BINARY_2 --home $CHAIN_DIR/$CHAINID_2 keys show rly2 --keyring-backend test -a) 100000000000stake  --home $CHAIN_DIR/$CHAINID_2

echo "Creating and collecting gentx..."
$BINARY gentx val1 7000000000stake --home $CHAIN_DIR/$CHAINID_1 --chain-id $CHAINID_1 --keyring-backend test
$BINARY_2 gentx val2 7000000000stake --home $CHAIN_DIR/$CHAINID_2 --chain-id $CHAINID_2 --keyring-backend test
$BINARY collect-gentxs --home $CHAIN_DIR/$CHAINID_1
$BINARY_2 collect-gentxs --home $CHAIN_DIR/$CHAINID_2

echo "Changing defaults and ports in app.toml and config.toml files..."
sed -i -e 's#"tcp://0.0.0.0:26656"#"tcp://0.0.0.0:'"$P2PPORT_1"'"#g' $CHAIN_DIR/$CHAINID_1/config/config.toml
sed -i -e 's#"tcp://127.0.0.1:26657"#"tcp://0.0.0.0:'"$RPCPORT_1"'"#g' $CHAIN_DIR/$CHAINID_1/config/config.toml
sed -i -e 's/timeout_commit = "5s"/timeout_commit = "1s"/g' $CHAIN_DIR/$CHAINID_1/config/config.toml
sed -i -e 's/timeout_propose = "3s"/timeout_propose = "1s"/g' $CHAIN_DIR/$CHAINID_1/config/config.toml
sed -i -e 's/index_all_keys = false/index_all_keys = true/g' $CHAIN_DIR/$CHAINID_1/config/config.toml
sed -i -e 's/enable = false/enable = true/g' $CHAIN_DIR/$CHAINID_1/config/app.toml
sed -i -e 's/swagger = false/swagger = true/g' $CHAIN_DIR/$CHAINID_1/config/app.toml
sed -i -e 's#"tcp://0.0.0.0:1317"#"tcp://0.0.0.0:'"$RESTPORT_1"'"#g' $CHAIN_DIR/$CHAINID_1/config/app.toml
sed -i -e 's#":8080"#":'"$ROSETTA_1"'"#g' $CHAIN_DIR/$CHAINID_1/config/app.toml


sed -i -e 's/timeout_commit = "5s"/timeout_commit = "1s"/g' $CHAIN_DIR/$CHAINID_2/config/config.toml
sed -i -e 's/timeout_propose = "3s"/timeout_propose = "1s"/g' $CHAIN_DIR/$CHAINID_2/config/config.toml
sed -i -e 's/index_all_keys = false/index_all_keys = true/g' $CHAIN_DIR/$CHAINID_2/config/config.toml
sed -i -e 's/enable = false/enable = true/g' $CHAIN_DIR/$CHAINID_2/config/app.toml
sed -i -e 's/swagger = false/swagger = true/g' $CHAIN_DIR/$CHAINID_2/config/app.toml
sed -i -e 's#"tcp://0.0.0.0:1317"#"tcp://0.0.0.0:'"$RESTPORT_2"'"#g' $CHAIN_DIR/$CHAINID_2/config/app.toml
sed -i -e 's#":8080"#":'"$ROSETTA_2"'"#g' $CHAIN_DIR/$CHAINID_2/config/app.toml


echo "Update genesis.json file with updated local params"
jq -r '.app_state.staking.params.unbonding_time |= "30s"' $CHAIN_DIR/$CHAINID_1/config/genesis.json > /tmp/genesis.json; mv /tmp/genesis.json $HOME/.persistenceCore/config/genesis.json
jq -r '.app_state.slashing.params.downtime_jail_duration |= "6s"' $CHAIN_DIR/$CHAINID_1/config/genesis.json > /tmp/genesis.json; mv /tmp/genesis.json $HOME/.persistenceCore/config/genesis.json
jq -r '.app_state.gov.deposit_params.max_deposit_period |= "30s"' $CHAIN_DIR/$CHAINID_1/config/genesis.json > /tmp/genesis.json; mv /tmp/genesis.json $HOME/.persistenceCore/config/genesis.json
jq -r '.app_state.gov.deposit_params.min_deposit[0].amount |= "10"' $CHAIN_DIR/$CHAINID_1/config/genesis.json > /tmp/genesis.json; mv /tmp/genesis.json $HOME/.persistenceCore/config/genesis.json
jq -r '.app_state.gov.voting_params.voting_period |= "30s"' $CHAIN_DIR/$CHAINID_1/config/genesis.json > /tmp/genesis.json; mv /tmp/genesis.json $HOME/.persistenceCore/config/genesis.json
jq -r '.app_state.gov.tally_params.quorum |= "0.000000000000000000"' $CHAIN_DIR/$CHAINID_1/config/genesis.json > /tmp/genesis.json; mv /tmp/genesis.json $HOME/.persistenceCore/config/genesis.json
jq -r '.app_state.gov.tally_params.threshold |= "0.000000000000000000"' $CHAIN_DIR/$CHAINID_1/config/genesis.json > /tmp/genesis.json; mv /tmp/genesis.json $HOME/.persistenceCore/config/genesis.json
jq -r '.app_state.gov.tally_params.veto_threshold |= "0.000000000000000000"' $CHAIN_DIR/$CHAINID_1/config/genesis.json > /tmp/genesis.json; mv /tmp/genesis.json $HOME/.persistenceCore/config/genesis.json


# Update host chain genesis to allow x/bank/MsgSend ICA tx execution
sed -i -e 's/\"allow_messages\":.*/\"allow_messages\": [\"\/cosmos.bank.v1beta1.MsgSend\", \"\/cosmos.staking.v1beta1.MsgDelegate\"]/g' $CHAIN_DIR/$CHAINID_2/config/genesis.json

# Set wasm as permissioned or permissionless based on environment variable
wasm_permission="Nobody"
if  $WASM_PERMISSIONLESS
then
  wasm_permission="Everybody"
fi
jq -r ".app_state.wasm.params.code_upload_access.permission |= \"${wasm_permission}\"" $CHAIN_DIR/$CHAINID_1/config/genesis.json > /tmp/genesis.json; mv /tmp/genesis.json $CHAIN_DIR/$CHAINID_1/config/genesis.json
