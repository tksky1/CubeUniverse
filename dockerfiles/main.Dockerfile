# CubeUniverse.main 程序入口Dockerfile 请勿直接部署
# docker build -t cubeuniverse-prepare -f dockerfiles/main.Dockerfile .
# docker tag cubeuniverse-prepare tksky1/cubeuniverse:0.1alpha
# docker push tksky1/cubeuniverse:0.1alpha
FROM golang:1.20
MAINTAINER tk_sky

COPY . .

ENV GO111MODULE=on \
GOPROXY=https://goproxy.cn \
GOPATH=""

RUN mkdir /app \
    && cd main && go mod download \
    && go build -o /app/main main.go  \
    && go clean -modcache && go clean -cache

WORKDIR /app
CMD ["/app/main"]
