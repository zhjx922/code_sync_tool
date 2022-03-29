# 代码同步小工具

> 因为解决不了xxx，所有开发了xxx

## 使用方法

```shell
go build -o code_sync_tool main.go
./code_sync_tool -c conf.ini
```

## 配置说明

```shell
#环境名称
[local_01]
#env 如果需要本地同步到远程，配置为：local，远程热加载可以选择为：server
env = local

#ssh
host = xxx.xxx.com
user = user
password = password
port = 22

#cmd 文件上传或者更新后执行的shell
cmd = ps aux | grep 'Master Process' | grep -v grep | awk '{print $2}' | xargs kill -USR1

#deployment
local_path = /Users/zhaojingxian/dev/php_server
deployment_path = /data/web/code_path

#ignore
ignore_prefix = .git,.idea
ignore_suffix = ~,.swap


[local_02]
#env
env = server

#cmd
cmd = ps aux | grep 'Master Process' | grep -v grep | awk '{print $2}' | xargs kill -USR1
cmd = ps aux | grep 'Master Process' | grep -v grep | awk '{print $2}' | xargs kill -USR2

#deployment
local_path = /data/web/code_path

#ignore
ignore_prefix = .git,.idea
ignore_suffix = ~,.swap
```

## Features

### 全局

* 按照前缀和后缀忽略监听文件

### 本地端模式

* 监听文件变动，自动通过SFTP上传到远程服务器
* 文件上传成功后，在远程服务器执行SHELL命令

### 服务端模式

* 监听文件变动，自动执行SHELL命令