# xiaomei（小美）
一个简单、实用的go语言web框架，注重日志、报警、部署、尽可能的自动化。

[![Go Report Card](https://goreportcard.com/badge/github.com/lovego/xiaomei)](https://goreportcard.com/report/github.com/lovego/xiaomei)
[![GoDoc](https://godoc.org/github.com/lovego/xiaomei?status.svg)](https://godoc.org/github.com/lovego/xiaomei)

命令相关文档
- [安装](#install)
- [概览](#overview)
- [生成项目](#new)
- [运行项目](#run)
- [基于Docker的部署](#deploy)
- [基于Nginx的接入层（负载均衡）](#access)
- [JSON格式的日志记录、实时收集](./server/log.md)

代码相关文档
- [最简单的过滤器](./server/filter.md)
- [Express风格的路由](./router)
- [请求](./request.md)
- [包含自动邮件报警的应答](./response.md)
- [会话](./session)
- [模版渲染](./renderer)
- [统一的配置](./config)
- [常见数据库连接](./config/db)

<a name="install"></a>
### 安装
```shell
go get github.com/lovego/xiaomei/xiaomei
```
执行如上`go get`命令即可将xiaomei安装到`$GOPATH/bin`目录中。
如果`$GOPATH/bin`已经在`$PATH`搜索路径中，你可以输入`xiaomei version`命令来检查xiaomei是否已经安装成功。如果输出类似"xiaomei version 18.7.13"的版本信息，就说明已经安装成功了。

<a name="overview"></a>
### 概览
在项目的开发流程中常用的命令如下：
```shell
xiaomei new example       # 生成项目
cd example                # 进入项目目录（后续命令可在项目内任何目录执行）
xiaomei app run           # 启动app服务器
xiaomei web run           # 启动web服务器
xiaomei deploy            # 部署到开发环境
xiaomei deploy qa         # 部署到QA环境
xiaomei deploy production # 部署到生产环境
```
xiaomei的完整命令行用法可以使用`xiaomei --help`来查看。

<a name="new"></a>
### 生成项目
xiaomei包含一个默认的项目模板，使用xiaomei new来生成。它包含了一个可立即运行、立即部署的"hello world"项目，基于这个基础来增加自己的功能即可。

<a name="deploy"></a>
### 基于Docker的部署
xiaomei所有的运行环境都是基于docker的，在开发环境的产出都是docker镜像，然后再将这些镜像部署到其他环境来提供服务。

现在xiaomei包含了这些镜像：
1. app镜像运行项目编译出的二进制文件，用来服务动态内容或者运行定时任务等。
2. web镜像运行nginx，它服务静态文件。
3. logc镜像运行logc工具，收集服务的日志，存储到ElasticSearch，供Kibana可视化展现。
4. godoc镜像运行godoc工具，从golang源码提供文档。
