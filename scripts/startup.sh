coreNode start>~/.coreNode/log &
sleep 10 ; coreClient rest-server --chain-id test>>~/.coreClient/log &
echo "
Node and Client started up. For logs:
tail -f ~/.coreNode/log
tail -f ~/.coreClient/log
"