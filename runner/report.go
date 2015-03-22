package runner

import (
	"time"
	"github.com/martin-helmich/distcrond/domain"
)

type RunReport struct {
	Time time.Time
	Duration time.Duration
	Success bool
	Output string
	Node *domain.Node
}
