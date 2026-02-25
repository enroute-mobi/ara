package model

type indexOneToMany struct {
	extractor    IndexableExtractor
	byIndexable  map[ModelId][]ModelId
	byIdentifier map[ModelId]ModelId
}

func NewIndex(extractor IndexableExtractor) *indexOneToMany {
	return &indexOneToMany{
		extractor:    extractor,
		byIndexable:  make(map[ModelId][]ModelId),
		byIdentifier: make(map[ModelId]ModelId),
	}
}

func (index *indexOneToMany) Index(model ModelInstance) {
	modelId := model.ModelId()
	indexable := index.extractor(model)

	currentIndexable, ok := index.byIdentifier[modelId]
	if ok {
		if currentIndexable != indexable {
			index.removeFromIndexable(currentIndexable, modelId)
		} else {
			return
		}
	}

	index.byIndexable[indexable] = append(index.byIndexable[indexable], modelId)
	index.byIdentifier[modelId] = indexable
}

func (index *indexOneToMany) Find(indexable ModelId) ([]ModelId, bool) {
	modelIds, ok := index.byIndexable[indexable]
	return modelIds, ok
}

func (index *indexOneToMany) FindOne(indexable ModelId) (ModelId, bool) {
	modelIds, ok := index.byIndexable[indexable]
	if len(modelIds) != 0 {
		return modelIds[0], ok
	}
	var id ModelId
	return id, ok
}

func (index *indexOneToMany) Delete(modelId ModelId) {
	currentIndexable, ok := index.byIdentifier[modelId]
	if !ok {
		return
	}

	index.removeFromIndexable(currentIndexable, modelId)
	delete(index.byIdentifier, modelId)
}

func (index *indexOneToMany) IndexableLength(indexable ModelId) int {
	return len(index.byIndexable[indexable])
}

func (index *indexOneToMany) removeFromIndexable(indexable, modelId ModelId) {
	if len(index.byIndexable[indexable]) == 0 {
		return
	}
	for i, indexedModelId := range index.byIndexable[indexable] {
		if indexedModelId == modelId {
			index.byIndexable[indexable] = append(index.byIndexable[indexable][:i], index.byIndexable[indexable][i+1:]...)
			if len(index.byIndexable[indexable]) == 0 {
				delete(index.byIndexable, indexable)
			}
			return
		}
	}
}
