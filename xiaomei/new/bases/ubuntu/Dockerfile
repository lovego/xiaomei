FROM ubuntu:16.04

COPY sources.list /etc/apt/sources.list

# setup locale
RUN locale-gen en_US.UTF-8
ENV LANG=en_US.UTF-8 LANGUAGE=en_US:en LC_ALL=en_US.UTF-8

# setup timezone
RUN ln -fs /usr/share/zoneinfo/Asia/Chongqing /etc/localtime \
  && dpkg-reconfigure -f noninteractive tzdata

# install less, wget, net-tools
RUN apt-get update \
  && apt-get install -y --no-install-recommends \
  less wget sudo ca-certificates net-tools iputils-ping dnsutils \
  && rm -rf /var/lib/apt/lists/*

# add user
RUN useradd -ms /bin/bash ubuntu && echo "ubuntu  ALL=NOPASSWD: ALL" >> /etc/sudoers

# set rlimit nofile
# RUN echo 'ubuntu - nofile unlimited' > /etc/security/limits.d/99-ubuntu.conf

