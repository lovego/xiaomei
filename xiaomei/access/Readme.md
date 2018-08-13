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

