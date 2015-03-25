package runner

import (
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	. "github.com/martin-helmich/distcrond/domain"
	logging "github.com/op/go-logging"
	"bytes"
	"errors"
	"fmt"
	"strings"
)

type SshExecutionStrategy struct {
	node *Node
	privateKey ssh.Signer
	clientConfig ssh.ClientConfig
}

func NewSshExecutionStrategy (node *Node) (*SshExecutionStrategy, error) {
	strat := new(SshExecutionStrategy)

	keyString, keyErr := ioutil.ReadFile(node.ConnectionOptions.SshKeyFile)
	if keyErr != nil {
		return nil, errors.New(fmt.Sprintf("Could not read private key file %s: %s", node.ConnectionOptions.SshKeyFile, keyErr))
	}

	privateKey, keyParseError := ssh.ParsePrivateKey(keyString)
	if keyParseError != nil {
		return nil, errors.New(fmt.Sprintf("Could not parse private key file %s: %s", node.ConnectionOptions.SshKeyFile, keyParseError))
	}

	strat.node         = node
	strat.privateKey   = privateKey
	strat.clientConfig = ssh.ClientConfig{
		User: node.ConnectionOptions.SshUser,
		Auth: []ssh.AuthMethod{ssh.PublicKeys(privateKey)},
	}

	return strat, nil
}

func (s *SshExecutionStrategy) ExecuteCommand(command Command, report *RunReportItem, logger *logging.Logger) error {
	var output bytes.Buffer

	//	keyString, keyErr := ioutil.ReadFile(s.node.ConnectionOptions.SshKeyFile)
	//	if keyErr != nil {
	//		return errors.New(fmt.Sprintf("Could not read private key file %s: %s", s.node.ConnectionOptions.SshKeyFile, keyErr))
	//	}
	//
	//	privateKey, keyParseError := ssh.ParsePrivateKey(keyString)
	//	if keyParseError != nil {
	//		return errors.New(fmt.Sprintf("Could not parse private key file %s: %s", s.node.ConnectionOptions.SshKeyFile, keyParseError))
	//	}
	//
	//	config := ssh.ClientConfig{
	//		User: s.node.ConnectionOptions.SshUser,
	//		Auth: []ssh.AuthMethod{ssh.PublicKeys(privateKey)},
	//	}

	client, clientErr := ssh.Dial("tcp", s.node.ConnectionOptions.SshHost, &s.clientConfig)
	if clientErr != nil {
		return errors.New(fmt.Sprintf("Could not open connection to %s: %s", s.node.Name, clientErr))
	}

	session, sesErr := client.NewSession()
	if sesErr != nil {
		return errors.New(fmt.Sprintf("Could not open SSH session on %s: %s", s.node.Name, sesErr))
	}

	defer session.Close()

	session.Stdout = &output

	originalArgs := command.Command()
	quotedArgs := make([]string, len(originalArgs))
	for i, c := range originalArgs {
		quotedArgs[i] = "'" + strings.Replace(c, "'", "\\'", -1) + "'"
	}

	logger.Debug("Executing %s on remote machine\n", quotedArgs)

	runErr := session.Run(strings.Join(quotedArgs, " "))

	report.Output = output.String()

	if runErr == nil {
		report.Success = true
	} else {
		report.Success = false
	}

	return nil
}
