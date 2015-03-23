package runner

import (
	"github.com/martin-helmich/distcrond/container"
	"github.com/martin-helmich/distcrond/domain"
	"github.com/martin-helmich/distcrond/storage"
	"errors"
	"fmt"
)

type JobRunner struct {
	nodes *container.NodeContainer
	storage storage.StorageBackend
}

func NewJobRunner(nodes *container.NodeContainer, storage storage.StorageBackend) *JobRunner {
	return &JobRunner{nodes, storage}
}

func (r *JobRunner) Run(job *domain.Job) error {
	logger := job.Logger
	nodes  := r.nodes.NodesForJob(job)

	if len(nodes) == 0 {
		return errors.New(fmt.Sprintf("No nodes available for job %s", job.Name))
	}

	done := make(chan bool, len(nodes))

	logger.Debug("Executing on %d nodes", len(nodes))
	for _, node := range nodes {
		go func(node *domain.Node) {
			logger.Debug("Executing on node %s\n", node.Name)

			report := domain.RunReport{}
			report.Initialize(job, node)

			strat, _ := GetStrategyForNode(node, job.Logger)
			if err := strat.ExecuteCommand(job.Command, &report); err != nil {
				logger.Error("%s", err)
			}

			logger.Debug("Done on %s\n", node.Name)
			logger.Info("Report: %s\n", report.Summary())

			go func() {
				if err := r.storage.SaveReport(job, &report); err != nil {
					logger.Error("%s", err)
				}
			}()

			done <- true
		}(node)
	}

	for i := 0; i < len(nodes); i ++ {
		<- done
	}

	logger.Info("%s: Done on all nodes", job.Name)

	return nil
}
