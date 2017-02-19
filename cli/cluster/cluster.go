package cluster

import (
	"errors"

	"github.com/bughou-go/xiaomei/config"
)

type Cluster struct {
	Managers     []*Node
	Workers      []*Node
	managerToken string
	workerToken  string
}

func Setup() error {
	if len(config.Cluster.Managers()) == 0 {
		return errors.New(`the cluster have no managers.`)
	}
	managers, workers, err := getClusterNodes()
	if err != nil {
		return err
	}
	return Cluster{Managers: managers, Workers: workers}.setup()
}

func (c Cluster) setup() error {
	if c.manager().Role == `` { // No Manager
		if err := c.Managers[0].init(); err != nil {
			return err
		}
	}
	if err := c.setupManagers(); err != nil {
		return err
	}
	return c.setupWorkers()
}

func (c Cluster) setupManagers() error {
	for _, m := range c.Managers {
		switch m.Role {
		case ``:
			if err := c.join(m, `manager`); err != nil {
				return err
			}
		case `worker`:
		}
	}
	return nil
}

func (c Cluster) setupWorkers() error {
	for _, w := range c.Workers {
		switch w.Role {
		case ``:
			if err := c.join(w, `worker`); err != nil {
				return err
			}
		case `manager`:
		}
	}
	return nil
}

func (c Cluster) join(node *Node, role string) error {
	if token, err := c.token(role); err != nil {
		return err
	} else {
		return node.join(role, token, c.manager().Config.Addr)
	}
}

func (c Cluster) token(role string) (string, error) {
	if role == `manager` {
		if c.managerToken != `` {
			return c.managerToken, nil
		}
	} else {
		if c.workerToken != `` {
			return c.workerToken, nil
		}
	}
	token, err := c.manager().token(role)
	if err == nil {
		if role == `manager` {
			c.managerToken = token
		} else {
			c.workerToken = token
		}
	}
	return token, err
}

func (c Cluster) manager() *Node {
	if m := findManager(c.Managers); m.Role != `` {
		return m
	}
	if m := findManager(c.Workers); m.Role != `` {
		return m
	}
	return &Node{}
}

func findManager(nodes []*Node) *Node {
	for _, node := range nodes {
		if node.Role == `manager` {
			return node
		}
	}
	return &Node{}
}
