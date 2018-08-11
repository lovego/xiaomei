<a name="overview"></a>
## 概览
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
