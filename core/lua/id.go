package lua

import "sync"

var incrementalId int
var incrementalIdMutex sync.Mutex

// id generates new unique incremental identifier
func id() int {
	incrementalIdMutex.Lock()
	defer incrementalIdMutex.Unlock()

	incrementalId++
	return incrementalId
}
