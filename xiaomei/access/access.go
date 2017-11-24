package access

import (
	"fmt"
	"strings"

	"github.com/lovego/utils/cmd"
	"github.com/lovego/xiaomei/xiaomei/cluster"
)

func accessPrint(env, svcName string) error {
	nginxConf, _, err := getNginxConf(env, svcName)
	if err != nil {
		return err
	}
	fmt.Print(nginxConf)
	return nil
}

func accessSetup(env, svcName, feature string) error {
	nginxConf, fileName, err := getNginxConf(env, svcName)
	if err != nil {
		return err
	}
	script := fmt.Sprintf(`
	sudo tee /etc/nginx/sites-enabled/%s.conf > /dev/null &&
	sudo mkdir -p /var/log/nginx/%s &&
	sudo nginx -t &&
	sudo service nginx reload
	`, fileName, fileName,
	)
	for _, node := range cluster.Get(env).GetNodes(feature) {
		if node.Labels[`access`] == `true` {
			if _, err := node.Run(
				cmd.O{Stdin: strings.NewReader(nginxConf)}, script,
			); err != nil {
				return err
			}
		}
	}
	return nil
}
