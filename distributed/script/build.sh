#!/bin/bash

# 本脚本在项目distributed目录下执行
# 需要go环境， go版本在go.mod中定义

go build -o ./build/interserver ./interface/main.go && echo -e "\033[33m Compile interface server successful ! \033[0m"
go build -o ./build/dataserver ./data/main.go && echo -e "\033[33m Compile data server successful ! \033[0m"

if ! [ -d ./build/config ]
then
  mkdir ./build/config
fi
cp ./config/data.yml ./build/config/
cp ./config/interface.yml ./build/config/

if ! [ -d ./build/var ]
then
  mkdir ./build/var
fi