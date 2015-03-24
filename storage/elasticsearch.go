package storage

import (
	"net/http"
	"github.com/martin-helmich/distcrond/domain"
	logging "github.com/op/go-logging"
	"fmt"
	"errors"
	"strings"
	"encoding/json"
)

type ElasticsearchBackend struct {
	host string
	port int
	index string
	uri string
	logger *logging.Logger
	client http.Client
}

func NewElasticsearchBackend(host string, port int, index string) *ElasticsearchBackend {
	logger, _ := logging.GetLogger("persistence_es")
	backend := ElasticsearchBackend{
		host: host,
		port: port,
		index: index,
		uri: fmt.Sprintf("http://%s:%d", host, port),
		logger: logger,
		client: http.Client{},
	}

	return &backend
}

func (e *ElasticsearchBackend) Connect() error {
	if err := e.request("", "GET", ""); err != nil {
		return errors.New(fmt.Sprintf("Elasticsearch backend at %s does not appear to be reachable.", e.uri))
	}
	return nil
}

func (e *ElasticsearchBackend) Disconnect() error {
	// HTTP is stateless, no disconnect needed. But other backends might.
	return nil
}

func (e *ElasticsearchBackend) SaveReport(report *domain.RunReport) error {
	body, _ := json.Marshal(report.ToJson())

	if err := e.requestDocument("reports", report.Id, "PUT", string(body)); err != nil {
		e.logger.Error(fmt.Sprintf("Error while persisting report %s: %s", report.Id, err))
		return err
	}

	e.logger.Debug("Persisted report: " + string(body))
	return nil
}

func (e *ElasticsearchBackend) requestDocument(docType string, docId string, method string, body string) error {
	path := fmt.Sprintf("%s/%s/%s", e.index, docType, docId)
	return e.request(path, method, body)
}

func (e *ElasticsearchBackend) request(path string, method string, body string) error {
	bufferReader := strings.NewReader(body)
	uri          := fmt.Sprintf("%s/%s", e.uri, path)

	request, reqErr := http.NewRequest(method, uri, bufferReader)
	if reqErr != nil {
		return reqErr
	}

	e.logger.Debug("Performing HTTP request to %s, %d bytes in body.", uri, len(body))

	resp, respErr := e.client.Do(request)
	if respErr != nil {
		return errors.New(fmt.Sprintf("Error while requesting %s: %s", uri, respErr))
	}

	if resp.StatusCode >= 300 {
		return errors.New(fmt.Sprintf("Unexpected status code %d while requesting %s", resp.StatusCode, uri))
	}

	return nil
}

func (p *ElasticsearchBackend) ReportsForJob(job *domain.Job) ([]domain.RunReportJson, error) {
	// TODO: Implement me!
	reports := make([]domain.RunReportJson, 0)
	return reports, nil
}
