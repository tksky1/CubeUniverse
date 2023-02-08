# CubeUniverse.UniverseBuilder 自动部署程序Dockerfile
# 自动部署CubeUniverse相关组件，由主程序入口调用，请勿直接部署
FROM golang:1.18
MAINTAINER tk_sky

COPY . .

ENV GO111MODULE=on \
GOPROXY=https://goproxy.cn \
GOPATH=""

RUN mkdir /app \
    && cd universeBuilder/main && go mod download \
    && go build -o /app/main main.go

WORKDIR /app
CMD ["/app/main"]
