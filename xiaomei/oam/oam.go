package oam

import (
	"fmt"

	"github.com/lovego/cmd"
	"github.com/lovego/xiaomei/xiaomei/cluster"
	"github.com/lovego/xiaomei/xiaomei/deploy/conf"
	"github.com/lovego/xiaomei/xiaomei/release"
)

func shell(svcName, env, feature string) error {
	_, err := cluster.Get(env).ServiceRun(svcName, feature, cmd.O{},
		fmt.Sprintf(
			"docker exec -it -e LINES=$(tput lines) -e COLUMNS=$(tput cols) "+
				"--detach-keys='ctrl-@' %s bash",
			conf.GetService(svcName, env).FirstContainerName(),
		),
	)
	return err
}

func logs(svcName, env, feature string) error {
	script := fmt.Sprintf(`
for name in $(docker ps -af name=^/%s --format '{{.Names}}'); do
	echo -e "\033[32m$name\033[0m"
	docker logs $name
	echo
done
`, release.ServiceName(svcName, env))
	return eachNodeRun(env, script, feature)
}

func operate(operation, svcName, env, feature string) error {
	if operation != `start` && operation != `stop` && operation != `restart` {
		return fmt.Errorf("invalid operation of %s", operation)
	}
	script := fmt.Sprintf(`
for name in $(docker ps -af name=^/%s --format '{{.Names}}'); do
	docker %s $name
done
`, release.ServiceName(svcName, env), operation)
	return eachNodeRun(env, script, feature)
}

func ps(svcName, env, feature string, watch bool) error {
	script := fmt.Sprintf(` docker ps -f name=^/%s`, release.ServiceName(svcName, env))
	if watch {
		script = WatchCmd() + script
	}
	return eachNodeRun(env, script, feature)
}

func WatchCmd() string {
	return `which watch >/dev/null || watch() {
  trap "echo; exit 0" INT
  while true; do
   clear
   "$@"
   sleep 2
  done
}
watch`
}

func eachNodeRun(env, script, feature string) error {
	for _, node := range cluster.Get(env).GetNodes(feature) {
		if _, err := node.Run(cmd.O{}, script); err != nil {
			return err
		}
	}
	return nil
}
