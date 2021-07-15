# vim: ft=dockerfile:

FROM registry.cn-beijing.aliyuncs.com/lovego/service

LABEL builder=xiaomei

COPY logc.yml logrotate.conf ./
RUN sudo chmod 644 logrotate.conf
RUN mkdir app-logs web-logs

CMD ["logc", "logc.yml"]

