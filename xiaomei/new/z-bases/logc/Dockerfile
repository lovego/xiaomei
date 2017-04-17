FROM 192.168.202.12:5000/xiaomei/baseimage

USER ubuntu
RUN mkdir -p /home/ubuntu/logc
WORKDIR /home/ubuntu

COPY logc /usr/local/bin/

ENTRYPOINT [ "logc" ]
