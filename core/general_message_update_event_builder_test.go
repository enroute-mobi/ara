package core

import (
	"io"
	"os"
	"testing"

	s "bitbucket.org/enroute-mobi/ara/core/settings"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/siri/sxml"
	"github.com/stretchr/testify/assert"
)

func Test_GeneralMessageUpdateEventBuilder_BuildGeneralMessageUpdateEvent(t *testing.T) {
	assert := assert.New(t)
	file, err := os.Open("testdata/long-general-message-response.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	content, err := io.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}

	response, _ := sxml.NewXMLGeneralMessageResponseFromContent(content)

	referentials := NewMemoryReferentials()
	referential := referentials.New("slug")
	referential.model = model.NewTestMemoryModel()
	referentials.Save(referential)

	partners := NewPartnerManager(referential)
	partner := partners.New("slug")
	settings := map[string]string{
		"remote_code_space": "remote_code_space",
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)
	partners.Save(partner)

	// StopPointRef
	code := model.NewCode("remote_code_space", "stopPointRef1")
	stopArea := referential.Model().StopAreas().New()
	stopArea.SetCode(code)
	stopArea.Save()
	stopAreaId := stopArea.Id()

	code2 := model.NewCode("remote_code_space", "stopPointRef2")
	stopArea2 := referential.Model().StopAreas().New()
	stopArea2.SetCode(code2)
	stopArea2.Save()
	stopArea2Id := stopArea2.Id()

	// LineRef
	code3 := model.NewCode("remote_code_space", "lineRef1")
	line := referential.Model().Lines().New()
	line.SetCode(code3)
	line.Save()
	lineId := line.Id()

	// Destinations
	code4 := model.NewCode("remote_code_space", "destinationRef1")
	destinationRef1 := referential.Model().StopAreas().New()
	destinationRef1.SetCode(code4)
	destinationRef1.Save()

	code5 := model.NewCode("remote_code_space", "destinationRef2")
	destinationRef2 := referential.Model().StopAreas().New()
	destinationRef2.SetCode(code5)
	destinationRef2.Save()

	// LineSections
	code6 := model.NewCode("remote_code_space", "lineSectionRef1")
	lineSectionRef1 := referential.Model().Lines().New()
	lineSectionRef1.SetCode(code6)
	lineSectionRef1.Save()

	code7 := model.NewCode("remote_code_space", "firstStop1")
	firstStop1 := referential.Model().StopAreas().New()
	firstStop1.SetCode(code7)
	firstStop1.Save()

	code8 := model.NewCode("remote_code_space", "lastStop1")
	lastStop1 := referential.Model().StopAreas().New()
	lastStop1.SetCode(code8)
	lastStop1.Save()

	code9 := model.NewCode("remote_code_space", "lineSectionRef2")
	lineSectionRef2 := referential.Model().Lines().New()
	lineSectionRef2.SetCode(code9)
	lineSectionRef2.Save()

	code10 := model.NewCode("remote_code_space", "firstStop2")
	firstStop2 := referential.Model().StopAreas().New()
	firstStop2.SetCode(code10)
	firstStop2.Save()

	code11 := model.NewCode("remote_code_space", "lastStop2")
	lastStop2 := referential.Model().StopAreas().New()
	lastStop2.SetCode(code11)
	lastStop2.Save()

	// Building
	builder := NewGeneralMessageUpdateEventBuilder(partner)
	events := NewCollectUpdateEvents()

	builder.buildGeneralMessageUpdateEvent(events, response.XMLGeneralMessages()[0], "producer")
	assert.Len(events.Situations, 1, "One event should have been created")

	event := events.Situations[0]
	assert.Equal("FRANCE", event.Format)
	assert.ElementsMatch([]string{"Commercial"}, event.Keywords)
	assert.Equal(model.ReportType("general"), event.ReportType)
	assert.Equal("test", event.Description.DefaultValue)
	assert.Nil(event.Summary)

	affects := event.Affects
	assert.Len(affects, 5, "Should have 5 affects: 3 affectedLines, 2 affectedStopAreas")

	// Affected Lines
	ok, affectedLine1 := event.TestFindAffectByLineId(lineId)
	assert.True(ok)

	// AffectedDestinations for LineRef1
	assert.Equal(destinationRef1.Id(), affectedLine1.AffectedDestinations[0].StopAreaId)
	assert.Equal(destinationRef2.Id(), affectedLine1.AffectedDestinations[1].StopAreaId)

	// AffectedRoutes for LineRef1
	assert.Equal("routeRef1", affectedLine1.AffectedRoutes[0].RouteRef)
	assert.Equal("routeRef2", affectedLine1.AffectedRoutes[1].RouteRef)

	// AffectedSections for LineSectionRef1
	ok, affectedLineSection1 := event.TestFindAffectByLineId(lineSectionRef1.Id())
	assert.True(ok)
	assert.Len(affectedLineSection1.AffectedSections, 1, "Should have 1 affectedSection for lineSectionRef1 ")
	assert.Equal(firstStop1.Id(), affectedLineSection1.AffectedSections[0].FirstStop)
	assert.Equal(lastStop1.Id(), affectedLineSection1.AffectedSections[0].LastStop)

	// AffectedSections for LineSectionRef2
	ok, affectedLineSection2 := event.TestFindAffectByLineId(lineSectionRef2.Id())
	assert.True(ok)
	assert.Len(affectedLineSection2.AffectedSections, 1, "Should have 1 affectedSection for lineSectionRef2")
	assert.Equal(firstStop2.Id(), affectedLineSection2.AffectedSections[0].FirstStop)
	assert.Equal(lastStop2.Id(), affectedLineSection2.AffectedSections[0].LastStop)

	// Affected StopAreas
	ok, _ = event.TestFindAffectByStopAreaId(stopAreaId)
	assert.True(ok)
	ok, _ = event.TestFindAffectByStopAreaId(stopArea2Id)
	assert.True(ok)
}

func Test_setReportType(t *testing.T) {
	assert := assert.New(t)
	partner := NewPartner()
	builder := NewGeneralMessageUpdateEventBuilder(partner)

	reportType := builder.setReportType("dummy")
	assert.Equal(model.SituationReportTypeGeneral, reportType)

	reportType = builder.setReportType("Perturbation")
	assert.Equal(model.SituationReportTypeIncident, reportType)
}

func Test_buildSituationAndDescriptionFromMessage(t *testing.T) {
	assert := assert.New(t)
	var TestCases = []struct {
		summary             *model.SituationTranslatedString
		messageType         string
		messageText         string
		expectedSummary     *model.SituationTranslatedString
		expectedDescription *model.SituationTranslatedString
		message             string
	}{
		{
			summary:     nil,
			messageType: "shortMessage",
			messageText: "a short message",
			expectedSummary: &model.SituationTranslatedString{
				DefaultValue: "a short message",
			},
			expectedDescription: nil,
			message:             "should set summary for shortMessage type",
		},
		{
			summary:         nil,
			messageType:     "longMessage",
			messageText:     "a long message",
			expectedSummary: nil,
			expectedDescription: &model.SituationTranslatedString{
				DefaultValue: "a long message",
			},
			message: "should set description for longMessage type",
		},
		{
			summary:         nil,
			messageType:     "dummy",
			messageText:     "A message with less than 160 characters and an emoji ⚠",
			expectedSummary: nil,
			expectedDescription: &model.SituationTranslatedString{
				DefaultValue: "A message with less than 160 characters and an emoji ⚠",
			},
			message: `for messageType other than shortMessage
should set description`,
		},
		{
			summary:     nil,
			messageType: "dummy",
			messageText: `Lorem ipsum dolor sit amet, consectetur adipiscing
 elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim
 veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.
Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore`,
			expectedSummary: nil,
			expectedDescription: &model.SituationTranslatedString{
				DefaultValue: `Lorem ipsum dolor sit amet, consectetur adipiscing
 elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim
 veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.
Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore`,
			},
			message: `for messageType other than shortMessage/longMessage
should set description if summary is not defined and text lenght > 160`,
		},
		{
			summary: &model.SituationTranslatedString{
				DefaultValue: "An existing summary ...",
			},
			messageType: "textOnly",
			messageText: "A message with less than 160 characters",
			expectedSummary: &model.SituationTranslatedString{
				DefaultValue: "An existing summary ...",
			},
			expectedDescription: &model.SituationTranslatedString{
				DefaultValue: "A message with less than 160 characters",
			},
			message: `When messageType is other than shortMessage
and summary is already defined, should keep existing summary and create description`,
		},
	}

	for _, tt := range TestCases {
		partner := NewPartner()
		builder := NewGeneralMessageUpdateEventBuilder(partner)

		event := &model.SituationUpdateEvent{}
		if tt.summary != nil {
			event.Summary = tt.summary
		}

		builder.buildSituationAndDescriptionFromMessage(
			tt.messageType,
			tt.messageText,
			event)
		assert.Equal(tt.expectedSummary, event.Summary, tt.message)
		assert.Equal(tt.expectedDescription, event.Description, tt.message)
	}
}

func Test_setAffectedStopArea(t *testing.T) {
	assert := assert.New(t)

	referentials := NewMemoryReferentials()
	referential := referentials.New("slug")
	referential.model = model.NewTestMemoryModel()
	referentials.Save(referential)

	partners := NewPartnerManager(referential)
	partner := partners.New("slug")
	settings := map[string]string{
		"remote_code_space": "remote_code_space",
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)
	partners.Save(partner)

	code := model.NewCode("remote_code_space", "stopPointRef1")
	stopArea := referential.Model().StopAreas().New()
	stopArea.SetCode(code)
	stopArea.Save()

	var TestCases = []struct {
		StopPointRef        string
		expectedEventAffect []model.Affect
		message             string
	}{
		{
			StopPointRef:        "dummy",
			expectedEventAffect: nil,
			message:             "Should not create an AffectedStopArea for unknown StopArea",
		},
		{
			StopPointRef: "stopPointRef1",
			expectedEventAffect: []model.Affect{
				&model.AffectedStopArea{
					StopAreaId: stopArea.Id(),
				},
			},
			message: "Should create an AffectedStopArea for known StopArea",
		},
	}
	for _, tt := range TestCases {
		event := &model.SituationUpdateEvent{}
		builder := NewGeneralMessageUpdateEventBuilder(partner)
		builder.setAffectedStopArea(event, tt.StopPointRef)
		assert.Equal(tt.expectedEventAffect, event.Affects, tt.message)
	}
}

func Test_setAffectedLine(t *testing.T) {
	assert := assert.New(t)

	referentials := NewMemoryReferentials()
	referential := referentials.New("slug")
	referential.model = model.NewTestMemoryModel()
	referentials.Save(referential)

	partners := NewPartnerManager(referential)
	partner := partners.New("slug")
	settings := map[string]string{
		"remote_code_space": "remote_code_space",
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)
	partners.Save(partner)

	code := model.NewCode("remote_code_space", "lineRef1")
	line := referential.Model().Lines().New()
	line.SetCode(code)
	line.Save()

	var TestCases = []struct {
		LineRef              string
		expectedAffectedLine *model.AffectedLine
		message              string
	}{
		{
			LineRef:              "dummy",
			expectedAffectedLine: nil,
			message:              "Should not create an AffectedLine for unknown Line",
		},
		{
			LineRef: "lineRef1",
			expectedAffectedLine: &model.AffectedLine{
				LineId: line.Id(),
			},
			message: "Should create an AffectedLine for known Line",
		},
	}

	for _, tt := range TestCases {
		affectedLines := make(map[model.LineId]*model.AffectedLine)
		builder := NewGeneralMessageUpdateEventBuilder(partner)
		builder.setAffectedLine(tt.LineRef, affectedLines)
		assert.Equal(tt.expectedAffectedLine, affectedLines[model.LineId(line.Id())], tt.message)
	}
}

func Test_setAffectedDestination(t *testing.T) {
	assert := assert.New(t)

	referentials := NewMemoryReferentials()
	referential := referentials.New("slug")
	referential.model = model.NewTestMemoryModel()
	referentials.Save(referential)

	partners := NewPartnerManager(referential)
	partner := partners.New("slug")
	settings := map[string]string{
		"remote_code_space": "remote_code_space",
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)
	partners.Save(partner)

	code := model.NewCode("remote_code_space", "destinationRef")
	stopArea := referential.Model().StopAreas().New()
	stopArea.SetCode(code)
	stopArea.Save()

	code2 := model.NewCode("remote_code_space", "lineRef")
	line := referential.Model().Lines().New()
	line.SetCode(code2)
	line.Save()

	var TestCases = []struct {
		StopPointRef         string
		expectedAffectedLine *model.AffectedLine
		message              string
	}{
		{
			StopPointRef: "dummy",
			expectedAffectedLine: &model.AffectedLine{
				LineId: line.Id(),
			},
			message: "Should not create an AffectedDestination for unknown StopArea",
		},
		{
			StopPointRef: "destinationRef",
			expectedAffectedLine: &model.AffectedLine{
				LineId: line.Id(),
				AffectedDestinations: []*model.AffectedDestination{
					&model.AffectedDestination{StopAreaId: stopArea.Id()},
				},
			},
			message: "Should create an AffectedDestination for known StopArea",
		},
	}

	for _, tt := range TestCases {
		builder := NewGeneralMessageUpdateEventBuilder(partner)
		affectedLines := make(map[model.LineId]*model.AffectedLine)
		affectedLine := model.NewAffectedLine()
		affectedLine.LineId = line.Id()
		affectedLines[line.Id()] = affectedLine
		builder.setAffectedDestination(line.Id(), tt.StopPointRef, affectedLines)
		assert.Equal(tt.expectedAffectedLine, affectedLine, tt.message)
	}
}

func Test_setAffectedSection(t *testing.T) {
	assert := assert.New(t)

	referentials := NewMemoryReferentials()
	referential := referentials.New("slug")
	referential.model = model.NewTestMemoryModel()
	referentials.Save(referential)

	partners := NewPartnerManager(referential)
	partner := partners.New("slug")
	settings := map[string]string{
		"remote_code_space": "remote_code_space",
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)
	partners.Save(partner)

	code := model.NewCode("remote_code_space", "firstStop")
	firstStop := referential.Model().StopAreas().New()
	firstStop.SetCode(code)
	firstStop.Save()

	code1 := model.NewCode("remote_code_space", "lastStop")
	lastStop := referential.Model().StopAreas().New()
	lastStop.SetCode(code1)
	lastStop.Save()

	code2 := model.NewCode("remote_code_space", "lineRef")
	line := referential.Model().Lines().New()
	line.SetCode(code2)
	line.Save()

	var TestCases = []struct {
		LineRef              string
		firstStop            string
		lastStop             string
		existingAffectedLine bool
		expectedAffectedLine *model.AffectedLine
		message              string
	}{

		{
			LineRef:              "DUMMY",
			firstStop:            "firstStop",
			lastStop:             "lastStop",
			existingAffectedLine: false,
			expectedAffectedLine: nil,
			message:              "Should not create an AffectedSection if lineRef does not exists",
		},
		{
			LineRef:              "lineRef",
			firstStop:            "DUMMY",
			lastStop:             "lastStop",
			existingAffectedLine: false,
			expectedAffectedLine: nil,
			message:              "Should not create an AffectedSection if firstStop does not exists",
		},
		{
			LineRef:              "lineRef",
			firstStop:            "firstStop",
			lastStop:             "DUMMY",
			existingAffectedLine: false,
			expectedAffectedLine: nil,
			message:              "Should not create an AffectedSection if lastStop does not exists",
		},
		{
			LineRef:              "lineRef",
			firstStop:            "firstStop",
			lastStop:             "lastStop",
			existingAffectedLine: false,
			expectedAffectedLine: &model.AffectedLine{
				LineId: line.Id(),
				AffectedSections: []*model.AffectedSection{
					&model.AffectedSection{FirstStop: firstStop.Id(), LastStop: lastStop.Id()},
				},
			},
			message: "Should create an AffectedSection if lineRef, firstStop and lastStop exists",
		},
		{
			LineRef:              "lineRef",
			firstStop:            "firstStop",
			lastStop:             "lastStop",
			existingAffectedLine: true,
			expectedAffectedLine: &model.AffectedLine{
				LineId: line.Id(),
				AffectedSections: []*model.AffectedSection{
					&model.AffectedSection{FirstStop: firstStop.Id(), LastStop: lastStop.Id()},
				},
				AffectedDestinations: []*model.AffectedDestination{
					&model.AffectedDestination{StopAreaId: firstStop.Id()},
				},
			},
			message: "Should add AffectedSection to existing AffectedLine",
		},
	}

	for _, tt := range TestCases {
		builder := NewGeneralMessageUpdateEventBuilder(partner)
		affectedLines := make(map[model.LineId]*model.AffectedLine)
		lineSection := LineSection{
			LineRef:   tt.LineRef,
			FirstStop: tt.firstStop,
			LastStop:  tt.lastStop,
		}

		if tt.existingAffectedLine {
			existingAffectedLine := &model.AffectedLine{
				LineId: line.Id(),
				AffectedDestinations: []*model.AffectedDestination{
					&model.AffectedDestination{StopAreaId: firstStop.Id()},
				},
			}
			affectedLines[line.Id()] = existingAffectedLine
		}

		builder.setAffectedSection(lineSection, affectedLines)

		if tt.expectedAffectedLine == nil {
			assert.Len(affectedLines, 0)
		} else {
			assert.Equal(tt.expectedAffectedLine, affectedLines[line.Id()], tt.message)
		}

	}
}
