package domain

import (
	"time"
	"fmt"
)

type RunReport struct {
	Time time.Time
	Duration time.Duration
	Success bool
	Output string
	Node *Node
}

func (r *RunReport) Initialize() {
	r.Time = time.Now()
}

func (r *RunReport) Summary() string {
	date, _ := r.Time.MarshalText()
	return fmt.Sprintf("On %s at %s (duration %s): %s, %d bytes of output", r.Node.Name, date, r.Duration.String(), r.successOrFail(), len(r.Output))
}

func (r *RunReport) successOrFail() string {
	if r.Success {
		return "success"
	} else {
		return "FAIL"
	}
}
