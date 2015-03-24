package storage

import (
	"github.com/martin-helmich/distcrond/domain"
	logging "github.com/op/go-logging"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type PlainFileStorageBackend struct {
	logDirectory string
	logger *logging.Logger
}

func NewPlainStorageBackend(logDirectory string) *PlainFileStorageBackend {
	logger, _ := logging.GetLogger("persistence_plain")
	return &PlainFileStorageBackend{logDirectory: logDirectory, logger: logger}
}

func (p *PlainFileStorageBackend) Connect() error {
	return nil
}

func (p *PlainFileStorageBackend) Disconnect() error {
	return nil
}

func (p *PlainFileStorageBackend) SaveReport(report *domain.RunReport) error {
	body, _ := json.MarshalIndent(report.ToJson(), "", "    ")
	time, _ := report.Time.Start.MarshalText()

	filename := fmt.Sprintf("%s/%s-%s-%s.json", p.logDirectory, time, report.Job.Name, report.Id)

	if err := ioutil.WriteFile(filename, body, os.ModePerm); err != nil {
		p.logger.Error(fmt.Sprintf("Error while persisting report %s: %s", report.Id, err))
		return err
	}

	p.logger.Debug("Persisted report: " + string(body))
	return nil
}
