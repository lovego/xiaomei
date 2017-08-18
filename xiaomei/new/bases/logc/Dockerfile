FROM 192.168.202.12:5000/xiaomei/base

RUN apt-get update \
  && apt-get install -y --no-install-recommends logrotate \
  && rm -rf /var/lib/apt/lists/*

COPY logc docker-kill /usr/local/bin/

WORKDIR /home/ubuntu

CMD [ "logc", "../logc.yml" ]
