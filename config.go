package main

import "flag"
import (
	"os"
	"errors"
	"fmt"
	"time"
)

const (
	STORAGE_PLAINFILES = "plain"
	STORAGE_ELASTICSEARCH = "es"
	STORAGE_REDIS = "redis"
)

type RuntimeConfig struct {
	// General configuration
	jobsDirectory string
	nodesDirectory string
	allowNoOwner bool
	storageBackend string
	healthCheckInterval time.Duration

	// Elasticsearch storage backend
	esHost string
	esPort int

	// Plainfiles storage backend
	pfPath string

	// Profiling configuration
	cpuprofile string
	memprofile string
}

func (c *RuntimeConfig) JobsDirectory() string {
	return c.jobsDirectory
}

func (c *RuntimeConfig) NodesDirectory() string {
	return c.nodesDirectory
}

func (c *RuntimeConfig) AllowNoOwner() bool {
	return c.allowNoOwner
}

func (c *RuntimeConfig) StorageBackend() string {
	return c.storageBackend
}

func (c *RuntimeConfig) ElasticSearchHost() string {
	return c.esHost
}

func (c *RuntimeConfig) ElasticSearchPort() int {
	return c.esPort
}

func (c *RuntimeConfig) LogDirectory() string {
	return c.pfPath
}

func (c *RuntimeConfig) CpuProfilingEnabled() bool {
	return c.cpuprofile != ""
}

func (c *RuntimeConfig) CpuProfilingTarget() string {
	return c.cpuprofile
}

func (c *RuntimeConfig) MemProfilingEnabled() bool {
	return c.memprofile != ""
}

func (c *RuntimeConfig) MemProfilingTarget() string {
	return c.memprofile
}

func (c *RuntimeConfig) HealthCheckInterval() time.Duration {
	return c.healthCheckInterval
}

func (c *RuntimeConfig) PopulateFromFlags() error {
	var healthCheckInterval string
	var err error

	flag.StringVar(&c.jobsDirectory, "jobsDirectory", "/etc/distcron/jobs.d", "Directory from which to load job definitions")
	flag.StringVar(&c.nodesDirectory, "nodesDirectory", "/etc/distcron/nodes.d", "Directory from which to load node definitions")
	flag.BoolVar(&c.allowNoOwner, "allowNoOwner", false, "Set to allow jobs to have no owners")
	flag.StringVar(&c.storageBackend, "storage", STORAGE_ELASTICSEARCH, "Which storage backend to use ('es' or 'plain')")
	flag.StringVar(&healthCheckInterval, "healthCheckInterval", "10s", "Interval in which to check node health")

	flag.StringVar(&c.esHost, "esHost", "localhost", "Elasticsearch host")
	flag.IntVar(&c.esPort, "esPort", 9200, "Elasticsearch port")

	flag.StringVar(&c.pfPath, "logDirectory", "/var/log/distcrond", "Directory to write log files to (for 'plain' storage backend')")

	flag.StringVar(&c.cpuprofile, "cpuprofile", "", "Write CPU profile to file")
	flag.StringVar(&c.memprofile, "memprofile", "", "Write Memory profile to file")

	flag.Parse()

	if c.healthCheckInterval, err = time.ParseDuration(healthCheckInterval); err != nil {
		return err
	}

	return nil
}

func (c *RuntimeConfig) IsValid() error {
	checkDir := func(dir string, purpose string) error {
		if _, err := os.Stat(dir); err != nil {
			if os.IsNotExist(err) {
				return errors.New(fmt.Sprintf("Non-existent %s %s!", purpose, dir))
			} else {
				return errors.New(fmt.Sprintf("Unknown error related to %s %s: %s", purpose, dir, err))
			}
		}
		return nil
	}

	if err := checkDir(c.jobsDirectory, "job configuration directory"); err != nil {
		return err
	}

	if err := checkDir(c.nodesDirectory, "node configuration directory"); err != nil {
		return err
	}

	switch c.storageBackend {
	case STORAGE_ELASTICSEARCH:
		if c.esHost != "" {
			return errors.New("No Elasticsearch host specified")
		}
	case STORAGE_PLAINFILES:
		if err := checkDir(c.pfPath, "log files target directory"); err != nil {
			return err
		}
	default:
		return errors.New(fmt.Sprintf("Unknown storage backend '%s', must be '" + STORAGE_ELASTICSEARCH + "'", c.storageBackend))
	}

	return nil
}
