package container

import (
	"errors"
	"fmt"
	"github.com/martin-helmich/distcrond/domain"
)

type JobContainer struct {
	jobs       []domain.Job
	jobsByName map[string]*domain.Job
}

func NewJobContainer(initialCapacity int) *JobContainer {
	container := new(JobContainer)
	container.jobs = make([]domain.Job, 0, initialCapacity)
	container.jobsByName = make(map[string]*domain.Job)
	return container
}

func (c *JobContainer) AddJob(job domain.Job) {
	c.jobs = append(c.jobs, job)
	c.jobsByName[job.Name] = &c.jobs[len(c.jobs)-1]
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

func (c *JobContainer) JobByName(n string) (*domain.Job, error) {
	if job, ok := c.jobsByName[n]; !ok {
		return nil, errors.New(fmt.Sprintf("No job with name '%s' is known", n))
	} else {
		return job, nil
	}
}
