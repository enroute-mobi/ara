package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// The OneToOne index is used only by Vehicles (ByVehicleJourney and
// ByNextStopVisit), so the tests exercise it with the production *Vehicle extractor.
func createTestOneToOneIndex() *indexOneToOne {
	return NewSimpleIndex(vehicleVjExtractor)
}

func Test_IndexOneToOne_FindOne(t *testing.T) {
	assert := assert.New(t)
	index := createTestOneToOneIndex()

	index.Index(&Vehicle{id: "v1", VehicleJourneyId: "vj-a"})

	id, ok := index.FindOne("vj-a")
	assert.True(ok, "should find the model after indexing")
	assert.Equal("v1", id)
}

// Delete is given a modelId, while the map is keyed by the indexable value.
// This guards the fix that makes Delete remove the entry for that model (the
// previous implementation deleted by the wrong key, so the entry was never
// removed).
func Test_IndexOneToOne_Delete(t *testing.T) {
	assert := assert.New(t)
	index := createTestOneToOneIndex()

	index.Index(&Vehicle{id: "v1", VehicleJourneyId: "vj-a"})
	index.Delete("v1")

	_, ok := index.FindOne("vj-a")
	assert.False(ok, "entry should be gone after Delete(modelId)")
}

// Re-indexing a model whose indexable changed must evict the previous key, so
// the index does not accumulate stale entries on reassignment.
func Test_IndexOneToOne_Index_EvictsPreviousKeyOnChange(t *testing.T) {
	assert := assert.New(t)
	index := createTestOneToOneIndex()

	v := &Vehicle{id: "v1", VehicleJourneyId: "vj-a"}
	index.Index(v)
	v.VehicleJourneyId = "vj-b"
	index.Index(v)

	_, ok := index.FindOne("vj-a")
	assert.False(ok, "stale key vj-a should be evicted when the model is re-indexed")
	id, ok := index.FindOne("vj-b")
	assert.True(ok, "current key vj-b should resolve")
	assert.Equal("v1", id)
}

// A model that shares an indexable value with another must not have its entry
// dropped when the other is deleted.
func Test_IndexOneToOne_Delete_KeepsSharedIndexableOfOtherModel(t *testing.T) {
	assert := assert.New(t)
	index := createTestOneToOneIndex()

	index.Index(&Vehicle{id: "v1", VehicleJourneyId: "vj-shared"})
	index.Index(&Vehicle{id: "v2", VehicleJourneyId: "vj-shared"}) // overwrites the value to v2

	index.Delete("v1")

	id, ok := index.FindOne("vj-shared")
	assert.True(ok, "the entry taken over by v2 must survive deleting v1")
	assert.Equal("v2", id)
}

// Deleting one model must not drop another model's entry that happens to share
// the index.
func Test_IndexOneToOne_Delete_KeepsOtherEntries(t *testing.T) {
	assert := assert.New(t)
	index := createTestOneToOneIndex()

	index.Index(&Vehicle{id: "v1", VehicleJourneyId: "vj-1"})
	index.Index(&Vehicle{id: "v2", VehicleJourneyId: "vj-2"})

	index.Delete("v1")

	_, ok := index.FindOne("vj-1")
	assert.False(ok, "deleted model's entry should be gone")
	id, ok := index.FindOne("vj-2")
	assert.True(ok, "other model's entry should remain")
	assert.Equal("v2", id)
}
