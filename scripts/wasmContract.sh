# genesis account / chain -id is test, commands to store, instantiate, execute a contract. this eg- hackatom, github.com/CosmWasm/cosmwasm/contracts/hackatom

assetClient tx wasm store /PATH_TO_WASM_COMTRACT/_.wasm --from test --gas 900000  -y --chain-id test

CODE_ID=$(assetClient query wasm list-code --chain-id test| jq .[-1].id)

assetClient keys add bob

INIT=$(jq -n --arg test $(assetClient keys show -a test) --arg bob $(assetClient keys show -a bob) '{"verifier":$test,"beneficiary":$bob}')

assetClient tx wasm instantiate $CODE_ID "$INIT" --from test --amount=50000stake  --label "escrow 1" -y --chain-id test

CONTRACT=$(assetClient query wasm list-contract-by-code $CODE_ID --chain-id test| jq -r .[0].address)

MINT='{"asset_mint":{"properties":"test5:7, test89:76"}}'

assetClient tx wasm execute $CONTRACT "$MINT" --from test -y --chain-id test

# issue asset normal
assetClient tx assetFactory mint --from test --properties test1:test1,test2:test2  --chain-id test
