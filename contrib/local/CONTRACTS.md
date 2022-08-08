# Deployment and Instantiation

See the list of codes that was uploaded to the testnet previously.
```
persistenceCore query wasm list-code --node https://rpc.testnet.persistence.one:443
```

You can set the `node` to the persistenceCore config and don't have to worry about passing that flag always

```
persistenceCore config node https://rpc.testnet.persistence.one:443
```

To upload the contract

```
RESP=$(persistenceCore tx gov submit-proposal wasm-store "path/to/the/compiled/wasm" \
  --title "title" \
  --description "description" \
  --deposit 10000stake \
  --run-as $VAL_ADDR \
  --instantiate-everybody "true" \
  --keyring-backend test \
  --from $VAL_ADDR --gas auto --fees 10000stake -y \
  --chain-id $CHAIN_ID \
  -b block -o json --gas-adjustment 1.5)
  
echo $RESP 
```

Now resp has the proposalID, extract the proposal_ID and vote on it

```
PROPOSAL_ID=$(echo "$RESP" | jq -r '.logs[0].events[] | select(.type == "submit_proposal") | .attributes[] | select(.key == "proposal_id") | .value')

persistenceCore tx gov vote $PROPOSAL_ID yes --from $VAL_ADDR --yes --chain-id $CHAIN_ID \
    --fees 500stake --gas auto --gas-adjustment 1.5 -b block --keyring-backend test -o json | jq
```

Now get the $CODE_ID by doing `wasm list-code` after the proposal passes and can proceed towards instantiation

```
persistenceCore tx wasm instantiate "$CODE_ID" "$INIT" --admin="$(persistenceCore keys show $VAL_ADDR -a --keyring-backend test)" \
  --from $VAL_ADDR --amount "10000stake" --label "local0.1.0" --gas-adjustment 1.5 --fees "10000stake" \
  --gas "auto" -y --chain-id $CHAIN_ID -b block --keyring-backend test -o json | jq
```

The `$INIT` variable has to be in a json structure with all the variables needed for the contract
example (init variable for the cw20 contract) :
```
INIT=$(cat <<EOF
{
  "name": "My first token",
  "symbol": "FRST",
  "decimals": 6,
  "initial_balances": [{
    "address": "$TEST1_KEY",
    "amount": "123456789000"
  }]
}
EOF
)
```
The contract is instantiated and can be interacted with now!!