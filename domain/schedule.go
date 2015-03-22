package domain

import (
	"time"
)

type ScheduleJson struct {
	Interval string
	Reference string
}

type Schedule struct {
	Interval time.Duration
	Reference *time.Time
}

func NewScheduleFromJson(json ScheduleJson) (Schedule, error) {
	interval, err := time.ParseDuration(json.Interval)
	if err != nil {
		return Schedule{}, err
	}

	schedule := Schedule {interval, nil}

	if len(json.Reference) > 0 {
		reference, refErr := time.Parse("15:04", json.Reference)
		if refErr != nil {
			return schedule, refErr
		}

		schedule.Reference = &reference
	}

	return schedule, nil
}

func (s Schedule) IsValid() error {
	return nil
}
