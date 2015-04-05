package runner

import (
	"github.com/martin-helmich/distcrond/storage"
	"github.com/martin-helmich/distcrond/container"
	"fmt"
	"github.com/martin-helmich/distcrond/domain"
	"errors"
	"time"
	"sync/atomic"
)

type AnyJobRunner GenericJobRunner

func NewAnyJobRunner(nodes *container.NodeContainer, storage storage.StorageBackend, health HealthChecker) JobRunner {
	return &AnyJobRunner{nodes: nodes, storage: storage, healthChecker: health}
}

func (r *AnyJobRunner) Run(job *domain.Job) error {
	logger := job.Logger
	nodes := r.nodes.NodeCandidatesForJob(job)

	if len(nodes) == 0 {
		return errors.New(fmt.Sprintf("No nodes available for job %s", job.Name))
	}

	logger.Debug("Executing on one of %d nodes", len(nodes))

	job.Lock.Lock()
	defer job.Lock.Unlock()

	report := domain.RunReport{}
	report.Initialize(job, 1)

	reportItem := &report.Items[0]
	reportItem.Time.Start = time.Now()
	reportItem.Success = false
	reportItem.Output = "Could not find any node to run job on"

	for _, node := range nodes {
		done := func() bool {
			atomic.AddInt32(&node.RunningJobs, 1)
			defer func() {
				atomic.AddInt32(&node.RunningJobs, -1)
				reportItem.Time.Stop = time.Now()
			}()

			logger.Debug("Executing on node %s\n", node.Name)
			reportItem.Node = node

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
					return false

				default:
					logger.Error("%s", err)
					reportItem.Success = false
					reportItem.Output = err.Error()
				}
			}

			logger.Debug("Done on %s\n", node.Name)
			return true
		}()

		if done {
			break
		}
	}

	report.Finalize()
	job.LastExecution = time.Now()

	logger.Info("Report: %s\n", reportItem.Summary())

	go func() {
		if err := r.storage.SaveReport(&report); err != nil {
			logger.Error("%s", err)
		}
	}()

	logger.Info("%s: Done on all nodes", job.Name)

	return nil
}
