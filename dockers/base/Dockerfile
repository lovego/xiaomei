FROM ubuntu:16.04

COPY sources.list /etc/apt/sources.list

# add user, generate locale, set timezone, install wget, set nofile
RUN useradd -ms /bin/bash ubuntu \
     && locale-gen en_US.UTF-8 \
     && ln -fs /usr/share/zoneinfo/Asia/Chongqing /etc/localtime \
     && dpkg-reconfigure -f noninteractive tzdata \
     && apt-get update \
     && apt-get install -y --no-install-recommends ca-certificates wget net-tools \
     && rm -rf /var/lib/apt/lists/*
      # echo 'ubuntu - nofile unlimited' > /etc/security/limits.d/99-ubuntu.conf

# set locale
ENV LANG=en_US.UTF-8 LANGUAGE=en_US:en LC_ALL=en_US.UTF-8

USER ubuntu
RUN mkdir /home/ubuntu/app /home/ubuntu/logs
WORKDIR /home/ubuntu/app

EXPOSE 3000
HEALTHCHECK --interval=1s --timeout=3s --retries=1 \
  CMD wget --spider http://localhost:3000/alive || exit 1

