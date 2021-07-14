FROM registry.cn-beijing.aliyuncs.com/lovego/service

ADD go1.12.4.linux-amd64.tar.gz /usr/local/
COPY xiaomei gospec godoc-start /usr/local/bin/

ENV PATH=/usr/local/go/bin:$PATH
ENV GOPATH=/home/ubuntu/go
# if cgo enabled, go test fails with 'exec: "gcc": executable file not found in $PATH'.
ENV CGO_ENABLED=0

USER ubuntu
RUN mkdir -p /home/ubuntu/go/src


