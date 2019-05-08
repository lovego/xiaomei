FROM registry.cn-beijing.aliyuncs.com/lovego/ubuntu

COPY go1.12.4.linux-amd64.tar.gz .
RUN tar -C /usr/local -xzf go1.12.4.linux-amd64.tar.gz && rm go1.12.4.linux-amd64.tar.gz
ENV PATH=/usr/local/go/bin:$PATH
ENV GOPATH=/home/ubuntu/go

USER ubuntu
RUN mkdir -p /home/ubuntu/go/src
WORKDIR /home/ubuntu/go/src

