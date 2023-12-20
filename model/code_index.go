package model

type CodeIndex struct {
	byCode       map[Code]ModelId
	byIdentifier map[ModelId]Codes
}

func NewCodeIndex() *CodeIndex {
	return &CodeIndex{
		byCode:       make(map[Code]ModelId),
		byIdentifier: make(map[ModelId]Codes),
	}
}

func (index *CodeIndex) Index(model ModelInstance) {
	currentIndexable, ok := index.byIdentifier[model.modelId()]
	if ok {
		for indexedCodeSpace, indexedCode := range currentIndexable {
			modelCode, ok := model.Code(indexedCodeSpace)
			if !ok || modelCode.Value() != indexedCode.Value() {
				delete(index.byCode, indexedCode)
			}
		}
	}

	if index.byIdentifier[model.modelId()] == nil {
		index.byIdentifier[model.modelId()] = make(Codes)
	}

	for _, code := range model.Codes() {
		index.byCode[code] = model.modelId()
		index.byIdentifier[model.modelId()][code.CodeSpace()] = code
	}
}

func (index *CodeIndex) Find(code Code) (ModelId, bool) {
	modelId, ok := index.byCode[code]
	return modelId, ok
}

func (index *CodeIndex) Delete(modelId ModelId) {
	currentIndexable, ok := index.byIdentifier[modelId]
	if !ok {
		return
	}

	for _, code := range currentIndexable {
		delete(index.byCode, code)
	}
	delete(index.byIdentifier, modelId)
}
