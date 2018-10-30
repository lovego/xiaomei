FROM ubuntu:16.04

COPY sources.list /etc/apt/sources.list

RUN apt-get update \
      && apt-get install -y --no-install-recommends \
      locales tzdata less wget sudo ca-certificates net-tools iputils-ping dnsutils \
      && rm -rf /var/lib/apt/lists/* \
      && locale-gen en_US.UTF-8 \
      && ln -fs /usr/share/zoneinfo/Asia/Chongqing /etc/localtime \
      && dpkg-reconfigure -f noninteractive tzdata \
      && useradd -ms /bin/bash ubuntu \
      && echo "ubuntu  ALL=NOPASSWD: ALL" >> /etc/sudoers

ENV LANG=en_US.UTF-8 LANGUAGE=en_US:en

# set rlimit nofile
# RUN echo 'ubuntu - nofile unlimited' > /etc/security/limits.d/99-ubuntu.conf

