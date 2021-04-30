#!/bin/bash


KAFKA_VERSION=2.8.0
KAFKA_FOLDER=kafka_2.13-"$KAFKA_VERSION"

if [ -d $KAFKA_FOLDER ]; then
  echo "File exists"
else
  echo "File does not exist, downloading $KAFKA_VERSION"
  wget http://mirrors.estointernet.in/apache/kafka/"$KAFKA_VERSION"/"$KAFKA_FOLDER".tgz
  tar -xzf "$KAFKA_FOLDER".tgz
  rm "$KAFKA_FOLDER".tgz
fi

cd "$KAFKA_FOLDER"

bin/zookeeper-server-start.sh config/zookeeper.properties &

bin/kafka-server-start.sh config/server.properties &

trap 'killall $BGPID; exit' SIGINT
sleep 1024 &
BGPID=${!}
sleep 1024
