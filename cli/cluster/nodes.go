package cluster

import (
	"encoding/json"
	"errors"
	"github.com/bughou-go/xiaomei/config"
)

type Node struct {
	Config config.Node
	Role   string // current role
}

func getClusterNodes() ([]Node, []Node, error) {
	if managers, workers, err := getNodesInfo(); err != nil {
		return nil, nil, err
	} else if err = checkIsInOneCluster(managers, workers); err != nil {
		return nil, nil, err
	} else {
		return mapNodes(managers), mapNodes(workers), nil
	}
}

// 检查所有节点只能属于同一个集群.
func checkIsInOneCluster(managers, workers []NodeInfo) error {
	m := make(map[string][]string)
	addToClusterIDMap(managers, m)
	addToClusterIDMap(workers, m)
	if len(m) <= 1 {
		return nil
	}
	s, err := json.MarshalIndent(m, ``, "\t")
	if err != nil {
		return err
	}
	return errors.New("found more than one cluster in target nodes:\n" + string(s))
}

func addToClusterIDMap(infos []NodeInfo, m map[string][]string) {
	for _, info := range infos {
		if info.Cluster.ID != `` {
			m[info.Cluster.ID] = append(m[info.Cluster.ID], info.Config.Addr)
		}
	}
}

func mapNodes(infos []NodeInfo) (result []Node) {
	for _, info := range infos {
		result = append(result, info.Node)
	}
	return
}
