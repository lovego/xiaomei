package godoc

import (
	"github.com/lovego/cmd"
)

func run() error {
	script := `
  killall godoc >/dev/null 2>&1
  nohup godoc -http=:1234 -index_interval=1s >/dev/null 2>&1 &
  `
	_, err := cmd.Run(cmd.O{}, `sh`, `-c`, script)
	return err
}

func deploy() error {
	script := `
  docker stop workspace-godoc >/dev/null 2>&1 && docker rm workspace-godoc
	docker run --name=workspace-godoc -d --restart=always \
    -e=GODOCPORT=1234 --publish=1234:1234 \
	  -v $(go env GOPATH):/home/ubuntu/go -v $(go env GOROOT):/usr/local/go \
	  hub.c.163.com/lovego/xiaomei/godoc
	`
	_, err := cmd.Run(cmd.O{}, `sh`, `-c`, script)
	return err
}

func rmDeploy() error {
	script := `docker stop workspace-godoc && docker rm workspace-godoc`
	_, err := cmd.Run(cmd.O{}, `sh`, `-c`, script)
	return err
}
