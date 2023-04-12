# CubeUniverse.backend Dockerfile
# CubeUniverse控制后端，存储控制组件，支持前端 由主程序入口调用，请勿直接部署
# docker build -t cubeuniverse-backend -f dockerfiles/backend.Dockerfile .
# docker tag cubeuniverse-backend tksky1/cubeuniverse-backend:0.1alpha
# docker push tksky1/cubeuniverse-backend:0.1alpha
FROM golang:1.20
MAINTAINER tk_sky

COPY . .

ENV GO111MODULE=on \
GOPROXY=https://goproxy.cn \
GOPATH=""

RUN mkdir /app \
    && cd control-backend/main && go mod download \
    && go build -o /app/main control.go routes.go \
    && go clean -modcache && go clean -cache

WORKDIR /app
CMD ["/app/main"]
