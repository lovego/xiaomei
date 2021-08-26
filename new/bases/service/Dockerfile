FROM ubuntu:20.04

COPY sources.list /etc/apt/sources.list
RUN apt-get update \
        && apt-get install -y --no-install-recommends locales tzdata sudo ca-certificates gnupg \
        openssh-client less wget lsof net-tools iputils-ping dnsutils logrotate \
        && rm -rf /var/lib/apt/lists/*

COPY nginx/sources.list /etc/apt/sources.list.d/nginx.list
# apt-key add depends on gnupg
RUN { wget -O- 'http://nginx.org/keys/nginx_signing.key' | apt-key add -; } \
        && apt-get update \
        && apt-get install -y --no-install-recommends nginx \
        && rm -rf /var/lib/apt/lists/*

RUN locale-gen en_US.UTF-8 \
        && ln -fs /usr/share/zoneinfo/Asia/Chongqing /etc/localtime \
        && dpkg-reconfigure -f noninteractive tzdata
ENV LANG=en_US.UTF-8 LANGUAGE=en_US:en

RUN  mkdir /etc/nginx/sites-enabled /etc/nginx/sites-available \
        && rm -rf /etc/nginx/conf.d /var/log/nginx/*.log
COPY nginx/nginx.conf nginx/proxy_params /etc/nginx/
COPY nginx/nginx-start /usr/local/bin/
COPY logc/logc logc/docker-kill /usr/local/bin/

# nginx use quit signal for graceful shutdown, so goa support it too.
STOPSIGNAL SIGQUIT
RUN useradd -ms /bin/bash ubuntu && echo "ubuntu  ALL=NOPASSWD: ALL" >> /etc/sudoers
USER ubuntu
WORKDIR /home/ubuntu
# set rlimit nofile
# RUN echo 'ubuntu - nofile unlimited' > /etc/security/limits.d/99-ubuntu.conf


