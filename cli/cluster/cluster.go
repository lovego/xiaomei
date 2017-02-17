package cluster

import (
	"errors"

	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils/cmd"
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
	if c.manager().Role == `` { // No Manager
		c.Managers[0].init()
	}
	if err := c.setupManagers(); err != nil {
		return err
	}
	return c.setupWorkers()
}

func (c Cluster) setupManagers() error {
	for i := 0; i < len(c.Managers); i++ {
		switch c.Managers[i].Role {
		case ``:
			if err := c.Managers[i].join(); err != nil {
				return err
			}
		case `worker`:
		}
	}
	return nil
}

func (c Cluster) setupWorkers() error {
	for i, w := range c.Workers {
		switch w.Role {
		case ``:
		case `manager`:
		}
	}
	return nil
}

func (c Cluster) join(node Node, role string) error {
	m := c.manager().token()
}

func (c Cluster) manager() Node {
	if m := findManager(c.Managers); m.Role != `` {
		return m
	}
	if m := findManager(c.Workers); m.Role != `` {
		return m
	}
	return Node{}
}

func findManager(nodes []Node) Node {
	for _, node := range nodes {
		if node.Role == `manager` {
			return node
		}
	}
	return Node{}
}
