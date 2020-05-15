make all

rm -rf ~/.coreNode
rm -rf ~/.coreClient

coreNode init --chain-id test test
coreClient keys add test
coreNode add-genesis-account test 100000000000000stake
coreNode gentx --name test --amount 1000000000stake
coreNode collect-gentxs