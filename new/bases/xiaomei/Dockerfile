FROM registry.cn-beijing.aliyuncs.com/lovego/ubuntu

RUN apt-get update && \
    apt-get install -y --no-install-recommends golang-1.10-go && \
    rm -rf /var/lib/apt/lists/*

COPY xiaomei /bin/

ENV PATH /usr/lib/go-1.10/bin:$PATH
ENV GOPATH /go

WORKDIR /go/src
