package container

import (
	"testing"
	"github.com/martin-helmich/distcrond/domain"
)

func assertThat(expr bool, e string, t *testing.T) {
	if expr == false {
		t.Error(e)
	}
}

func TestCountReturnsNodeCount(t *testing.T) {
	node := domain.Node{}
	nodes := NewNodeContainer(10)

	nodes.AddNode(node)

	if nodes.Count() != 1 {
		t.Error("Node count is not 1")
	}
}

func TestNodesForJobReturnsAllNodesFromRolePattern(t *testing.T) {
	nodes := []domain.Node{
		domain.Node{Name: "n1", Roles: []string{"web"}},
		domain.Node{Name: "n2", Roles: []string{"web"}},
		domain.Node{Name: "n3", Roles: []string{"db"}},
	}

	c := NewNodeContainer(3)
	c.AddNode(nodes[0])
	c.AddNode(nodes[1])
	c.AddNode(nodes[2])

	job := new(domain.Job)
	job.Policy.Hosts = domain.POLICY_ALL
	job.Policy.Roles = []string{"web"}

	selectedNodes := c.NodesForJob(job)

	if len(selectedNodes) != 2 {
		t.Error("Wrong node count")
		return
	}

	if selectedNodes[0].Name == "n1" && selectedNodes[1].Name == "n2" {
		t. Log("Correct nodes returned")
		return
	}

	if selectedNodes[1].Name == "n1" && selectedNodes[0].Name == "n2" {
		t. Log("Correct nodes returned")
		return
	}

	t.Error("Incorrect nodes returned")
}

func TestNodesForJobReturnsAnyNodeFromRolePattern(t *testing.T) {
	nodes := []domain.Node{
		domain.Node{Name: "n1", Roles: []string{"web"}},
		domain.Node{Name: "n2", Roles: []string{"web"}},
		domain.Node{Name: "n3", Roles: []string{"db"}},
	}

	c := NewNodeContainer(3)
	c.AddNode(nodes[0])
	c.AddNode(nodes[1])
	c.AddNode(nodes[2])

	job := new(domain.Job)
	job.Policy.Hosts = domain.POLICY_ANY
	job.Policy.Roles = []string{"web"}

	selectedNodes := c.NodesForJob(job)

	if len(selectedNodes) != 1 {
		t.Error("Wrong node count")
		return
	}

	if selectedNodes[0].Name == "n1"  {
		t.Log("Correct nodes returned")
		return
	}

	if selectedNodes[0].Name == "n2" {
		t.Log("Correct nodes returned")
		return
	}

	t.Error("Incorrect nodes returned")
}

func TestNodesForJobReturnsAllNodesFromNodeList(t *testing.T) {
	nodes := []domain.Node{
		domain.Node{Name: "n1", Roles: []string{"web"}},
		domain.Node{Name: "n2", Roles: []string{"web"}},
		domain.Node{Name: "n3", Roles: []string{"db"}},
	}

	c := NewNodeContainer(3)
	c.AddNode(nodes[0])
	c.AddNode(nodes[1])
	c.AddNode(nodes[2])

	job := new(domain.Job)
	job.Policy.Hosts = domain.POLICY_ALL
	job.Policy.HostList = []string{"n1", "n2"}

	selectedNodes := c.NodesForJob(job)

	if len(selectedNodes) != 2 {
		t.Error("Wrong node count")
		return
	}

	if selectedNodes[0].Name == "n1" && selectedNodes[1].Name == "n2" {
		t. Log("Correct nodes returned")
		return
	}

	if selectedNodes[1].Name == "n1" && selectedNodes[0].Name == "n2" {
		t. Log("Correct nodes returned")
		return
	}

	t.Error("Incorrect nodes returned")
}

func TestNodesForJobReturnsAnyNodesFromNodeList(t *testing.T) {
	nodes := []domain.Node{
		domain.Node{Name: "n1", Roles: []string{"web"}},
		domain.Node{Name: "n2", Roles: []string{"web"}},
		domain.Node{Name: "n3", Roles: []string{"db"}},
	}

	c := NewNodeContainer(3)
	c.AddNode(nodes[0])
	c.AddNode(nodes[1])
	c.AddNode(nodes[2])

	job := new(domain.Job)
	job.Policy.Hosts = domain.POLICY_ANY
	job.Policy.HostList = []string{"n1", "n2"}

	selectedNodes := c.NodesForJob(job)

	if len(selectedNodes) != 1 {
		t.Error("Wrong node count")
		return
	}

	if selectedNodes[0].Name == "n1" {
		t. Log("Correct nodes returned")
		return
	}

	if selectedNodes[0].Name == "n2" {
		t. Log("Correct nodes returned")
		return
	}

	t.Error("Incorrect nodes returned")
}

func TestRoleOverlappingHostsAreReturnedOnlyOnce(t *testing.T) {
	nodes := []domain.Node{
		domain.Node{Name: "n1", Roles: []string{"foo", "web"}},
		domain.Node{Name: "n2", Roles: []string{"foo", "db"}},
	}

	c := NewNodeContainer(2)
	c.AddNode(nodes[0])
	c.AddNode(nodes[1])

	job := new(domain.Job)
	job.Policy.Hosts = domain.POLICY_ALL
	job.Policy.Roles = []string{"web", "foo"}

	selectedNodes := c.NodesForJob(job)

	assertThat(len(selectedNodes) == 2, "Wrong node count", t)
}

func TestNodeByNameReturnsNodeByName(t *testing.T) {
	node := domain.Node{Name: "n1", Roles: []string{"web"}}
	nodes := NewNodeContainer(1)
	nodes.AddNode(node)

	foundNode, err := nodes.NodeByName("n1")
	assertThat(err == nil, "Unexpected error", t)
	assertThat(foundNode != nil, "Node not found", t)
	assertThat(foundNode.Name == "n1", "Wrong node returned", t)
}

func TestNodeByNameReturnsPointerToSameNode(t *testing.T) {
	node := domain.Node{Name: "n1", Roles: []string{"web"}}
	nodes := NewNodeContainer(1)
	nodes.AddNode(node)

	foundNode1, _ := nodes.NodeByName("n1")
	foundNode2, _ := nodes.NodeByName("n1")

	assertThat(foundNode1 == foundNode2, "Different pointers returned", t)
}

func TestNodesForJobReturnsUniformPointerToNode(t *testing.T) {
	nodes := []domain.Node{
		domain.Node{Name: "n1", Roles: []string{"web"}},
		domain.Node{Name: "n2", Roles: []string{"db"}},
	}

	c := NewNodeContainer(2)
	c.AddNode(nodes[0])
	c.AddNode(nodes[1])

	job := new(domain.Job)
	job.Policy.Hosts = domain.POLICY_ANY
	job.Policy.HostList = []string{"n1"}

	selectedNodes := c.NodesForJob(job)
	realNode, _ := c.NodeByName("n1")

	if len(selectedNodes) != 1 {
		t.Error("Wrong node count")
		return
	}

	assertThat(selectedNodes[0] == realNode, "Different pointers returned", t)
}
