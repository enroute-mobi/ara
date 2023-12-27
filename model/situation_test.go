package model

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Situation_Id(t *testing.T) {
	situation := Situation{
		id: "6ba7b814-9dad-11d1-0-00c04fd430c8",
	}

	if situation.Id() != "6ba7b814-9dad-11d1-0-00c04fd430c8" {
		t.Errorf("situation.Id() returns wrong value, got: %s, required: %s", situation.Id(), "6ba7b814-9dad-11d1-0-00c04fd430c8")
	}
}

func Test_Situation_MarshalJSON(t *testing.T) {
	assert := assert.New(t)
	situation := Situation{
		id:     "6ba7b814-9dad-11d1-0-00c04fd430c8",
		Origin: "test",
	}

	situation.Description = &SituationTranslatedString{
		DefaultValue: "Joyeux Noel",
	}
	situation.Summary = &SituationTranslatedString{
		DefaultValue: "Noel",
	}

	affectStopArea := NewAffectedStopArea()
	affectStopArea.StopAreaId = "259344234"
	situation.Affects = append(situation.Affects, affectStopArea)

	affectLine := NewAffectedLine()
	affectLine.LineId = "222"
	affectDestinationId := "333"

	affectedDestination := &AffectedDestination{StopAreaId: StopAreaId(affectDestinationId)}
	affectLine.AffectedDestinations = append(affectLine.AffectedDestinations, affectedDestination)
	affectedSectionFirstStopId := "firstStop"
	affectedSectionLastStopId := "lastStop"
	affectedSection := &AffectedSection{
		FirstStop: StopAreaId(affectedSectionFirstStopId),
		LastStop:  StopAreaId(affectedSectionLastStopId),
	}
	affectLine.AffectedSections = append(affectLine.AffectedSections, affectedSection)
	affectedRoute := &AffectedRoute{RouteRef: "Route:66:LOC"}
	affectLine.AffectedRoutes = append(affectLine.AffectedRoutes, affectedRoute)
	situation.Affects = append(situation.Affects, affectLine)

	expected := `{
"Origin":"test",
"ValidityPeriods": null,
"PublicationWindows": null,
"Affects":[
{"Type":"StopArea","StopAreaId":"259344234"},
{"Type":"Line","LineId":"222",
"AffectedDestinations":[{"StopAreaId":"333"}],
"AffectedSections":[{"FirstStop":"firstStop","LastStop":"lastStop"}],
"AffectedRoutes":[{"RouteRef":"Route:66:LOC"}]}
],
"Description":{"DefaultValue":"Joyeux Noel"},
"Summary":{"DefaultValue":"Noel"},
"Id":"6ba7b814-9dad-11d1-0-00c04fd430c8"}`

	jsonBytes, err := situation.MarshalJSON()
	assert.Nil(err)
	assert.JSONEq(expected, string(jsonBytes))
}

func Test_Situation_UnmarshalJSON(t *testing.T) {
	assert := assert.New(t)
	text := `{
"Origin":"test",
"Codes": { "reflex": "FR:77491:ZDE:34004:STIF", "hastus": "sqypis" },
"Affects":[
{"Type":"StopArea","StopAreaId":"259344234"},
{"Type":"Line","LineId":"222","AffectedDestinations":[{"StopAreaId":"333"}],
"AffectedSections":[{"FirstStop":"firstStop","LastStop":"lastStop"}],
"AffectedRoutes":[{"RouteRef":"Route:66:LOC"}]}
],
"Description":{"DefaultValue":"Joyeux Noel"},
"Summary":{"DefaultValue":"Noel"},
"Id":"6ba7b814-9dad-11d1-0-00c04fd430c8"}`

	situation := &Situation{}
	err := json.Unmarshal([]byte(text), &situation)
	assert.Nil(err)

	expectedCodes := []Code{
		NewCode("reflex", "FR:77491:ZDE:34004:STIF"),
		NewCode("hastus", "sqypis"),
	}

	expectedSmmary := &SituationTranslatedString{
		DefaultValue: "Noel",
	}
	expectedDescription := &SituationTranslatedString{
		DefaultValue: "Joyeux Noel",
	}

	expectedAffectedStopArea := NewAffectedStopArea()
	expectedAffectedStopArea.StopAreaId = "259344234"

	expectedAffectedLine := NewAffectedLine()
	expectedAffectedLine.LineId = "222"
	affectedDestination := &AffectedDestination{StopAreaId: StopAreaId("333")}
	expectedAffectedLine.AffectedDestinations = append(expectedAffectedLine.AffectedDestinations, affectedDestination)

	affectedSection := &AffectedSection{
		FirstStop: StopAreaId("firstStop"),
		LastStop:  StopAreaId("lastStop"),
	}
	expectedAffectedLine.AffectedSections = append(expectedAffectedLine.AffectedSections, affectedSection)

	expectedAffectedRoute := &AffectedRoute{RouteRef: "Route:66:LOC"}
	expectedAffectedLine.AffectedRoutes = append(expectedAffectedLine.AffectedRoutes, expectedAffectedRoute)

	assert.Equal(expectedSmmary, situation.Summary)
	assert.Equal(expectedDescription, situation.Description)
	assert.Len(situation.Affects, 2)
	assert.Equal(expectedAffectedStopArea, situation.Affects[0])
	assert.Equal(expectedAffectedLine, situation.Affects[1])

	for _, expectedCode := range expectedCodes {
		code, found := situation.Code(expectedCode.CodeSpace())
		assert.True(found)
		assert.Equal(expectedCode, code)
	}
}

func Test_Situation_Save(t *testing.T) {
	model := NewMemoryModel()
	situation := model.Situations().New()
	code := NewCode("codeSpace", "value")
	situation.SetCode(code)

	if situation.model != model {
		t.Errorf("New situation model should be MemorySituation model")
	}

	ok := situation.Save()
	if !ok {
		t.Errorf("situation.Save() should succeed")
	}
	_, ok = model.Situations().Find(situation.Id())
	if !ok {
		t.Errorf("New situation should be found in MemorySituation")
	}
}

func Test_Situation_Code(t *testing.T) {
	situation := Situation{
		id: "6ba7b814-9dad-11d1-0-00c04fd430c8",
	}
	situation.codes = make(Codes)
	code := NewCode("codeSpace", "value")
	situation.SetCode(code)

	foundCode, ok := situation.Code("codeSpace")
	if !ok {
		t.Errorf("Code should return true if Code exists")
	}
	if foundCode.Value() != code.Value() {
		t.Errorf("Code should return a correct Code:\n got: %v\n want: %v", foundCode, code)
	}

	_, ok = situation.Code("wrongkind")
	if ok {
		t.Errorf("Code should return false if Code doesn't exist")
	}

	if len(situation.Codes()) != 1 {
		t.Errorf("Codes should return an array with set Codes, got: %v", situation.Codes())
	}
}

func Test_MemorySituations_New(t *testing.T) {
	situations := NewMemorySituations()

	situation := situations.New()
	if situation.Id() != "" {
		t.Errorf("New situation identifier should be an empty string, got: %s", situation.Id())
	}
}

func Test_MemorySituations_Save(t *testing.T) {
	situations := NewMemorySituations()

	situation := situations.New()

	if success := situations.Save(&situation); !success {
		t.Errorf("Save should return true")
	}

	if situation.Id() == "" {
		t.Errorf("New situation identifier shouldn't be an empty string")
	}
}

func Test_MemorySituations_Find_NotFound(t *testing.T) {
	situations := NewMemorySituations()
	_, ok := situations.Find("6ba7b814-9dad-11d1-0-00c04fd430c8")
	if ok {
		t.Errorf("Find should return false when Situation isn't found")
	}
}

func Test_MemorySituations_Find(t *testing.T) {
	situations := NewMemorySituations()

	existingSituation := situations.New()
	situations.Save(&existingSituation)

	situationId := existingSituation.Id()

	situation, ok := situations.Find(situationId)
	if !ok {
		t.Errorf("Find should return true when situation is found")
	}
	if situation.Id() != situationId {
		t.Errorf("Find should return a situation with the given Id")
	}
}

func Test_MemorySituations_FindAll(t *testing.T) {
	situations := NewMemorySituations()

	for i := 0; i < 5; i++ {
		existingSituation := situations.New()
		situations.Save(&existingSituation)
	}

	foundSituations := situations.FindAll()

	if len(foundSituations) != 5 {
		t.Errorf("FindAll should return all situations")
	}
}

func Test_MemorySituations_Delete(t *testing.T) {
	situations := NewMemorySituations()
	existingSituation := situations.New()
	code := NewCode("codeSpace", "value")
	existingSituation.SetCode(code)
	situations.Save(&existingSituation)

	situations.Delete(&existingSituation)

	_, ok := situations.Find(existingSituation.Id())
	if ok {
		t.Errorf("Deleted situation should not be findable")
	}
}
