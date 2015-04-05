package runner

import (
	"github.com/martin-helmich/distcrond/domain"
	"github.com/martin-helmich/distcrond/container"
	"github.com/martin-helmich/distcrond/storage"
)

type JobRunner interface {
	Run(job *domain.Job) error
}

type GenericJobRunner struct {
	nodes *container.NodeContainer
	storage storage.StorageBackend
}

type DispatchingRunner struct {
	allRunner JobRunner
	anyRunner JobRunner
}

func NewDispatchingRunner(nodes *container.NodeContainer, storage storage.StorageBackend) *DispatchingRunner {
	return &DispatchingRunner{
		allRunner: NewAllJobRunner(nodes, storage),
		anyRunner: NewAnyJobRunner(nodes, storage),
	}
}

func (d *DispatchingRunner) Run(job *domain.Job) error {
	switch job.Policy.Hosts {
	case domain.POLICY_ALL:
		return d.allRunner.Run(job)

	case domain.POLICY_ANY:
		return d.anyRunner.Run(job)
	}

	return nil
}
