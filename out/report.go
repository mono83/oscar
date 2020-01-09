package out

import (
	"encoding/json"
	"sort"
	"sync"
	"time"

	"github.com/mono83/oscar/events"
)

// Report contains aggregated data about tests
type Report struct {
	m       sync.Mutex
	self    *ReportNode
	current *ReportNode

	Statistics struct {
		TotalEvents int
		Assertions  int
	}
}

// OnEvent receives test event and registers it in report
func (r *Report) OnEvent(e *events.Emitted) {
	r.m.Lock()
	defer r.m.Unlock()

	r.Statistics.TotalEvents++

	if r.current == nil {
		r.self = &ReportNode{Type: "Report", Name: "Report", Events: []events.Emitted{}}
		r.current = r.self
	}

	// Processing event
	events.EventRouter{
		Log: func(l events.LogEvent, em *events.Emitted) {
			r.IfFound(em.OwnerID, func(node *ReportNode) {
				node.Events.Add(*em)
			})
		},
		Start: func(s events.Start, em *events.Emitted) {
			r.IfFound(em.OwnerID, func(node *ReportNode) {
				node.StartAt = em.Time
				if r.self.StartAt.IsZero() {
					r.self.StartAt = em.Time
				}
			})
		},
		Finish: func(_ events.Finish, em *events.Emitted) {
			r.IfFound(em.OwnerID, func(node *ReportNode) {
				node.FinishAt = em.Time
				if r.self.FinishAt.Before(em.Time) {
					r.self.FinishAt = em.Time
				}
			})
		},
		Remote: func(rr events.RemoteRequest, em *events.Emitted) {
			r.IfFound(em.OwnerID, func(node *ReportNode) {
				node.Events.Add(*em)
			})
		},
		Assert: func(a events.AssertDone, em *events.Emitted) {
			r.Statistics.Assertions++
			r.IfFound(em.OwnerID, func(node *ReportNode) {
				node.Events.Add(*em)
			})
		},
		Failure: func(f events.Failure, em *events.Emitted) {
			err := f.Error()
			r.IfFound(em.OwnerID, func(node *ReportNode) {
				node.Error = &err
			})
		},
		Skip: func(s events.Skip, em *events.Emitted) {
			err := s.Error()
			r.IfFound(em.OwnerID, func(node *ReportNode) {
				node.Error = &err
				node.IsSkip = true
			})
		},
		Var: func(v events.SetVar, em *events.Emitted) {
			r.IfFound(em.OwnerID, func(node *ReportNode) {
				if len(node.Variables) == 0 {
					node.Variables = map[string]string{}
				}

				node.Variables[v.Key] = v.Value
			})
		},
		Sleep: func(s events.Sleep, em *events.Emitted) {
			r.IfFound(em.OwnerID, func(node *ReportNode) {
				node.Sleep += time.Duration(s)
			})
		},
		RegistrationIn: func(rin events.RegistrationBegin, em *events.Emitted) {
			node := &ReportNode{
				parent: r.current,
				ID:     rin.ID,
				Type:   rin.Type,
				Name:   rin.Name,
				Events: []events.Emitted{},
			}

			r.current.Elements = append(r.current.Elements, node)
			r.current = node
		},
		RegistrationOut: func(rout events.RegistrationEnd, em *events.Emitted) {
			if r.current.parent != nil {
				r.current = r.current.parent
			}
		},
	}.OnEvent(e)
}

// Suites returns suites collection
func (r *Report) Suites() []*ReportNode {
	return r.self.Elements
}

// Find searches for node by ID and returns it, if found
func (r *Report) Find(id int) *ReportNode {
	if r.self != nil {
		return r.self.Find(id)
	}

	return nil
}

// IfFound searched for node by ID and invokes callback if found
func (r *Report) IfFound(id int, callback func(*ReportNode)) {
	if f := r.Find(id); f != nil && callback != nil {
		callback(f)
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
