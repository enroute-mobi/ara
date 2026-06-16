package model

type indexOneToOne struct {
	extractor   IndexableExtractor
	byIndexable map[string]string
	byModelId   map[string]string
}

func NewSimpleIndex(extractor IndexableExtractor) *indexOneToOne {
	return &indexOneToOne{
		extractor:   extractor,
		byIndexable: make(map[string]string),
		byModelId:   make(map[string]string),
	}
}

func (index *indexOneToOne) Index(model ModelInstance) {
	modelId := model.ModelId()
	indexable := index.extractor(model)

	// Evict the previous key for this model when its indexable changed, so the
	// index doesn't retain a stale entry after a reassignment. The byModelId
	// reverse map makes this O(1). Guard on the value so we never drop an entry
	// that another model has since taken over.
	if previous, ok := index.byModelId[modelId]; ok && previous != indexable {
		if index.byIndexable[previous] == modelId {
			delete(index.byIndexable, previous)
		}
	}

	index.byIndexable[indexable] = modelId
	index.byModelId[modelId] = indexable
}

func (index *indexOneToOne) FindOne(indexable string) (string, bool) {
	modelId, ok := index.byIndexable[indexable]
	return modelId, ok
}

func (index *indexOneToOne) Find(indexable string) ([]string, bool) {
	modelId, ok := index.byIndexable[indexable]
	return []string{modelId}, ok
}

func (index *indexOneToOne) Delete(modelId string) {
	indexable, ok := index.byModelId[modelId]
	if !ok {
		return
	}
	// Only drop the indexable entry if it still points at this model, so a model
	// that shares the indexable value isn't removed by mistake.
	if index.byIndexable[indexable] == modelId {
		delete(index.byIndexable, indexable)
	}
	delete(index.byModelId, modelId)
}

func (index *indexOneToOne) IndexableLength(indexable string) int {
	return len(index.byIndexable[indexable])
}
