package model

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Macro_CreateCode_StopArea(t *testing.T) {
	assert := assert.New(t)

	model := NewTestMemoryModel()
	manager := NewMacroManager()
	attributes := `{"source_code_space": "sae", "target_code_space": "regional", "target_pattern": "prefix:%{value}:suffix"}`

	sm := &SelectMacro{
		Id:              "id2",
		ReferentialSlug: "referential",
		ContextId:       sql.NullString{String: "", Valid: false},
		Position:        0,
		Type:            CreateCode,
		ModelType:       sql.NullString{String: "StopArea", Valid: true},
		Hook:            sql.NullString{String: "AfterCreate", Valid: true},
		Attributes:      sql.NullString{String: attributes, Valid: true},
	}

	cb := &contextBuilder{
		childrenId: "",
		macro:      nil,
		updaters:   []*SelectMacro{sm},
	}

	builder := &macroBuilder{
		manager:        manager,
		initialContext: []*contextBuilder{cb},
		contexes:       make(map[string]*contextBuilder),
	}

	err := builder.buildMacros()
	if len(err) != 0 {
		t.Fatal("Macro should be created: ", err)
	}
	model.macros = manager

	code1 := NewCode("sae", "test1")

	sa := model.StopAreas().New()
	sa.SetCode(code1)
	sa.Save()

	code2 := NewCode("sae", "test2")
	regionalCode := NewCode("regional", "test")

	sa2 := model.StopAreas().New()
	sa2.SetCode(code2)
	sa2.SetCode(regionalCode)
	sa2.Save()

	code3 := NewCode("sae", "test3")

	updateManager := newUpdateManager(model)

	event1 := &StopAreaUpdateEvent{
		Code: code1,
		Name: "Test 1",
	}
	event2 := &StopAreaUpdateEvent{
		Code: code2,
		Name: "Test 2",
	}
	event3 := &StopAreaUpdateEvent{
		Code: code3,
		Name: "Test 3",
	}

	updateManager.Update(event1)
	updateManager.Update(event2)
	updateManager.Update(event3)

	updatedSA1, ok := model.StopAreas().FindByCode(code1)
	assert.True(ok)
	_, ok = updatedSA1.Code("regional")
	assert.False(ok)

	updatedSA2, ok := model.StopAreas().FindByCode(code2)
	assert.True(ok)
	foundRegionalCode, _ := updatedSA2.Code("regional")
	assert.Equal(regionalCode.Value(), foundRegionalCode.Value())

	updatedSA3, ok := model.StopAreas().FindByCode(code3)
	assert.True(ok)
	foundRegionalCode, _ = updatedSA3.Code("regional")
	assert.Equal("prefix:test3:suffix", foundRegionalCode.Value())
}

func Test_Macro_CreateCode_Line(t *testing.T) {
	assert := assert.New(t)

	model := NewTestMemoryModel()
	manager := NewMacroManager()
	attributes := `{"source_code_space": "sae", "target_code_space": "regional", "target_pattern": "prefix:%{value}:suffix"}`

	sm := &SelectMacro{
		Id:              "id2",
		ReferentialSlug: "referential",
		ContextId:       sql.NullString{String: "", Valid: false},
		Position:        0,
		Type:            CreateCode,
		ModelType:       sql.NullString{String: "Line", Valid: true},
		Hook:            sql.NullString{String: "AfterCreate", Valid: true},
		Attributes:      sql.NullString{String: attributes, Valid: true},
	}

	cb := &contextBuilder{
		childrenId: "",
		macro:      nil,
		updaters:   []*SelectMacro{sm},
	}

	builder := &macroBuilder{
		manager:        manager,
		initialContext: []*contextBuilder{cb},
		contexes:       make(map[string]*contextBuilder),
	}

	err := builder.buildMacros()
	if len(err) != 0 {
		t.Fatal("Macro should be created: ", err)
	}
	model.macros = manager

	code1 := NewCode("sae", "test1")

	sa := model.Lines().New()
	sa.SetCode(code1)
	sa.Save()

	code2 := NewCode("sae", "test2")
	regionalCode := NewCode("regional", "test")

	sa2 := model.Lines().New()
	sa2.SetCode(code2)
	sa2.SetCode(regionalCode)
	sa2.Save()

	code3 := NewCode("sae", "test3")

	updateManager := newUpdateManager(model)

	event1 := &LineUpdateEvent{
		Code: code1,
		Name: "Test 1",
	}
	event2 := &LineUpdateEvent{
		Code: code2,
		Name: "Test 2",
	}
	event3 := &LineUpdateEvent{
		Code: code3,
		Name: "Test 3",
	}

	updateManager.Update(event1)
	updateManager.Update(event2)
	updateManager.Update(event3)

	updatedSA1, ok := model.Lines().FindByCode(code1)
	assert.True(ok)
	_, ok = updatedSA1.Code("regional")
	assert.False(ok)

	updatedSA2, ok := model.Lines().FindByCode(code2)
	assert.True(ok)
	foundRegionalCode, _ := updatedSA2.Code("regional")
	assert.Equal(regionalCode.Value(), foundRegionalCode.Value())

	updatedSA3, ok := model.Lines().FindByCode(code3)
	assert.True(ok)
	foundRegionalCode, _ = updatedSA3.Code("regional")
	assert.Equal("prefix:test3:suffix", foundRegionalCode.Value())
}
