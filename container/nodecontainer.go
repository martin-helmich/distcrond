package container

import (
	"errors"
	"math/rand"
	"github.com/martin-helmich/distcrond/domain"
	"github.com/martin-helmich/distcrond/logging"
	"fmt"
)

type NodeContainer struct {
	nodes []domain.Node
	nodesByName map[string]*domain.Node
	nodesByRole map[string][]*domain.Node
}

func NewNodeContainer(initialCapacity int) *NodeContainer {
	container := new(NodeContainer)
	container.nodes = make([]domain.Node, 0, initialCapacity)
	container.nodesByName = make(map[string]*domain.Node)
	container.nodesByRole = make(map[string][]*domain.Node)
	return container
}

func (c *NodeContainer) AddNode(node domain.Node) {
	c.nodes = append(c.nodes, node)
	c.nodesByName[node.Name] = &c.nodes[len(c.nodes)-1]

	for _, role := range node.Roles {
		if _, ok := c.nodesByRole[role]; ok == false {
			c.nodesByRole[role] = make([]*domain.Node, 0, 3)
		}
		c.nodesByRole[role] = append(c.nodesByRole[role], &c.nodes[len(c.nodes)-1])
	}
}

func (c *NodeContainer) Count() int {
	return len(c.nodes)
}

func (c *NodeContainer) NodeByName(name string) (*domain.Node, error) {
	if node, ok := c.nodesByName[name]; ok {
		return node, nil
	} else {
		return nil, errors.New(fmt.Sprintf("No node with name %s is known", name))
	}
}

func (c *NodeContainer) NodesForJob(job *domain.Job) []*domain.Node {
	potentialNodes := c.potentialNodesForJob(job)

	logging.Debug("Found %d potential nodes for job %s: %s", len(potentialNodes), job.Name, potentialNodes)

	if job.Policy.Hosts == domain.POLICY_ALL {
		return potentialNodes
	} else {
		idx := rand.Int() % len(potentialNodes)
		return potentialNodes[idx:idx+1]
	}
}

func (c *NodeContainer) potentialNodesForJob(job *domain.Job) []*domain.Node {
	nodes := make([]*domain.Node, 0, len(c.nodes))

	if len(job.Policy.HostList) > 0 {
		for _, name := range job.Policy.HostList {
			if node, ok := c.nodesByName[name]; ok {
				nodes = append(nodes, node)
			}
		}
	} else if len(job.Policy.Roles) > 0 {
		knownHosts := make(map[string]bool)
		for _, role := range job.Policy.Roles {
			if hostsWithRole, ok := c.nodesByRole[role]; ok {
				for _, node := range hostsWithRole {
					if _, known := knownHosts[node.Name]; known == false {
						nodes = append(nodes, node)
						knownHosts[node.Name] = true
					}
				}
			}
		}
	}

	return nodes
}
