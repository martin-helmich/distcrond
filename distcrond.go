package main

import (
	"os"
	"os/signal"
	"github.com/martin-helmich/distcrond/domain"
	"github.com/martin-helmich/distcrond/reader"
	"github.com/martin-helmich/distcrond/container"
	"github.com/martin-helmich/distcrond/scheduler"
	"github.com/martin-helmich/distcrond/runner"
	"github.com/martin-helmich/distcrond/logging"
	"github.com/martin-helmich/distcrond/storage"
	"runtime/pprof"
)

var runtimeConfig *RuntimeConfig

type JobContainer struct {
	Jobs []domain.Job
}

func main() {
	runtimeConfig = new(RuntimeConfig)
	runtimeConfig.PopulateFromFlags()

	logging.Setup()
	log := logging.Logger

	if err := runtimeConfig.IsValid(); err != nil {
		log.Fatal(err)
	}

	if runtimeConfig.CpuProfilingEnabled() {
		f, err := os.Create(runtimeConfig.CpuProfilingTarget())
		if err != nil {
			log.Fatal(err)
		}

		log.Info("Starting CPU profiling")

		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	nodeContainer := container.NewNodeContainer(5)
	jobContainer := container.NewJobContainer(10)

	nodesLoaded, jobsLoaded := make(chan bool), make(chan bool)

	go func() {
		nodeReader := reader.NewNodeReader(nodeContainer)
		if err := nodeReader.ReadFromDirectory(runtimeConfig.NodesDirectory()); err != nil {
			log.Fatal(err)
		}
		nodesLoaded <- true
	}()

	go func() {
		jobReader := reader.NewJobReader(runtimeConfig, jobContainer)
		if err := jobReader.ReadFromDirectory(runtimeConfig.JobsDirectory()); err != nil {
			log.Fatal(err)
		}
		jobsLoaded <- true
	}()

	<-nodesLoaded
	<-jobsLoaded

	storageBackend, err := storage.BuildStorageBackend(runtimeConfig)
	if err != nil {
		log.Fatal(err)
	}

	if err := storageBackend.Connect(); err != nil {
		log.Fatal(err)
	}

	jobRunner := runner.NewJobRunner(nodeContainer, storageBackend)
	jobScheduler := scheduler.NewScheduler(jobContainer, nodeContainer, jobRunner)
	go jobScheduler.Run()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		<- c
		log.Notice("Received SIGINT")
		jobScheduler.Abort()
	}()

	<-jobScheduler.Done
}
