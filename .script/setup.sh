make all

rm -rf ~/.persistenceNode
rm -rf ~/.persistenceClient

mkdir ~/.persistenceNode
mkdir ~/.persistenceClient

persistenceNode init --chain-id test test
persistenceClient keys add test --recover<<<"y
wage thunder live sense resemble foil apple course spin horse glass mansion midnight laundry acoustic rhythm loan scale talent push green direct brick please"
persistenceNode add-genesis-account test 100000000000000stake
persistenceNode gentx --name test --amount 1000000000stake
persistenceNode collect-gentxs