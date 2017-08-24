package deploy

import (
	"fmt"

	"github.com/lovego/xiaomei/utils/cmd"
	"github.com/lovego/xiaomei/xiaomei/cluster"
	"github.com/lovego/xiaomei/xiaomei/release"
)

func logs(svcName, feature string) error {
	script := fmt.Sprintf(`
for name in $(docker ps -af name=%s_%s --format '{{.Names}}'); do
	echo -e "\033[32m$name\033[0m"
	docker logs $name
	echo
done
`, release.DeployName(), svcName)
	return eachNodeRun(script, feature)
}

func rm(svcName, feature string) error {
	script := fmt.Sprintf(`
for name in $(docker ps -af name=%s_%s --format '{{.Names}}'); do
	docker stop $name >/dev/null 2>&1 && docker rm $name
done
`, release.DeployName(), svcName)
	return eachNodeRun(script, feature)
}

func ps(svcName, feature string, watch bool) error {
	script := fmt.Sprintf(`docker ps -f name=%s_%s`, release.DeployName(), svcName)
	if watch {
		script = `watch ` + script
	}
	return eachNodeRun(script, feature)
}

func eachNodeRun(script, feature string) error {
	for _, node := range cluster.Nodes(feature) {
		if _, err := node.Run(cmd.O{}, script); err != nil {
			return err
		}
	}
	return nil
}
