# xiaomei（小美）
一个简单、实用的go语言web框架，注重报警、日志、部署、尽可能的自动化。

[![Build Status](https://travis-ci.org/lovego/xiaomei.svg?branch=master)](https://travis-ci.org/lovego/xiaomei)
[![Coverage Status](https://coveralls.io/repos/github/lovego/xiaomei/badge.svg?branch=master)](https://coveralls.io/github/lovego/xiaomei?branch=master)
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
- [包含自动报警的应答](./response.md)
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
xiaomei new example          # 生成项目
cd example                   # 进入项目目录（后续命令可在项目内任何目录执行）

xiaomei app run              # 在本机启动app服务器
xiaomei web run              # 在本机启动web服务器

xiaomei deploy               # 部署所有服务到开发环境
xiaomei access -s            # 设置开发环境的Nginx接入层

xiaomei deploy qa            # 部署所有服务到QA环境
xiaomei access -s qa         # 设置QA环境的Nginx接入层

xiaomei deploy production    # 部署所有服务到生产环境
xiaomei access -s production # 设置生产环境的Nginx接入层
```
xiaomei的完整命令行用法可以使用`xiaomei --help`来查看。

<a name="new"></a>
### 生成项目
```shell
xiaomei new example
```
执行以上命令，在当前目录生成了一个名为example的项目。项目目录结构如下：
```
example
├── main.go                        # 入口文件
├── filter                         # 过滤器，所有HTTP请求要先经过过滤，才会进入routes
│   └── filter.go                  # 登录、权限、跨域等全局检查，都可以写在这里
├── routes                         # 路由目录
│   ├── example-api-doc.md         # API文档样例
│   └── routes.go                  # 路由代码，路由要尽量“瘦”，只做参数传递等，不做业务逻辑。
├── helpers                        # routes和filter使用的辅助方法
├── models                         # 模型目录，所有的业务逻辑都应该写在这里
├── tasks                          # 定时任务目录
│   └── tasks.go
├── release                        # 发布目录
│   ├── access                     # 接入层配置
│   │   ├── access.conf.tmpl       # 项目的nginx配置
│   │   └── godoc.conf.tmpl        # godoc服务的nginx配置
│   ├── clusters.yml               # 各个环境的机器配置
│   ├── deploy.yml                 # 各个服务的部署配置
│   ├── img-app                    # 应用服务器镜像
│   │   ├── Dockerfile             # 生成应用服务器镜像的Dockerfile
│   │   ├── config                 # 应用服务器的配置
│   │   │   ├── config.yml         # 框架必需的配置
│   │   │   └── envs               # 各个环境的自定义配置
│   │   │       ├── dev.yml        # 开发环境
│   │   │       ├── test.yml       # 单元测试环境
│   │   │       ├── demo.yml       # 演示环境
│   │   │       ├── qa.yml         # 质量保证环境
│   │   │       └── production.yml # 生产环境
│   │   ├── log                    # 日志目录
│   │   └── views                  # 视图目录
│   │       └── layout             # 布局模板目录
│   │           └── default.tmpl   # 默认布局模版
│   ├── img-logc                   # 日志采集器镜像
│   │   ├── Dockerfile             # 生成日志采集器镜像的Dockerfile
│   │   ├── logc.yml               # 日志采集配置文件
│   │   └── logrotate.conf         # 日志轮转配置文件
│   └── img-web                    # web服务器镜像
│       ├── Dockerfile             # 生成web服务器镜像的Dockerfile
│       ├── public                 # web服务的公开静态文件
│       │   └── index.html          
│       └── web.conf.tmpl          # nginx配置文件
└── vendor                         # 固化第三方依赖包的vendor目录
```
它包含了一个可立即运行、立即部署的"hello world"项目，基于这个基础来增加自己的功能即可。

<a name="run"></a>
### 运行项目
```shell
xiaomei app run
```
在项目内执行以上命令，就可以运行应用服务器，截图如下：

<a name="deploy"></a>
### 基于Docker的部署
xiaomei所有的运行环境都是基于docker的，在开发环境的产出都是docker镜像，然后再将这些镜像部署到其他环境来提供服务。

现在xiaomei包含了这些镜像：
1. app镜像运行项目编译出的二进制文件，用来服务动态内容或者运行定时任务等。
2. web镜像运行nginx，它服务静态文件。
3. logc镜像运行logc工具，收集服务的日志，存储到ElasticSearch，供Kibana可视化展现。
4. godoc镜像运行godoc工具，从golang源码提供文档。
