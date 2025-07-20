package hub

import (
	"sync"

	"github.com/DeepAung/gradient/grader-server/proto"
)

var (
	mu      sync.Mutex
	results = make(map[string]<-chan proto.StatusType)
)

func CreateResult(resultId string, resultCh <-chan proto.StatusType) {
	mu.Lock()
	results[resultId] = resultCh
	mu.Unlock()
}

func DeleteResult(resultId string) {
	mu.Lock()
	delete(results, resultId)
	mu.Unlock()
}

func PopResult(resultId string) (<-chan proto.StatusType, bool) {
	mu.Lock()
	resultCh, ok := results[resultId]
	mu.Unlock()

	if ok {
		DeleteResult(resultId)
	}

	return resultCh, ok
}
