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
docker push registry.example.com/example/app:dev-180809-142750
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

