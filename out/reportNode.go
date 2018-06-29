package out

import (
	"time"
)

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
	Error      *string
	IsSkip     bool
	Sleep      time.Duration
	Logs       []ReportLogLine
	Remotes    []ReportRemoteRequest
	Variables  map[string]string
}

// Find searches for node by ID (including self) and returns it, if found
func (r *ReportNode) Find(id int) *ReportNode {
	if r.ID == id {
		return r
	}

	if len(r.Elements) > 0 {
		for _, e := range r.Elements {
			if f := e.Find(id); f != nil {
				return f
			}
		}
	}

	return nil
}

// GetParentID returns parent identifier, if node has parent
func (r ReportNode) GetParentID() *int {
	if r.parent != nil {
		return &r.parent.ID
	}

	return nil
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

// CountAssertionsSuccessRecursive returns success assertions count made by this and child nodes
func (r ReportNode) CountAssertionsSuccessRecursive() (success int) {
	success = r.Assertions

	for _, child := range r.Elements {
		success += child.CountAssertionsSuccessRecursive()
	}

	return
}

// CountAssertionsFailedRecursive returns failed assertions count made by this and child nodes
func (r ReportNode) CountAssertionsFailedRecursive() (failed int) {
	if r.Error != nil {
		failed = 1
	}

	for _, child := range r.Elements {
		failed += child.CountAssertionsFailedRecursive()
	}

	return
}

// CountAssertionsSkippedRecursive returns amount of test skips made by this and child nodes
func (r ReportNode) CountAssertionsSkippedRecursive() (skipped int) {
	if r.Error != nil && r.IsSkip {
		skipped = 1
	}

	for _, child := range r.Elements {
		skipped += child.CountAssertionsSkippedRecursive()
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

// RemoteRequestsRecursive return slice of remote requests, obtained recursively from current node and sub nodes
func (r ReportNode) RemoteRequestsRecursive() []ReportRemoteRequest {
	var requests []ReportRemoteRequest

	if len(r.Remotes) > 0 {
		requests = append(requests, r.Remotes...)
	}

	for _, child := range r.Elements {
		requests = append(requests, child.RemoteRequestsRecursive()...)
	}

	return requests
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
	URI     string
	Elapsed time.Duration
	Success bool
}
