package model

type indexOneToOne struct {
	extractor   IndexableExtractor
	byIndexable map[string]string
}

func NewSimpleIndex(extractor IndexableExtractor) *indexOneToOne {
	return &indexOneToOne{
		extractor:   extractor,
		byIndexable: make(map[string]string),
	}
}

func (index *indexOneToOne) Index(model ModelInstance) {
	modelId := model.ModelId()
	indexable := index.extractor(model)

	index.byIndexable[indexable] = modelId
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
	// byIndexable is keyed by the indexable value, not the modelId, so we must
	// remove the entries whose value is this modelId (there may be more than one
	// if the indexable changed without a prior Delete).
	for indexable, id := range index.byIndexable {
		if id == modelId {
			delete(index.byIndexable, indexable)
		}
	}
}

func (index *indexOneToOne) IndexableLength(indexable string) int {
	return len(index.byIndexable[indexable])
}
