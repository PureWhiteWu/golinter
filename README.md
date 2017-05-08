# golinter
### 前言

这是一个用go写的代码风格检测的服务器

目的是为了统一接口，让代码风格检测更加方便

欢迎大家能各种提issue和pr，感激不尽

目前支持的语言：

* java
* cpp
* python
* javascript
* 更多语言将被支持

### 使用说明

目前仅支持通过post提交代码并获取返回的结果

格式如下：

```Son
{
  "language": "java",
  "source": "source code"
}
```

服务器收到后如果没有出错，会返回如下格式的一个json：

```json
{
  "error_num": 2, 
  "errors": [
    "error1",
    "error2"
  ] 
}
```

使用命令如下：

```shell
go run server.go dispatch.go
```

端口是48722（不要问我为啥是这个端口）

当然也可以在本地自行编译完打包成二进制运行。
