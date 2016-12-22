package setup

import (
	"fmt"
	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils"
	"github.com/bughou-go/xiaomei/utils/cmd"
	"path"
)

func SetupHosts() {
	filePath := path.Join(config.App.Root(), `config/shell/setup-hosts.sh`)
	if !utils.IsFile(filePath) {
		fmt.Println(`no such file: ` + filePath)
		return
	}
	cmd.Run(cmd.O{Panic: true}, filePath)

	fmt.Println(`setup hosts ok.`)
}

/*
#!/bin/bash

dir="$(dirname $0)/.."
target=/etc/hosts

config=$(<"$target")
begin_hosts_end='#begin'$'\n'$(<"$dir/hosts.txt")$'\n''#end'

if [[ "$config" == *"$begin_hosts_end"* ]]; then
  exit
elif [[ "$config" == *'#begin'*'#end'* ]]; then
  echo "${config/'#begin'*'#end'/$begin_hosts_end}" | sudo tee "$target" > /dev/null
else
  echo $'\n'"$begin_hosts_end" | sudo tee --append "$target" > /dev/null
fi
*/
