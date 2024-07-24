package model

import (
	"bitbucket.org/enroute-mobi/ara/gtfs"
	"encoding/json"

	"github.com/stretchr/testify/assert"
	"testing"
	"time"
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

	situation.InternalTags = []string{"tag1"}
	// Affects
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

	// Consequences
	periodStartTime, _ := time.Parse(time.RFC3339, "2016-09-22T07:58:34+02:00")
	periodEndTime, _ := time.Parse(time.RFC3339, "2017-09-22T10:11:34+02:00")
	period := &TimeRange{
		StartTime: periodStartTime,
		EndTime:   periodEndTime,
	}
	var periods []*TimeRange
	periods = append(periods, period)
	consequence := &Consequence{Periods: periods}
	situation.Consequences = append(situation.Consequences, consequence)

	expected := `{
"Origin":"test",
"ValidityPeriods": null,
"PublicationWindows": null,
"InternalTags":["tag1"],
"Affects":[
{"Type":"StopArea","StopAreaId":"259344234"},
{"Type":"Line","LineId":"222",
"AffectedDestinations":[{"StopAreaId":"333"}],
"AffectedSections":[{"FirstStop":"firstStop","LastStop":"lastStop"}],
"AffectedRoutes":[{"RouteRef":"Route:66:LOC"}]}
],
"Consequences":[
{"Periods":[{"StartTime":"2016-09-22T07:58:34+02:00","EndTime":"2017-09-22T10:11:34+02:00"}]}
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
"InternalTags":["tag1"],
"Affects":[
{"Type":"StopArea","StopAreaId":"259344234"},
{"Type":"Line","LineId":"222","AffectedDestinations":[{"StopAreaId":"333"}],
"AffectedSections":[{"FirstStop":"firstStop","LastStop":"lastStop"}],
"AffectedRoutes":[{"RouteRef":"Route:66:LOC"}]}
],
"Consequences":[
{"Periods":[{"StartTime":"2016-09-22T07:58:34+02:00","EndTime":"2017-09-22T10:11:34+02:00"}]}
],
"Description":{"DefaultValue":"Joyeux Noel"},
"Summary":{"DefaultValue":"Noel"},
"Id":"6ba7b814-9dad-11d1-0-00c04fd430c8"}`

	apiSituation := &APISituation{}
	apiSituation.codes = make(Codes)

	err := json.Unmarshal([]byte(text), &apiSituation)
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

	assert.Equal(expectedSmmary, apiSituation.Summary)
	assert.Equal([]string{"tag1"}, apiSituation.InternalTags)
	assert.Equal(expectedDescription, apiSituation.Description)
	assert.Len(apiSituation.Affects, 2)
	assert.Equal(expectedAffectedStopArea, apiSituation.Affects[0])
	assert.Equal(expectedAffectedLine, apiSituation.Affects[1])

	for _, expectedCode := range expectedCodes {
		code, found := apiSituation.Code(expectedCode.CodeSpace())
		assert.True(found)
		assert.Equal(expectedCode, code)
	}
}

func Test_Situation_Save(t *testing.T) {
	model := NewTestMemoryModel()
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

func Test_Validate_Empty(t *testing.T) {
	assert := assert.New(t)
	situations := NewMemorySituations()
	situation := situations.New()
	apiSituation := situation.Definition()

	assert.False(apiSituation.Validate())

	assert.Equal([]string{"Can't be empty"}, apiSituation.Errors.Get("CodeSpace"))
	assert.Equal([]string{"Can't be empty"}, apiSituation.Errors.Get("SituationNumber"))
	assert.Equal([]string{"Can't be empty"}, apiSituation.Errors.Get("Affects"))
	assert.Equal([]string{"Can't be empty"}, apiSituation.Errors.Get("ValidityPeriods"))
}

func Test_Validate_Empty_Summary(t *testing.T) {
	assert := assert.New(t)
	situations := NewMemorySituations()
	situation := situations.New()
	situation.Summary = &SituationTranslatedString{
		DefaultValue: "",
	}
	apiSituation := situation.Definition()

	assert.False(apiSituation.Validate())
	assert.Equal([]string{"Can't be empty"}, apiSituation.Errors.Get("Summary"))
}

func Test_Validate_SituationAlreadyExists(t *testing.T) {
	assert := assert.New(t)
	situations := NewMemorySituations()
	situation := situations.New()
	apiSituation := situation.Definition()
	apiSituation.SituationNumber = "test"
	apiSituation.ExistingSituationCode = true
	apiSituation.Validate()

	assert.Equal([]string{"Is already in use"}, apiSituation.Errors.Get("SituationNumber"))
}

func Test_Validate_Sanitize_Summary(t *testing.T) {
	assert := assert.New(t)
	situations := NewMemorySituations()
	situation := situations.New()
	situation.Summary = &SituationTranslatedString{
		DefaultValue: "<script>alert('Boo!');</script>",
	}
	apiSituation := situation.Definition()
	apiSituation.Validate()
	assert.Equal("", apiSituation.Summary.DefaultValue, "Shoud saninitze the summary")
}

func Test_Validate_Sanitize_Description(t *testing.T) {
	assert := assert.New(t)
	situations := NewMemorySituations()
	situation := situations.New()
	situation.Description = &SituationTranslatedString{
		DefaultValue: "<script>alert('Boo!');</script>",
	}
	apiSituation := situation.Definition()
	apiSituation.Validate()
	assert.Equal("", apiSituation.Description.DefaultValue, "Shoud saninitze the summary")
}

func Test_Validate_ValidityPeriods_Without_StartTime(t *testing.T) {
	assert := assert.New(t)
	situations := NewMemorySituations()
	situation := situations.New()

	timeLayout := "2006/01/02-15:04:05"
	testTime, _ := time.Parse(timeLayout, "2007/01/02-15:04:05")
	period := &TimeRange{
		EndTime: testTime,
	}
	situation.ValidityPeriods = append(situation.ValidityPeriods, period)

	apiSituation := situation.Definition()
	apiSituation.Validate()

	assert.Equal([]string{"Can't be empty"}, apiSituation.Errors.Get("ValidityPeriods"))
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

func Test_AffectFromProto(t *testing.T) {
	assert := assert.New(t)
	model := NewTestMemoryModel()

	stopArea := model.StopAreas().New()
	code := NewCode("external", "A")
	stopArea.SetCode(code)
	stopArea.Save()
	stopAreaA := "A"

	line := model.Lines().New()
	code = NewCode("external", "1")
	line.SetCode(code)
	line.Save()
	lineValue := "1"

	unknownModel := "WRONG"

	var TestCases = []struct {
		entity                 *gtfs.EntitySelector
		remoteCodeSpace        string
		valid                  bool
		expectedAffect         Affect
		expectedMonitoringRefs []string
		expectedLineRefs       []string
		message                string
	}{
		{
			entity: &gtfs.EntitySelector{
				StopId: &stopAreaA,
			},
			remoteCodeSpace: "external",
			valid:           true,
			expectedAffect: &AffectedStopArea{
				StopAreaId: stopArea.Id(),
			},
			expectedMonitoringRefs: []string{"A"},
			expectedLineRefs:       []string{},
			message: `EntitySelector with only a StopId
should create an affectedStopArea`,
		},
		{
			entity: &gtfs.EntitySelector{
				RouteId: &lineValue,
			},
			remoteCodeSpace: "external",
			valid:           true,
			expectedAffect: &AffectedLine{
				LineId: line.Id(),
			},
			expectedMonitoringRefs: []string{},
			expectedLineRefs:       []string{"1"},
			message: `EntitySelector with only a LineId
should create an affectedLine`,
		},
		{
			entity: &gtfs.EntitySelector{
				StopId:  &stopAreaA,
				RouteId: &lineValue,
			},
			remoteCodeSpace: "external",
			valid:           true,
			expectedAffect: &AffectedStopArea{
				StopAreaId: stopArea.Id(),
				LineIds:    []LineId{line.Id()},
			},
			expectedMonitoringRefs: []string{"A"},
			expectedLineRefs:       []string{"1"},
			message: `EntitySelector with valid StopId and
LineId should create an affectedStopArea with LineId in LineIds`,
		},
		{
			entity: &gtfs.EntitySelector{
				StopId:  &stopAreaA,
				RouteId: &unknownModel,
			},
			remoteCodeSpace: "external",
			valid:           true,
			expectedAffect: &AffectedStopArea{
				StopAreaId: stopArea.Id(),
			},
			expectedMonitoringRefs: []string{"A"},
			expectedLineRefs:       []string{},
			message: `EntitySelector with valid StopId and unknown
LineId should create an affectedStopArea without LineIds`,
		},
		{
			entity: &gtfs.EntitySelector{
				StopId:  &unknownModel,
				RouteId: &lineValue,
			},
			remoteCodeSpace: "external",
			valid:           false,
			expectedAffect:  nil,
			message: `EntitySelector with unknow StopId should
not create any affect`,
		},
		{
			entity:          &gtfs.EntitySelector{},
			remoteCodeSpace: "external",
			valid:           false,
			expectedAffect:  nil,
			message: `EntitySelector empty should not create
any affect`,
		},
	}

	for _, tt := range TestCases {
		affect, collectedRefs, err := AffectFromProto(tt.entity, tt.remoteCodeSpace, model)
		if !tt.valid {
			assert.Error(err)
			continue
		}
		assert.Nil(err)
		assert.Equalf(tt.expectedAffect, affect, tt.message)
		assert.Equal(tt.expectedMonitoringRefs, GetReferencesSlice(collectedRefs.MonitoringRefs))
		assert.Equal(tt.expectedLineRefs, GetReferencesSlice(collectedRefs.LineRefs))
	}
}

func GetReferencesSlice(refs map[string]struct{}) []string {
	refSlice := make([]string, len(refs))
	i := 0
	for ref := range refs {
		refSlice[i] = ref
		i++
	}
	return refSlice
}

func Test_AffectToProto(t *testing.T) {
	assert := assert.New(t)
	model := NewTestMemoryModel()

	stopArea := model.StopAreas().New()
	code := NewCode("external", "A")
	stopArea.SetCode(code)
	stopArea.Save()
	stopAreaValue := "A"

	line := model.Lines().New()
	code = NewCode("external", "1")
	line.SetCode(code)
	line.Save()
	lineValue := "1"

	wrongStopArea := model.StopAreas().New()
	code = NewCode("WRONG", "B")
	wrongStopArea.SetCode(code)
	wrongStopArea.Save()

	wrongLine := model.Lines().New()
	wrongLine.SetCode(code)
	wrongLine.Save()

	var TestCases = []struct {
		affect                Affect
		remoteCodeSpace       string
		valid                 bool
		expectedStopId        *string
		expectedRouteId       *string
		expectedCollectedRefs *AffectRefs
		message               string
	}{
		{
			affect: &AffectedLine{
				LineId: line.Id(),
			},
			remoteCodeSpace: "external",
			valid:           true,
			expectedStopId:  nil,
			expectedRouteId: &lineValue,
			message:         `AffectedLine with valid line should create RouteId`,
		},
		{
			affect: &AffectedStopArea{
				StopAreaId: stopArea.Id(),
			},
			remoteCodeSpace: "external",
			valid:           true,
			expectedStopId:  &stopAreaValue,
			expectedRouteId: nil,
			message:         `AffectedStopArea with valid stopArea should create StopId`,
		},
		{
			affect: &AffectedStopArea{
				StopAreaId: stopArea.Id(),
				LineIds:    []LineId{line.Id()},
			},
			remoteCodeSpace: "external",
			valid:           true,
			expectedStopId:  &stopAreaValue,
			expectedRouteId: &lineValue,
			message: `AffectedStopArea with valid StopArea and LineIds should
create StopId and RouteId`,
		},
		{
			affect: &AffectedStopArea{
				StopAreaId: wrongStopArea.Id(),
			},
			remoteCodeSpace: "external",
			valid:           false,
			message:         `AffectedStopArea with unknown stopArea should be invalid`,
		},
		{
			affect: &AffectedStopArea{
				StopAreaId: stopArea.Id(),
				LineIds:    []LineId{wrongLine.Id()},
			},
			remoteCodeSpace: "external",
			valid:           true,
			expectedStopId:  &stopAreaValue,
			expectedRouteId: nil,
			message: `AffectedStopArea with valid stopArea and unknwon line
should create StopId only`,
		},
	}

	for _, tt := range TestCases {
		entitySelector, _, err := AffectToProto(tt.affect, tt.remoteCodeSpace, model)
		if !tt.valid {
			assert.Error(err)
			continue
		}

		assert.Nil(err)
		if tt.expectedStopId == nil {
			assert.Nil(entitySelector[0].StopId)
		}

		if tt.expectedRouteId == nil {
			assert.Nil(entitySelector[0].RouteId)
		}

		if tt.expectedRouteId != nil && tt.expectedStopId != nil {
			assert.Equal(tt.expectedStopId, entitySelector[0].StopId)
			assert.Equal(tt.expectedRouteId, entitySelector[0].RouteId)
		}
	}
}
