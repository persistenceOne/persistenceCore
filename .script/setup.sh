#!/bin/bash
.script/reset.sh

test_mnemonic="wage thunder live sense resemble foil apple course spin horse glass mansion midnight laundry acoustic rhythm loan scale talent push green direct brick please"

persistenceCore init test --chain-id test
echo $test_mnemonic | persistenceCore keys add test --recover --keyring-backend test
persistenceCore add-genesis-account test 100000000000000stake --keyring-backend test
persistenceCore gentx test 10000000stake --chain-id test --keyring-backend test
persistenceCore collect-gentxs