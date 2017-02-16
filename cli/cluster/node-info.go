package cluster

import (
	"encoding/json"
	"errors"

	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils/cmd"
)

type NodeInfo struct {
	Config      config.Node
	Role        string
	ClusterID   string
	NodeID      string
	ManagersIDs []string
}

type swarmInfo struct {
	NodeID  string
	Cluster struct {
		ID string
	}
	RemoteManagers []struct {
		NodeID string
	}
}

func getNodesInfo() ([]NodeInfo, []NodeInfo, error) {
	for _, manager := range config.Cluster.Managers() {
	}
	return nil, nil
}

func getNodeInfo(nodeConfig config.Node) (*nodeInfo, error) {
	output, err := cmd.Run(cmd.O{Output: true}, `ssh`, nodeConfig.SshAddr(),
		`docker`, `system`, `info`, `-f`, `{{ json .Swarm }}`,
	)
	if err != nil {
		return nil, err
	}
	info := swarmInfo{}
	if err := json.Unmarshal([]byte(output), &info); err != nil {
		return nil, err
	}
	return makeNodeInfo(nodeConfig, info), nil
}

func makeNodeInfo(nodeConfig config.Node, info swarmInfo) NodeInfo {
	node := NodeInfo{Config: nodeConfig}
	if info.NodeId == `` {
		return node // not in swarm mode
	}
	node.NodeID = info.NodeId
	if info.Cluster.ID != `` {
		node.Role = `manager` // only manager stores ClusterID
		node.ClusterID = info.Cluster.ID
	} else {
		node.Role = `worker`
		ids := []string{}
		for _, m := range info.RemoteManagers {
			ids = append(ids, m.NodeID)
		}
		node.ManagersIDs = ids
	}
	return node
}
