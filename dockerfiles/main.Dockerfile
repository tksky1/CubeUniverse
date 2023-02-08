# CubeUniverse.main 程序入口Dockerfile
# docker build -t cubeuniverse-prepare -f dockerfiles/main.Dockerfile .
# docker tag cubeuniverse-prepare tksky1/cubeuniverse:dev0.1
# docker push tksky1/cubeuniverse:dev0.1
FROM golang:1.18
MAINTAINER tk_skyc

COPY . .

ENV GO111MODULE=on \
GOPROXY=https://goproxy.cn \
GOPATH=""

RUN mkdir /app \
    && cd main && go mod download \
    && go build -o /app/main main.go

WORKDIR /app
CMD ["/app/main"]
