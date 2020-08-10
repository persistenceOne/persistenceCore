
#add chain id in config
assetClient config chain-id test

#set env variables
NONCE="$RANDOM"
SLEEP=6
PASSWD="123123123"

#Create users
assetClient keys add main$NONCE <<< "$PASSWD"
assetClient keys add eve$NONCE <<< "$PASSWD"
assetClient keys add bob$NONCE <<< "$PASSWD"
assetClient keys add alice$NONCE <<< "$PASSWD"

#name users with their addresses
TEST=$(assetClient keys show -a test <<< "$PASSWD")
MAIN=$(assetClient keys show -a main$NONCE <<< "$PASSWD")
ALICE=$(assetClient keys show -a alice$NONCE <<< "$PASSWD")
BOB=$(assetClient keys show -a bob$NONCE <<< "$PASSWD")
EVE=$(assetClient keys show -a eve$NONCE <<< "$PASSWD")

#Load coins in main
assetClient tx send $TEST $MAIN 10000persistence,10000stake,10000atom -y <<< "$PASSWD
$PASSWD"
sleep $SLEEP
#send coins in users
assetClient tx send $MAIN $EVE 100persistence,100stake,100atom -y <<< "$PASSWD
$PASSWD"
sleep $SLEEP
assetClient tx send $EVE $BOB 10persistence,10stake,10atom -y <<< "$PASSWD
$PASSWD"
assetClient tx send $MAIN $ALICE 100persistence,100stake,100atom -y <<< "$PASSWD
$PASSWD"
sleep $SLEEP

#recursively send coins
assetClient tx send $MAIN $EVE 100persistence,100stake,100atom -y <<< "$PASSWD
$PASSWD"
assetClient tx send $EVE $BOB 50persistence,50stake,50atom -y <<< "$PASSWD
$PASSWD"
assetClient tx send $BOB $ALICE 5persistence,5stake,5atom -y <<< "$PASSWD
$PASSWD"
assetClient tx send $ALICE $MAIN 1persistence,1stake,1atom -y <<< "$PASSWD
$PASSWD"
sleep $SLEEP

# identities issue, provision, unprovision
ID_1=identity1$NONCE
ID_2=identity2$NONCE
ID_3=identity3$NONCE
assetClient tx identities issue -y --from $MAIN --to $MAIN  --properties $ID_1:$ID_1 <<< "$PASSWD
$PASSWD"
assetClient tx identities issue -y --from $BOB --to $BOB  --properties $ID_2:$ID_2 <<< "$PASSWD
$PASSWD"
assetClient tx identities issue -y --from $ALICE --to $ALICE  --properties $ID_3:$ID_3 <<< "$PASSWD
$PASSWD"
sleep $SLEEP
MAIN_ID=test...$(echo $(assetClient q identities identities) | awk -v var=$ID_1 '{for(i=1;i<=NF;i++)if($i=="hashid:" && $(i+14)==var)print $(i+2)}')
BOB_ID=test...$(echo $(assetClient q identities identities) | awk -v var=$ID_2 '{for(i=1;i<=NF;i++)if($i=="hashid:" && $(i+14)==var)print $(i+2)}')
ALICE_ID=test...$(echo $(assetClient q identities identities) | awk -v var=$ID_3 '{for(i=1;i<=NF;i++)if($i=="hashid:" && $(i+14)==var)print $(i+2)}')
#provision identities
assetClient tx identities provision -y --from $MAIN --to $EVE --identityID $MAIN_ID <<< "$PASSWD
$PASSWD"
sleep $SLEEP
assetClient tx identities unprovision -y --from $MAIN --to $EVE --identityID $MAIN_ID <<< "$PASSWD
$PASSWD"
sleep $SLEEP
assetClient query identities identities

#assets mint, mutate burn
ASSET_P1=assets1$NONCE
ASSET_P2=assets2$NONCE
ASSET_P3=assets3$NONCE
ASSET_P4=assets4$NONCE
ASSET_P5=assets5$NONCE
ASSET_P6=assets6$NONCE
ASSET_P7=assets7$NONCE
assetClient tx assets mint -y --from $MAIN --fromID $MAIN_ID --toID $MAIN_ID --properties $ASSET_P1:$ASSET_P1 <<< "$PASSWD
$PASSWD"
assetClient tx assets mint -y --from $BOB --fromID $BOB_ID --toID $BOB_ID --properties $ASSET_P2:$ASSET_P2 <<< "$PASSWD
$PASSWD"
sleep $SLEEP
MAIN_ASSET_1=test...$(echo $(assetClient q assets assets) | awk -v var=$ASSET_P1 '{for(i=1;i<=NF;i++)if($i=="hashid:" && $(i+15)==var)print $(i+2)}')
BOB_ASSET_1=test...$(echo $(assetClient q assets assets) | awk -v var=$ASSET_P2 '{for(i=1;i<=NF;i++)if($i=="hashid:" && $(i+15)==var)print $(i+2)}')
assetClient tx assets mutate -y --from $MAIN --fromID $MAIN_ID --assetID $MAIN_ASSET_1 --properties $ASSET_P1:mutated$ASSET_P1 <<< "$PASSWD
$PASSWD"
sleep $SLEEP
assetClient tx assets burn -y --from $MAIN --fromID $MAIN_ID --assetID $MAIN_ASSET_1 <<< "$PASSWD
$PASSWD"
sleep $SLEEP
assetClient tx assets mint -y --from $MAIN --fromID $MAIN_ID --toID $MAIN_ID --properties $ASSET_P1:$ASSET_P1 <<< "$PASSWD
$PASSWD"
sleep $SLEEP

assetClient tx assets mint -y --from $MAIN --fromID $MAIN_ID --toID $MAIN_ID --properties $ASSET_P3:$ASSET_P3 <<< "$PASSWD
$PASSWD"
assetClient tx assets mint -y --from $ALICE --fromID $ALICE_ID --toID $ALICE_ID --properties $ASSET_P4:$ASSET_P4 <<< "$PASSWD
$PASSWD"
sleep $SLEEP
MAIN_ASSET_2=test...$(echo $(assetClient q assets assets) | awk -v var=$ASSET_P3 '{for(i=1;i<=NF;i++)if($i=="hashid:" && $(i+15)==var)print $(i+2)}')
ALICE_ASSET_1=test...$(echo $(assetClient q assets assets) | awk -v var=$ASSET_P4 '{for(i=1;i<=NF;i++)if($i=="hashid:" && $(i+15)==var)print $(i+2)}')

assetClient query assets assets

assetClient query splits splits

#order make and cancel
assetClient tx orders make --from $MAIN --fromID $MAIN_ID --toID $BOB_ID --makerSplit 1 --makerSplitID $MAIN_ASSET_1 --exchangeRate="1" --takerSplitID $BOB_ASSET_1 -y <<< "$PASSWD
$PASSWD"
sleep $SLEEP
MAIN_BOB_ORDER_1=test..$(echo $(assetClient q orders orders) | awk -v var=$MAIN_ASSET_1 '{for(i=1;i<=NF;i++)if($i=="hashid:"){for(j=1;j<=i+40;j++)if($j==var){print $(i+2)}}}')
assetClient tx orders cancel --from $MAIN --orderID $MAIN_BOB_ORDER_1 -y <<< "$PASSWD
$PASSWD"
sleep $SLEEP

#order make and take private
assetClient tx orders make --from $MAIN --fromID $MAIN_ID --toID $BOB_ID --makerSplit 1 --makerSplitID $MAIN_ASSET_1 --exchangeRate="1" --takerSplitID $BOB_ASSET_1 -y <<< "$PASSWD
$PASSWD"
sleep $SLEEP
MAIN_BOB_ORDER_1=test..$(echo $(assetClient q orders orders) | awk -v var=$MAIN_ASSET_1 '{for(i=1;i<=NF;i++)if($i=="hashid:"){for(j=1;j<=i+40;j++)if($j==var){print $(i+2)}}}')
assetClient tx orders take --from $BOB --orderID $MAIN_BOB_ORDER_1 --takerSplit 1 --fromID $BOB_ID -y <<< "$PASSWD
$PASSWD"
sleep $SLEEP

#order make and take public
assetClient tx orders make --from $MAIN --fromID $MAIN_ID --makerSplit 1 --makerSplitID $MAIN_ASSET_2 --exchangeRate="1" --takerSplitID $ALICE_ASSET_1 -y <<< "$PASSWD
$PASSWD"
sleep $SLEEP
MAIN_ORDER_2=test..$(echo $(assetClient q orders orders) | awk -v var=$MAIN_ASSET_2 '{for(i=1;i<=NF;i++)if($i=="hashid:"){for(j=1;j<=i+40;j++)if($j==var){print $(i+2)}}}')
assetClient tx orders take  --from $ALICE --orderID $MAIN_ORDER_2 --takerSplit 1 --fromID $ALICE_ID -y <<< "$PASSWD
$PASSWD"
sleep $SLEEP

#splits send
assetClient tx splits send -y --from $MAIN --fromID $MAIN_ID --toID $BOB_ID --ownableID $BOB_ASSET_1 --split "1" <<< "$PASSWD
$PASSWD"
assetClient tx splits send -y --from $BOB --fromID $BOB_ID --toID $MAIN_ID --ownableID $MAIN_ASSET_1 --split "1" <<< "$PASSWD
$PASSWD"
assetClient tx splits send -y --from $ALICE --fromID $ALICE_ID --toID $MAIN_ID --ownableID $MAIN_ASSET_2 --split "1" <<< "$PASSWD
$PASSWD"
sleep $SLEEP
assetClient tx splits send -y --from $MAIN --fromID $MAIN_ID --toID $ALICE_ID --ownableID $ALICE_ASSET_1 --split "1" <<< "$PASSWD
$PASSWD"
sleep $SLEEP

##wraping coins
assetClient tx splits wrap -y --from $MAIN --fromID $MAIN_ID --coins 50stake <<< "$PASSWD
$PASSWD"
assetClient tx splits wrap -y --from $BOB --fromID $BOB_ID --coins 50persistence <<< "$PASSWD
$PASSWD"
assetClient tx splits wrap -y --from $ALICE --fromID $ALICE_ID --coins 50atom <<< "$PASSWD
$PASSWD"
sleep $SLEEP

# orders maker asset taker split
assetClient tx orders make --from $MAIN --fromID $MAIN_ID --makerSplit 1 --makerSplitID $MAIN_ASSET_1 --exchangeRate="2.25" --takerSplitID persistence -y <<< "$PASSWD
$PASSWD"
sleep $SLEEP
MAIN_ORDER_1=test..$(echo $(assetClient q orders orders) | awk -v var=$MAIN_ASSET_1 '{for(i=1;i<=NF;i++)if($i=="hashid:"){for(j=1;j<=i+40;j++)if($j==var){print $(i+2)}}}')
assetClient tx orders take --from $BOB --fromID $BOB_ID --orderID $MAIN_ORDER_1 --takerSplit 5 -y <<< "$PASSWD
$PASSWD"
sleep $SLEEP

# orders maker split taker asset
assetClient tx orders make --from $MAIN --fromID $MAIN_ID --makerSplit 10 --makerSplitID stake --exchangeRate="0.1" --takerSplitID $MAIN_ASSET_1 -y <<< "$PASSWD
$PASSWD"
sleep $SLEEP
MAIN_ORDER_1=test..$(echo $(assetClient q orders orders) | awk -v var=$MAIN_ASSET_1 '{for(i=1;i<=NF;i++)if($i=="hashid:"){for(j=1;j<=i+60;j++)if($j==var){print $(i+2)}}}')
assetClient tx orders take --from $BOB --fromID $BOB_ID --orderID $MAIN_ORDER_1 --takerSplit 1 -y <<< "$PASSWD
$PASSWD"
sleep $SLEEP

# orders maker split taker split
assetClient tx orders make --from $MAIN --fromID $MAIN_ID --makerSplit 10 --makerSplitID stake --exchangeRate="0.7" --takerSplitID atom -y <<< "$PASSWD
$PASSWD"
sleep $SLEEP
MAIN_ORDER_2=test..$(echo $(assetClient q orders orders) | awk -v var=MAIN_ID '{for(i=1;i<=NF;i++)if($i=="hashid:"){for(j=1;j<=i+40;j++)if($j==var){print $(i+2)}}}')
assetClient tx orders take --from $ALICE --fromID $ALICE_ID --orderID $MAIN_ORDER_2 --takerSplit 9 -y <<< "$PASSWD
$PASSWD"
sleep $SLEEP
