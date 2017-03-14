FROM 192.168.202.12:5000/xiaomei/baseimage

USER ubuntu
RUN mkdir /home/ubuntu/logs

EXPOSE 3000
HEALTHCHECK --interval=3s --timeout=3s --retries=20 \
  CMD wget -q --spider http://localhost:3000/_alive || exit 1

