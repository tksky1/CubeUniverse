# CubeUniverse.UniverseBuilder 自动部署程序Dockerfile
# 自动部署CubeUniverse相关组件，由主程序入口调用，请勿直接部署
# docker build -t cubeuniverse-builder -f dockerfiles/universeBuilder.Dockerfile .
# docker tag cubeuniverse-builder tksky1/cubeuniverse-builder:0.1alpha
# docker push tksky1/cubeuniverse-builder:0.1alpha
FROM golang:1.20
MAINTAINER tk_sky

RUN cd / && mkdir CubeUniverse
COPY . /CubeUniverse

ENV GO111MODULE=on \
GOPROXY=https://goproxy.cn \
GOPATH=""

RUN cd /CubeUniverse && mkdir app/ \
    && cd universeBuilder/main && go mod download \
    && go build -o /CubeUniverse/app/main universeBuilder.go \
    && go clean -modcache && go clean -cache

WORKDIR /CubeUniverse/app
CMD ["/CubeUniverse/app/main"]
