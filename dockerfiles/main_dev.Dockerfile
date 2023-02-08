# CubeUniverse.main 程序入口-调试版本的Dockerfile
# /dev-tmp/makeMain.sh

FROM golang:1.18
MAINTAINER tk_sky

COPY . .

ENV GO111MODULE=on \
GOPROXY=https://goproxy.cn \
GOPATH=""

CMD ["dev-tmp/main"]
