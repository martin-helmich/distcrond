package domain

import (
	"time"
	"fmt"
	"github.com/twinj/uuid"
)

type DurationJson struct {
	Milliseconds float64 `json:"milliseconds"`
	String string `json:"string"`
}

type TimePairJson struct {
	Start string `json:"start"`
	Stop string `json:"stop"`
}

type RunReportJson struct {
	Job string `json:"job"`
	Time TimePairJson `json:"time"`
	Duration DurationJson `json:"duration"`
	Success bool `json:"success"`
	Items []RunReportItemJson `json:"items"`
}

type RunReportItemJson struct {
	Time TimePairJson `json:"time"`
	Duration DurationJson `json:"duration"`
	Success bool `json:"success"`
	Output string `json:"output"`
	Node string `json:"node"`
}


// Run Report
// ==========

type TimePair struct {
	Start time.Time
	Stop time.Time
}

func (t *TimePair) ToJson() TimePairJson {
	var js TimePairJson

	if start, err := t.Start.MarshalText(); err == nil {
		js.Start = string(start)
	}

	if stop, err := t.Stop.MarshalText(); err == nil {
		js.Stop = string(stop)
	}

	return js
}

type RunReport struct {
	Id string
	Job *Job
	Time TimePair
	Items []RunReportItem
}

func (r *RunReport) Initialize(job *Job, nodeCount int) {
	r.Id = uuid.Formatter(uuid.NewV4(), uuid.CleanHyphen)
	r.Time.Start = time.Now()
	r.Job = job
	r.Items = make([]RunReportItem, nodeCount)

	for i := 0; i < nodeCount; i ++ {
		r.Items[i].Id = uuid.Formatter(uuid.NewV4(), uuid.CleanHyphen)
	}
}

func (r *RunReport) Finalize() {
	r.Time.Stop = time.Now()
}

func (r *RunReport) successOrFail() string {
	if r.Success() {
		return "success"
	} else {
		return "FAIL"
	}
}

func (r *RunReport) Success() bool {
	success := true
	for _, i := range r.Items {
		success = success && i.Success
	}
	return success
}

func (r *RunReport) Duration() time.Duration {
	return r.Time.Stop.Sub(r.Time.Start)
}

func (r *RunReport) ToJson() RunReportJson {
	items := make([]RunReportItemJson, len(r.Items))
	for i, item := range r.Items {
		items[i] = item.ToJson()
	}

	dur := r.Duration()

	return RunReportJson{
		Job: r.Job.Name,
		Time: r.Time.ToJson(),
		Duration: DurationJson{
			Milliseconds: float64(dur.Nanoseconds()) / float64(time.Millisecond),
			String: dur.String(),
		},
		Success: r.Success(),
		Items: items,
	}
}


// Run Report Item
// ===============

type RunReportItem struct {
	Id string
	Time TimePair
	Success bool
	Output string
	Node *Node
}

func (i *RunReportItem) Summary() string {
	date, _ := i.Time.Start.MarshalText()
	return fmt.Sprintf("On %s at %s (duration %s): %s, %d bytes of output", i.Node.Name, date, i.Duration().String(), i.successOrFail(), len(i.Output))
}

func (i *RunReportItem) successOrFail() string {
	if i.Success {
		return "success"
	} else {
		return "FAIL"
	}
}

func (i *RunReportItem) Duration() time.Duration {
	return i.Time.Stop.Sub(i.Time.Start)
}

func (i *RunReportItem) ToJson() RunReportItemJson {
	dur := i.Duration()
	return RunReportItemJson{
		Node: i.Node.Name,
		Time: i.Time.ToJson(),
		Duration: DurationJson{
			Milliseconds: float64(dur.Nanoseconds()) / float64(time.Millisecond),
			String: dur.String(),
		},
		Success: i.Success,
		Output: i.Output,
	}
}
