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

	objectid3 := model.NewObjectID("remote_objectid_kind", "lineRef1")
	line := referential.Model().Lines().New()
	line.SetObjectID(objectid3)
	line.Save()
	lineId := line.Id()

	objectid4 := model.NewObjectID("remote_objectid_kind", "destinationRef1")
	destinationRef1 := referential.Model().StopAreas().New()
	destinationRef1.SetObjectID(objectid4)
	destinationRef1.Save()

	objectid5 := model.NewObjectID("remote_objectid_kind", "destinationRef2")
	destinationRef2 := referential.Model().StopAreas().New()
	destinationRef2.SetObjectID(objectid5)
	destinationRef2.Save()

	builder := NewGeneralMessageUpdateEventBuilder(partner)
	events := &[]*model.SituationUpdateEvent{}

	builder.buildGeneralMessageUpdateEvent(events, response.XMLGeneralMessages()[0], "producer")
	assert.Len(*events, 1, "One event should have been created")

	event := (*events)[0]
	if event.Format != "FRANCE" {
		t.Errorf("Wrong Format, expected: FRANCE, got: %v", event.Format)
	}

	assert.ElementsMatch([]string{"Commercial"}, event.Keywords)
	assert.Equal(model.ReportType("general"), event.ReportType)
	assert.Equal("test", event.Description.DefaultValue)
	assert.Nil(event.Summary)

	affects := event.Affects
	assert.Len(affects, 3)

	// Affected Lines
	assert.Equal("Line", affects[0].GetType())
	assert.Equal(model.ModelId(lineId), affects[0].GetId())
	assert.Equal(destinationRef1.Id(), affects[0].(*model.AffectedLine).AffectedDestinations[0].StopAreaId)
	assert.Equal(destinationRef2.Id(), affects[0].(*model.AffectedLine).AffectedDestinations[1].StopAreaId)
	// PLEASE CHANGE ME !!!!!
	assert.Len(affects[0].(*model.AffectedLine).AffectedSections, 0)

	// Affected StopAreas
	assert.Equal("StopArea", affects[1].GetType())
	assert.Equal(model.ModelId(stopAreaId), affects[1].GetId())
	assert.Equal("StopArea", affects[2].GetType())
	assert.Equal(model.ModelId(stopArea2Id), affects[2].GetId())

	// if len(event.SituationAttributes.LineSections) != 2 {
	// 	t.Fatalf("Wrong number of LineSections, expected: 2, got: %v", len(event.SituationAttributes.LineSections))
	// }
	// firstLineSection := *event.SituationAttributes.LineSections[0]
	// if ref, _ := firstLineSection.Get("FirstStop"); ref.ObjectId.Value() != "firstStop1" {
	// 	t.Errorf("Wrong first LineSection FirstStop: %v", ref)
	// }
	// if ref, _ := firstLineSection.Get("LastStop"); ref.ObjectId.Value() != "lastStop1" {
	// 	t.Errorf("Wrong first LineSection LastStop: %v", ref)
	// }
	// if ref, _ := firstLineSection.Get("LineRef"); ref.ObjectId.Value() != "lineSectionRef1" {
	// 	t.Errorf("Wrong first LineSection LineRef: %v", ref)
	// }
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
					Type:       "StopArea",
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
					Type:   "Line",
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
				Type:   "Line",
				LineId: line.Id(),
			},
			message: "Should not create an AffectedDestination for unknown StopArea",
		},
		{
			StopPointRef: "destinationRef",
			expectedAffectedLine: &model.AffectedLine{
				Type:   "Line",
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
