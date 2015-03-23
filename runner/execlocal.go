package runner

import (
	"time"
	"os/exec"
	"github.com/martin-helmich/distcrond/domain"
	"bytes"
)

type LocalExecutionStrategy struct {
	node *domain.Node
	logger interface {Debug(string, ...interface {})}
}

func (s *LocalExecutionStrategy) ExecuteCommand(command domain.Command, report *domain.RunReportItem) error {
	var output bytes.Buffer
	var start time.Time
	var cmd *exec.Cmd

	start = time.Now()

	args := command.Command()

	s.logger.Debug("Executing %s on local machine", args)

	//cmd = exec.Command("/bin/sh", "-c", command)
	cmd = &exec.Cmd{
		Path: args[0],
		Args: args,
	}
	cmd.Stdout = &output

	err := cmd.Run()

	report.Duration = time.Now().Sub(start)
	report.Output = output.String()
	report.Node = s.node

	if err == nil {
		report.Success = true
	} else {
		report.Success = false
	}

	return nil
}
