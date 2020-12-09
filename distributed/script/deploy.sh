#!/bin/bash

## 本部署脚本不会一并部署rabbitmq、es等
## 仅部署本项目

## 确保在用户目录下有这么几个文件夹
ssh ubuntu@119.29.5.95 "cd ~ ; mkdir -p dos/build ; mkdir -p dos/var"

### 复制运行时有用内容
scp -r distributed/build/config ubuntu@119.29.5.95:~/dos/build/
scp distributed/build/dataserver distributed/build/interserver ubuntu@119.29.5.95:~/dos/build/
scp distributed/script/start_min.sh distributed/script/stop_min.sh ubuntu@119.29.5.95:~/dos/

### 登陆后执行  eeof 可以自定义, 下面指令遇到eeof即停
###
ssh ubuntu@119.29.5.95 << eeof

cd dos
bash stop_min.sh
bash start_min.sh
exit

eeof
echo "done!"


