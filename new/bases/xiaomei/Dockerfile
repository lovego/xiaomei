FROM registry.cn-beijing.aliyuncs.com/lovego/service

RUN sudo apt-get update && sudo apt-get install -y --no-install-recommends gcc libc6-dev \
        && sudo rm -rf /var/lib/apt/lists/*

ADD go1.16.9.linux-amd64.tar.gz /usr/local/
COPY xiaomei gospec godoc-start /usr/local/bin/

ENV PATH=/usr/local/go/bin:$PATH
ENV GOPATH=/home/ubuntu/go

USER ubuntu
RUN mkdir -p /home/ubuntu/go/src


