# Deployment and Instantiation

These are the steps to upload and instantiate a contract

- Submit the upload proposal
- Vote on the proposal, wait for it to pass
- Submit proposal for instantiating the contract
- Vote on the proposal, wait for it to pass
- The contract is now instantiated and ready to be interacted with.

See the list of codes that was uploaded to the testnet previously.
```
persistenceCore query wasm list-code --node https://rpc.testnet.persistence.one:443
```

You can set the `node` to the persistenceCore config and don't have to worry about passing that flag always

```
persistenceCore config node https://rpc.testnet.persistence.one:443
```

To upload the contract via proposal

```
RESP=$(persistenceCore tx gov submit-proposal wasm-store "<path/to/the/compiled/wasm>" \
  --title "title" \
  --description "description" \
  --deposit 10000uxprt \
  --run-as $TEST_KEY \
  --instantiate-everybody "true" \
  --keyring-backend test \
  --from $TEST_KEY --gas auto --fees 10000uxprt -y \
  --chain-id test-core-1 \
  -b block -o json --gas-adjustment 1.1)
  
echo $RESP 
```
The `$TEST_KEY` can be any valid persistenceAddress. Make sure it has some test tokens.

Now `$RESP` has the proposalID, extract the proposal_ID and vote on it

```
PROPOSAL_ID=$(echo "$RESP" | jq -r '.logs[0].events[] | select(.type == "submit_proposal") | .attributes[] | select(.key == "proposal_id") | .value')

persistenceCore tx gov vote $PROPOSAL_ID yes --from $TEST_KEY --yes --chain-id test-core-1 \
    --fees 500uxprt --gas auto --gas-adjustment 1.1 -b block --keyring-backend test -o json | jq
```

The contract is instantiated via a gov-proposal 

Get the $CODE_ID by doing `wasm list-code` after the proposal passes and proceed towards instantiation

```
RESP=$(persistenceCore tx gov submit-proposal instantiate-contract $CODE_ID "$INIT" \
  --admin="$TEST_KEY" \
  --from $TEST_KEY \
  --deposit 10000uxprt \
  --label "label" \
  --title "title" \
  --description "description" \
  --gas-adjustment 1.1 \
  --fees "10000uxprt" \
  --gas "auto" \
  --run-as $TEST_KEY \
  -y --chain-id test-core-1 -b block -o json)
```

The `$TEST_KEY` can be any valid persistenceAddress. Make sure it has some test tokens.

The `$INIT` variable has to be in a json structure with all the variables needed for the contract
example (init variable for the cw20 contract) :
```
INIT=$(cat <<EOF
{
  "name": "My first token",
  "symbol": "FRST",
  "decimals": 6,
  "initial_balances": [{
    "address": "$TEST_KEY",
    "amount": "123456789000"
  }]
}
EOF
)

```


Extract the proposal_ID and vote on it.

After the proposal passes, the contract will be instantiated and can be interacted with!!