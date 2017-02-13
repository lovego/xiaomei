FROM ubuntu:16.04

COPY sources.list /etc/apt/sources.list

# add user, generate locale, set timezone, install curl, set nofile
RUN useradd -ms /bin/bash ubuntu && \
      locale-gen en_US.UTF-8 && \
      ln -fs /usr/share/zoneinfo/Asia/Chongqing /etc/localtime && \
      dpkg-reconfigure -f noninteractive tzdata && \
      apt-get update && \
      apt-get install -y --no-install-recommends curl && \
      rm -rf /var/lib/apt/lists/*
      # echo 'ubuntu - nofile unlimited' > /etc/security/limits.d/99-ubuntu.conf

# set locale
ENV LANG     en_US.UTF-8
ENV LANGUAGE en_US:en
ENV LC_ALL   en_US.UTF-8 

USER ubuntu
WORKDIR /home/ubuntu
EXPOSE 3000

HEALTHCHECK --interval=3s --timeout=3s --retries=1 \
  CMD curl -f http://localhost:3000/alive || exit 1
