# 调试用脚本
cd ../main
go build -o ../dev-tmp/main main.go
cd ..
docker build -t main-dev -f dockerfiles/main_dev.Dockerfile .
docker save main-dev -o ./dev-tmp/main-dev.tar
scp ./dev-tmp/main-dev.tar 192.168.79.12:/home/node
scp ./dev-tmp/main-dev.tar 192.168.79.13:/home/node