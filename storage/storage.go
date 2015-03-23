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
}

type StorageBackend interface {
	Connect() error
	Disconnect() error
	SaveReport(report *domain.RunReport) error
}

func BuildStorageBackend(config StorageBackendConfiguration) (StorageBackend, error) {
	switch {
	case config.StorageBackend() == "es":
		return NewElasticsearchBackend(
			config.ElasticSearchHost(),
			config.ElasticSearchPort(),
			"distcrond",
		), nil
	default:
		return &ElasticsearchBackend{}, errors.New(fmt.Sprintf("Unknown storage backend type: '%s'", config.StorageBackend()))
	}
}
