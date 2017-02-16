package deploy

import (
	"encoding/json"
	"errors"

	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils/cmd"
)

func Setup() error {
	if len(config.Cluster.Managers()) == 0 {
		return errors.New(`the cluster have no managers.`)
	}
	for _, m := range config.Cluster.Managers() {

	}
	return nil
}

// 检查所有节点，当前的集群数要么为0要么为1.
func CheckNodes() error {
}

func setupManagers() {
}

func setupWorkers() {
}
