package scheduler

import (
	"log"
	"time"
	"github.com/martin-helmich/distcrond/container"
	"github.com/martin-helmich/distcrond/domain"
	"github.com/martin-helmich/distcrond/runner"
)

type Scheduler struct {
	jobContainer *container.JobContainer
	nodeContainer *container.NodeContainer
	runner *runner.JobRunner
	abort chan bool

	Done chan bool
}

func NewScheduler(jobs *container.JobContainer, nodes *container.NodeContainer, runner *runner.JobRunner) *Scheduler {
	return &Scheduler {
		jobs,
		nodes,
		runner,
		make(chan bool),
		make(chan bool),
	}
}

func (s *Scheduler) Abort() {
	s.abort <- true
}

func (s *Scheduler) Run() {
	log.Println("Starting scheduler")

	var count int = s.jobContainer.Count()
	var semaphores []chan bool = make([]chan bool, count)
	var tickers []*time.Ticker = make([]*time.Ticker, count)

	withLock := func(f func(), i int) {
		semaphores[i] <- true
		f()
		<-semaphores[i]
	}

	for i := 0; i < count; i ++ {
		job := s.jobContainer.Get(i)
		semaphores[i] = make(chan bool, 1)
		go func(job *domain.Job, i int) {
			tickers[i] = time.NewTicker(job.Schedule.Interval)
			for t := range tickers[i].C {
				withLock(func() {
					log.Printf("Executing job %s at %s", job.Name, t)
					s.runner.Run(job)
					job.LastExecution = time.Now()
				}, i)
			}
		}(job, i)
	}

	select {
	case <- s.abort:
		log.Println("Aborting")

		log.Println("Stopping tickers...")
		for i := 0; i < count; i ++ {
			tickers[i].Stop()
		}

		log.Println("Waiting for running jobs...")
		for i := 0; i < count; i ++ {
			semaphores[i] <- true
		}

		log.Println("Done")
		s.Done <- true
	}
}
