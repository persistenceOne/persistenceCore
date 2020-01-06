rm -rf ~/.hubNode
rm -rf ~/.hubClient
make all
hubNode init test --chain-id test
hubClient keys add test<<!
qweqweqwe
qweqweqwe
!
hubNode add-genesis-account test 100000000000000stake
hubNode gentx --name test --amount 1000000000stake<<!
qweqweqwe
!
hubNode collect-gentxs