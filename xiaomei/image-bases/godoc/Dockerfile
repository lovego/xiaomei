FROM 192.168.202.12:5000/xiaomei/base

COPY godoc /usr/local/bin/
RUN mkdir -p /usr/local/go/doc /usr/local/go/lib /usr/local/go/project /usr/local/go/help
COPY pkg.png /usr/local/go/doc/gopher/

USER ubuntu
RUN mkdir -p /home/ubuntu/go/src
WORKDIR /home/ubuntu/go/src

# godoc exit fatal error if no package.  so make a placeholder package "z".
RUN mkdir z && echo 'package z' > /home/ubuntu/go/src/z/z.go

ENV GOPATH=/home/ubuntu/go
CMD godoc -http=:${GODOCPORT:-7000} -index_interval=-1s
