package access

import (
	"fmt"
	"strings"

	"github.com/lovego/xiaomei/utils/cmd"
	"github.com/lovego/xiaomei/xiaomei/cluster"
)

func accessPrint(svcName string) error {
	nginxConf, _, err := getNginxConf(svcName)
	if err != nil {
		return err
	}
	fmt.Print(nginxConf)
	return nil
}

func accessSetup(svcName string) error {
	nginxConf, fileName, err := getNginxConf(svcName)
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
	for _, node := range cluster.Nodes() {
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
