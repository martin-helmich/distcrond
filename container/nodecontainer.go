package container

import (
	"errors"
	"fmt"
	"github.com/martin-helmich/distcrond/domain"
	"github.com/martin-helmich/distcrond/logging"
	"math/rand"
)

type NodeContainer struct {
	nodes       []domain.Node
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

func (c *NodeContainer) Get(i int) *domain.Node {
	return &c.nodes[i]
}

func (c *NodeContainer) NodeByName(name string) (*domain.Node, error) {
	if node, ok := c.nodesByName[name]; ok {
		return node, nil
	} else {
		return nil, errors.New(fmt.Sprintf("No node with name %s is known", name))
	}
}

func (c *NodeContainer) NodeCandidatesForJob(job *domain.Job) []*domain.Node {
	nodes := c.potentialNodesForJob(job, true)

	// Fisher-Yates shuffle
	for i, _ := range nodes {
		j := rand.Intn(i + 1)
		nodes[i], nodes[j] = nodes[j], nodes[i]
	}

	return nodes
}

func (c *NodeContainer) NodesForJob(job *domain.Job) []*domain.Node {
	switch job.Policy.Hosts {
	case domain.POLICY_ALL:
		return c.potentialNodesForJob(job, false)

	case domain.POLICY_ANY:
		potentialNodes := c.potentialNodesForJob(job, true)
		logging.Debug("Found %d potential nodes for job %s: %s", len(potentialNodes), job.Name, potentialNodes)

		idx := rand.Int() % len(potentialNodes)
		return potentialNodes[idx : idx+1]

	default:
		logging.Error("Invalid job policy: %s", job.Policy.Hosts)
		return make([]*domain.Node, 0)
	}
}

func (c *NodeContainer) NodesWithStatus(status domain.NodeStatus) []*domain.Node {
	return c.NodesByFilter(func(n *domain.Node) bool {
		return n.Status == status
	})
}

func (c *NodeContainer) NodesByFilter(filter (func(*domain.Node) bool)) []*domain.Node {
	var nodes []*domain.Node = make([]*domain.Node, 0, len(c.nodes))

	for _, node := range c.nodes {
		if filter(&node) {
			nodes = append(nodes, &node)
		}
	}

	return nodes
}

func (c *NodeContainer) potentialNodesForJob(job *domain.Job, onlyHealthyNodes bool) []*domain.Node {
	var filter func(*domain.Node) bool = func(_ *domain.Node) bool { return true }
	var nodes []*domain.Node = make([]*domain.Node, 0, len(c.nodes))

	if onlyHealthyNodes {
		filter = func(node *domain.Node) bool {
			return node.Status == domain.STATUS_UP
		}
	}

	if len(job.Policy.HostList) > 0 {
		for _, name := range job.Policy.HostList {
			if node, ok := c.nodesByName[name]; ok && filter(node) {
				nodes = append(nodes, node)
			}
		}
	} else if len(job.Policy.Roles) > 0 {
		knownHosts := make(map[string]bool)
		for _, role := range job.Policy.Roles {
			if hostsWithRole, ok := c.nodesByRole[role]; ok {
				for _, node := range hostsWithRole {
					if _, known := knownHosts[node.Name]; known == false && filter(node) {
						nodes = append(nodes, node)
						knownHosts[node.Name] = true
					}
				}
			}
		}
	}

	return nodes
}
