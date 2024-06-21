# Go+Rails远程文件管理

## 项目简介

Go语言编写客户端来监控客户机指定目录的文件变化，并通过Websocket把文件的实时变化信息提交到服务端。

Ruby on Rails开发后端管理，实时查看客户机上指定文件夹中的文件信息，可以远程对文件改名和删除等操作。

## 配置

### 服务端

启动服务端

```shell
cd hub
bundle install
rails s
```

### 客户端

编译客户端

```shell
cd edge
go build -o edge
```

交叉编译

```shell
# 编译Windows可执行文件
GOOS=windows GOARCH=amd64 go build -o edge.exe
```

默认监视目录为`./watch_files`，可以通过修改源代码中`watcher.Add("./watch_files")`更改监视目录

```go
err = watcher.Add("./watch_files") // 换成你要监视的目录
```

设置客户端key

在服务端后台新建客户端，获取key，然后在客户端可执行文件所在目录下创建`key.txt`文件，文件内容为key值

## CHANGES

- 2024-06-21: 实现监视目录添加和删除文件时文件信息上报到服务器
