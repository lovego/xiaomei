FROM ubuntu:16.04

# set locale
RUN locale-gen en_US.UTF-8  
ENV LANG en_US.UTF-8  
ENV LANGUAGE en_US:en  
ENV LC_ALL en_US.UTF-8 

# set timezone
RUN ln -fs /usr/share/zoneinfo/Asia/Chongqing /etc/localtime && \
  dpkg-reconfigure -f noninteractive tzdata


COPY sources.list /etc/apt/sources.list

RUN apt-get update && \
  apt-get install -y --no-install-recommends lsof && \
  rm -rf /var/lib/apt/lists/*
