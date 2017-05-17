FROM 192.168.202.12:5000/xiaomei/base

COPY sources.list /etc/apt/sources.list.d/nginx.list

RUN wget -O- 'http://nginx.org/keys/nginx_signing.key' | apt-key add - \
      && apt-get update && apt-get install -y --no-install-recommends nginx \
      && rm -rf /var/lib/apt/lists/* /etc/nginx/conf.d \
      && mkdir /etc/nginx/sites-enabled /etc/nginx/sites-available

COPY nginx.conf proxy_params /etc/nginx/
COPY nginx-start /usr/local/bin/

CMD [ "nginx-start" ]
EXPOSE 80 443

