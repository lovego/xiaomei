FROM 192.168.202.12:5000/xiaomei/baseimage

USER ubuntu
RUN mkdir /home/ubuntu/logs

EXPOSE 3000
HEALTHCHECK --interval=1s --timeout=3s --retries=1 \
  CMD wget --spider http://localhost:3000/alive || exit 1

