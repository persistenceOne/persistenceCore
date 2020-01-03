hubNode init test --chain-id test
hubClient keys add test
hubNode add-genesis-account test 100000000000000stake
hubNode gentx --name test --amount 1000000000stake
hubNode collect-gentxs