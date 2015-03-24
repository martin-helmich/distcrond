package storage

import (
	"errors"
	"github.com/martin-helmich/distcrond/domain"
	"fmt"
)

type StorageBackendConfiguration interface {
	StorageBackend() string

	ElasticSearchHost() string
	ElasticSearchPort() int

	LogDirectory() string
}

type StorageBackend interface {
	Connect() error
	Disconnect() error
	SaveReport(report *domain.RunReport) error
}

func BuildStorageBackend(config StorageBackendConfiguration) (StorageBackend, error) {
	switch config.StorageBackend() {
	case "es":
		return NewElasticsearchBackend(
			config.ElasticSearchHost(),
			config.ElasticSearchPort(),
			"distcrond",
		), nil
	case "plain":
		return NewPlainStorageBackend(config.LogDirectory()), nil
	default:
		return &ElasticsearchBackend{}, errors.New(fmt.Sprintf("Unknown storage backend type: '%s'", config.StorageBackend()))
	}
}
