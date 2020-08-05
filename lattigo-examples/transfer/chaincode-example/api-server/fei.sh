#!/bin/bash

# 构建
go get github.com/go-kit/kit@v0.8.0
#bee run

# 打包
bee pack
# 解压到/tmp/api-server-bank1
tar -xvzf api-server.tar.gz -C /tmp/api-server-bank1

# 启动api-server
cd /tmp/api-server-bank1
./api-server&

# 解压到/tmp/api-server-bank2
