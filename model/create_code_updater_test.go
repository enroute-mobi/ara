package model

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Macro_CreateCode_StopArea(t *testing.T) {
	assert := assert.New(t)

	model := newTestModel(t)
	manager := model.Macros().(*macroManager)
	attributes := `{"source_code_space": "internal", "target_code_space": "external", "target_pattern": "prefix:%{value}:suffix"}`

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

	cb := &macroContextBuilder{
		childrenId: "",
		macro:      nil,
		updaters:   []*SelectMacro{sm},
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

	code1 := NewCode("internal", "test1")

	sa := model.StopAreas().New()
	sa.SetCode(code1)
	sa.Save()

	code2 := NewCode("internal", "test2")
	regionalCode := NewCode("external", "test")

	sa2 := model.StopAreas().New()
	sa2.SetCode(code2)
	sa2.SetCode(regionalCode)
	sa2.Save()

	code3 := NewCode("internal", "test3")

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
	_, ok = updatedSA1.Code("external")
	assert.False(ok)

	updatedSA2, ok := model.StopAreas().FindByCode(code2)
	assert.True(ok)
	foundRegionalCode, _ := updatedSA2.Code("external")
	assert.Equal(regionalCode.Value(), foundRegionalCode.Value())

	updatedSA3, ok := model.StopAreas().FindByCode(code3)
	assert.True(ok)
	foundRegionalCode, _ = updatedSA3.Code("external")
	assert.Equal("prefix:test3:suffix", foundRegionalCode.Value())
}

func Test_Macro_CreateCode_Line(t *testing.T) {
	assert := assert.New(t)

	model := newTestModel(t)
	manager := model.Macros().(*macroManager)
	attributes := `{"source_code_space": "internal", "target_code_space": "external", "target_pattern": "prefix:%{value}:suffix"}`

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

	cb := &macroContextBuilder{
		childrenId: "",
		macro:      nil,
		updaters:   []*SelectMacro{sm},
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

	code1 := NewCode("internal", "test1")

	line := model.Lines().New()
	line.SetCode(code1)
	line.Save()

	code2 := NewCode("internal", "test2")
	regionalCode := NewCode("external", "test")

	line2 := model.Lines().New()
	line2.SetCode(code2)
	line2.SetCode(regionalCode)
	line2.Save()

	code3 := NewCode("internal", "test3")

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

	updatedLine, ok := model.Lines().FindByCode(code1)
	assert.True(ok)
	if ok {
		_, ok = updatedLine.Code("external")
		assert.False(ok)
	}

	updatedLine2, ok := model.Lines().FindByCode(code2)
	assert.True(ok)
	if ok {
		foundRegionalCode, _ := updatedLine2.Code("external")
		assert.Equal(regionalCode.Value(), foundRegionalCode.Value())
	}

	updatedLine3, ok := model.Lines().FindByCode(code3)
	assert.True(ok)
	if ok {
		foundRegionalCode, _ := updatedLine3.Code("external")
		assert.Equal("prefix:test3:suffix", foundRegionalCode.Value())
	}
}
