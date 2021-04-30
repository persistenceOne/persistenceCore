#!/bin/bash

KAFKA_VERSION=2.8.0
KAFKA_FOLDER=kafka_2.13-"$KAFKA_VERSION"

cd "$KAFKA_FOLDER"

bin/zookeeper-server-start.sh config/zookeeper.properties &

bin/kafka-server-start.sh config/server.properties &

sleep 2s
declare -a topics=$(bin/kafka-topics.sh --list --zookeeper localhost:2181)
echo "====================================================="

for topic in "${topics[@]}"; do
  echo "$topic"
  bin/kafka-topics.sh --delete --zookeeper localhost:2181 --topic "$topic"
done
echo "marked for deletion"

bin/kafka-server-stop.sh
sleep 2s

bin/zookeeper-server-stop.sh

cd ..
