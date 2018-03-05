FROM hub.c.163.com/lovego/xiaomei/ubuntu

USER ubuntu
WORKDIR /home/ubuntu

EXPOSE 3000
HEALTHCHECK --interval=1m --timeout=3s \
  CMD wget -q --spider http://localhost:${GOPORT:-3000}/_alive || exit 1

