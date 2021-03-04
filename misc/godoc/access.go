package godoc

import (
	"fmt"
	"strings"

	"github.com/lovego/cmd"
)

const nginxConf = `# vim: set ft=nginx:

server {
  listen  80;
  server_name godoc.dev;
  charset utf-8;

  location / {
    proxy_pass http://0.0.0.0:1234;
    include    proxy_params;
  }
  access_log /var/log/nginx/godoc.dev/access.log;
  error_log  /var/log/nginx/godoc.dev/access.err;
}
`

func accessPrint() error {
	fmt.Print(nginxConf)
	return nil
}

func accessSetup() error {
	script := `
set -e
sudo tee /etc/nginx/sites-enabled/godoc.dev.conf > /dev/null
sudo mkdir -p /var/log/nginx/godoc.dev
nginx -s reload
`
	_, err := cmd.Run(
		cmd.O{Stdin: strings.NewReader(nginxConf)}, `bash`, `-c`, script,
	)
	return err
}
