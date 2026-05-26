package model

type indexOneToMany struct {
	extractor    IndexableExtractor
	byIndexable  map[string][]string
	byIdentifier map[string]string
}

func NewIndex(extractor IndexableExtractor) *indexOneToMany {
	return &indexOneToMany{
		extractor:    extractor,
		byIndexable:  make(map[string][]string),
		byIdentifier: make(map[string]string),
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

func (index *indexOneToMany) Find(indexable string) ([]string, bool) {
	modelIds, ok := index.byIndexable[indexable]
	return modelIds, ok
}

func (index *indexOneToMany) FindOne(indexable string) (string, bool) {
	modelIds, ok := index.byIndexable[indexable]
	if len(modelIds) != 0 {
		return modelIds[0], ok
	}
	var id string
	return id, ok
}

func (index *indexOneToMany) Delete(modelId string) {
	currentIndexable, ok := index.byIdentifier[modelId]
	if !ok {
		return
	}

	index.removeFromIndexable(currentIndexable, modelId)
	delete(index.byIdentifier, modelId)
}

func (index *indexOneToMany) IndexableLength(indexable string) int {
	return len(index.byIndexable[indexable])
}

func (index *indexOneToMany) removeFromIndexable(indexable, modelId string) {
	if len(index.byIndexable[indexable]) == 0 {
		return
	}
	for i, indexedstring := range index.byIndexable[indexable] {
		if indexedstring == modelId {
			index.byIndexable[indexable] = append(index.byIndexable[indexable][:i], index.byIndexable[indexable][i+1:]...)
			if len(index.byIndexable[indexable]) == 0 {
				delete(index.byIndexable, indexable)
			}
			return
		}
	}
}
