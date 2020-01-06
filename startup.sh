hubNode start>~/.hubNode/log &
sleep 10 ; hubClient rest-server --chain-id test>>~/.hubClient/log &
echo "
Node and Client started up. For logs:
tail -f ~/.hubNode/log
tail -f ~/.hubClient/log
"