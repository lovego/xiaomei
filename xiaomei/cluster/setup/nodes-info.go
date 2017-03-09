package setup

import (
	"encoding/json"
	"errors"

	"github.com/bughou-go/xiaomei/utils/cmd"
	"github.com/bughou-go/xiaomei/xiaomei/release"
)

type NodeInfo struct {
	Node
	Cluster struct {
		ID string
	}
	NodeID         string
	RemoteManagers []struct {
		NodeID string
	}
}

func (info *NodeInfo) fetch() error {
	output, err := cmd.Run(cmd.O{Output: true}, `ssh`, info.Config.SshAddr(),
		`docker`, `system`, `info`, `-f`, `{{ json .Swarm }}`,
	)
	if err != nil {
		return err
	}
	if err := json.Unmarshal([]byte(output), info); err != nil {
		return err
	}
	if info.NodeID != `` { // in swarm mode
		if info.Cluster.ID != `` {
			info.Role = `manager` // only manager stores ClusterID
		} else {
			info.Role = `worker`
		}
	}
	return nil
}

func getNodesInfo(clusterConf release.Cluster) ([]NodeInfo, []NodeInfo, error) {
	managers, err1 := fetchNodesInfo(clusterConf.Managers)
	if err1 != nil {
		return nil, nil, err1
	}
	workers, err2 := fetchNodesInfo(clusterConf.Workers)
	if err2 != nil {
		return nil, nil, err2
	}

	if err := setupEmptyClusterID(managers, workers); err != nil {
		return nil, nil, err
	}
	return managers, workers, nil
}

func fetchNodesInfo(nodes []release.Node) (result []NodeInfo, err error) {
	for _, n := range nodes {
		info := NodeInfo{Node: Node{Config: n}}
		if err := info.fetch(); err != nil {
			return nil, err
		}
		result = append(result, info)
	}
	return
}

func setupEmptyClusterID(managers, workers []NodeInfo) error {
	m := map[string]string{}
	mapNodeID2ClusterID(managers, m)
	mapNodeID2ClusterID(workers, m)
	if err := setupClusterID(managers, m); err != nil {
		return err
	}
	return setupClusterID(workers, m)
}

func mapNodeID2ClusterID(infos []NodeInfo, m map[string]string) {
	for _, info := range infos {
		if info.Cluster.ID != `` {
			m[info.NodeID] = info.Cluster.ID
		}
	}
}

func setupClusterID(infos []NodeInfo, m map[string]string) error {
	for _, info := range infos {
		if info.NodeID != `` && info.Cluster.ID == `` && !findAndSetClusterID(info, m) {
			return errors.New(`no ClusterID found for ` + info.Config.Addr)
		}
	}
	return nil
}

func findAndSetClusterID(info NodeInfo, m map[string]string) bool {
	for _, manger := range info.RemoteManagers {
		if clusterID := m[manger.NodeID]; clusterID != `` {
			info.Cluster.ID = clusterID
			return true
		}
	}
	return false
}
