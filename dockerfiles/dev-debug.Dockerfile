# CubeUniverse开发用调试版本，不要使用
FROM golang:1.18
MAINTAINER tk_sky
RUN rm -f dev-tmp/builder-dev.tar
COPY . .

CMD ["dev-tmp/main"]
