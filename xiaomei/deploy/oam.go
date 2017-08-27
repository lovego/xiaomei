package deploy

import (
	"fmt"

	"github.com/lovego/xiaomei/utils/cmd"
	"github.com/lovego/xiaomei/xiaomei/cluster"
	"github.com/lovego/xiaomei/xiaomei/release"
)

func logs(env, svcName, feature string) error {
	script := fmt.Sprintf(`
for name in $(docker ps -af name=%s --format '{{.Names}}'); do
	echo -e "\033[32m$name\033[0m"
	docker logs $name
	echo
done
`, release.ServiceName(env, svcName))
	return eachNodeRun(env, script, feature)
}

func rm(env, svcName, feature string) error {
	script := fmt.Sprintf(`
for name in $(docker ps -af name=%s --format '{{.Names}}'); do
	docker stop $name >/dev/null 2>&1 && docker rm $name
done
`, release.ServiceName(env, svcName))
	return eachNodeRun(env, script, feature)
}

func ps(env, svcName, feature string, watch bool) error {
	script := fmt.Sprintf(`docker ps -f name=%s`, release.ServiceName(env, svcName))
	if watch {
		script = `watch ` + script
	}
	return eachNodeRun(env, script, feature)
}

func eachNodeRun(env, script, feature string) error {
	for _, node := range cluster.Get(env).GetNodes(feature) {
		if _, err := node.Run(cmd.O{}, script); err != nil {
			return err
		}
	}
	return nil
}
