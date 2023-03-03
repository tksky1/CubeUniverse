# CubeUniverse开发用调试版本，不要使用
FROM golang:1.18
MAINTAINER tk_sky

COPY . .

CMD ["dev-tmp/main"]
