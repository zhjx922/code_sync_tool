# 代码同步小工具

> 因为解决不了xxx，所有开发了xxx

## 使用方法

```shell
go build -o code_sync_tool main.go
./code_sync_tool -c conf.ini
```

## Features

### 全局

* 按照前缀和后缀忽略监听文件

### 本地端模式

* 监听文件变动，自动通过SFTP上传到远程服务器
* 文件上传成功后，在远程服务器执行SHELL命令

### 服务端模式

* 监听文件变动，自动执行SHELL命令