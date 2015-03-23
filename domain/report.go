package domain

import (
	"time"
	"fmt"
	"github.com/twinj/uuid"
	"encoding/json"
)

type DurationJson struct {
	Milliseconds float64 `json:"milliseconds"`
	String string `json:"string"`
}

type RunReportJson struct {
	Job string `json:"job"`
	Time string `json:"time"`
	Duration DurationJson `json:"duration"`
	Success bool `json:"success"`
	Output string `json:"output"`
	Node string `json:"node"`
}

type RunReport struct {
	Id string
	Job *Job
	Time time.Time
	Duration time.Duration
	Success bool
	Output string
	Node *Node
}

func (r *RunReport) Initialize(job *Job, node *Node) {
	r.Id = uuid.Formatter(uuid.NewV4(), uuid.CleanHyphen)
	r.Time = time.Now()
	r.Job = job
	r.Node = node
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

func (r *RunReport) ToJson() string {
	jsonTime, _ := r.Time.MarshalText()
	jsonReport := RunReportJson{
		Time: string(jsonTime),
		Duration: DurationJson{
			Milliseconds: float64(r.Duration.Nanoseconds()) / float64(time.Millisecond),
			String: r.Duration.String(),
		},
		Success: r.Success,
		Output: r.Output,
		Node: r.Node.Name,
	}

	out, _ := json.Marshal(&jsonReport)
	return string(out)
}
