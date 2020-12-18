#add chain id in config
assetClient config chain-id test

#set env variables
NONCE="$RANDOM"
SLEEP=6
PASSWD="123123123"
KEYRING="--keyring-backend os"
MODE="-b sync"

#Create users
ACCOUNT_NAME_1=account1$NONCE
ACCOUNT_NAME_2=account2$NONCE
ACCOUNT_NAME_3=account3$NONCE
ACCOUNT_NAME_4=account4$NONCE
assetClient keys add $ACCOUNT_NAME_1 $KEYRING
assetClient keys add $ACCOUNT_NAME_2 $KEYRING
assetClient keys add $ACCOUNT_NAME_3 $KEYRING
assetClient keys add $ACCOUNT_NAME_4 $KEYRING

#name users with their addresses
TEST=$(assetClient keys show -a test $KEYRING)
ACCOUNT_1=$(assetClient keys show -a $ACCOUNT_NAME_1 $KEYRING)
ACCOUNT_2=$(assetClient keys show -a $ACCOUNT_NAME_2 $KEYRING)
ACCOUNT_3=$(assetClient keys show -a $ACCOUNT_NAME_3 $KEYRING)
ACCOUNT_4=$(assetClient keys show -a $ACCOUNT_NAME_4 $KEYRING)

#Load coins in main
assetClient tx send $TEST $ACCOUNT_1 10000stake -y $KEYRING $MODE
sleep $SLEEP
#send coins in users
assetClient tx send $ACCOUNT_1 $ACCOUNT_3 110stake -y $KEYRING $MODE
sleep $SLEEP
assetClient tx send $ACCOUNT_3 $ACCOUNT_4 10stake -y $KEYRING $MODE
assetClient tx send $ACCOUNT_1 $ACCOUNT_2 100stake -y $KEYRING $MODE
sleep $SLEEP

#recursively send coins
assetClient tx send $ACCOUNT_1 $ACCOUNT_3 100stake -y $KEYRING $MODE
assetClient tx send $ACCOUNT_3 $ACCOUNT_2 50stake -y $KEYRING $MODE
assetClient tx send $ACCOUNT_2 $ACCOUNT_4 5stake -y $KEYRING $MODE
assetClient tx send $ACCOUNT_4 $ACCOUNT_2 5stake -y $KEYRING $MODE
sleep $SLEEP

# identities nub, define, issue, provision, unprovision
NUB_ID_1=nubID1$NONCE
NUB_ID_2=nubID2$NONCE
NUB_ID_3=nubID3$NONCE
assetClient tx identities nub -y --from $ACCOUNT_1 --nubID $NUB_ID_1 $KEYRING $MODE
assetClient tx identities nub -y --from $ACCOUNT_2 --nubID $NUB_ID_2 $KEYRING $MODE
assetClient tx identities nub -y --from $ACCOUNT_3 --nubID $NUB_ID_3 $KEYRING $MODE
sleep $SLEEP
ACCOUNT_1_NUB_ID=$(echo $(assetClient q identities identities) | awk -v var="$ACCOUNT_1" '{for(i=1;i<=NF;i++)if($i==var)print $(i-6)"|"$(i-3)}')
ACCOUNT_2_NUB_ID=$(echo $(assetClient q identities identities) | awk -v var="$ACCOUNT_2" '{for(i=1;i<=NF;i++)if($i==var)print $(i-6)"|"$(i-3)}')
ACCOUNT_3_NUB_ID=$(echo $(assetClient q identities identities) | awk -v var="$ACCOUNT_3" '{for(i=1;i<=NF;i++)if($i==var)print $(i-6)"|"$(i-3)}')

IDENTITY_DEFINE_IMMUTABLE_1_ID="identityDefineImmutable1$NONCE"
IDENTITY_DEFINE_IMMUTABLE_1="$IDENTITY_DEFINE_IMMUTABLE_1_ID:S|"
IDENTITY_DEFINE_IMMUTABLE_META_1_ID="identityDefineImmutableMeta1$NONCE"
IDENTITY_DEFINE_IMMUTABLE_META_1="$IDENTITY_DEFINE_IMMUTABLE_META_1_ID:I|identityDefineImmutableMeta1$NONCE"
IDENTITY_DEFINE_MUTABLE_1_ID="identityDefineMutable1$NONCE"
IDENTITY_DEFINE_MUTABLE_1="$IDENTITY_DEFINE_MUTABLE_1_ID:D|"
IDENTITY_DEFINE_MUTABLE_META_1_ID="identityDefineMutableMeta1$NONCE"
IDENTITY_DEFINE_MUTABLE_META_1="$IDENTITY_DEFINE_MUTABLE_META_1_ID:H|"
assetClient tx identities define -y --from $ACCOUNT_1 --fromID $ACCOUNT_1_NUB_ID \
 --immutableTraits "$IDENTITY_DEFINE_IMMUTABLE_1" \
 --immutableMetaTraits "$IDENTITY_DEFINE_IMMUTABLE_META_1" \
 --mutableTraits "$IDENTITY_DEFINE_MUTABLE_1" \
 --mutableMetaTraits "$IDENTITY_DEFINE_MUTABLE_META_1" $KEYRING $MODE

sleep $SLEEP
IDENTITY_DEFINE_CLASSIFICATION=$(echo $(assetClient q classifications classifications) | awk -v var="$IDENTITY_DEFINE_IMMUTABLE_META_1_ID" '{for(i=1;i<=NF;i++)if($i==var)print $(i-10)"."$(i-7)}')

assetClient tx identities issue -y --from $ACCOUNT_1 --fromID $ACCOUNT_1_NUB_ID --classificationID $IDENTITY_DEFINE_CLASSIFICATION --to $ACCOUNT_1 \
 --immutableProperties "$IDENTITY_DEFINE_IMMUTABLE_1""stringValue" \
 --immutableMetaProperties "$IDENTITY_DEFINE_IMMUTABLE_META_1" \
 --mutableProperties "$IDENTITY_DEFINE_MUTABLE_1""1.01" \
 --mutableMetaProperties "$IDENTITY_DEFINE_MUTABLE_META_1""123" \
 $KEYRING $MODE

sleep $SLEEP
IDENTITY_ISSUE_ACCOUNT_1=$(echo $(assetClient q identities identities) | awk -v var="$IDENTITY_DEFINE_CLASSIFICATION" '{for(i=1;i<=NF;i++)if($i==var)print $i"|"$(i+3)}')
assetClient tx identities provision -y --from $ACCOUNT_1 --to $ACCOUNT_4 --identityID $IDENTITY_ISSUE_ACCOUNT_1  $KEYRING $MODE
sleep $SLEEP
assetClient tx identities unprovision -y --from $ACCOUNT_1 --to $ACCOUNT_4 --identityID $IDENTITY_ISSUE_ACCOUNT_1 $KEYRING $MODE
sleep $SLEEP


#metas reveal
assetClient tx metas reveal -y --from $ACCOUNT_1 --metaFact "S|stringValue$NONCE"
assetClient tx metas reveal -y --from $ACCOUNT_2 --metaFact "I|identityValue$NONCE"
assetClient tx metas reveal -y --from $ACCOUNT_3 --metaFact "D|0.101010$NONCE"
assetClient tx metas reveal -y --from $ACCOUNT_4 --metaFact "H|1$NONCE"
sleep $SLEEP

#assets mint, mutate burn
ASSET_DEFINE_IMMUTABLE_1_ID="assetDefineImmutable1$NONCE"
ASSET_DEFINE_IMMUTABLE_1="$ASSET_DEFINE_IMMUTABLE_1_ID:S|"
ASSET_DEFINE_IMMUTABLE_META_1_ID="assetDefineImmutableMeta1$NONCE"
ASSET_DEFINE_IMMUTABLE_META_1="$ASSET_DEFINE_IMMUTABLE_META_1_ID:I|assetDefineImmutableMeta1$NONCE"
ASSET_DEFINE_MUTABLE_1_ID="assetDefineMutable1$NONCE"
ASSET_DEFINE_MUTABLE_1="$ASSET_DEFINE_MUTABLE_1_ID:D|"
ASSET_DEFINE_MUTABLE_META_1_ID="assetDefineMutableMeta1$NONCE"
ASSET_DEFINE_MUTABLE_META_1="$ASSET_DEFINE_MUTABLE_META_1_ID:H|"
assetClient tx assets define -y --from $ACCOUNT_1 --fromID $ACCOUNT_1_NUB_ID \
 --immutableTraits "$ASSET_DEFINE_IMMUTABLE_1" \
 --immutableMetaTraits "$ASSET_DEFINE_IMMUTABLE_META_1" \
 --mutableTraits "$ASSET_DEFINE_MUTABLE_1" \
 --mutableMetaTraits "$ASSET_DEFINE_MUTABLE_META_1" $KEYRING $MODE


ASSET_DEFINE_IMMUTABLE_2_ID="assetDefineImmutable2$NONCE"
ASSET_DEFINE_IMMUTABLE_2="$ASSET_DEFINE_IMMUTABLE_2_ID:S|"
ASSET_DEFINE_IMMUTABLE_META_2_ID="assetDefineImmutableMeta2$NONCE"
ASSET_DEFINE_IMMUTABLE_META_2="$ASSET_DEFINE_IMMUTABLE_META_2_ID:I|assetDefineImmutableMeta$NONCE"
ASSET_DEFINE_MUTABLE_2_ID="assetDefineMutable2$NONCE"
ASSET_DEFINE_MUTABLE_2="$ASSET_DEFINE_MUTABLE_2_ID:D|"
ASSET_DEFINE_MUTABLE_META_2_ID="assetDefineMutableMeta2$NONCE"
ASSET_DEFINE_MUTABLE_META_2="$ASSET_DEFINE_MUTABLE_META_2_ID:H|"
assetClient tx assets define -y --from $ACCOUNT_2 --fromID $ACCOUNT_2_NUB_ID \
 --immutableTraits "$ASSET_DEFINE_IMMUTABLE_2" \
 --immutableMetaTraits "$ASSET_DEFINE_IMMUTABLE_META_2" \
 --mutableTraits "$ASSET_DEFINE_MUTABLE_2" \
 --mutableMetaTraits "$ASSET_DEFINE_MUTABLE_META_2" $KEYRING $MODE

sleep $SLEEP
ASSET_DEFINE_CLASSIFICATION_1=$(echo $(assetClient q classifications classifications) | awk -v var="$ASSET_DEFINE_IMMUTABLE_META_1_ID" '{for(i=1;i<=NF;i++)if($i==var)print $(i-10)"."$(i-7)}')
assetClient tx assets mint -y --from $ACCOUNT_1 --fromID $ACCOUNT_1_NUB_ID --classificationID $ASSET_DEFINE_CLASSIFICATION_1 --toID $ACCOUNT_1_NUB_ID \
 --immutableProperties "$ASSET_DEFINE_IMMUTABLE_1""stringValue" \
 --immutableMetaProperties "$ASSET_DEFINE_IMMUTABLE_META_1" \
 --mutableProperties "$ASSET_DEFINE_MUTABLE_1""1.01" \
 --mutableMetaProperties "$ASSET_DEFINE_MUTABLE_META_1""123" \
 $KEYRING $MODE

ASSET_DEFINE_CLASSIFICATION_2=$(echo $(assetClient q classifications classifications) | awk -v var="$ASSET_DEFINE_IMMUTABLE_META_2_ID" '{for(i=1;i<=NF;i++)if($i==var)print $(i-10)"."$(i-7)}')
assetClient tx assets mint -y --from $ACCOUNT_2 --fromID $ACCOUNT_2_NUB_ID --classificationID $ASSET_DEFINE_CLASSIFICATION_2 --toID $ACCOUNT_2_NUB_ID \
 --immutableProperties "$ASSET_DEFINE_IMMUTABLE_2""stringValue" \
 --immutableMetaProperties "$ASSET_DEFINE_IMMUTABLE_META_2" \
 --mutableProperties "$ASSET_DEFINE_MUTABLE_2""1.01" \
 --mutableMetaProperties "$ASSET_DEFINE_MUTABLE_META_2""123" \
 $KEYRING $MODE

sleep $SLEEP
ASSET_MINT_1=$(echo $(assetClient q assets assets) | awk -v var="$ASSET_DEFINE_CLASSIFICATION_1" '{for(i=1;i<=NF;i++)if($i==var)print $i"|"$(i+3)}')
ASSET_MINT_2=$(echo $(assetClient q assets assets) | awk -v var="$ASSET_DEFINE_CLASSIFICATION_2" '{for(i=1;i<=NF;i++)if($i==var)print $i"|"$(i+3)}')

assetClient tx assets mutate -y --from $ACCOUNT_1 --fromID $ACCOUNT_1_NUB_ID --assetID $ASSET_MINT_1 \
 --mutableProperties "$ASSET_DEFINE_MUTABLE_1""1.012" \
 --mutableMetaProperties "$ASSET_DEFINE_MUTABLE_META_1""1234" $KEYRING $MODE
assetClient tx assets burn -y --from $ACCOUNT_2 --fromID $ACCOUNT_2_NUB_ID --assetID $ASSET_MINT_2 $KEYRING $MODE
sleep $SLEEP
#remint asset2
assetClient tx assets mint -y --from $ACCOUNT_2 --fromID $ACCOUNT_2_NUB_ID --classificationID $ASSET_DEFINE_CLASSIFICATION_2 --toID $ACCOUNT_2_NUB_ID \
 --immutableProperties "$ASSET_DEFINE_IMMUTABLE_2""stringValue" \
 --immutableMetaProperties "$ASSET_DEFINE_IMMUTABLE_META_2" \
 --mutableProperties "$ASSET_DEFINE_MUTABLE_2""1.01" \
 --mutableMetaProperties "$ASSET_DEFINE_MUTABLE_META_2""123" \
 $KEYRING $MODE
sleep $SLEEP

##wraping unwrapping send coins
assetClient tx splits wrap -y --from $ACCOUNT_1 --fromID $ACCOUNT_1_NUB_ID --coins 20stake $KEYRING $MODE
assetClient tx splits wrap -y --from $ACCOUNT_2 --fromID $ACCOUNT_2_NUB_ID --coins 20stake $KEYRING $MODE
assetClient tx splits wrap -y --from $ACCOUNT_3 --fromID $ACCOUNT_3_NUB_ID --coins 20stake $KEYRING $MODE
sleep $SLEEP
assetClient tx splits unwrap -y --from $ACCOUNT_1 --fromID $ACCOUNT_1_NUB_ID --ownableID stake --split 1 $KEYRING $MODE
assetClient tx splits unwrap -y --from $ACCOUNT_2 --fromID $ACCOUNT_2_NUB_ID --ownableID stake --split 1 $KEYRING $MODE
assetClient tx splits unwrap -y --from $ACCOUNT_3 --fromID $ACCOUNT_3_NUB_ID --ownableID stake --split 1 $KEYRING $MODE
sleep $SLEEP
assetClient tx splits send -y --from $ACCOUNT_3 --fromID $ACCOUNT_3_NUB_ID --toID $ACCOUNT_3_NUB_ID --ownableID stake --split 1 $KEYRING $MODE

#order make take cancel
ORDER_MUTABLE_META_TRAITS="takerID:I|,exchangeRate:D|,expiry:H|,makerOwnableSplit:D|"
ORDER_DEFINE_IMMUTABLE_1_ID="orderDefineImmutable1$NONCE"
ORDER_DEFINE_IMMUTABLE_1="$ORDER_DEFINE_IMMUTABLE_1_ID:S|"
ORDER_DEFINE_IMMUTABLE_META_1_ID="orderDefineImmutableMeta1$NONCE"
ORDER_DEFINE_IMMUTABLE_META_1="$ORDER_DEFINE_IMMUTABLE_META_1_ID:I|orderDefineImmutableMeta1$NONCE"
ORDER_DEFINE_MUTABLE_1_ID="orderDefineMutable1$NONCE"
ORDER_DEFINE_MUTABLE_1="$ORDER_DEFINE_MUTABLE_1_ID:D|"
ORDER_DEFINE_MUTABLE_META_1_ID="orderDefineMutableMeta1$NONCE"
ORDER_DEFINE_MUTABLE_META_1="$ORDER_DEFINE_MUTABLE_META_1_ID:H|"
assetClient tx orders define -y --from $ACCOUNT_1 --fromID $ACCOUNT_1_NUB_ID \
 --immutableTraits "$ORDER_DEFINE_IMMUTABLE_1" \
 --immutableMetaTraits "$ORDER_DEFINE_IMMUTABLE_META_1" \
 --mutableTraits "$ORDER_DEFINE_MUTABLE_1" \
 --mutableMetaTraits "$ORDER_DEFINE_MUTABLE_META_1"",takerID:I|,exchangeRate:D|,expiry:H|,makerOwnableSplit:D|" $KEYRING $MODE

ORDER_DEFINE_IMMUTABLE_2_ID="orderDefineImmutable2$NONCE"
ORDER_DEFINE_IMMUTABLE_2="$ORDER_DEFINE_IMMUTABLE_2_ID:S|"
ORDER_DEFINE_IMMUTABLE_META_2_ID="orderDefineImmutableMeta2$NONCE"
ORDER_DEFINE_IMMUTABLE_META_2="$ORDER_DEFINE_IMMUTABLE_META_2_ID:I|orderDefineImmutableMeta2$NONCE"
ORDER_DEFINE_MUTABLE_2_ID="orderDefineMutable2$NONCE"
ORDER_DEFINE_MUTABLE_2="$ORDER_DEFINE_MUTABLE_2_ID:D|"
ORDER_DEFINE_MUTABLE_META_2_ID="orderDefineMutableMeta2$NONCE"
ORDER_DEFINE_MUTABLE_META_2="$ORDER_DEFINE_MUTABLE_META_2_ID:H|"
assetClient tx orders define -y --from $ACCOUNT_2 --fromID $ACCOUNT_2_NUB_ID \
 --immutableTraits "$ORDER_DEFINE_IMMUTABLE_2" \
 --immutableMetaTraits "$ORDER_DEFINE_IMMUTABLE_META_2" \
 --mutableTraits "$ORDER_DEFINE_MUTABLE_2" \
 --mutableMetaTraits "$ORDER_DEFINE_MUTABLE_META_2"",takerID:I|,exchangeRate:D|,expiry:H|,makerOwnableSplit:D|" $KEYRING $MODE

sleep $SLEEP
ORDER_DEFINE_CLASSIFICATION_1=$(echo $(assetClient q classifications classifications) | awk -v var="$ORDER_DEFINE_IMMUTABLE_META_1_ID" '{for(i=1;i<=NF;i++)if($i==var)print $(i-10)"."$(i-7)}')
assetClient tx orders make -y --from $ACCOUNT_1 --fromID $ACCOUNT_1_NUB_ID --classificationID $ORDER_DEFINE_CLASSIFICATION_1 --toID $ACCOUNT_1_NUB_ID \
 --makerOwnableID "$ASSET_MINT_1" --makerOwnableSplit "0.000000000000000001" --takerOwnableID stake\
 --immutableProperties "$ORDER_DEFINE_IMMUTABLE_1""stringValue" \
 --immutableMetaProperties "$ORDER_DEFINE_IMMUTABLE_META_1" \
 --mutableProperties "$ORDER_DEFINE_MUTABLE_1""1.01" \
 --mutableMetaProperties "$ORDER_DEFINE_MUTABLE_META_1""123,takerID:I|,exchangeRate:D|1" \
 $KEYRING $MODE

ORDER_DEFINE_CLASSIFICATION_2=$(echo $(assetClient q classifications classifications) | awk -v var="$ORDER_DEFINE_IMMUTABLE_META_2_ID" '{for(i=1;i<=NF;i++)if($i==var)print $(i-10)"."$(i-7)}')
assetClient tx orders make -y --from $ACCOUNT_2 --fromID $ACCOUNT_2_NUB_ID --classificationID $ORDER_DEFINE_CLASSIFICATION_2 --toID $ACCOUNT_2_NUB_ID \
 --makerOwnableID "$ASSET_MINT_2" --makerOwnableSplit "0.000000000000000001" --takerOwnableID "$ASSET_MINT_1"\
 --immutableProperties "$ORDER_DEFINE_IMMUTABLE_2""stringValue" \
 --immutableMetaProperties "$ORDER_DEFINE_IMMUTABLE_META_2" \
 --mutableProperties "$ORDER_DEFINE_MUTABLE_2""1.01" \
 --mutableMetaProperties "$ORDER_DEFINE_MUTABLE_META_2""123,takerID:I|,exchangeRate:D|1" \
 $KEYRING $MODE

sleep $SLEEP
ORDER_MAKE_1_ID=$(echo $(assetClient q orders orders) | awk -v var="$ORDER_DEFINE_IMMUTABLE_META_1_ID" '{for(i=1;i<=NF;i++)if($i==var)print $(i-19)"*"$(i-16)"*"$(i-13)"*"$(i-10)"*"$(i-7)}')
assetClient tx orders cancel -y --from $ACCOUNT_1 --fromID $ACCOUNT_1_NUB_ID --orderID "$ORDER_MAKE_1_ID" $KEYRING $MODE
ORDER_MAKE_2_ID=$(echo $(assetClient q orders orders) | awk -v var="$ORDER_DEFINE_IMMUTABLE_META_2_ID" '{for(i=1;i<=NF;i++)if($i==var)print $(i-19)"*"$(i-16)"*"$(i-13)"*"$(i-10)"*"$(i-7)}')
sleep $SLEEP
assetClient tx orders take -y --from $ACCOUNT_1 --fromID $ACCOUNT_1_NUB_ID --orderID "$ORDER_MAKE_2_ID" --takerOwnableSplit "0.000000000000000001" $KEYRING $MODE

#Maintainers deputize