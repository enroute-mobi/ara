package model

/*
This struct is meant to be embedded in managers to handle multiple indexes
without fearing to forget one in the save or delete methods

Warning: This is an unsafe package that:
  - doesn't check if a particular index exists
  - doesn't protect its maps with a mutex

It's meant for internal use only and these matters needs to be checked elsewere
*/
const (
	OneToOne = iota
	OneToMany

	ByLine
	ByStopArea
	ByVehicleJourney
	ByStopVisit
	ByParent
	ByReferent
)

type IndexableExtractor func(ModelInstance) ModelId

type Index interface {
	Index(ModelInstance)
	Find(ModelId) ([]ModelId, bool)
	FindOne(ModelId) (ModelId, bool)
	Delete(ModelId)
	IndexableLength(ModelId) int
}

type IndexHandler struct {
	i  map[int]Index
	ci *codeIndex
}

// InitHandler creates a CodeIndex by default but it can be disabled
func (is *IndexHandler) InitIndexes(hasCodes ...bool) {
	is.i = make(map[int]Index)
	if len(hasCodes) == 0 || hasCodes[0] {
		is.ci = NewCodeIndex()
	}
}

func (is *IndexHandler) AddIndex(name, kind int, extractor IndexableExtractor) {
	switch kind {
	case OneToOne:
		is.i[name] = NewSimpleIndex(extractor)
	case OneToMany:
		is.i[name] = NewIndex(extractor)
	}
}

func (is *IndexHandler) Index(m ModelInstance) {
	for _, v := range is.i {
		v.Index(m)
	}
	if is.ci != nil {
		is.ci.Index(m)
	}
}

func (is *IndexHandler) GetIndex(n int) Index {
	return is.i[n]
}

func (is *IndexHandler) FindBy(n int, i ModelId) ([]ModelId, bool) {
	return is.i[n].Find(i)
}

func (is *IndexHandler) FindOneBy(n int, i ModelId) (ModelId, bool) {
	return is.i[n].FindOne(i)
}

func (is *IndexHandler) IndexableLength(n int, i ModelId) int {
	return is.i[n].IndexableLength(i)
}

func (is *IndexHandler) ByCode() *codeIndex {
	return is.ci
}

func (is *IndexHandler) Deindex(mid ModelId) {
	for _, v := range is.i {
		v.Delete(mid)
	}
	if is.ci != nil {
		is.ci.Delete(mid)
	}
}
