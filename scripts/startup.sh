rm -rf ~/.coreLog
mkdir ~/.coreLog
coreNode start>~/.coreLog/coreNode &
sleep 10 ; coreClient rest-server --chain-id test>>~/.coreLog/coreClient &
echo "
Node and Client started up. For logs:
tail -f ~/.coreLog/coreNode
tail -f ~/.coreLog/coreClient
Save logs before restarting.
"