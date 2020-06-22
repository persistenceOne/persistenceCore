assetClient tx wasm store ../../CosmWasm/wasmd/x/wasm/internal/keeper/testdata/contract.wasm --from test --gas 600000  -y --chain-id test

CODE_ID=$(assetClient query wasm list-code --chain-id test| jq .[-1].id)

assetClient keys add bob

INIT=$(jq -n --arg test $(assetClient keys show -a test) --arg bob $(assetClient keys show -a bob) '{"verifier":$test,"beneficiary":$bob}')

assetClient tx wasm instantiate $CODE_ID "$INIT" --from test --amount=50000stake  --label "escrow 1" -y --chain-id test

CONTRACT=$(assetClient query wasm list-contract-by-code $CODE_ID --chain-id test| jq -r .[0].address)

APPROVE='{"asset_mint":{"properties":"test1:1, test2:4"}}'

assetClient tx wasm execute $CONTRACT "$APPROVE" --from test -y --chain-id test

# issue assset

assetClient tx assetFactory mint --from test --properties test1:test1,test2:test2  --chain-id test