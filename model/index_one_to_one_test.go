package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func createTestOneToOneIndex() *indexOneToOne {
	extractor := func(instance ModelInstance) string {
		return string((instance.(*StopVisit)).VehicleJourneyId)
	}
	return NewSimpleIndex(extractor)
}

func Test_IndexOneToOne_FindOne(t *testing.T) {
	assert := assert.New(t)
	index := createTestOneToOneIndex()

	index.Index(&StopVisit{id: "stopVisitId", VehicleJourneyId: "dummy"})

	id, ok := index.FindOne("dummy")
	assert.True(ok, "should find the model after indexing")
	assert.Equal("stopVisitId", id)
}

// Delete is keyed by modelId, while the map is keyed by the indexable value.
// This guards the fix that makes Delete remove the entry whose value is the
// modelId (the previous implementation deleted by the wrong key, so the entry
// was never removed).
func Test_IndexOneToOne_Delete(t *testing.T) {
	assert := assert.New(t)
	index := createTestOneToOneIndex()

	index.Index(&StopVisit{id: "stopVisitId", VehicleJourneyId: "dummy"})
	index.Delete("stopVisitId")

	_, ok := index.FindOne("dummy")
	assert.False(ok, "entry should be gone after Delete(modelId)")
}

// If the indexable changes without a prior Delete, the OneToOne index can hold
// several keys pointing at the same modelId. Delete must remove all of them.
func Test_IndexOneToOne_Delete_RemovesAllStaleKeys(t *testing.T) {
	assert := assert.New(t)
	index := createTestOneToOneIndex()

	sv := &StopVisit{id: "stopVisitId", VehicleJourneyId: "vj-a"}
	index.Index(sv)
	sv.VehicleJourneyId = "vj-b"
	index.Index(sv) // leaves the "vj-a" -> stopVisitId entry behind

	index.Delete("stopVisitId")

	_, ok := index.FindOne("vj-a")
	assert.False(ok, "stale key vj-a should be removed by Delete")
	_, ok = index.FindOne("vj-b")
	assert.False(ok, "current key vj-b should be removed by Delete")
}

// Deleting one model must not drop another model's entry that happens to share
// the index.
func Test_IndexOneToOne_Delete_KeepsOtherEntries(t *testing.T) {
	assert := assert.New(t)
	index := createTestOneToOneIndex()

	index.Index(&StopVisit{id: "sv1", VehicleJourneyId: "vj-1"})
	index.Index(&StopVisit{id: "sv2", VehicleJourneyId: "vj-2"})

	index.Delete("sv1")

	_, ok := index.FindOne("vj-1")
	assert.False(ok, "deleted model's entry should be gone")
	id, ok := index.FindOne("vj-2")
	assert.True(ok, "other model's entry should remain")
	assert.Equal("sv2", id)
}
