package server

import (
	"github.com/martin-helmich/distcrond/domain"
	"net/http"
	"github.com/julienschmidt/httprouter"
	"encoding/json"
	"github.com/martin-helmich/distcrond/container"
	"fmt"
	"time"
)

type JobHandler SubHandler

type JobOwnerResource struct {
	Name string `json:"name"`
	EmailAddress string `json:"email_address"`
}

type RoleExecutionPolicyResource struct {
	Hosts string `json:"select"`
	Roles []string `json:"roles"`
}

type HostsExecutionPolicyResource struct {
	Hosts string `json:"select"`
	HostList []string `json:"hosts"`
}

type DateResource struct {
	Timestamp int64 `json:"timestamp"`
	String string `json:"string"`
}

type JobResource struct {
	Name string `json:"name"`
	Href string `json:"href"`
	Links [1]LinkResource `json:"links"`
	Description string `json:"description"`
	Owners []JobOwnerResource `json:"owners"`
	Policy interface {} `json:"execution_policy"`
	Schedule string `json:"execution_schedule"`
	Command []string `json:"command"`
	LastExecution *DateResource `json:"last_execution"`
	NextExecution *DateResource `json:"next_execution"`
}

func (h *JobHandler) resourceFromJob(job *domain.Job, res *JobResource, host string) {
	job.Lock.RLock()
	defer job.Lock.RUnlock()

	res.Name = job.Name
	res.Href = fmt.Sprintf("http://%s/jobs/%s", host, job.Name)
	res.Description = job.Description

	res.Schedule = job.ScheduleSpec

	res.Links[0].Href = fmt.Sprintf("http://%s/jobs/%s/reports", host, job.Name)
	res.Links[0].Rel = "reports"

	res.Command = job.Command.Command()

	res.Owners = make([]JobOwnerResource, len(job.Owners))
	for i, owner := range job.Owners {
		res.Owners[i].Name = owner.Name
		res.Owners[i].EmailAddress = owner.EmailAddress
	}

	if !job.LastExecution.IsZero() {
		res.LastExecution = &DateResource{
			Timestamp: job.LastExecution.UnixNano(),
			String: job.LastExecution.String(),
		}
	} else {
		res.LastExecution = nil
	}

	next := job.Schedule.Next(time.Now())
	res.NextExecution = &DateResource{
		Timestamp: next.UnixNano(),
		String: next.String(),
	}

	if len(job.Policy.Roles) > 0 {
		res.Policy = RoleExecutionPolicyResource{job.Policy.Hosts, job.Policy.Roles}
	} else {
		res.Policy = HostsExecutionPolicyResource{job.Policy.Hosts, job.Policy.HostList}
	}
}

func (h *JobHandler) JobList(resp http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	var jobs *container.JobContainer = h.server.jobs

	jobCount := jobs.Count()
	jobResources := make([]JobResource, jobCount)
	for i := 0; i < jobCount; i ++ {
		job := jobs.Get(i)
		h.resourceFromJob(job, &jobResources[i], req.Host)
	}

	jsonBody, _ := json.MarshalIndent(jobResources, "", "  ")

	resp.Header().Set("Content-Type", "application/json")
	resp.Write(jsonBody)
}

func (h *JobHandler) JobSingle(resp http.ResponseWriter, req *http.Request, params httprouter.Params) {
	var jobs *container.JobContainer = h.server.jobs

	if job, err := jobs.JobByName(params.ByName("job")); err != nil {
		resp.WriteHeader(404)
	} else {
		res := JobResource{}
		h.resourceFromJob(job, &res, req.Host)

		jsonBody, _ := json.MarshalIndent(res, "", "  ")

		resp.Header().Set("Content-Type", "application/json")
		resp.Write(jsonBody)
	}
}
