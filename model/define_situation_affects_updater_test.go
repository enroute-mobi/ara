package model

import (
	"database/sql"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Macro_DefineSituationAffects(t *testing.T) {
	assert := assert.New(t)

	model := NewTestMemoryModel()
	manager := NewMacroManager()

	sm := &SelectMacro{
		Id:              "id2",
		ReferentialSlug: "referential",
		ContextId:       sql.NullString{String: "", Valid: false},
		Position:        0,
		Type:            DefineSituationAffects,
		ModelType:       sql.NullString{String: "Situation", Valid: true},
		Hook:            sql.NullString{String: "", Valid: false},
		Attributes:      sql.NullString{String: "{}", Valid: true},
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

	asa := AffectedStopArea{
		StopAreaId: "said",
	}

	al := AffectedLine{
		LineId: "lid",
	}

	al2 := AffectedLine{
		LineId: "lid",
	}

	code := NewCode("codeSpace", "value")

	c1 := &Consequence{
		Affects: []Affect{asa},
	}
	c2 := &Consequence{
		Affects: []Affect{al2},
	}

	s := model.Situations().New()
	s.SetCode(code)
	s.Save()

	updateManager := newUpdateManager(model)

	event := &SituationUpdateEvent{
		SituationCode: code,
		Version:       1,
		Affects:       []Affect{al},
		Consequences:  []*Consequence{c1, c2},
	}

	updateManager.Update(event)

	updatedSituation, ok := model.Situations().FindByCode(code)
	assert.True(ok)

	assert.Equal(2, len(updatedSituation.Affects))
}
