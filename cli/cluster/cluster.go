package cluster

import (
	"errors"

	"github.com/bughou-go/xiaomei/config"
)

type Cluster struct {
	Managers []Node
	Workers  []Node
}

func Setup() error {
	if len(config.Cluster.Managers()) == 0 {
		return errors.New(`the cluster have no managers.`)
	}
	managers, workers, err := getClusterNodes()
	if err != nil {
		return err
	}
	return Cluster{managers, workers}.setup()
}

func (c Cluster) setup() error {
	if err := c.setupManagers(); err != nil {
		return err
	}
	return c.setupWorkers()
}

func (c Cluster) setupManagers() error {
	for _, m := range c.Managers {
		switch m.Role {
		case ``:
			c.managerAddr()
		case `worker`:
		}
	}
	return nil
}

func (c Cluster) setupWorkers() error {
	for _, w := range c.Workers {
		switch w.Role {
		case ``:
		case `manager`:
		}
	}
	return nil
}

func (c Cluster) managerAddr() string {
	if m := findManager(c.Managers); m.Role != `` {
		return m.Config.Addr
	}
	if m := findManager(c.Workers); m.Role != `` {
		return m.Config.Addr
	}
	return ``
}

func findManager(nodes []Node) Node {
	for _, node := range nodes {
		if node.Role == `manager` {
			return node
		}
	}
	return Node{}
}
