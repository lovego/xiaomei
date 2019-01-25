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

