package reader

import (
	"os"
	"errors"
	"fmt"
	"strings"
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"github.com/martin-helmich/distcrond/domain"
	"github.com/martin-helmich/distcrond/logging"
	"github.com/martin-helmich/distcrond/runner"
)

type NodeReader struct {
	receiver NodeReceiver
}

type NodeReceiver interface {
	AddNode(domain.Node)
}

func NewNodeReader(receiver NodeReceiver) *NodeReader {
	reader := new(NodeReader)
	reader.receiver = receiver

	return reader
}

func (r NodeReader) ReadFromDirectory(directory string) error {
	logging.Info("Reading node configuration")

	var walk filepath.WalkFunc = func(path string, file os.FileInfo, err error) error {
		if file.IsDir() {
			return nil
		}

		if file.Name()[0] == '.' {
			logging.Debug("Skipping %s", path)
			return nil
		}

		fileContents, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		wrapError := func(err error) error {
			return errors.New(fmt.Sprintf("Error parsing file %s: %s", file.Name(), err))
		}

		nodeJson := domain.NodeJson{}
		nodeName := strings.Replace(file.Name(), ".json", "", 1)

		if err = json.Unmarshal(fileContents, &nodeJson); err != nil {
			return wrapError(err)
		}

		node, err := domain.NewNodeFromJson(nodeName, nodeJson)
		if err != nil {
			return err
		}

		if str, strErr := runner.GetStrategyForNode(&node); strErr == nil {
			node.ExecutionStrategy = str
		} else {
			return strErr
		}

		logging.Debug("Read node from %s: %s", file.Name(), node)

		if validErr := node.IsValid(); validErr != nil {
			return wrapError(validErr)
		}

		r.receiver.AddNode(node)

		return nil
	}

	if err := filepath.Walk(directory, walk); err != nil {
		return err
	}

	return nil
}
