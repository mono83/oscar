package out

import (
	"encoding/json"
	"github.com/mono83/oscar/events"
	"sort"
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
		if r.self.StartAt.IsZero() {
			r.self.StartAt = now
		}

		r.current.Elements = append(r.current.Elements, node)
		r.current = node
	} else if _, ok := e.(events.Finish); ok {
		r.current.FinishAt = now
		r.self.FinishAt = now
		if r.current.parent != nil {
			r.current = r.current.parent
		}
	} else if rr, ok := e.(events.RemoteRequest); ok {
		r.current.Remotes = append(r.current.Remotes, ReportRemoteRequest{
			Time:    now,
			Type:    rr.Type,
			URI:     rr.URI,
			Elapsed: rr.Elapsed,
			Success: rr.Success,
		})
	} else if a, ok := e.(events.AssertDone); ok {
		if a.Error == nil {
			r.current.Assertions++
		} else if r.current.Error == nil {
			errmsg := a.Error.Error()
			r.current.Error = &errmsg
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

// Suites returns suites collection
func (r *Report) Suites() []*ReportNode {
	return r.self.Elements
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

// TopRemoteRequests returns remote requests, that consumes most time
func (r *Report) TopRemoteRequests(count int, mode string) []AggregatedReportRemoteRequest {
	rr := r.self.RemoteRequestsRecursive()
	if len(rr) == 0 {
		return nil
	}

	// Aggregating
	agg := map[string]*AggregatedReportRemoteRequest{}
	for _, req := range rr {
		key := req.Type + req.URI
		if v, ok := agg[key]; ok {
			v.Elapsed = append(v.Elapsed, req.Elapsed)
		} else {
			agg[key] = &AggregatedReportRemoteRequest{
				Type:    req.Type,
				URI:     req.URI,
				Elapsed: []time.Duration{req.Elapsed},
			}
		}
	}

	// Converting to slice
	var sl []AggregatedReportRemoteRequest
	for _, v := range agg {
		sl = append(sl, *v)
	}

	// Sorting by total
	sort.Slice(sl, func(i, j int) bool {
		if mode == "avg" {
			return sl[i].Avg().Nanoseconds() > sl[j].Avg().Nanoseconds()
		} else if mode == "max" {
			return sl[i].Max().Nanoseconds() > sl[j].Max().Nanoseconds()
		}

		return sl[i].Total().Nanoseconds() > sl[j].Total().Nanoseconds()
	})

	if len(sl) > count {
		sl = sl[0:count]
	}

	return sl
}

// AggregatedReportRemoteRequest contains aggregated data for time spent on remote requests
type AggregatedReportRemoteRequest struct {
	Type    string
	URI     string
	Elapsed []time.Duration
}

// Total calculates total time, spent on single URL
func (a AggregatedReportRemoteRequest) Total() time.Duration {
	var t time.Duration
	for _, e := range a.Elapsed {
		t += e
	}
	return t
}

// Count returns amount of requests, made on single URL
func (a AggregatedReportRemoteRequest) Count() int {
	return len(a.Elapsed)
}

// Avg calculates average time across requests to single URL
func (a AggregatedReportRemoteRequest) Avg() time.Duration {
	return time.Duration(a.Total().Nanoseconds() / int64(a.Count()))
}

// Min returns minimal time, spent on URL
func (a AggregatedReportRemoteRequest) Min() time.Duration {
	min := a.Elapsed[0]
	if len(a.Elapsed) > 0 {
		for _, x := range a.Elapsed {
			if x < min {
				min = x
			}
		}
	}
	return min
}

// Max returns max time, spent on URL
func (a AggregatedReportRemoteRequest) Max() time.Duration {
	min := a.Elapsed[0]
	if len(a.Elapsed) > 0 {
		for _, x := range a.Elapsed {
			if x > min {
				min = x
			}
		}
	}
	return min
}

// LoadJSON loads report from saved JSON
func LoadJSON(bts []byte) (*Report, error) {
	var r ReportNode
	if err := json.Unmarshal(bts, &r); err != nil {
		return nil, err
	}

	x := &Report{
		self:    &r,
		current: &r,
	}

	linkParentRecursively(x.self)
	return x, nil
}

func linkParentRecursively(n *ReportNode) {
	for _, e := range n.Elements {
		e.parent = n
		linkParentRecursively(e)
	}
}
