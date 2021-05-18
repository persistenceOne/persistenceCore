./.setup.sh

persistenceCore init test --chain-id test
persistenceCore keys add test --recover<<<"y
wage thunder live sense resemble foil apple course spin horse glass mansion midnight laundry acoustic rhythm loan scale talent push green direct brick please"
persistenceCore add-genesis-account test 100000000000000stake
persistenceCore gentx test 10000000stake --chain-id test
persistenceCore collect-gentxs

persistenceCore pStake chain.json "" --ports "localhost:9092"

bin/zookeeper-server-start.sh config/zookeeper.properties

bin/kafka-server-start.sh config/server.properties


persistenceCore pStake chain.json "wage thunder live sense resemble foil apple course spin horse glass mansion midnight laundry acoustic rhythm loan scale talent push green direct brick please" --ports="localhost:9092" --ethSleepTime 2000 --tmSleepTime  2000 --tmStart 1 --ethStart 4772131 --ethPrivateKey ce0b4f52b909ed065181c5295632e398018bbe9d8be9aaab30ca0831dfef905c