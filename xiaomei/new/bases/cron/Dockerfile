FROM 192.168.202.12:5000/xiaomei/base

RUN apt-get update \
  && apt-get install -y --no-install-recommends cron logrotate\
  && rm -rf /etc/logrotate.d/* /var/lib/apt/lists/*\
  && echo 'include /etc/logrotate.d' > /etc/logrotate.conf

COPY docker-kill /usr/local/bin/
CMD [ "cron", "-f" ]
