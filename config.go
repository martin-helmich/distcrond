package main

import "flag"
import (
	"os"
	"errors"
	"fmt"
)

type RuntimeConfig struct {
	jobsDirectory string
	nodesDirectory string
	allowNoOwner bool
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

func (c *RuntimeConfig) PopulateFromFlags() {
	flag.StringVar(&c.jobsDirectory, "jobsDirectory", "/etc/distcron/jobs.d", "Directory from which to load job definitions")
	flag.StringVar(&c.nodesDirectory, "nodesDirectory", "/etc/distcron/nodes.d", "Directory from which to load node definitions")
	flag.BoolVar(&c.allowNoOwner, "allowNoOwner", false, "Set to allow jobs to have no owners")
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

	return nil
}
