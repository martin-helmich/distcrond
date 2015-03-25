package domain

import (
	"fmt"
	"errors"
	"time"
	logging "github.com/op/go-logging"
	"sync"
)

type JobValidationConfig interface {
	AllowNoOwner() bool
}

type JobJson struct {
	Description string `json:"description"`
	Owners []JobOwnerJson `json:"owners"`
	Policy ExecutionPolicyJson `json:"policy"`
	Schedule ScheduleJson `json:"schedule"`
	ShellCommand string `json:"shell_command"`
	Command []string `json:"command"`
	Environment map[string]string `json:"environment"`
}

type Job struct {
	// Domain properties
	Name string
	Description string
	Owners []JobOwner
	Policy ExecutionPolicy
	Schedule Schedule
	Command Command
	LastExecution time.Time
	Environment map[string]string

	// Auxiliary properties
	Logger *logging.Logger
	Lock sync.RWMutex
}

func NewJobFromJson(name string, json JobJson) (Job, error) {
	owners := make([]JobOwner, len(json.Owners))
	for i, ownerJson := range(json.Owners) {
		if owner, err := NewJobOwnerFromJson(ownerJson); err == nil {
			owners[i] = owner
		} else {
			return Job{}, err
		}
	}

	if (len(json.Command) > 0 && len(json.ShellCommand) > 0) || (len(json.Command) == 0 && len(json.ShellCommand) == 0) {
		return Job{}, errors.New("Exactly one of 'ShellCommand' or 'Command' must be specified")
	}

	var command Command
	if len(json.Command) > 0 {
		command = ExecCommand{json.Command}
	} else {
		command = ShellCommand{json.ShellCommand}
	}

	policy, pErr := NewExecutionPolicyFromJson(json.Policy)
	if pErr != nil {
		return Job{}, pErr
	}

	schedule, sErr := NewScheduleFromJson(json.Schedule)
	if sErr != nil {
		return Job{}, sErr
	}

	logger, lErr := logging.GetLogger(name)
	if lErr != nil {
		return Job{}, lErr
	}

	return Job {
		Name: name,
		Description: json.Description,
		Owners: owners,
		Policy: policy,
		Schedule: schedule,
		Command: command,
		Environment: json.Environment,
		Logger: logger,
	}, nil
}

func (j Job) IsValid(config JobValidationConfig) error {
	if len(j.Name) == 0 {
		return errors.New("Job name must not be empty")
	}

	if len(j.Owners) == 0 && config.AllowNoOwner() == false {
		return errors.New("Job must have specified at least one owner")
	}

	if err := j.Policy.IsValid(); err != nil {
		return errors.New(fmt.Sprintf("Invalid execution policy: %s", err))
	}

	if err := j.Schedule.IsValid(); err != nil {
		return errors.New(fmt.Sprintf("Invalid schedule: %s", err))
	}

	for i, owner := range(j.Owners) {
		if err := owner.IsValid(); err != nil {
			return errors.New(fmt.Sprintf("Invalid owner %d: %s", i, err))
		}
	}

	return nil
}
