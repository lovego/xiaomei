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
- [概览](./xiaomei)
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
- [JSON化的日志记录、实时收集](#logging)

<a name="new"></a>
## 生成项目
```shell
MacBook:~/go/src/example$ xiaomei new example
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
## 运行项目
```shell
MacBook:~/go/src/example$ xiaomei app exec
2018/08/09 09:16:58 compile the app server binary.
2018/08/09 09:17:00 check app code spec.
2018/08/09 09:17:00 starting.(dev)
2018/08/09 09:17:00 started.(:3000)
```
在项目内执行以上命令，就可以运行应用服务器，该命令编译项目为可执行文件，然后执行它。
如果本机已经安装好了Docker，也可以执行如下命令来运行项目：

```shell
MacBook:~/go/src/example$ xiaomei app run
2018/08/09 09:11:15 compile the app server binary.
2018/08/09 09:11:17 check app code spec.
2018/08/09 09:11:17 building app image.
docker build --pull --file=Dockerfile --tag=registry.example.com/example/app .
... # 此处省略若干构建镜像过程中的输出
Successfully built 3f8020e2c72a
Successfully tagged registry.example.com/example/app:latest
2018/08/09 09:11:25 starting.(dev)
2018/08/09 09:11:25 started.(:3001)
```
该命令编译项目为可执行文件，然后构建并运行Docker镜像。因为部署时也使用Docker镜像的方式，
因此`run`比`exec`更加接近真实的部署环境，但是因为每次都要构建镜像，所以比`exec`慢一些。`run`跟部署一样，需要注意[非Linux环境宿主机访问限制](#host-network)。

<a name="deploy"></a>
## 基于Docker的部署
xiaomei的部署是基于docker的，首先基于项目构建docker镜像，再将镜像推送到registry，
然后在目标部署机器上拉取并运行这些镜像。
release/clusters.yml是部署机器的配置文件，release/deploy.yml每台机器上部署的服务的配置文件。

```shell
MacBook:~/go/src/example$ xiaomei app deploy
2018/08/09 14:27:50 time tag: 180809-142750
2018/08/09 14:27:50 compile the app server binary.
2018/08/09 14:27:50 check app code spec.
2018/08/09 14:27:50 building app image.
docker build --pull --file=Dockerfile --tag=registry.example.com/example/app:dev-180809-142750 .
... # 此处省略若干构建镜像过程中的输出
Successfully built 3f8020e2c72a
Successfully tagged registry.example.com/example/app:dev-180809-142750
2018/08/09 14:27:57 pushing app image.
docker push registry.hztl3.com/example/app:dev-180809-142750
... # 此处省略若干推送镜像过程中的输出
dev-180809-142750: digest: sha256:a1fd5e7a60d529883a15c46adefa1a0aa1de4ba901d4f47b85315e4b9ab1e1c7 size: 2819
2018/08/09 14:27:59 deploying ubuntu@127.0.0.1
example-dev-app.3001
2018/08/09 14:27:59 starting.(dev)
2018/08/09 14:27:59 started.(:3001)
example-dev-app.4001
2018/08/09 14:28:00 starting.(dev)
2018/08/09 14:28:00 started.(:4001)
```
如上是部署到开发环境的命令及输出。

### <a name="host-network">非Linux环境宿主机访问限制</a>
需要注意，在Linux环境下使用了`--network=host`的Docker网络模式，因此可以在容器内直接使用localhost或127.0.0.1来访问宿主机。在非Linux环境下，Docker不支持`--network=host`，所以则不能使用localhost或127.0.0.1来访问宿主机,，而需要使用宿主机的其他IP来访问。

<a name="access"></a>
## 基于Nginx的接入层（负载均衡）
xiaomei使用Nginx作为接入层和负载均衡。release/access.conf.tmpl是Nginx配置的模板文件，可以根据你的需求修改。xiaomei根据该模板生成Nginx配置。 使用如下命令可查看生成的nginx配置。
```
MacBook:~/go/src/example$ xiaomei access
upstream example-dev-app {
  server 127.0.0.1:3001;
  server 127.0.0.1:4001;
  keepalive 1;
}
upstream example-dev-web {
  server 127.0.0.1:8001 fail_timeout=3m;
  keepalive 1;
}
server {
  listen 80;
  server_name example.dev.example.com;
  
  location = / {
    proxy_pass http://example-dev-web;
  }
  location ~ \.(html|js|css|png|gif|jpg|svg|ico|woff|woff2|ttf|eot|map|json)$ {
    proxy_pass http://example-dev-web;
  }
  location / {
    proxy_pass http://example-dev-app;
  }

  proxy_http_version 1.1;
  proxy_set_header Connection "";
  proxy_set_header Host $http_host;
  proxy_set_header X-Real-IP $remote_addr;
  proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
  proxy_set_header X-Forwarded-Proto $scheme;
  proxy_connect_timeout 3s;

  access_log /var/log/nginx/example.dev.example.com/access.log;
  error_log  /var/log/nginx/example.dev.example.com/access.err;
}
```
使用`xiaomei access -s`命令则可以将Nginx配置写到接入层机器的`/etc/nginx/sites-enabled/<domain>`这个文件内，并且重新加载Nginx配置。
1. release/clusters.yml文件中`labels.access`为`true`的机器就是需要配置接入层的机器。
2. 重新加载Nginx配置通过执行 `sudo systemctl reload nginx` 或 `sudo service nginx reload` 命令来完成，因此`xiaomei access -s`是只支持Linux系统的。
3. Ubuntu的Nginx的主配置文件默认包含`include /etc/nginx/sites-enabled/*;`这条配置，所以/etc/nginx/sites-enabled目录下的所有配置文件都会生效。其他Linux发行版，如果没有这条配置，需要自行添加。
4. 其中`<domain>`代表项目的域名，在release/img-app/config/config.yml配置文件中配置。

```
ubuntu@ubuntu:~/go/src/example$ xiaomei access -s
2018/08/11 17:41:00 ubuntu@127.0.0.1
nginx: the configuration file /etc/nginx/nginx.conf syntax is ok
nginx: configuration file /etc/nginx/nginx.conf test is successful
```

<a name="logging"></a>
## JSON格式的日志记录、实时收集
```
{
  "agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.99 Safari/537.36",
  "at": "2018-08-11T18:11:36.052356321+08:00",
  "duration": 0.067402,
  "host": "localhost:3000",
  "ip": "::1",
  "level": "info",
  "machineName": "MacBook",
  "method": "GET",
  "path": "/",
  "query": {
    "name": [
      "value"
    ]
  },
  "rawQuery": "name=value",
  "refer": "",
  "reqBodySize": 0,
  "resBodySize": 23,
  "status": 200,
  "tags": {
    "hello": "world"
  }
}
```
