package model

import (
	"database/sql"
	"testing"

	sattr "bitbucket.org/enroute-mobi/ara/siri/siri_attributes"
	"github.com/stretchr/testify/assert"
)

func Test_MacroBuilder_Ok(t *testing.T) {
	assert := assert.New(t)

	manager := NewMacroManager()

	sm := &SelectMacro{
		Id:              "id",
		ReferentialSlug: "referential",
		ContextId:       sql.NullString{String: "", Valid: false},
		Position:        0,
		Type:            "IfAttribute",
		ModelType:       sql.NullString{String: "VehicleJourney", Valid: true},
		Hook:            sql.NullString{String: "", Valid: false},
		Attributes:      sql.NullString{String: "{\"attribute_name\": \"DirectionName\", \"value\": \"Aller\"}", Valid: true},
	}

	cb := &macroContextBuilder{
		childrenId: "",
		macro:      sm,
	}

	builder := &macroBuilder{
		manager:        manager,
		initialContext: []*macroContextBuilder{cb},
		contexes:       make(map[string]*macroContextBuilder),
	}

	err := builder.buildMacros()
	assert.Len(err, 0)
}

func Test_MacroBuilder_NOk(t *testing.T) {
	assert := assert.New(t)

	manager := NewMacroManager()

	sm := &SelectMacro{
		Id:              "id",
		ReferentialSlug: "referential",
		ContextId:       sql.NullString{String: "", Valid: false},
		Position:        0,
		Type:            "IfAttribute",
		ModelType:       sql.NullString{String: "VehicleJourney", Valid: true},
		Hook:            sql.NullString{String: "", Valid: false},
		Attributes:      sql.NullString{String: "{\"value\": \"Aller\"}", Valid: true},
	}

	cb := &macroContextBuilder{
		childrenId: "",
		macro:      sm,
	}

	builder := &macroBuilder{
		manager:        manager,
		initialContext: []*macroContextBuilder{cb},
		contexes:       make(map[string]*macroContextBuilder),
	}

	err := builder.buildMacros()
	assert.Len(err, 1)
}

func Test_Macro_UpdateVehicleJourney(t *testing.T) {
	assert := assert.New(t)

	model := NewTestMemoryModel()
	manager := NewMacroManager()

	smc := &SelectMacro{
		Id:              "id",
		ReferentialSlug: "referential",
		ContextId:       sql.NullString{String: "", Valid: false},
		Position:        0,
		Type:            "IfAttribute",
		ModelType:       sql.NullString{String: "VehicleJourney", Valid: true},
		Hook:            sql.NullString{String: "", Valid: false},
		Attributes:      sql.NullString{String: "{\"attribute_name\": \"DirectionName\", \"value\": \"Aller\"}", Valid: true},
	}

	smu := &SelectMacro{
		Id:              "id2",
		ReferentialSlug: "referential",
		ContextId:       sql.NullString{String: "id", Valid: true},
		Position:        0,
		Type:            "SetAttribute",
		ModelType:       sql.NullString{String: "VehicleJourney", Valid: true},
		Hook:            sql.NullString{String: "", Valid: false},
		Attributes:      sql.NullString{String: "{\"attribute_name\": \"DirectionType\", \"value\": \"Outbound\"}", Valid: true},
	}

	cb := &macroContextBuilder{
		childrenId: "",
		macro:      smc,
		updaters:   []*SelectMacro{smu},
	}

	builder := &macroBuilder{
		manager:        manager,
		initialContext: []*macroContextBuilder{cb},
		contexes:       make(map[string]*macroContextBuilder),
	}

	err := builder.buildMacros()
	if len(err) != 0 {
		t.Fatal("Macro should be created: ", err)
	}
	model.macros = manager

	code := NewCode("codeSpace", "value")
	sa := model.StopAreas().New()
	sa.SetCode(code)
	sa.Save()

	l := model.Lines().New()
	l.SetCode(code)
	l.Save()

	vj := model.VehicleJourneys().New()
	vj.SetCode(code)
	vj.LineId = l.Id()
	vj.Save()

	updateManager := newUpdateManager(model)

	event := &VehicleJourneyUpdateEvent{
		Code:       code,
		LineCode:   code,
		attributes: NewAttributes(),
	}
	event.attributes.Set(sattr.DirectionName, "Aller")

	updateManager.Update(event)
	updatedVehicleJourney, ok := model.VehicleJourneys().FindByCode(code)
	assert.True(ok)

	assert.Equal(updatedVehicleJourney.DirectionType, "Outbound")
}
