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

