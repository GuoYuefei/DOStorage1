#!/bin/bash

# 本脚本在项目根目录下执行
# 需要go环境， go版本在go.mod中定义

go build -o ./build/interserver ./interface/main.go
go build -o ./build/dataserver ./data/main.go

mkdir ./build/config
cp ./config/data.yml ./build/config/
cp ./config/interface.yml ./build/config/