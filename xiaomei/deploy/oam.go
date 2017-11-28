package deploy

import (
	"fmt"

	"github.com/lovego/cmd"
	"github.com/lovego/xiaomei/xiaomei/cluster"
	"github.com/lovego/xiaomei/xiaomei/deploy/conf"
	"github.com/lovego/xiaomei/xiaomei/release"
)

func shell(svcName, env, feature string) error {
	_, err := cluster.Get(env).ServiceRun(svcName, feature, cmd.O{},
		`docker exec -it --detach-keys='ctrl-@' `+conf.GetService(svcName, env).FirstContainerName()+` bash`,
	)
	return err
}

func logs(svcName, env, feature string) error {
	script := fmt.Sprintf(`
for name in $(docker ps -af name=%s --format '{{.Names}}'); do
	echo -e "\033[32m$name\033[0m"
	docker logs $name
	echo
done
`, release.ServiceName(svcName, env))
	return eachNodeRun(env, script, feature)
}

func restart(svcName, env, feature string) error {
	script := fmt.Sprintf(`
for name in $(docker ps -af name=%s --format '{{.Names}}'); do
	docker restart $name
done
`, release.ServiceName(svcName, env))
	return eachNodeRun(env, script, feature)
}

func rmDeploy(svcName, env, feature string) error {
	script := fmt.Sprintf(`
for name in $(docker ps -af name=%s --format '{{.Names}}'); do
	docker stop $name >/dev/null 2>&1 && docker rm $name
done
`, release.ServiceName(svcName, env))
	return eachNodeRun(env, script, feature)
}

func ps(svcName, env, feature string, watch bool) error {
	script := fmt.Sprintf(`docker ps -f name=%s`, release.ServiceName(svcName, env))
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
