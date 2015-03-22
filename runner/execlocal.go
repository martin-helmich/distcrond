package runner

import (
	"time"
	"os/exec"
	"github.com/martin-helmich/distcrond/domain"
	"bytes"
	"log"
)

type LocalExecutionStrategy struct {
	node *domain.Node
}

func (s *LocalExecutionStrategy) ExecuteCommand(command domain.Command, report *RunReport) error {
	var output bytes.Buffer
	var start time.Time
	var cmd *exec.Cmd

	start = time.Now()

	args := command.Command()

	log.Printf("Executing %s on local machine\n", args)

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
