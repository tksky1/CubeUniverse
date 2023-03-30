#!/bin/sh

/data/kafka/bin/kafka-storage.sh format \
                    --config /data/kafka/config/kraft/server.properties \
                    --cluster-id $(/data/kafka/bin/kafka-storage.sh random-uuid)

/data/kafka/bin/kafka-server-start.sh /data/kafka/config/kraft/server.properties