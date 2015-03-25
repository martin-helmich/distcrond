package runner

import (
	"os/exec"
	"github.com/martin-helmich/distcrond/domain"
	"bytes"
)

type LocalExecutionStrategy struct {
	node *domain.Node
}

func (s *LocalExecutionStrategy) ExecuteCommand(job *domain.Job, report *domain.RunReportItem) error {
	var output bytes.Buffer
	var cmd *exec.Cmd

	args := job.Command.Command()

	job.Logger.Debug("Executing %s on local machine", args)

	env := make([]string, len(job.Environment))
	i   := 0
	for key, value := range job.Environment {
		env[i] = key + "=" + value
		i++
	}

	cmd = &exec.Cmd{
		Path: args[0],
		Args: args,
		Env: env,
	}
	cmd.Stdout = &output

	err := cmd.Run()

	report.Output = output.String()

	if err == nil {
		report.Success = true
	} else {
		report.Success = false
	}

	return nil
}
