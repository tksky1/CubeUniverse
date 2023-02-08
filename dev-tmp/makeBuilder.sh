# 调试用脚本
cd ../universeBuilder
go build -o ../dev-tmp/main main/universeBuilder.go
cd ..
docker build -t builder-dev -f dockerfiles/universeBuilder_dev.Dockerfile .
docker save builder-dev -o ./dev-tmp/builder-dev.tar
scp ./dev-tmp/builder-dev.tar 192.168.79.12:/home/node
scp ./dev-tmp/builder-dev.tar 192.168.79.13:/home/node