package runner

import (
	"os/exec"
	"github.com/martin-helmich/distcrond/domain"
	logging "github.com/op/go-logging"
	"bytes"
)

type LocalExecutionStrategy struct {
	node *domain.Node
}

func (s *LocalExecutionStrategy) ExecuteCommand(command domain.Command, report *domain.RunReportItem, logger *logging.Logger) error {
	var output bytes.Buffer
	var cmd *exec.Cmd

	args := command.Command()

	logger.Debug("Executing %s on local machine", args)

	cmd = &exec.Cmd{
		Path: args[0],
		Args: args,
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
