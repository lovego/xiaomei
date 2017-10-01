# xiaomei 小而美的go语言web框架。

[![GoDoc](https://godoc.org/github.com/lovego/xiaomei?status.svg)](https://godoc.org/github.com/lovego/xiaomei)

## 安装
```
  go get github.com/lovego/xiaomei/xiaomei
```

## 使用
```
  xiaomei new example       # 生成项目
  xiaomei app run           # 启动app服务器
  xiaomei web run           # 启动web服务器
  xiaomei deploy            # 部署到开发环境
  xiaomei deploy qa         # 部署到QA环境
  xiaomei deploy production # 部署到生产环境
  xiaomei --help            # 完整的xiaomei命令文档
```

## 介绍
  xiaomei包含两个部分：1. app服务器，2. 基于docker的开发、部署工具。

### 一、app服务器

1. Router 支持基于字符串和正则表达式的路由（express风格）。

2. Renderer 支持layout和partial的模板渲染（Rails风格）。

3. Request、Response 它们封装了http.Request、http.ResponseWriter、以及Renderer，以提供模板渲染等功能。

### 二、基于docker的开发、部署工具

xiaomei包含一个名为xiaomei工具，用来支持开发、部署、运维。

xiaomei包含一个默认的项目模板，使用xiaomei new来生成。它包含了一个可立即运行、立即部署的"hello world"项目，基于这个基础来增加自己的功能即可。

xiaomei所有的运行环境都是基于docker的，在开发环境的产出都是docker镜像，然后再将这些镜像部署到其他环境来提供服务。

现在xiaomei包含了这些镜像：
1. app镜像运行项目编译出的二进制文件，用来服务动态内容或者运行定时任务等。
2. web镜像运行nginx，它服务静态文件。
3. godoc镜像运行godoc工具，从golang源码提供文档。
4. access镜像运行nginx，根据域名将请求转发到不同项目的服务。
5. logc镜像运行logc工具，收集服务的日志，存储到ElasticSearch，供可视化展现。

