package runner

import (
	"errors"
	. "github.com/martin-helmich/distcrond/domain"
	"fmt"
)

func GetStrategyForNode(node *Node) (ExecutionStrategy, error) {
	switch {
	case node.ConnectionType == CONN_LOCAL:
		return &LocalExecutionStrategy{node}, nil

	case node.ConnectionType == CONN_SSH:
		if str, err := NewSshExecutionStrategy(node); err != nil {
			return nil, err
		} else {
			return str, nil
		}

	default:
		return &NullExecutionStrategy{}, errors.New(fmt.Sprintf("Unknown connection type for node %s: %s", node.Name, node.ConnectionType))
	}
}

type NullExecutionStrategy struct {}

func (n *NullExecutionStrategy) ExecuteCommand(_ *Job, _ *RunReportItem) error {
	return nil
}
