package godoc

import (
	"github.com/lovego/cmd"
)

func deploy() error {
	_, err := cmd.Run(cmd.O{}, `sh`, `-c`, `
  docker stop workspace-godoc >/dev/null 2>&1 && docker rm workspace-godoc
	docker run --name=workspace-godoc -d --restart=always \
	--network=host -e=GODOCPORT=1234 \
	-v $(go env GOPATH):/home/ubuntu/go -v $(go env GOROOT):/usr/local/go \
	hub.c.163.com/lovego/xiaomei/godoc
	`)
	return err
}
