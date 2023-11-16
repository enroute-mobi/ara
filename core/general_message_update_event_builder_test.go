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
	referential.model = model.NewMemoryModel()
	referentials.Save(referential)

	partners := NewPartnerManager(referential)
	partner := partners.New("slug")
	settings := map[string]string{
		"remote_objectid_kind": "remote_objectid_kind",
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)
	partners.Save(partner)

	// StopPointRef
	objectid := model.NewObjectID("remote_objectid_kind", "stopPointRef1")
	stopArea := referential.Model().StopAreas().New()
	stopArea.SetObjectID(objectid)
	stopArea.Save()
	stopAreaId := stopArea.Id()

	objectid2 := model.NewObjectID("remote_objectid_kind", "stopPointRef2")
	stopArea2 := referential.Model().StopAreas().New()
	stopArea2.SetObjectID(objectid2)
	stopArea2.Save()
	stopArea2Id := stopArea2.Id()

	// LineRef
	objectid3 := model.NewObjectID("remote_objectid_kind", "lineRef1")
	line := referential.Model().Lines().New()
	line.SetObjectID(objectid3)
	line.Save()
	lineId := line.Id()

	// Destinations
	objectid4 := model.NewObjectID("remote_objectid_kind", "destinationRef1")
	destinationRef1 := referential.Model().StopAreas().New()
	destinationRef1.SetObjectID(objectid4)
	destinationRef1.Save()

	objectid5 := model.NewObjectID("remote_objectid_kind", "destinationRef2")
	destinationRef2 := referential.Model().StopAreas().New()
	destinationRef2.SetObjectID(objectid5)
	destinationRef2.Save()

	// LineSections
	objectid6 := model.NewObjectID("remote_objectid_kind", "lineSectionRef1")
	lineSectionRef1 := referential.Model().Lines().New()
	lineSectionRef1.SetObjectID(objectid6)
	lineSectionRef1.Save()

	objectid7 := model.NewObjectID("remote_objectid_kind", "firstStop1")
	firstStop1 := referential.Model().StopAreas().New()
	firstStop1.SetObjectID(objectid7)
	firstStop1.Save()

	objectid8 := model.NewObjectID("remote_objectid_kind", "lastStop1")
	lastStop1 := referential.Model().StopAreas().New()
	lastStop1.SetObjectID(objectid8)
	lastStop1.Save()

	objectid9 := model.NewObjectID("remote_objectid_kind", "lineSectionRef2")
	lineSectionRef2 := referential.Model().Lines().New()
	lineSectionRef2.SetObjectID(objectid9)
	lineSectionRef2.Save()

	objectid10 := model.NewObjectID("remote_objectid_kind", "firstStop2")
	firstStop2 := referential.Model().StopAreas().New()
	firstStop2.SetObjectID(objectid10)
	firstStop2.Save()

	objectid11 := model.NewObjectID("remote_objectid_kind", "lastStop2")
	lastStop2 := referential.Model().StopAreas().New()
	lastStop2.SetObjectID(objectid11)
	lastStop2.Save()

	// Building
	builder := NewGeneralMessageUpdateEventBuilder(partner)
	events := &[]*model.SituationUpdateEvent{}

	builder.buildGeneralMessageUpdateEvent(events, response.XMLGeneralMessages()[0], "producer")
	assert.Len(*events, 1, "One event should have been created")

	event := (*events)[0]
	assert.Equal("FRANCE", event.Format)
	assert.ElementsMatch([]string{"Commercial"}, event.Keywords)
	assert.Equal(model.ReportType("general"), event.ReportType)
	assert.Equal("test", event.Description.DefaultValue)
	assert.Nil(event.Summary)

	affects := event.Affects
	assert.Len(affects, 5, "Should have 5 affetcs: 3 affctedLines, 2 affectedStopAreas")

	// Affected Lines
	assert.Equal(model.SituationType("Line"), affects[0].GetType())
	assert.Equal(model.ModelId(lineId), affects[0].GetId(), "Should be Id of lineRef1")

	assert.Equal(destinationRef1.Id(), affects[0].(*model.AffectedLine).AffectedDestinations[0].StopAreaId)
	assert.Equal(destinationRef2.Id(), affects[0].(*model.AffectedLine).AffectedDestinations[1].StopAreaId)

	// AffectedSections
	assert.Len(affects[0].(*model.AffectedLine).AffectedSections, 0, "Should have no affected section for lineRef1")
	assert.Len(affects[1].(*model.AffectedLine).AffectedSections, 1, "Should have 1 affectedSection for lineSectionRef1 ")
	assert.Len(affects[2].(*model.AffectedLine).AffectedSections, 1, "Should have 1 affecteSection for lineSection2")

	affectedSectionLineSection1 := affects[1].(*model.AffectedLine).AffectedSections[0]
	assert.Equal(firstStop1.Id(), affectedSectionLineSection1.FirstStop)
	assert.Equal(lastStop1.Id(), affectedSectionLineSection1.LastStop)

	affectedSectionLineSection2 := affects[2].(*model.AffectedLine).AffectedSections[0]
	assert.Equal(firstStop2.Id(), affectedSectionLineSection2.FirstStop)
	assert.Equal(lastStop2.Id(), affectedSectionLineSection2.LastStop)

	// Affected StopAreas
	assert.Equal(model.SituationType("StopArea"), affects[3].GetType())
	assert.Equal(model.ModelId(stopAreaId), affects[3].GetId())
	assert.Equal(model.SituationType("StopArea"), affects[4].GetType())
	assert.Equal(model.ModelId(stopArea2Id), affects[4].GetId())
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
			summary:     nil,
			messageType: "dummy",
			messageText: "A message < 160 characters",
			expectedSummary: &model.SituationTranslatedString{
				DefaultValue: "A message < 160 characters",
			},
			expectedDescription: nil,
			message: `for messageType other than shortMessage/longMessage
should set summary if summary is not defined and text lenght < 160`,
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
			messageText: "A message < 160 characters",
			expectedSummary: &model.SituationTranslatedString{
				DefaultValue: "An existing summary ...",
			},
			expectedDescription: &model.SituationTranslatedString{
				DefaultValue: "A message < 160 characters",
			},
			message: `When messageType is other than shortMessage/longMessage
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
	referential.model = model.NewMemoryModel()
	referentials.Save(referential)

	partners := NewPartnerManager(referential)
	partner := partners.New("slug")
	settings := map[string]string{
		"remote_objectid_kind": "remote_objectid_kind",
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)
	partners.Save(partner)

	objectid := model.NewObjectID("remote_objectid_kind", "stopPointRef1")
	stopArea := referential.Model().StopAreas().New()
	stopArea.SetObjectID(objectid)
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
	referential.model = model.NewMemoryModel()
	referentials.Save(referential)

	partners := NewPartnerManager(referential)
	partner := partners.New("slug")
	settings := map[string]string{
		"remote_objectid_kind": "remote_objectid_kind",
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)
	partners.Save(partner)

	objectid := model.NewObjectID("remote_objectid_kind", "lineRef1")
	line := referential.Model().Lines().New()
	line.SetObjectID(objectid)
	line.Save()

	var TestCases = []struct {
		LineRef             string
		expectedEventAffect []model.Affect
		message             string
	}{
		{
			LineRef:             "dummy",
			expectedEventAffect: nil,
			message:             "Should not create an AffectedLine for unknown Line",
		},
		{
			LineRef: "lineRef1",
			expectedEventAffect: []model.Affect{
				&model.AffectedLine{
					LineId: line.Id(),
				},
			},
			message: "Should create an AffectedLine for known Line",
		},
	}
	for _, tt := range TestCases {
		event := &model.SituationUpdateEvent{}
		builder := NewGeneralMessageUpdateEventBuilder(partner)
		builder.setAffectedLine(event, tt.LineRef)
		assert.Equal(tt.expectedEventAffect, event.Affects, tt.message)
	}
}

func Test_setAffectedDestination(t *testing.T) {
	assert := assert.New(t)

	referentials := NewMemoryReferentials()
	referential := referentials.New("slug")
	referential.model = model.NewMemoryModel()
	referentials.Save(referential)

	partners := NewPartnerManager(referential)
	partner := partners.New("slug")
	settings := map[string]string{
		"remote_objectid_kind": "remote_objectid_kind",
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)
	partners.Save(partner)

	objectid := model.NewObjectID("remote_objectid_kind", "destinationRef")
	stopArea := referential.Model().StopAreas().New()
	stopArea.SetObjectID(objectid)
	stopArea.Save()

	objectid2 := model.NewObjectID("remote_objectid_kind", "lineRef")
	line := referential.Model().Lines().New()
	line.SetObjectID(objectid2)
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
		event := &model.SituationUpdateEvent{}
		builder := NewGeneralMessageUpdateEventBuilder(partner)
		affectedLine := model.NewAffectedLine()
		affectedLine.LineId = line.Id()
		builder.setAffectedDestination(event, tt.StopPointRef, affectedLine)
		assert.Equal(tt.expectedAffectedLine, affectedLine, tt.message)
	}
}

func Test_setAffectedSection(t *testing.T) {
	assert := assert.New(t)

	referentials := NewMemoryReferentials()
	referential := referentials.New("slug")
	referential.model = model.NewMemoryModel()
	referentials.Save(referential)

	partners := NewPartnerManager(referential)
	partner := partners.New("slug")
	settings := map[string]string{
		"remote_objectid_kind": "remote_objectid_kind",
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)
	partners.Save(partner)

	objectid := model.NewObjectID("remote_objectid_kind", "firstStop")
	firstStop := referential.Model().StopAreas().New()
	firstStop.SetObjectID(objectid)
	firstStop.Save()

	objectid1 := model.NewObjectID("remote_objectid_kind", "lastStop")
	lastStop := referential.Model().StopAreas().New()
	lastStop.SetObjectID(objectid1)
	lastStop.Save()

	objectid2 := model.NewObjectID("remote_objectid_kind", "lineRef")
	line := referential.Model().Lines().New()
	line.SetObjectID(objectid2)
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
		event := &model.SituationUpdateEvent{}
		builder := NewGeneralMessageUpdateEventBuilder(partner)
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
			event.Affects = append(event.Affects, existingAffectedLine)
		}

		builder.setAffectedSection(event, lineSection)

		if tt.expectedAffectedLine == nil {
			assert.Nil(event.Affects)
		} else {
			assert.Equal(tt.expectedAffectedLine, event.Affects[0].(*model.AffectedLine), tt.message)
		}

	}
}
