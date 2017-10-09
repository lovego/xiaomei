package workspace_godoc

import (
	"fmt"
	"os"

	"github.com/lovego/xiaomei/utils/cmd"
)

func deploy() error {
	goRoot := os.Getenv(`GOROOT`)
	if goRoot == `` {
		goRoot = `/usr/local/go`
	}
	script := fmt.Sprintf(`
  docker stop workspace-godoc >/dev/null 2>&1 && docker rm workspace-godoc
	docker run --name=workspace-godoc -d --restart=always \
	--network=host -e=GODOCPORT=1234 \
	-v %s:/home/ubuntu/go -v %s:/usr/local/go \
	hub.c.163.com/lovego/xiaomei/godoc
	`, os.Getenv(`GOPATH`), goRoot)
	_, err := cmd.Run(cmd.O{}, `sh`, `-c`, script)
	return err
}
