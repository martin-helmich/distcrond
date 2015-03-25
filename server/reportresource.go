package server

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
//	"encoding/json"
	"fmt"
	"encoding/json"
)

type ReportHandler SubHandler

func (h *ReportHandler) ReportsByJob(resp http.ResponseWriter, req *http.Request, param httprouter.Params) {
	job, err := h.server.jobs.JobByName(param.ByName("job"));
	if err != nil {
		h.server.logger.Warning(fmt.Sprintf("Job %s reports requested, but not found", param.ByName("job")))
		resp.WriteHeader(404)
		return
	}

	reports, sErr := h.server.store.ReportsForJob(job)
	if sErr != nil {
		h.server.logger.Error(fmt.Sprintf("Reports for job %s could not be loaded: %s", job.Name, sErr))
		resp.WriteHeader(500)
		return
	}

	resp.Header().Set("Content-Type", "application/json")

	body, _ := json.Marshal(reports)
	resp.Write(body)
}
