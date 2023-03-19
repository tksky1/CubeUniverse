# CubeUniverse开发用调试版本，不要使用
FROM golang:1.18
MAINTAINER tk_sky
ENV GO111MODULE=on \
GOPROXY=https://goproxy.cn \
GOPATH=""
RUN rm -f dev-tmp/*.tar
COPY . .

CMD ["dev-tmp/main"]
