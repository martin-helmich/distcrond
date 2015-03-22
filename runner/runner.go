package runner

import (
	"log"
	"github.com/martin-helmich/distcrond/container"
	"github.com/martin-helmich/distcrond/domain"
	"errors"
	"fmt"
)

type JobRunner struct {
	nodes *container.NodeContainer
}

func NewJobRunner(nodes *container.NodeContainer) *JobRunner {
	return &JobRunner{nodes}
}

func (r *JobRunner) Run(job *domain.Job) error {
	nodes := r.nodes.NodesForJob(job)

	if len(nodes) == 0 {
		return errors.New(fmt.Sprintf("No nodes available for job %s", job.Name))
	}

	done := make(chan bool, len(nodes))

	log.Printf("%s: Executing on %d nodes", job.Name, len(nodes))
	for _, node := range nodes {
		go func(node *domain.Node) {
			log.Printf("%s: Executing on node %s\n", job.Name, node.Name)

			report := RunReport{}

			strat, _ := GetStrategyForNode(node)
			strat.ExecuteCommand(job.Command, &report)

			log.Printf("%s: Done on %s\n", job.Name, node.Name)
			log.Printf("%s: Report: %s\n", job.Name, report)

			done <- true
		}(node)
	}

	for i := 0; i < len(nodes); i ++ {
		<- done
	}

	log.Printf("%s: Done on all nodes", job.Name)

	return nil
}
