package main

//import "log"
import "os"
import (
	"os/signal"
	"github.com/martin-helmich/distcrond/domain"
	"github.com/martin-helmich/distcrond/reader"
	"github.com/martin-helmich/distcrond/container"
	"github.com/martin-helmich/distcrond/scheduler"
	"github.com/martin-helmich/distcrond/runner"
	"github.com/martin-helmich/distcrond/logging"
//	logging "github.com/op/go-logging"
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
		os.Exit(1)
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

	jobRunner := runner.NewJobRunner(nodeContainer)
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
