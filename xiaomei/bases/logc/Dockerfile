FROM 192.168.202.12:5000/xiaomei/base

USER ubuntu
RUN mkdir -p /home/ubuntu/logc

COPY logc /usr/local/bin/

CMD [ "logc", "../logc.yml" ]
