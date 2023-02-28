# CubeUniverse.UniverseBuilder 自动部署程序Dockerfile 调试版本
# 自动部署CubeUniverse相关组件，由主程序入口调用，请勿直接部署
# /dev-tmp/makeBuilder.sh
FROM golang:1.18
MAINTAINER tk_sky

COPY . .

CMD ["dev-tmp/main"]
