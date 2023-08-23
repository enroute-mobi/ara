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
	partner := referential.Partners().New("slug")

	settings := map[string]string{
		"remote_objectid_kind": "remote_objectid_kind",
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)
	builder := NewGeneralMessageUpdateEventBuilder(partner)

	events := &[]*model.SituationUpdateEvent{}

	builder.buildGeneralMessageUpdateEvent(events, response.XMLGeneralMessages()[0], "producer")
	assert.Len(*events, 1, "One event should have been created")

	event := (*events)[0]
	if event.SituationAttributes.Format != "FRANCE" {
		t.Errorf("Wrong Format, expected: FRANCE, got: %v", event.SituationAttributes.Format)
	}

	assert.ElementsMatch([]string{"Commercial"}, event.Keywords)
	assert.Equal(model.ReportType("general"), event.ReportType)
	assert.Equal("test", event.Description.DefaultValue)
	assert.Nil(event.Summary)

	if len(event.SituationAttributes.References) != 12 {
		t.Fatalf("Wrong number of References, expected: 12, got: %v", len(event.SituationAttributes.References))
	}
	if event.SituationAttributes.References[0].ObjectId.Value() != "lineRef1" {
		t.Errorf("Wrong first LineRef: %v", event.SituationAttributes.References[0])
	}
	if event.SituationAttributes.References[2].ObjectId.Value() != "stopPointRef1" {
		t.Errorf("Wrong first StopPointRef: %v", event.SituationAttributes.References[2])
	}
	if event.SituationAttributes.References[4].ObjectId.Value() != "journeyPatternRef1" {
		t.Errorf("Wrong first JourneyPatternRef: %v", event.SituationAttributes.References[4])
	}
	if event.SituationAttributes.References[6].ObjectId.Value() != "destinationRef1" {
		t.Errorf("Wrong first DestinationRef: %v", event.SituationAttributes.References[6])
	}
	if event.SituationAttributes.References[8].ObjectId.Value() != "routeRef1" {
		t.Errorf("Wrong first RouteRef: %v", event.SituationAttributes.References[8])
	}
	if event.SituationAttributes.References[10].ObjectId.Value() != "groupOfLineRef1" {
		t.Errorf("Wrong first GroupOfLinesRef: %v", event.SituationAttributes.References[10])
	}

	if len(event.SituationAttributes.LineSections) != 2 {
		t.Fatalf("Wrong number of LineSections, expected: 2, got: %v", len(event.SituationAttributes.LineSections))
	}
	firstLineSection := *event.SituationAttributes.LineSections[0]
	if ref, _ := firstLineSection.Get("FirstStop"); ref.ObjectId.Value() != "firstStop1" {
		t.Errorf("Wrong first LineSection FirstStop: %v", ref)
	}
	if ref, _ := firstLineSection.Get("LastStop"); ref.ObjectId.Value() != "lastStop1" {
		t.Errorf("Wrong first LineSection LastStop: %v", ref)
	}
	if ref, _ := firstLineSection.Get("LineRef"); ref.ObjectId.Value() != "lineSectionRef1" {
		t.Errorf("Wrong first LineSection LineRef: %v", ref)
	}
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
