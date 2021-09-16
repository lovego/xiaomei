package oam

import (
	"fmt"

	"github.com/lovego/cmd"
	"github.com/lovego/xiaomei/access"
	"github.com/lovego/xiaomei/release"
)

func shell(svcName, env, feature string) error {
	_, err := release.GetCluster(env).ServiceRun(svcName, feature, cmd.O{},
		fmt.Sprintf(
			"docker exec -it -e LINES=$(tput lines) -e COLUMNS=$(tput cols) "+
				"--detach-keys='ctrl-@' %s bash",
			release.GetService(env, svcName).FirstContainerName(),
		),
	)
	return err
}

func logs(svcName, env, feature, options string) error {
	script := fmt.Sprintf(`
for name in $(docker ps -af name='%s' --format '{{.Names}}'); do
	echo -e "\033[32m$name\033[0m"
	docker logs %s $name
	echo
done
`, release.ContainerNameRegexp(svcName, env), options)
	return eachNodeRun(env, script, feature)
}

func operate(operation, svcName, env, feature string) error {
	if operation != `start` && operation != `stop` && operation != `restart` {
		return fmt.Errorf("invalid operation: %s", operation)
	}
	script := fmt.Sprintf(`
for name in $(docker ps -af name='%s' --format '{{.Names}}'); do
	docker %s $name
done
`, release.ContainerNameRegexp(svcName, env), operation)
	for _, node := range release.GetCluster(env).GetNodes(feature) {
		if _, err := node.Run(cmd.O{}, script); err != nil {
			return err
		}
		if access.HasAccess(node.Services(env, svcName)) {
			if err := access.ReloadNginx(env, feature); err != nil {
				return err
			}
		}
	}
	return nil
}

func ps(svcName, env, feature string, watch bool) error {
	script := fmt.Sprintf(` docker ps -af name='%s'`, release.ContainerNameRegexp(svcName, env))
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
	for _, node := range release.GetCluster(env).GetNodes(feature) {
		if _, err := node.Run(cmd.O{}, script); err != nil {
			return err
		}
	}
	return nil
}
