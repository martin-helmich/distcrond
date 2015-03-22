package domain

import (
	"errors"
	"fmt"
)

const (
	CONN_LOCAL = "local"
	CONN_SSH = "ssh"
)

type ConnectionOptions struct {
	SshHost string
	SshUser string
	SshKeyFile string
}

func (o ConnectionOptions) SetDefaults(forType string) {
	if forType == CONN_SSH {
		o.SshHost = "localhost"
		o.SshUser = "root"
		o.SshKeyFile = "~/.ssh/id_rsa"
	}
}

func (o ConnectionOptions) IsValid(forType string) error {
	if forType == CONN_SSH {
		if len(o.SshHost) == 0 {
			return errors.New("SSH host is empty")
		}

		if len(o.SshUser) == 0 {
			return errors.New("SSH user is empty")
		}

		if len(o.SshKeyFile) == 0 {
			return errors.New("SSH key is empty")
		}
	} else if forType == CONN_LOCAL {

	}

	return nil
}

type Node struct {
	Name string
	Roles []string
	ConnectionType string
	ConnectionOptions ConnectionOptions
}

func (n Node) IsValid() error {
	if n.ConnectionType != CONN_LOCAL && n.ConnectionType != CONN_SSH {
		return errors.New("Invalid connection type (must be either " + CONN_LOCAL + " or " + CONN_SSH + ").")
	}

	if len(n.Name) == 0 {
		return errors.New("Name is empty")
	}

	if err := n.ConnectionOptions.IsValid(n.ConnectionType); err != nil {
		return errors.New(fmt.Sprintf("Invalid connection options: %s", err))
	}

	return nil
}
