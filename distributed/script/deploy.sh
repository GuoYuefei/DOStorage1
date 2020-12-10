#!/bin/bash

## 本部署脚本不会一并部署rabbitmq、es等
## 仅部署本项目
## 且部署时在部署机子上重新构建，所以需要go、git环境  // 使用这种方法也是因为架构，系统文件版本不同会导致二进制文件无法运行

## 确保在用户目录下有这么几个文件夹
ssh ubuntu@119.29.5.95 "cd ~ ; mkdir -p dos/build/var ; mkdir -p ~/DOStorage1/"

### 复制内容
scp -r ./** ubuntu@119.29.5.95:~/DOStorage1/

### 登陆后执行  eeof 可以自定义, 下面指令遇到eeof即停
###
ssh ubuntu@119.29.5.95 << eeof

cd DOStorage1/distributed
rm -rf build && echo "delete build folder, then will build"
bash script/build.sh
echo "如果有，则关闭之前运行的程序"
cd ~/dos/ && bash stop_min.sh
cd -
cp -r build ~/dos/
cp script/start_min.sh script/stop_min.sh ~/dos/
cd ~
rm -rf DOStorage1
cd dos

echo "开启程序"
bash start_min.sh

exit

eeof
echo "done!"

exit 0


