package out

import (
	"encoding/json"
	"github.com/mono83/oscar/events"
	"sync"
	"time"
)

// Report contains aggregated data about tests
type Report struct {
	m       sync.Mutex
	self    *ReportNode
	current *ReportNode
}

// OnEvent receives test event and registers it in report
func (r *Report) OnEvent(e interface{}) {
	now := time.Now().UTC()

	r.m.Lock()
	defer r.m.Unlock()

	if r.current == nil {
		r.self = &ReportNode{Type: "Report", Name: "Report"}
		r.current = r.self
	}

	if l, ok := e.(events.LogEvent); ok {
		r.current.Logs = append(r.current.Logs, ReportLogLine{
			Time:    now,
			Level:   l.Level,
			Message: l.Pattern,
		})
	} else if s, ok := e.(events.Start); ok {
		node := &ReportNode{
			parent:  r.current,
			ID:      s.ID,
			Type:    s.Type,
			Name:    s.Name,
			StartAt: now,
		}

		r.current.Elements = append(r.current.Elements, node)
		r.current = node
	} else if _, ok := e.(events.Finish); ok {
		r.current.FinishAt = now
		if r.current.parent != nil {
			r.current = r.current.parent
		}
	} else if rr, ok := e.(events.RemoteRequest); ok {
		r.current.Remotes = append(r.current.Remotes, ReportRemoteRequest{
			Time:    now,
			Type:    rr.Type,
			Elapsed: rr.Elapsed,
			Success: rr.Success,
		})
	} else if a, ok := e.(events.AssertDone); ok {
		if a.Error == nil {
			r.current.Assertions++
		} else if r.current.Error == nil {
			r.current.Error = a.Error
		}
	} else if v, ok := e.(events.SetVar); ok {
		if len(r.current.Variables) == 0 {
			r.current.Variables = map[string]string{}
		}

		r.current.Variables[v.Key] = v.Value
	} else if s, ok := e.(events.Sleep); ok {
		r.current.Sleep += time.Duration(s)
	}
}

// Flatten converts report data from tree to linear
func (r *Report) Flatten() []ReportNode {
	if r.self == nil {
		return nil
	}
	return r.self.Flatten()
}

// JSON returns report in JSON format
func (r *Report) JSON() string {
	bts, err := json.MarshalIndent(r.self, "", "  ")
	if err == nil {
		return string(bts)
	}

	return ""
}

// ReportNode is a report node
type ReportNode struct {
	ID   int
	Type string
	Name string

	StartAt  time.Time
	FinishAt time.Time

	parent   *ReportNode
	Elements []*ReportNode

	Assertions int
	Error      error
	Sleep      time.Duration
	Logs       []ReportLogLine
	Remotes    []ReportRemoteRequest
	Variables  map[string]string
}

// CountAssertions returns assertions count, made by this node
func (r ReportNode) CountAssertions() (success int, failed int) {
	success = r.Assertions
	if r.Error != nil {
		failed = 1
	}

	return
}

// CountAssertionsRecursive returns assertions count made by this and child nodes
func (r ReportNode) CountAssertionsRecursive() (success int, failed int) {
	success = r.Assertions
	if r.Error != nil {
		failed = 1
	}

	for _, child := range r.Elements {
		cs, cf := child.CountAssertionsRecursive()
		success += cs
		failed += cf
	}

	return
}

// CountRemoteRequestsRecursive returns amount of remote requests, done by this node and its childs
func (r ReportNode) CountRemoteRequestsRecursive() int {
	count := len(r.Remotes)

	for _, child := range r.Elements {
		count += child.CountRemoteRequestsRecursive()
	}

	return count
}

// ElapsedRemoteRecursive returns total elapsed time, spent by current node and it's childs on remote requests
func (r ReportNode) ElapsedRemoteRecursive() time.Duration {
	var total time.Duration

	for _, rr := range r.Remotes {
		total += rr.Elapsed
	}

	for _, child := range r.Elements {
		total += child.ElapsedRemoteRecursive()
	}

	return total
}

// ElapsedSleepRecursive returns total elapsed time, spent by current node and it's childs on sleeps
func (r ReportNode) ElapsedSleepRecursive() time.Duration {
	total := r.Sleep

	for _, child := range r.Elements {
		total += child.ElapsedSleepRecursive()
	}

	return total
}

// Elapsed returns time, spent by this node
func (r ReportNode) Elapsed() time.Duration {
	if r.StartAt.IsZero() || r.FinishAt.IsZero() {
		return time.Duration(0)
	}

	return r.FinishAt.Sub(r.StartAt)
}

// Flatten converts node data from tree to linear
func (r ReportNode) Flatten() []ReportNode {
	var nodes []ReportNode
	nodes = append(nodes, r)
	for _, child := range r.Elements {
		nodes = append(nodes, child.Flatten()...)
	}

	return nodes
}

// ReportLogLine contains logging data with event time
type ReportLogLine struct {
	Time    time.Time
	Level   byte
	Message string
}

// ReportRemoteRequest contains remote request event data with event time
type ReportRemoteRequest struct {
	Time    time.Time
	Type    string
	Elapsed time.Duration
	Success bool
}
