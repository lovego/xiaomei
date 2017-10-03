FROM hub.c.163.com/lovego/xiaomei/base

COPY godoc godoc-start /usr/local/bin/
COPY pkg.png /usr/local/go/doc/gopher/
RUN mkdir -p /usr/local/go/lib /usr/local/go/project /usr/local/go/help
COPY favicon.ico /usr/local/go/

USER ubuntu
RUN mkdir -p /home/ubuntu/go/src
WORKDIR /home/ubuntu/go/src

# godoc exit fatal error if no package.  so make a placeholder package "z".
RUN mkdir z && echo 'package z' > /home/ubuntu/go/src/z/z.go

ENV GOPATH=/home/ubuntu/go
CMD [ "godoc-start" ]
