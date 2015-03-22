package runner

import (
	"errors"
	"github.com/martin-helmich/distcrond/domain"
	"fmt"
)

type ExecutionStrategy interface {
	ExecuteCommand(command domain.Command, report *RunReport) error
}

func GetStrategyForNode(node *domain.Node) (ExecutionStrategy, error) {
	if node.ConnectionType == domain.CONN_LOCAL {
		return &LocalExecutionStrategy{node}, nil
	} else if node.ConnectionType == domain.CONN_SSH {
		return &SshExecutionStrategy{node}, nil
	} else {
		return &NullExecutionStrategy{}, errors.New(fmt.Sprintf("Unknown connection type for node %s: %s", node.Name, node.ConnectionType))
	}
}

type NullExecutionStrategy struct {}

func (n *NullExecutionStrategy) ExecuteCommand(_ domain.Command, _ *RunReport) error {
	return nil
}
