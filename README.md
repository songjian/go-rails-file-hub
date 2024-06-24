# Go+Rails远程文件管理

## 项目简介

Go语言编写客户端来监控客户机指定目录的文件变化，并通过Websocket把文件的实时变化信息提交到服务端。

Ruby on Rails开发后端管理，实时查看客户机上指定文件夹中的文件信息，可以远程对文件改名和删除等操作。

## 配置

### 客户端

客户端使用Go语言编写，需要安装Go环境。

#### 客户端目录结构

```shell
tree edge
edge
├── edge.go
├── go.mod
├── go.sum
├── key.txt
└── watch_files

1 directory, 4 files
```

* `edge.go`：客户端源代码
* `go.mod`：Go模块文件
* `go.sum`：Go模块文件
* `key.txt`：客户端key, 用于验证客户端身份
* `watch_files`：监视目录

#### 安装依赖

```shell
go get github.com/fsnotify/fsnotify
go get github.com/gorilla/websocket
```

#### 编译

编译Linux可执行文件：

```shell
go build -o edge
```

编译Windows可执行文件，需要交叉编译：

```shell
GOOS=windows GOARCH=amd64 go build -o edge.exe
```

默认监视目录为`./watch_files`，可以通过修改源代码中`watcher.Add("./watch_files")`更改监视目录

```go
err = watcher.Add("./watch_files") // 换成你要监视的目录
```

#### 客户端key

客户端key在服务端后台Client中新建客户端，获取key后，然后在客户端可执行文件所在目录下创建`key.txt`文件，key值写入文件中。

### 服务端

服务端使用Ruby on Rails开发，需要安装Ruby环境。

#### 服务端目录

服务端代码在`hub`目录下。

#### 启动服务端

安装依赖：

```shell
bundle install
```

启动服务：

```shell
rails s
```

## CHANGES

### 2024-06-21

* 实现监视目录添加和删除文件时文件信息上报到服务器

### 2024-06-24

* 增加远程删除文件功能
* 增加客户端在线状态显示
