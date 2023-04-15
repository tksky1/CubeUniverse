# CubeUniverse.operator Dockerfile
# 监控调度CubeUniverse集群，维持集群健康，由主程序入口调用，请勿直接部署
# docker build -t cubeuniverse-operator -f dockerfiles/operator.Dockerfile .
# docker tag cubeuniverse-operator tksky1/cubeuniverse-operator:0.1alpha
# docker push tksky1/cubeuniverse-operator:0.1alpha
FROM golang:1.20
MAINTAINER tk_sky

RUN cd / && mkdir CubeUniverse
COPY . /CubeUniverse

ENV GO111MODULE=on \
GOPROXY=https://goproxy.cn \
GOPATH=""

RUN cd /CubeUniverse && mkdir app/ \
    && cd universeOperator/main && go mod download \
    && go build -o /CubeUniverse/app/main mainCubeOperator.go routes.go \
    && go clean -modcache && go clean -cache

WORKDIR /CubeUniverse/app
CMD ["/CubeUniverse/app/main"]
