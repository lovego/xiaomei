FROM 192.168.202.12:5000/xiaomei/baseimage

RUN apt-get update \
  && apt-get install -y --no-install-recommends nginx \
  && rm -rf /var/lib/apt/lists/*

COPY nginx.conf /etc/nginx/nginx.conf

CMD ["nginx"]
EXPOSE 80 443

WORKDIR /var/www

