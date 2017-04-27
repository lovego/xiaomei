package simple

import (
	"fmt"

	"github.com/lovego/xiaomei/utils/cmd"
	"github.com/lovego/xiaomei/xiaomei/cluster"
	"github.com/lovego/xiaomei/xiaomei/release"
)

func (d driver) Logs(svcName string) error {
	script := fmt.Sprintf(`
for name in $(docker ps -af name=%s_%s. --format '{{.Names}}'); do
	echo -e "\033[32m$name\033[0m"
	docker logs $name
	echo
done
`, release.Name(), svcName)
	return eachNodeRun(script)
}

func (d driver) RmDeploy(svcName string) error {
	script := fmt.Sprintf(`
for name in $(docker ps -af name=%s_%s. --format '{{.Names}}'); do
	docker stop $name >/dev/null 2>&1 && docker rm $name
done
`, release.Name(), svcName)
	return eachNodeRun(script)
}

func (d driver) Ps(svcName string, watch bool, options []string) error {
	return eachNodeRun(getPsScript(svcName, watch))
}

func getPsScript(svcName string, watch bool) string {
	script := fmt.Sprintf(`docker ps -f name=%s_%s`, release.Name(), svcName)
	if watch {
		script = `watch ` + script
	}
	return script
}

func eachNodeRun(script string) error {
	for _, node := range cluster.Nodes() {
		if _, err := node.Run(cmd.O{}, script); err != nil {
			return err
		}
	}
	return nil
}
