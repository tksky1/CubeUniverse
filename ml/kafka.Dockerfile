FROM ubuntu/kafka:latest
MAINTAINER ahz-r3v
COPY entrypoint.sh /usr/local/bin/
COPY server.properties /opt/kafka/config/
RUN chmod +x /opt/kafka/bin/zookeeper-server-start.sh
RUN chmod +x /usr/local/bin/entrypoint.sh
ENV ZOOKEEPER_HOST=127.0.0.1 \
TZ=UTC

EXPOSE 9092
EXPOSE 2181

# CMD ["/bin/bash", "zookeeper-server-start.sh","/etc/kafka/zookeeper.properties", "&&", "/bin/bash", "kafka-server-start.sh", "/etc/kafka/server.properties", "--override","zookeeper.connect=127.0.0.1:2181"]
# CMD ["zookeeper-server-start.sh", "/etc/kafka/zookeeper.properties"]
# ENTRYPOINT [ "zookeeper-server-start.sh" ]