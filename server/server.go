package server

import (
	"fmt"
	"net/http"
	logging "github.com/op/go-logging"
	"github.com/martin-helmich/distcrond/container"
	"github.com/julienschmidt/httprouter"
	"github.com/martin-helmich/distcrond/storage"
)

type LinkResource struct {
	Href string `json:"href"`
	Rel string `json:"rel"`
}

type RestServer struct {
	server http.Server
	mux http.Handler

	nodes *container.NodeContainer
	jobs *container.JobContainer
	store storage.StorageBackend
	logger *logging.Logger
}

type SubHandler struct {
	server *RestServer
}

func NewRestServer(port int, nodes *container.NodeContainer, jobs *container.JobContainer, store storage.StorageBackend, logger *logging.Logger) *RestServer {
	server := new(RestServer)
	server.nodes = nodes
	server.jobs = jobs
	server.logger = logger
	server.store = store

	nodehandler := NodeHandler{server}
	jobhandler := JobHandler{server}
	reporthandler := ReportHandler{server}

	router := httprouter.New()
	router.GET("/nodes", nodehandler.NodeList)
	router.GET("/nodes/:node", nodehandler.NodeSingle)
	router.GET("/jobs", jobhandler.JobList)
	router.GET("/jobs/:job", jobhandler.JobSingle)
	router.GET("/jobs/:job/reports", reporthandler.ReportsByJob)

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
