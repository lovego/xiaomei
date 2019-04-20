FROM registry.cn-beijing.aliyuncs.com/lovego/ubuntu

RUN apt-get update \
  && apt-get install -y --no-install-recommends logrotate \
  && rm -rf /var/lib/apt/lists/*

COPY logc logc-start docker-kill /usr/local/bin/

WORKDIR /home/ubuntu
CMD [ "logc-start" ]
