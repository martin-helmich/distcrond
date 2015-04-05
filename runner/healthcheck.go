package runner

import (
	"time"
	"github.com/martin-helmich/distcrond/domain"
	logging "github.com/op/go-logging"
)

type HealthCheckerConfiguration interface {
	HealthCheckInterval() time.Duration
}

type HealthChecker interface {
	ScheduleHealthCheck(node *domain.Node)
}

type healthCheckerImpl struct {
	interval time.Duration
	logger *logging.Logger
}

func NewHealthChecker(config HealthCheckerConfiguration) HealthChecker {
	logger, _ := logging.GetLogger("healthcheck")
	return &healthCheckerImpl{config.HealthCheckInterval(), logger}
}

func (h *healthCheckerImpl) ScheduleHealthCheck(node *domain.Node) {
	h.logger.Info("Scheduled health check for node %s", node.Name)

	ticker := time.NewTicker(h.interval)
	go func() {
		strat := node.ExecutionStrategy
		for {
			select {
			case <-ticker.C:
				if err := strat.HealthCheck(); err == nil {
					h.logger.Notice("Node %s is up again", node.Name)
					ticker.Stop()

					func() {
						node.Lock.Lock()
						defer node.Lock.Unlock()

						node.Status = domain.STATUS_UP
					}()
				} else {
					h.logger.Warning("Node %s is still down", node.Name)
				}
			}
		}
	}()
}
