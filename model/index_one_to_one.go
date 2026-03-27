package model

type indexOneToOne struct {
	extractor   IndexableExtractor
	byIndexable map[ModelId]ModelId
}

func NewSimpleIndex(extractor IndexableExtractor) *indexOneToOne {
	return &indexOneToOne{
		extractor:   extractor,
		byIndexable: make(map[ModelId]ModelId),
	}
}

func (index *indexOneToOne) Index(model ModelInstance) {
	modelId := model.ModelId()
	indexable := index.extractor(model)

	index.byIndexable[indexable] = modelId
}

func (index *indexOneToOne) FindOne(indexable ModelId) (ModelId, bool) {
	modelId, ok := index.byIndexable[indexable]
	return modelId, ok
}

func (index *indexOneToOne) Find(indexable ModelId) ([]ModelId, bool) {
	modelId, ok := index.byIndexable[indexable]
	return []ModelId{modelId}, ok
}

func (index *indexOneToOne) Delete(modelId ModelId) {
	delete(index.byIndexable, modelId)
}

func (index *indexOneToOne) IndexableLength(indexable ModelId) int {
	return len(index.byIndexable[indexable])
}
