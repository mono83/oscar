package events

import "time"

// EmittedBuffer contains multiple events
type EmittedBuffer []Emitted

// Add adds new item to events buffer
func (b *EmittedBuffer) Add(e Emitted) {
	*b = append(*b, e)
}

// AddIf adds item to buffer if provided callback returns true
func (b *EmittedBuffer) AddIf(e Emitted, c func(Emitted) bool) {
	if c != nil && c(e) {
		b.Add(e)
	}
}

// RemoteRequests returns list of remote requests
func (b *EmittedBuffer) RemoteRequests() []RemoteRequest {
	var rrs []RemoteRequest
	for _, v := range *b {
		if rr, ok := v.Data.(RemoteRequest); ok {
			rrs = append(rrs, rr)
		}
	}
	return rrs
}

// RemoteRequestsCount returns count of remote requests been
// held inside buffer
func (b *EmittedBuffer) RemoteRequestsCount() (count int) {
	for _, v := range *b {
		if _, ok := v.Data.(RemoteRequest); ok {
			count++
		}
	}
	return
}

// RemoteRequestsElapsed returns total elapsed time in remote requests
func (b *EmittedBuffer) RemoteRequestsElapsed() (elapsed time.Duration) {
	for _, v := range *b {
		if rr, ok := v.Data.(RemoteRequest); ok {
			elapsed += rr.Elapsed
		}
	}
	return
}
