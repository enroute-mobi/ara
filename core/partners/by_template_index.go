package partners

/* This is a very specific index, we can make it more generic by copying model.ByTemplateIndex */

type ByTemplateIndex struct {
	byTemplate   map[Id][]Id
	byIdentifier map[Id]Id
}

func NewByTemplateIndex() *ByTemplateIndex {
	return &ByTemplateIndex{
		byTemplate:   make(map[Id][]Id),
		byIdentifier: make(map[Id]Id),
	}
}

func (index *ByTemplateIndex) Index(p Partner) {
	pt := p.FromTemplate()
	if pt == "" {
		return
	}

	id := p.Id()

	currentByTemplate, ok := index.byIdentifier[id]
	if ok {
		if currentByTemplate != pt {
			index.removeFromIndexable(currentByTemplate, id)
		} else {
			return
		}
	}

	index.byTemplate[pt] = append(index.byTemplate[pt], id)
	index.byIdentifier[id] = pt
}

func (index *ByTemplateIndex) Find(pt Id) []Id {
	return index.byTemplate[pt]
}

func (index *ByTemplateIndex) Delete(id Id) {
	currentByTemplate, ok := index.byIdentifier[id]
	if !ok {
		return
	}

	index.removeFromIndexable(currentByTemplate, id)
	delete(index.byIdentifier, id)
}

func (index *ByTemplateIndex) IndexableLength(pt Id) int {
	return len(index.byTemplate[pt])
}

func (index *ByTemplateIndex) removeFromIndexable(pt, id Id) {
	if len(index.byTemplate[pt]) == 0 {
		return
	}
	for i, indexedId := range index.byTemplate[pt] {
		if indexedId == id {
			index.byTemplate[pt] = append(index.byTemplate[pt][:i], index.byTemplate[pt][i+1:]...)
			if len(index.byTemplate[pt]) == 0 {
				delete(index.byTemplate, pt)
			}
			return
		}
	}
}
