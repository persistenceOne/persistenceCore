persistenceNode start >~/.persistenceNode/log &
sleep 10
persistenceClient rest-server --chain-id test $1 $2>~/.persistenceClient/log &
echo "
Node and Client started up. For logs:
tail -f ~/.persistenceNode/log
tail -f ~/.persistenceClient/log
"
