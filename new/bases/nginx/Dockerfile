FROM registry.cn-beijing.aliyuncs.com/lovego/ubuntu

COPY sources.list /etc/apt/sources.list.d/nginx.list

RUN wget -O- 'http://nginx.org/keys/nginx_signing.key' | apt-key add - \
      && apt-get update && apt-get install -y --no-install-recommends nginx \
      && rm -rf /var/lib/apt/lists/* /etc/nginx/conf.d /var/log/nginx/*.log \
      && mkdir /etc/nginx/sites-enabled /etc/nginx/sites-available

COPY nginx.conf proxy_params /etc/nginx/
COPY nginx-start /usr/local/bin/

STOPSIGNAL SIGQUIT
CMD [ "nginx-start" ]
WORKDIR /var/log/nginx

