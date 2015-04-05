package runner

import (
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	. "github.com/martin-helmich/distcrond/domain"
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

func (s *SshExecutionStrategy) quote(c string) string {
	return "'" + strings.Replace(c, "'", "\\'", -1) + "'"
}

func (s *SshExecutionStrategy) HealthCheck() error {
	client, clientErr := ssh.Dial("tcp", s.node.ConnectionOptions.SshHost, &s.clientConfig)
	if clientErr != nil {
		return NewNodeDownError(s.node, "Could not open TCP connection", clientErr)
	}

	session, sesErr := client.NewSession()
	if sesErr != nil {
		return NewNodeDownError(s.node, "Could not start SSH session", sesErr)
	}

	defer session.Close()

	return nil
}

func (s *SshExecutionStrategy) ExecuteCommand(job *Job, report *RunReportItem) error {
	var output bytes.Buffer

	client, clientErr := ssh.Dial("tcp", s.node.ConnectionOptions.SshHost, &s.clientConfig)
	if clientErr != nil {
		return NewNodeDownError(s.node, "Could not open TCP connection", clientErr)
	}

	session, sesErr := client.NewSession()
	if sesErr != nil {
		return NewNodeDownError(s.node, "Could not start SSH session", sesErr)
	}

	defer session.Close()

	session.Stdout = &output

	originalArgs := job.Command.Command()
	quotedArgs := make([]string, len(originalArgs))
	for i, c := range originalArgs {
		quotedArgs[i] = s.quote(c)
	}

	job.Logger.Debug("Executing %s on remote machine", quotedArgs)

	cmdStrings := make([]string, 0, len(job.Environment) + 1)
	for key, value := range job.Environment {
		if err := session.Setenv(key, value); err != nil {
			job.Logger.Warning("Could not remotely set environment variable %s to %s. Check your 'AcceptEnv' server setting.", key, value)
			cmdStrings = append(cmdStrings, "export " + key + "=" + s.quote(value))
		}
	}

	cmdStrings = append(cmdStrings, strings.Join(quotedArgs, " "))
	cmd := strings.Join(cmdStrings, " ; ")

	job.Logger.Debug("Actually running \"%s\"", cmd)

	runErr := session.Run(strings.Join(cmdStrings, " ; "))

	report.Output = output.String()

	if runErr == nil {
		report.Success = true
	} else {
		report.Success = false
	}

	return nil
}
