package hooks

type Type uint8

const (
	/* Warning: There's two kinds of Hooks, simple and complex
	Simple Hooks needs to be sorted and are returned in sequence
	(when we request AfterCreate, we also get AfterSave)
	Complex controls can't be handled with generic methods,
	they need to be defined after the simple ones
	*/
	AfterCreate Type = iota
	AfterSave
	/* AfterAllStopVisitSave is a complex hook that can only be applied to StopVisits
	and creates a control on their VehicleJourney
	*/
	AfterAllStopVisitSave

	Total               = 3
	TotalSimpleControls = 2
)

var Hook = map[string]Type{
	"AfterCreate":           AfterCreate,
	"AfterSave":             AfterSave,
	"AfterAllStopVisitSave": AfterAllStopVisitSave,
}
