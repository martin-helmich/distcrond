package domain

import (
	"fmt"
	"errors"
	"time"
)

type JobValidationConfig interface {
	AllowNoOwner() bool
}

type JobJson struct {
	Description string
	Owners []JobOwnerJson
	Policy ExecutionPolicyJson
	Schedule ScheduleJson
	ShellCommand string
	Command []string
}

type Job struct {
	Name string
	Description string
	Owners []JobOwner
	Policy ExecutionPolicy
	Schedule Schedule
	Command Command
	LastExecution time.Time
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

	return Job {
		Name: name,
		Description: json.Description,
		Owners: owners,
		Policy: policy,
		Schedule: schedule,
		Command: command,
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
