package storage

import (
	"github.com/martin-helmich/distcrond/domain"
	logging "github.com/op/go-logging"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync/atomic"
	"time"
)

type PlainFileStorageBackend struct {
	logDirectory string
	logger *logging.Logger
	counter int64
//	reports []domain.RunReportJson
//	timer *time.Timer
}

func NewPlainStorageBackend(logDirectory string) *PlainFileStorageBackend {
	logger, _ := logging.GetLogger("persistence_plain")
	return &PlainFileStorageBackend{
		logDirectory: logDirectory,
		logger: logger,
		counter: 0,
//		make([]domain.RunReportJson, 0, 64),
	}
}

func (p *PlainFileStorageBackend) Connect() error {
	return nil
}

func (p *PlainFileStorageBackend) Disconnect() error {
	return nil
}

func (p *PlainFileStorageBackend) SaveReport(report *domain.RunReport) error {
	jsonReport := report.ToJson()
	body, _ := json.MarshalIndent(jsonReport, "", "    ")
	time, _ := report.Time.Start.MarshalText()

	filename := fmt.Sprintf("%s/%s-%s-%s.json", p.logDirectory, time, report.Job.Name, report.Id)

	if err := ioutil.WriteFile(filename, body, os.ModePerm); err != nil {
		p.logger.Error(fmt.Sprintf("Error while persisting report %s: %s", report.Id, err))
		return err
	}

	atomic.AddInt64(&p.counter, 1)
//	p.reports = append(p.reports, jsonReport)

	p.logger.Debug("Persisted report: " + string(body))
	return nil
}

func (p *PlainFileStorageBackend) ReportsForJob(job *domain.Job) ([]domain.RunReportJson, error) {
	start := time.Now()
	reports := make([]domain.RunReportJson, 0, atomic.LoadInt64(&p.counter))

	var walk filepath.WalkFunc = func(path string, file os.FileInfo, _ error) error {
		if file.IsDir() || file.Name()[0] == '.' {
			return nil
		}

		if content, err := ioutil.ReadFile(path); err != nil {
			return err
		} else {
			report := domain.RunReportJson{}
			if jErr := json.Unmarshal(content, &report); jErr != nil {
				return jErr
			}

			if report.Job == job.Name {
				reports = append(reports, report)
			}
			return nil
		}
	}

	if err := filepath.Walk(p.logDirectory, walk); err != nil {
		return nil, err
	}

	p.logger.Debug("Took %s for loading reports for job %s", time.Now().Sub(start).String(), job.Name)

	return reports, nil
}
