# xiaomei（小美）
一个简单、实用的go语言web框架，注重日志、报警、部署、尽量自动化。

[![Go Report Card](https://goreportcard.com/badge/github.com/lovego/xiaomei)](https://goreportcard.com/report/github.com/lovego/xiaomei)
[![GoDoc](https://godoc.org/github.com/lovego/xiaomei?status.svg)](https://godoc.org/github.com/lovego/xiaomei)

快速入门
- [概览](#overview)
- [安装](#install)
- [生成项目](#new)
- [运行项目](#run)
- [部署项目](#deploy)

参考文档
- [过滤器](./server/filter.md)
- [路由](./router)
- [请求](./request.md)
- [应答](./response.md)
- [会话](./session)
- [模版渲染](./renderer)
- [配置](./config)
- [数据库连接](./config/db)
- [日志](./server/log.md)
- [报警](./alarm.md)

### 概览
开发流程中常用的命令如下：
```shell
xiaomei new example       # 生成项目
cd example                # 进入项目目录
xiaomei app run           # 启动app服务器
xiaomei web run           # 启动web服务器
xiaomei deploy            # 部署到开发环境
xiaomei deploy qa         # 部署到QA环境
xiaomei deploy production # 部署到生产环境
xiaomei --help            # 完整的xiaomei命令文档
```

### 安装
```shell
go get github.com/lovego/xiaomei/xiaomei
```
执行如上`go get`命令即可将xiaomei安装到`$GOPATH/bin`目录中。
如果`$GOPATH/bin`已经在`$PATH`搜索路径中，你可以输入`xiaomei version`命令来检查xiaomei是否已经安装成功。如果输出类似"xiaomei version 18.7.13"的版本信息，就说明已经安装成功了。


## 介绍
xiaomei包含两个部分：1. app服务器，2. 基于docker的开发、部署工具。

### 一、app服务器

1. Router 支持基于字符串和正则表达式的路由（Express风格）。

2. Renderer 支持layout和partial的模板渲染（Rails风格）。

3. Request、Response 它们封装了http.Request、http.ResponseWriter、以及Renderer，以提供模板渲染等功能。

4. Session 支持会话读写，且内置了基于加密cookie的会话。

5. config  统一的配置结构和解析。

### 二、基于docker的开发、部署工具

xiaomei包含一个名为xiaomei工具，用来支持开发、部署、运维。

xiaomei包含一个默认的项目模板，使用xiaomei new来生成。它包含了一个可立即运行、立即部署的"hello world"项目，基于这个基础来增加自己的功能即可。

xiaomei所有的运行环境都是基于docker的，在开发环境的产出都是docker镜像，然后再将这些镜像部署到其他环境来提供服务。

现在xiaomei包含了这些镜像：
1. app镜像运行项目编译出的二进制文件，用来服务动态内容或者运行定时任务等。
2. web镜像运行nginx，它服务静态文件。
3. logc镜像运行logc工具，收集服务的日志，存储到ElasticSearch，供Kibana可视化展现。
4. godoc镜像运行godoc工具，从golang源码提供文档。

