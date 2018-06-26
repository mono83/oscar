package lua

import "sync"

var incrementalID int
var incrementalIDMutex sync.Mutex

// id generates new unique incremental identifier
func id() int {
	incrementalIDMutex.Lock()
	defer incrementalIDMutex.Unlock()

	incrementalID++
	return incrementalID
}
