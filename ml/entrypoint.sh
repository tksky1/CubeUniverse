#!/bin/bash
ERR_MSG="E: If either ZOOKEEPER_HOST or ZOOKEEPER_PORT is defined, both must be defined"
if [ -n "${ZOOKEEPER_HOST}" ] || [ -n "${ZOOKEEPER_PORT}" ]; then
  if [ -z "${ZOOKEEPER_HOST}" ]; then
    echo "E: ZOOKEEPER_HOST is not defined"
    echo "${ERR_MSG}"
    exit 1
  elif [ -z "${ZOOKEEPER_PORT}" ]; then
    echo "E: ZOOKEEPER_PORT is not defined"
    echo "${ERR_MSG}"
    exit 1
  fi
  exec /opt/kafka/bin/zookeeper-server-start.sh /etc/kafka/zookeeper.properties &
  exec /opt/kafka/bin/kafka-server-start.sh /opt/kafka/config/server.properties 
else
  exec /opt/kafka/bin/kafka-server-start.sh "${@}"
fi

# #!/bin/bash
# ERR_MSG="E: If either ZOOKEEPER_HOST or ZOOKEEPER_PORT is defined, both must be defined"
# if [ -n "${ZOOKEEPER_HOST}" ] || [ -n "${ZOOKEEPER_PORT}" ]; then
#   if [ -z "${ZOOKEEPER_HOST}" ]; then
#     echo "E: ZOOKEEPER_HOST is not defined"
#     echo "${ERR_MSG}"
#     exit 1
#   elif [ -z "${ZOOKEEPER_PORT}" ]; then
#     echo "E: ZOOKEEPER_PORT is not defined"
#     echo "${ERR_MSG}"
#     exit 1
#   fi
#   exec /opt/kafka/bin/zookeeper-server-start.sh /etc/kafka/zookeeper.properties
#   exec /opt/kafka/bin/kafka-server-start.sh "${@}" --override zookeeper.connect="${ZOOKEEPER_HOST}":"${ZOOKEEPER_PORT}"
# else
#   exec /opt/kafka/bin/kafka-server-start.sh "${@}"
# fi