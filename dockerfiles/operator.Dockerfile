# CubeUniverse.operator Dockerfile
# 监控调度CubeUniverse集群，维持集群健康
FROM golang:1.18
MAINTAINER tk_sky

COPY . .

ENV GO111MODULE=on \
GOPROXY=https://goproxy.cn \
GOPATH=""

RUN mkdir /app \
    && cd universeOperator/main && go mod download \
    && go build -o /app/main main.go

WORKDIR /app
CMD ["/app/main"]
