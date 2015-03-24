package server

import (
	"net/http"
	"encoding/json"
	"github.com/martin-helmich/distcrond/domain"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"sync/atomic"
)

type NodeHandler SubHandler

type NodeResource struct {
	Name string `json:"name"`
	Href string `json:"href"`
	Roles []string `json:"roles"`
	Status string `json:"status"`
	RunningJobs int32 `json:"running_jobs"`
}

func (h *NodeHandler) resourceFromNode(node *domain.Node, res *NodeResource, host string) {
	res.Name = node.Name
	res.Href = fmt.Sprintf("http://%s/nodes/%s", host, node.Name)
	res.Roles = node.Roles
	res.RunningJobs = atomic.LoadInt32(&node.RunningJobs)

	switch node.Status {
	case domain.STATUS_DOWN:
		res.Status = "down"
	default:
		res.Status = "up"
	}
}

func (h *NodeHandler) NodeList(resp http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	nodeCount := h.server.nodes.Count()
	nodeResources := make([]NodeResource, nodeCount)
	for i := 0; i < nodeCount; i ++ {
		node := h.server.nodes.Get(i)
		h.resourceFromNode(node, &nodeResources[i], req.Host)
	}

	jsonBody, _ := json.MarshalIndent(nodeResources, "", "  ")

	resp.Header().Set("Content-Type", "application/json")
	resp.Write(jsonBody)
}

func (h *NodeHandler) NodeSingle(resp http.ResponseWriter, req *http.Request, param httprouter.Params) {
	if node, err := h.server.nodes.NodeByName(param.ByName("node")); err != nil {
		resp.WriteHeader(404)
	} else {
		res := NodeResource{}
		h.resourceFromNode(node, &res, req.Host)

		jsonBody, _ := json.MarshalIndent(res, "", "  ")

		resp.Header().Set("Content-Type", "application/json")
		resp.Write(jsonBody)
	}
}
