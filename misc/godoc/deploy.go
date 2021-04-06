package godoc

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/lovego/cmd"
	"github.com/lovego/xiaomei/release"
)

func run() error {
	if nullDevice, err := os.OpenFile(`/dev/null`, os.O_WRONLY, 0); err != nil {
		return err
	} else if !cmd.Ok(cmd.O{Stdout: nullDevice}, `which`, `godoc`) {
		if err := release.GoGetByProxy("golang.org/x/tools/..."); err != nil {
			return err
		}
	}
	script := `
killall godoc >/dev/null 2>&1
nohup godoc -http=:1234 -index_interval=1s >/dev/null 2>&1 &
`
	_, err := cmd.Run(cmd.O{}, `sh`, `-c`, script)
	if err == nil {
		printAddr()
	}
	return err
}

func deploy() error {
	script := os.Expand(`
docker stop workspace-godoc >/dev/null 2>&1 && docker rm workspace-godoc
docker run --name=workspace-godoc -d --restart=always -e=GODOCPORT=1234 --publish=1234:1234 \
	-v $($GoCmd env GOPATH):/home/ubuntu/go -v $($GoCmd env GOROOT):/usr/local/go \
	registry.cn-beijing.aliyuncs.com/lovego/godoc
`, func(name string) string {
		return release.GoCmd()
	})

	_, err := cmd.Run(cmd.O{}, `sh`, `-c`, script)
	if err == nil {
		printAddr()
	}
	return err
}

func rmDeploy() error {
	script := `docker stop workspace-godoc && docker rm workspace-godoc`
	_, err := cmd.Run(cmd.O{}, `sh`, `-c`, script)
	return err
}

func printAddr() {
	fmt.Println("started at " + color.GreenString(`http://localhost:1234/`))
}
