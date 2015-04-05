package runner

import (
	"github.com/martin-helmich/distcrond/container"
	"github.com/martin-helmich/distcrond/domain"
	"github.com/martin-helmich/distcrond/storage"
	"errors"
	"fmt"
	"time"
	"sync/atomic"
)

type AllJobRunner GenericJobRunner

func NewAllJobRunner(nodes *container.NodeContainer, storage storage.StorageBackend, health HealthChecker) JobRunner {
	return &AllJobRunner{nodes: nodes, storage: storage, healthChecker: health}
}

func (r *AllJobRunner) Run(job *domain.Job) error {
	logger := job.Logger
	nodes  := r.nodes.NodesForJob(job)

	if len(nodes) == 0 {
		return errors.New(fmt.Sprintf("No nodes available for job %s", job.Name))
	}

	done := make(chan bool, len(nodes))
	logger.Debug("Executing on %d nodes", len(nodes))

	job.Lock.Lock()
	defer job.Lock.Unlock()

	report := domain.RunReport{}
	report.Initialize(job, len(nodes))

	for i, node := range nodes {
		go func(node *domain.Node, reportItem *domain.RunReportItem) {
			logger.Debug("Executing on node %s\n", node.Name)

			reportItem.Node = node
			reportItem.Time.Start = time.Now()
			atomic.AddInt32(&node.RunningJobs, 1)

			strat := node.ExecutionStrategy
			if err := strat.ExecuteCommand(job, reportItem); err != nil {
				switch err.(type) {
				case NodeDownError:
					func() {
						node.Lock.Lock();
						defer node.Lock.Unlock()

						logger.Warning("Node %s is down.", node.Name)
						node.Status = domain.STATUS_DOWN
					}()
					r.healthChecker.ScheduleHealthCheck(node)
				}

				logger.Error("%s", err)
				reportItem.Success = false
				reportItem.Output = err.Error()
			}

			atomic.AddInt32(&node.RunningJobs, -1)
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
	job.LastExecution = time.Now()

	go func() {
		if err := r.storage.SaveReport(&report); err != nil {
			logger.Error("%s", err)
		}
	}()

	logger.Info("%s: Done on all nodes", job.Name)

	return nil
}
