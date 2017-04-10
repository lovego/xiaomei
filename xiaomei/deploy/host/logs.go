package host

import (
	"fmt"

	"github.com/lovego/xiaomei/utils/cmd"
	"github.com/lovego/xiaomei/xiaomei/cluster"
	"github.com/lovego/xiaomei/xiaomei/release"
)

func (d driver) Logs(svcName string) error {
	for _, node := range cluster.GetCluster().Nodes() {
		if nodeLog(svcName, node) {
			return nil
		}
	}
	return nil
}

func nodeLog(svcName string, node cluster.Node) bool {
	script := fmt.Sprintf(`
for name in $(docker ps -af name=%s_%s. --format '{{.Names}}'); do
	echo -e "\033[32m$name\033[0m"
	docker logs $name
	echo
done
`, release.Name(), svcName)
	_, err := node.Run(cmd.O{}, script)
	return err == nil
}
