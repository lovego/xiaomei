# xiaomei（小美）
一个简单、实用的go语言web框架，注重部署、日志、报警、尽可能的自动化。

[![Build Status](https://travis-ci.org/lovego/xiaomei.svg?branch=master)](https://travis-ci.org/lovego/xiaomei)
[![Coverage Status](https://coveralls.io/repos/github/lovego/xiaomei/badge.svg?branch=master)](https://coveralls.io/github/lovego/xiaomei?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/lovego/xiaomei)](https://goreportcard.com/report/github.com/lovego/xiaomei)
[![GoDoc](https://godoc.org/github.com/lovego/xiaomei?status.svg)](https://godoc.org/github.com/lovego/xiaomei)

## 安装
```shell
go get github.com/lovego/xiaomei/xiaomei
```
执行如上`go get`命令即可将xiaomei安装到`$GOPATH/bin`目录中。
如果`$GOPATH/bin`已经在`$PATH`搜索路径中，你可以输入`xiaomei version`命令来检查xiaomei是否已经安装成功。如果输出类似"xiaomei version 18.7.13"的版本信息，就说明已经安装成功了。

现在xiaomei包含了三个服务：
1. app服务：运行项目编译出的二进制文件，用来服务HTTP请求、执行定时任务等。
2. web服务：运行nginx，用来服务项目内的静态文件。
3. logc服务：运行logc日志收集工具，存储到ElasticSearch，供Kibana可视化展现。

命令相关文档
- [命令概览](./xiaomei)
- [生成项目](./xiaomei/new)
- [运行项目](./xiaomei/run)
- [基于Docker的部署](./xiaomei/deploy)
- [基于Nginx的接入层（负载均衡）](./xiaomei/access)

代码相关文档
- [最简单的过滤器](./server/filter.md)
- [Express风格的路由](./router)
- [请求](./request.md)
- [应答（自动报警）](./response.md)
- [会话](./session)
- [模版渲染](./renderer)
- [统一的配置](./config)
- [常见数据库连接](./config/db)
- [HTTP服务器](./server)


