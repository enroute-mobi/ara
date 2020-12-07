package model

type IndexableExtractor func(ModelInstance) ModelId

type Index struct {
	extractor    IndexableExtractor
	byIndexable  map[ModelId][]ModelId
	byIdentifier map[ModelId]ModelId
}

func NewIndex(extractor IndexableExtractor) *Index {
	return &Index{
		extractor:    extractor,
		byIndexable:  make(map[ModelId][]ModelId),
		byIdentifier: make(map[ModelId]ModelId),
	}
}

func (index *Index) Index(model ModelInstance) {
	modelId := model.modelId()
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

func (index *Index) Find(indexable ModelId) ([]ModelId, bool) {
	modelIds, ok := index.byIndexable[indexable]
	return modelIds, ok
}

func (index *Index) Delete(modelId ModelId) {
	currentIndexable, ok := index.byIdentifier[modelId]
	if !ok {
		return
	}

	index.removeFromIndexable(currentIndexable, modelId)
	delete(index.byIdentifier, modelId)
}

func (index *Index) removeFromIndexable(indexable, modelId ModelId) {
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
