# xiaomei（小美）
一个简单、实用的go语言web框架，注重部署、日志、报警、尽可能的自动化。

[![Build Status](https://travis-ci.org/lovego/xiaomei.svg?branch=master)](https://travis-ci.org/lovego/xiaomei)
[![Coverage Status](https://img.shields.io/coveralls/github/lovego/xiaomei/master.svg)](https://coveralls.io/github/lovego/xiaomei?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/lovego/xiaomei)](https://goreportcard.com/report/github.com/lovego/xiaomei)
[![GoDoc](https://godoc.org/github.com/lovego/xiaomei?status.svg)](https://godoc.org/github.com/lovego/xiaomei)

## 安装
```shell
go get github.com/lovego/xiaomei
```
执行如上`go get`命令即可将xiaomei安装到`$GOPATH/bin`目录中。
如果`$GOPATH/bin`已经在`$PATH`搜索路径中，你可以输入`xiaomei version`命令来检查xiaomei是否已经安装成功。如果输出类似"xiaomei version 18.7.13"的版本信息，就说明已经安装成功了。

现在xiaomei包含了三个服务：
1. app服务：运行项目编译出的二进制文件，用来服务HTTP请求、执行定时任务等。
2. web服务：运行nginx，用来服务项目内的静态文件。
3. logc服务：运行logc日志收集工具，存储到ElasticSearch，供Kibana可视化展现。

## 文档
- [命令概览](#overview)
- [生成项目](./new)
- [运行项目](./services/run)
- [基于Docker的部署](./services/deploy)
- [基于Nginx的接入层（负载均衡）](./access)

<a name="overview"></a>
## 概览

在项目的开发流程中常用的命令如下：
```shell
xiaomei new example          # 生成项目
cd example                   # 进入项目目录（后续命令可在项目内任何目录执行）

xiaomei app run              # 在本机启动app服务器
xiaomei web run              # 在本机启动web服务器

xiaomei deploy                  # 部署所有服务到开发环境
xiaomei access setup            # 设置开发环境的Nginx接入层

xiaomei deploy qa               # 部署所有服务到QA环境
xiaomei access setup qa         # 设置QA环境的Nginx接入层

xiaomei deploy production       # 部署所有服务到生产环境
xiaomei access setup production # 设置生产环境的Nginx接入层
```

xiaomei的完整命令行用法可以使用`xiaomei --help`来查看。

