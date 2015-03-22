package container

import (
	"log"
	"github.com/martin-helmich/distcrond/domain"
)

type JobContainer struct {
	jobs []domain.Job
}

func NewJobContainer(initialCapacity int) *JobContainer {
	container := new(JobContainer)
	container.jobs = make([]domain.Job, 0, initialCapacity)
	return container
}

func (c *JobContainer) AddJob(job domain.Job) {
	c.jobs = append(c.jobs, job)
	log.Printf("Added job %s\n", job.Name)
}

func (c *JobContainer) Count() int {
	return len(c.jobs)
}

func (c *JobContainer) All() []domain.Job {
	return c.jobs
}

func (c *JobContainer) Get(i int) *domain.Job {
	return &c.jobs[i]
}
