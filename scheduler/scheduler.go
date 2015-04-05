package scheduler

import (
	"time"
	"github.com/martin-helmich/distcrond/container"
	. "github.com/martin-helmich/distcrond/domain"
	"github.com/martin-helmich/distcrond/runner"
	"github.com/martin-helmich/distcrond/logging"
	"sync/atomic"
)

type Scheduler struct {
	jobContainer *container.JobContainer
	nodeContainer *container.NodeContainer
	runner runner.JobRunner
	abort chan bool

	Done chan bool
}

func NewScheduler(jobs *container.JobContainer, nodes *container.NodeContainer, runner runner.JobRunner) *Scheduler {
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

func (s *Scheduler) nextRunDate(job *Job, now time.Time) time.Time {
	reference := job.Schedule.Reference
	todayReference := time.Date(now.Year(), now.Month(), now.Day(), reference.Hour(), reference.Minute(), reference.Second(), reference.Nanosecond(), reference.Location())

	if todayReference.Before(now) {
		for todayReference.Before(now) {
			todayReference = todayReference.Add(job.Schedule.Interval)
		}
	} else {
		for todayReference.After(now) {
			todayReference = todayReference.Add(-job.Schedule.Interval)
		}
		todayReference = todayReference.Add(job.Schedule.Interval)
	}

	return todayReference
}

func (s *Scheduler) Run() {
	logging.Info("Starting scheduler")

	var count      int               = s.jobContainer.Count()
	var semaphores []chan bool       = make([]chan bool, count)
	var tickers    chan *time.Ticker = make(chan *time.Ticker, count)
	var now        time.Time         = time.Now()

	var startedTickers int64 = 0

	withLock := func(f func(), i int) {
		semaphores[i] <- true
		f()
		<-semaphores[i]
	}

	var start []time.Time = make([]time.Time, count)
	for i := 0; i < count; i ++ {
		start[i] = s.nextRunDate(s.jobContainer.Get(i), now)
	}

	for i := 0; i < count; i ++ {
		job := s.jobContainer.Get(i)
		semaphores[i] = make(chan bool, 1)
		go func(job *Job, i int) {
			wait := start[i].Sub(now)
			logging.Debug("Next execution of %s scheduled for %s, waiting %s", job.Name, start[i].String(), wait.String())
			<- time.After(wait)

			atomic.AddInt64(&startedTickers, 1)

			ticker := time.NewTicker(job.Schedule.Interval)
			tickers <- ticker

			logging.Debug("Started timer for %s", job.Name)

			for t := range ticker.C {
				withLock(func() {
					logging.Notice("Executing job %s at %s", job.Name, t)
					s.runner.Run(job)
				}, i)
			}
		}(job, i)
	}

	select {
	case <- s.abort:
		logging.Notice("Aborting")

		logging.Debug("Stopping tickers...")
		for i := atomic.LoadInt64(&startedTickers); i > 0; i -- {
			(<- tickers).Stop()
		}

		logging.Notice("Waiting for running jobs...")
		for i := 0; i < count; i ++ {
			semaphores[i] <- true
		}

		logging.Debug("Done")
		s.Done <- true
	}
}
