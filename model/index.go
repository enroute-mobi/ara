package model

type ModelId string

//type StopAreaId ModelId
//type LineId ModelId
// ...

type ModelInstance interface{}
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

func (index *Index) Index(modelId ModelId, model ModelInstance) {
	indexable := index.extractor(model)

	currentIndexable, ok := index.byIdentifier[modelId]
	if ok && currentIndexable != indexable {
		index.removeFromIndexable(currentIndexable, modelId)
	}

	index.byIndexable[indexable] = append(index.byIndexable[indexable], modelId)
	index.byIdentifier[modelId] = indexable
}

func (index *Index) Find(indexable ModelId) ([]ModelId, bool) {
	modelId, ok := index.byIndexable[indexable]
	return modelId, ok
}

func (index *Index) Delete(modelId ModelId) {
	currentIndexable, ok := index.byIdentifier[modelId]
	if !ok {
		return
	}

	index.removeFromIndexable(currentIndexable, modelId)
	delete(index.byIdentifier, modelId)
}
