FROM ubuntu:16.04

COPY locale /etc/default/locale
COPY sources.list /etc/apt/sources.list

RUN apt-get update &&\
 apt-get install -y --no-install-recommends lsof &&\
 rm -rf /var/lib/apt/lists/*
