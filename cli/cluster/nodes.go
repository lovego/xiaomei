package cluster

import (
	"encoding/json"
	"errors"

	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils/cmd"
)

type Node struct {
	Config config.Node
	Role   string // current role
}

func (n *Node) init() error {
	_, err := cmd.Run(cmd.O{Print: true}, `ssh`,
		append([]string{n.Config.SshAddr(), `docker`, `swarm`, `init`}, n.addrFlags()...)...,
	)
	if err == nil {
		n.Role = `manager`
	}
	return err
}

func (n *Node) join(role, token, addr string) error {
	args := []string{n.Config.SshAddr(), `docker`, `swarm`, `join`}
	if role == `manager` {
		args = append(args, n.addrFlags()...)
	}
	args = append(args, `--token`, token, addr)
	_, err := cmd.Run(cmd.O{Print: true}, `ssh`, args...)
	if err == nil {
		n.Role = role
	}
	return err
}

func (n Node) token(role string) (string, error) {
	if n.Role == `` {
		return ``, nil
	}
	return cmd.Run(cmd.O{Output: true}, `ssh`, n.Config.SshAddr(),
		`docker`, `swarm`, `join-token`, `-q`, role,
	)
}

func (n Node) addrFlags() []string {
	var addr string
	if n.Config.ListenAddr != `` {
		addr = n.Config.ListenAddr
	} else {
		addr = n.Config.Addr
	}
	return []string{`--advertise-addr`, addr, `--listen-addr`, addr}
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
