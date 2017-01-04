package deploy

import (
	"fmt"

	"github.com/bughou-go/xiaomei/config"
	"github.com/fatih/color"
)

func Deploy(commit, serverFilter string) error {
	tag, err := setupDeployTag(commit)
	if err != nil {
		return err
	}

	servers := config.Servers.Matched(serverFilter)
	updated := make(map[string]bool)
	for _, server := range servers {
		sshAddr := server.SshAddr()
		color.Cyan(sshAddr)
		if !updated[server.Addr] {
			updateCode(sshAddr, tag)
			copyBinary(sshAddr, tag)
			updated[server.Addr] = true
		}
		setupTasks(sshAddr, tag, server.Tasks)
	}
	fmt.Printf("deployed %d servers!\n", len(servers))
	return ClearTags()
}
