package main

import "flag"
import (
	"os"
	"errors"
	"fmt"
)

const (
	STORAGE_ELASTICSEARCH = "es"
	STORAGE_REDIS = "redis"
)

type RuntimeConfig struct {
	jobsDirectory string
	nodesDirectory string
	allowNoOwner bool
	storageBackend string
	esHost string
	esPort int
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

func (c *RuntimeConfig) PopulateFromFlags() {
	flag.StringVar(&c.jobsDirectory, "jobsDirectory", "/etc/distcron/jobs.d", "Directory from which to load job definitions")
	flag.StringVar(&c.nodesDirectory, "nodesDirectory", "/etc/distcron/nodes.d", "Directory from which to load node definitions")
	flag.BoolVar(&c.allowNoOwner, "allowNoOwner", false, "Set to allow jobs to have no owners")
	flag.StringVar(&c.storageBackend, "storage", STORAGE_ELASTICSEARCH, "Which storage backend to use ('es' or 'redis')")
	flag.StringVar(&c.esHost, "esHost", "localhost", "Elasticsearch host")
	flag.IntVar(&c.esPort, "esPort", 9200, "Elasticsearch port")
	flag.Parse()
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

	if c.storageBackend != STORAGE_ELASTICSEARCH {
		return errors.New(fmt.Sprintf("Unknown storage backend '%s', must be '" + STORAGE_ELASTICSEARCH + "'", c.storageBackend))
	}

	return nil
}
