package runner

import (
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"time"
	"github.com/martin-helmich/distcrond/domain"
	"bytes"
	"log"
	"errors"
	"fmt"
	"strings"
)

type SshExecutionStrategy struct {
	node *domain.Node
}

func (s *SshExecutionStrategy) ExecuteCommand(command domain.Command, report *RunReport) error {
	var output bytes.Buffer
	var start time.Time

	key, keyErr := ioutil.ReadFile(s.node.ConnectionOptions.SshKeyFile)
	if keyErr != nil {
		return errors.New(fmt.Sprintf("Could not read private key file %s: %s", s.node.ConnectionOptions.SshKeyFile, keyErr))
	}

	//pk, _ := ssh.ParsePublicKey(key)
	pk, _ := ssh.ParsePrivateKey(key)
//	signer, _ := ssh.NewSignerFromKey(pk)

	config := ssh.ClientConfig{
		User: s.node.ConnectionOptions.SshUser,
		Auth: []ssh.AuthMethod{ssh.PublicKeys(pk)},
	}

	client, clientErr := ssh.Dial("tcp", s.node.ConnectionOptions.SshHost, &config)
	if clientErr != nil {
		return errors.New(fmt.Sprintf("Could not open connection to %s: %s", s.node.Name, clientErr))
	}

	session, sesErr := client.NewSession()
	if sesErr != nil {
		return errors.New(fmt.Sprintf("Could not open SSH session on %s: %s", s.node.Name, sesErr))
	}

	defer session.Close()

	session.Stdout = &output

	start = time.Now()

	originalArgs := command.Command()
	quotedArgs := make([]string, len(originalArgs))
	for i, c := range originalArgs {
		quotedArgs[i] = "'" + strings.Replace(c, "'", "\\'", -1) + "'"
	}

	log.Printf("Executing %s on remote machine\n", quotedArgs)

	//cmd = exec.Command("/bin/sh", "-c", command)
//	cmd = &exec.Cmd{
//		Path: args[0],
//		Args: args,
//	}
//	cmd.Stdout = &output

//	err := cmd.Run()
	runErr := session.Run(strings.Join(quotedArgs, " "))

	report.Duration = time.Now().Sub(start)
	report.Output = output.String()
	report.Node = s.node

	if runErr == nil {
		report.Success = true
	} else {
		report.Success = false
	}

	return nil
}
