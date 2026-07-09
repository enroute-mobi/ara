package core

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ARA-1956: AlreadySend previously read under RLock then re-acquired Lock to
// write, leaving a window where two concurrent calls for the same subID could
// both observe "not yet sent" and both return false — causing a duplicate
// DeleteSubscription request. The fix holds a single write Lock for the whole
// method, making the check-and-register atomic.
func Test_DeletedSubscriptions_AlreadySend_AtomicUnderConcurrency(t *testing.T) {
	ds := NewDeletedSubscriptions()
	const goroutines = 50

	results := make([]bool, goroutines)
	var wg sync.WaitGroup

	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		i := i
		go func() {
			defer wg.Done()
			results[i] = ds.AlreadySend("RATPCapIDF:Subscription::001:LOC")
		}()
	}
	wg.Wait()

	falseCount := 0
	for _, r := range results {
		if !r {
			falseCount++
		}
	}
	assert.Equal(t, 1, falseCount, "exactly one goroutine should register the deletion; %d did (TOCTOU would cause > 1)", falseCount)
}
