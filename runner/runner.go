package runner

import (
	"github.com/martin-helmich/distcrond/container"
	"github.com/martin-helmich/distcrond/domain"
	"github.com/martin-helmich/distcrond/storage"
	"errors"
	"fmt"
	"time"
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

	report := domain.RunReport{}
	report.Initialize(job, len(nodes))

	for i, node := range nodes {
		go func(node *domain.Node, reportItem *domain.RunReportItem) {
			logger.Debug("Executing on node %s\n", node.Name)

			reportItem.Time.Start = time.Now()

			strat := node.ExecutionStrategy

			if err := strat.ExecuteCommand(job.Command, reportItem, job.Logger); err != nil {
				logger.Error("%s", err)
			}

			reportItem.Time.Stop = time.Now()

			logger.Debug("Done on %s\n", node.Name)
			logger.Info("Report: %s\n", reportItem.Summary())

			done <- true
		}(node, &report.Items[i])
	}

	for i := 0; i < len(nodes); i ++ {
		<- done
	}

	report.Finalize()

	go func() {
		if err := r.storage.SaveReport(&report); err != nil {
			logger.Error("%s", err)
		}
	}()

	logger.Info("%s: Done on all nodes", job.Name)

	return nil
}
