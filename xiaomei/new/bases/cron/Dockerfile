FROM 192.168.202.12:5000/xiaomei/base

RUN apt-get update \
  && apt-get install -y --no-install-recommends cron \
  && rm -rf /var/lib/apt/lists/*

CMD [ "cron", "-f" ]
