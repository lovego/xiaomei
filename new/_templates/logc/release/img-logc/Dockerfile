# vim: ft=dockerfile:

FROM registry.cn-beijing.aliyuncs.com/lovego/logc:2

LABEL builder=xiaomei

COPY logc.yml logrotate.conf .
RUN chmod 644 logrotate.conf
RUN mkdir app-logs web-logs

