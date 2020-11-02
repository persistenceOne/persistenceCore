
killall assetNode
killall assetClient
echo "
Node and Client Shutdown."

sleep 4

rm -rf ~/.assetNode
rm -rf ~/.assetClient

mkdir ~/.assetNode
mkdir ~/.assetClient

assetNode init --chain-id test test
assetClient keys add test --recover<<<"y
wage thunder live sense resemble foil apple course spin horse glass mansion midnight laundry acoustic rhythm loan scale talent push green direct brick please"
assetNode add-genesis-account test 100000000000000stake
assetNode gentx --name test --amount 1000000000stake
assetNode collect-gentxs


assetNode start >~/.assetNode/log &
sleep 10
assetClient rest-server --chain-id test -b block >~/.assetClient/log &

echo "
Node and Client started up. For logs:
tail -f ~/.assetNode/log
tail -f ~/.assetClient/log
"

npm run test:awesome
