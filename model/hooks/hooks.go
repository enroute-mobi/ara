package hooks

type Type uint8

const (
	// Warning: Hooks needs to be sorted
	AfterCreate Type = iota
	AfterSave

	Total = 2
)

var Hook = map[string]Type{
	"AfterCreate": AfterCreate,
	"AfterSave":   AfterSave,
}
