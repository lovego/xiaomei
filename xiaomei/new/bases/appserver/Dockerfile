FROM hub.c.163.com/lovego/xiaomei/base

USER ubuntu
WORKDIR /home/ubuntu

EXPOSE 3000
HEALTHCHECK --interval=3s --timeout=3s --retries=20 \
  CMD wget -q --spider http://localhost:${GOPORT:-3000}/_alive || exit 1

