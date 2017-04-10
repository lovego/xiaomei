FROM 192.168.202.12:5000/xiaomei/baseimage

RUN apt-get update \
  && apt-get install -y --no-install-recommends nginx \
  && rm -rf /var/lib/apt/lists/* \
    /etc/nginx/sites-enabled/default \
    /etc/nginx/sites-available/default

COPY nginx.conf /etc/nginx/nginx.conf
# COPY nginx-start /usr/local/bin/

CMD ["nginx"]
EXPOSE 80 443

