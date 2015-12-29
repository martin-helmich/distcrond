package server

import (
	"fmt"
	"net/http"
	logging "github.com/op/go-logging"
	"github.com/martin-helmich/distcrond/container"
	"github.com/julienschmidt/httprouter"
	"github.com/martin-helmich/distcrond/storage"
	"time"
	"encoding/json"
)

type LinkResource struct {
	Href string `json:"href"`
	Rel string `json:"rel"`
}

type RootResource struct {
	Links []LinkResource `json:"links"`
}

type RestServer struct {
	server http.Server
	mux http.Handler

	nodes *container.NodeContainer
	jobs *container.JobContainer
	store storage.StorageBackend
	logger *logging.Logger

	root *RootResource
}

type SubHandler struct {
	server *RestServer
}

func (s *RestServer) decorate(handler httprouter.Handle) httprouter.Handle {
	return func(resp http.ResponseWriter, req *http.Request, par httprouter.Params) {
		start := time.Now()
		handler(resp, req, par)
		dur := time.Now().Sub(start)

		s.logger.Info("%s %s %s", req.Method, req.URL.Path, dur.String())
	}
}

func (h *RestServer) RootHandler(resp http.ResponseWriter, req *http.Request, param httprouter.Params) {
	root := h.root
	root.Links[0].Href = fmt.Sprintf("http://%s/jobs", req.Host)
	root.Links[1].Href = fmt.Sprintf("http://%s/nodes", req.Host)

	resp.Header().Set("Content-Type", "application/json")

	body, _ := json.Marshal(root)
	resp.Write(body)
}

func (h *RestServer) buildRootResource() {
	h.root = new(RootResource)
	h.root.Links = []LinkResource {
		LinkResource{"/jobs", "jobs"},
		LinkResource{"/nodes", "nodes"},
	}
}

func NewRestServer(port int, nodes *container.NodeContainer, jobs *container.JobContainer, store storage.StorageBackend, logger *logging.Logger) *RestServer {
	server := new(RestServer)
	server.nodes = nodes
	server.jobs = jobs
	server.logger = logger
	server.store = store
	server.buildRootResource()

	nodehandler := NodeHandler{server}
	jobhandler := JobHandler{server}
	reporthandler := ReportHandler{server}

	router := httprouter.New()
	router.GET("/", server.decorate(server.RootHandler))
	router.GET("/nodes", server.decorate(nodehandler.NodeList))
	router.GET("/nodes/:node", server.decorate(nodehandler.NodeSingle))
	router.GET("/jobs", server.decorate(jobhandler.JobList))
	router.GET("/jobs/:job", server.decorate(jobhandler.JobSingle))
	router.GET("/jobs/:job/reports", server.decorate(reporthandler.ReportsByJob))

	server.mux = router
	server.server = http.Server{
		Addr: fmt.Sprintf(":%d", port),
	}

//	server.mux.HandleFunc("/foo", func(resp http.ResponseWriter, req *http.Request) {
//		resp.Write([]byte("Hallo Welt!"))
//	})
//
//	server.mux.Handle("/nodes/", &NodeHandler{server})

	server.server.Handler = server.mux

	return server
}

func (s *RestServer) Start() {
	s.logger.Notice("Starting REST API server")

	if err := s.server.ListenAndServe(); err != nil {
		s.logger.Fatal(err)
	}
}
