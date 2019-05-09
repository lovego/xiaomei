FROM registry.cn-beijing.aliyuncs.com/lovego/golang

COPY xiaomei /usr/local/bin/
# if cgo enabled, go test fails with 'exec: "gcc": executable file not found in $PATH'.
ENV CGO_ENABLED=0


