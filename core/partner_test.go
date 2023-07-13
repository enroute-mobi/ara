package core

import (
	"reflect"
	"strconv"
	"testing"
	"time"

	"bitbucket.org/enroute-mobi/ara/clock"
	s "bitbucket.org/enroute-mobi/ara/core/settings"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/uuid"
	"github.com/stretchr/testify/assert"
)

func createTestPartnerManager() *PartnerManager {
	referentials := NewMemoryReferentials()
	referential := referentials.New(ReferentialSlug("referential"))
	referential.collectManager = NewTestCollectManager()
	referentials.Save(referential)
	return (referential.partners).(*PartnerManager)
}

func Test_Partner_Id(t *testing.T) {
	partner := &Partner{
		id: "6ba7b814-9dad-11d1-0-00c04fd430c8",
	}

	if expected := PartnerId("6ba7b814-9dad-11d1-0-00c04fd430c8"); partner.Id() != expected {
		t.Errorf("Partner.Id() returns wrong value, got: %s, required: %s", partner.Id(), expected)
	}
}

func Test_Partner_Slug(t *testing.T) {
	partner := &Partner{
		slug: "partner",
	}

	if expected := PartnerSlug("partner"); partner.Slug() != expected {
		t.Errorf("Partner.Slug() returns wrong value, got: %s, required: %s", partner.Id(), expected)
	}
}

func Test_Partner_OperationnalStatus(t *testing.T) {
	partner := NewPartner()

	if expected := OPERATIONNAL_STATUS_UNKNOWN; partner.PartnerStatus.OperationnalStatus != expected {
		t.Errorf("partner.PartnerStatus.OperationnalStatus returns wrong status, got: %s, required: %s", partner.PartnerStatus.OperationnalStatus, expected)
	}
}

func Test_Partner_OperationnalStatus_PushCollector(t *testing.T) {
	partners := createTestPartnerManager()
	partner := partners.New("slug")
	settings := map[string]string{
		"local_credential":     "loc",
		"remote_objectid_kind": "_internal",
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)
	partner.ConnectorTypes = []string{"push-collector"}
	partners.Save(partner)

	// No Connectors
	_, err := partner.CheckStatus()
	if err == nil {
		t.Fatalf("should have an error when partner doesn't have any connectors")
	}

	// Push collector but old collect
	partner.RefreshConnectors()

	ps, err := partner.CheckStatus()
	if err != nil {
		t.Fatalf("should not have an error during checkstatus: %v", err)
	}
	if expected := OPERATIONNAL_STATUS_DOWN; ps.OperationnalStatus != expected {
		t.Errorf("partner.PartnerStatus.OperationnalStatus returns wrong status, got: %s, required: %s", partner.PartnerStatus.OperationnalStatus, expected)
	}

	// Push collector but recent collect
	partner.alternativeStatusCheck.LastCheck = clock.DefaultClock().Now()

	ps, err = partner.CheckStatus()
	if err != nil {
		t.Fatalf("should not have an error during checkstatus: %v", err)
	}
	if expected := OPERATIONNAL_STATUS_UP; ps.OperationnalStatus != expected {
		t.Errorf("partner.PartnerStatus.OperationnalStatus returns wrong status, got: %s, required: %s", partner.PartnerStatus.OperationnalStatus, expected)
	}
}

func Test_Partner_OperationnalStatus_GtfsCollector(t *testing.T) {
	partners := createTestPartnerManager()
	partner := partners.New("slug")

	settings := map[string]string{
		"local_credential":     "loc",
		"remote_objectid_kind": "_internal",
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)
	partner.ConnectorTypes = []string{GTFS_RT_REQUEST_COLLECTOR}
	partners.Save(partner)

	// No Connectors
	_, err := partner.CheckStatus()
	if err == nil {
		t.Fatalf("should have an error when partner doesn't have any connectors")
	}

	// Push collector but old collect
	partner.RefreshConnectors()

	ps, err := partner.CheckStatus()
	if err != nil {
		t.Fatalf("should not have an error during checkstatus: %v", err)
	}
	if expected := OPERATIONNAL_STATUS_DOWN; ps.OperationnalStatus != expected {
		t.Errorf("partner.PartnerStatus.OperationnalStatus returns wrong status, got: %s, required: %s", partner.PartnerStatus.OperationnalStatus, expected)
	}

	// Push collector but recent collect
	partner.alternativeStatusCheck.LastCheck = clock.DefaultClock().Now()
	partner.alternativeStatusCheck.Status = OPERATIONNAL_STATUS_UNKNOWN

	ps, err = partner.CheckStatus()
	if err != nil {
		t.Fatalf("should not have an error during checkstatus: %v", err)
	}
	if expected := OPERATIONNAL_STATUS_UNKNOWN; ps.OperationnalStatus != expected {
		t.Errorf("partner.PartnerStatus.OperationnalStatus returns wrong status, got: %s, required: %s", partner.PartnerStatus.OperationnalStatus, expected)
	}
}

func Test_Partner_SubcriptionCancel(t *testing.T) {
	partners := createTestPartnerManager()
	partner := partners.New("slug")

	settings := map[string]string{
		"remote_url":           "une url",
		"remote_objectid_kind": "_internal",
		s.PARTNER_MAX_RETRY:    "1",
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)
	partner.ConnectorTypes = []string{"siri-stop-monitoring-subscription-collector"}

	partner.subscriptionManager = NewMemorySubscriptions(partner)
	partners.Save(partner)

	partner.subscriptionManager.SetUUIDGenerator(uuid.NewFakeUUIDGenerator())
	referential := partner.Referential()

	stopArea := referential.Model().StopAreas().New()
	stopArea.CollectedAlways = false
	objectid := model.NewObjectID("_internal", "coicogn2")
	stopArea.SetObjectID(objectid)
	stopArea.Save()

	stopVisit := referential.Model().StopVisits().New()
	objectid = model.NewObjectID("_internal", "stopvisit1")
	stopVisit.SetObjectID(objectid)
	stopVisit.StopAreaId = stopArea.Id()
	stopVisit.Collected(time.Now())
	stopVisit.Save()

	objId := model.NewObjectID("_internal", "coicogn2")
	ref := model.Reference{
		ObjectId: &objId,

		Type: "StopArea",
	}

	subscription := partner.Subscriptions().FindOrCreateByKind(StopMonitoringCollect)
	subscription.CreateAndAddNewResource(ref)
	subscription.Save()

	partner.CancelSubscriptions()
	if len(partner.Subscriptions().FindAll()) != 0 {
		t.Errorf("Subscriptions should not be found \n")
	}
}

func Test_Partner_MarshalJSON(t *testing.T) {
	partner := &Partner{
		id: "6ba7b814-9dad-11d1-0-00c04fd430c8",
		PartnerStatus: PartnerStatus{
			OperationnalStatus: OPERATIONNAL_STATUS_UNKNOWN,
		},
		slug:           "partner",
		ConnectorTypes: []string{},
	}
	partner.PartnerSettings = s.NewEmptyPartnerSettings(partner.UUIDGenerator)
	expected := `{"Id":"6ba7b814-9dad-11d1-0-00c04fd430c8","Slug":"partner","PartnerStatus":{"OperationnalStatus":"unknown","RetryCount":0,"ServiceStartedAt":"0001-01-01T00:00:00Z"},"ConnectorTypes":[],"Settings":{}}`
	jsonBytes, err := partner.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}

	jsonString := string(jsonBytes)
	if jsonString != expected {
		t.Errorf("Partner.MarshalJSON() returns wrong json:\n got: %s\n want: %s", jsonString, expected)
	}

	partner.Name = "PartnerName"
	expected = `{"Id":"6ba7b814-9dad-11d1-0-00c04fd430c8","Slug":"partner","Name":"PartnerName","PartnerStatus":{"OperationnalStatus":"unknown","RetryCount":0,"ServiceStartedAt":"0001-01-01T00:00:00Z"},"ConnectorTypes":[],"Settings":{}}`
	jsonBytes, err = partner.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}

	jsonString = string(jsonBytes)
	if jsonString != expected {
		t.Errorf("Partner.MarshalJSON() returns wrong json:\n got: %s\n want: %s", jsonString, expected)
	}
}

func Test_Partner_Save(t *testing.T) {
	partners := createTestPartnerManager()
	partner := partners.New("partner")

	if partner.manager != partners {
		t.Errorf("New partner manager should be partners")
	}

	ok := partner.Save()
	if !ok {
		t.Errorf("partner.Save() should succeed")
	}
	partner = partners.Find(partner.Id())
	if partner == nil {
		t.Errorf("New Partner should be found in Partners manager")
	}
}

func Test_Partner_RefreshConnectors(t *testing.T) {
	partner := &Partner{connectors: make(map[string]Connector)}
	partner.RefreshConnectors()
	if partner.CheckStatusClient() != nil {
		t.Errorf("Partner CheckStatus client should be nil, got: %v", reflect.TypeOf(partner.CheckStatusClient()))
	}

	partner.ConnectorTypes = []string{"siri-check-status-client"}
	partner.RefreshConnectors()
	if _, ok := partner.CheckStatusClient().(*SIRICheckStatusClient); !ok {
		t.Errorf("Partner CheckStatus client should be SIRICheckStatusClient, got: %v", reflect.TypeOf(partner.CheckStatusClient()))
	}

	partner.ConnectorTypes = []string{"test-check-status-client"}
	partner.RefreshConnectors()
	if _, ok := partner.CheckStatusClient().(*TestCheckStatusClient); !ok {
		t.Errorf("Partner CheckStatus client should be TestCheckStatusClient, got: %v", reflect.TypeOf(partner.CheckStatusClient()))
	}
}
func Test_CanCollect(t *testing.T) {
	assert := assert.New(t)
	var TestCases = []struct {
		collectIncludeStopAreas       []string
		collectExcludeStopAreas       []string
		collectUseDiscoveredStopAreas bool
		discoveredStopArea            string
		collectIncludeLines           []string
		collectExcludeLines           []string
		collectUseDiscoveredLines     bool
		discoveredLine                string
		lineIds                       []string
		expectedOutput                bool
		testName                      int
	}{
		{
			collectIncludeStopAreas:       []string{},
			collectExcludeStopAreas:       []string{},
			collectUseDiscoveredStopAreas: false,
			collectIncludeLines:           []string{},
			collectExcludeLines:           []string{},
			collectUseDiscoveredLines:     false,
			lineIds:                       []string{"dummy"},
			expectedOutput:                true,
			testName:                      1,
		},
		{
			collectIncludeStopAreas:       []string{},
			collectExcludeStopAreas:       []string{},
			collectUseDiscoveredStopAreas: false,
			collectIncludeLines:           []string{"dummy"},
			collectExcludeLines:           []string{},
			collectUseDiscoveredLines:     false,
			lineIds:                       []string{"dummy"},
			expectedOutput:                true,
			testName:                      2,
		},
		{
			collectIncludeStopAreas:       []string{},
			collectExcludeStopAreas:       []string{},
			collectUseDiscoveredStopAreas: false,
			collectIncludeLines:           []string{},
			collectExcludeLines:           []string{"dummy"},
			collectUseDiscoveredLines:     false,
			lineIds:                       []string{"dummy"},
			expectedOutput:                false,
			testName:                      3,
		},
		{
			collectIncludeStopAreas:       []string{"dummy"},
			collectExcludeStopAreas:       []string{},
			collectUseDiscoveredStopAreas: false,
			collectIncludeLines:           []string{},
			collectExcludeLines:           []string{},
			collectUseDiscoveredLines:     false,
			lineIds:                       []string{},
			expectedOutput:                true,
			testName:                      4,
		},
		{
			collectIncludeStopAreas:       []string{},
			collectExcludeStopAreas:       []string{"dummy"},
			collectUseDiscoveredStopAreas: false,
			collectIncludeLines:           []string{},
			collectExcludeLines:           []string{},
			collectUseDiscoveredLines:     false,
			lineIds:                       []string{},
			expectedOutput:                false,
			testName:                      5,
		},
	}

	for _, tt := range TestCases {
		partner := NewPartner()

		settings := map[string]string{}
		// StopArea
		if len(tt.collectIncludeStopAreas) != 0 {
			settings[s.COLLECT_INCLUDE_STOP_AREAS] = tt.collectIncludeStopAreas[0]
		}

		if len(tt.collectExcludeStopAreas) != 0 {
			settings[s.COLLECT_EXCLUDE_STOP_AREAS] = tt.collectExcludeStopAreas[0]
		}

		if tt.collectUseDiscoveredStopAreas {
			settings[s.COLLECT_USE_DISCOVERED_SA] = "true"
		}

		if tt.discoveredStopArea != "" {
			partner.discoveredStopAreas[tt.discoveredStopArea] = struct{}{}
		}

		// Line
		if len(tt.collectIncludeLines) != 0 {
			settings[s.COLLECT_INCLUDE_LINES] = tt.collectIncludeLines[0]
		}

		if len(tt.collectExcludeLines) != 0 {
			settings[s.COLLECT_EXCLUDE_LINES] = tt.collectExcludeLines[0]
		}

		if tt.collectUseDiscoveredLines {
			settings[s.COLLECT_USE_DISCOVERED_LINES] = "true"
		}

		if tt.discoveredLine != "" {
			partner.discoveredLines[tt.discoveredLine] = struct{}{}
		}

		lineIds := make(map[string]struct{})
		if len(tt.lineIds) != 0 {
			lineIds[tt.lineIds[0]] = struct{}{}
		}

		// test
		partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)
		output := partner.CanCollect("dummy", lineIds)

		assert.Equal(tt.expectedOutput, output, tt.testName)
	}
}

func Test_CanCollectStopArea(t *testing.T) {
	assert := assert.New(t)

	var TestCases = []struct {
		collectIncludeStopAreas       []string
		collectExcludeStopAreas       []string
		collectUseDiscoveredStopAreas bool
		discoveredStopArea            string
		expectedOutput                s.CollectStatus
		testName                      int
	}{
		{
			collectIncludeStopAreas:       []string{},
			collectExcludeStopAreas:       []string{},
			collectUseDiscoveredStopAreas: false,
			expectedOutput:                s.COLLECT_UNKNOWN,
			testName:                      1,
		},
		{
			collectIncludeStopAreas:       []string{"dummy"},
			collectExcludeStopAreas:       []string{},
			collectUseDiscoveredStopAreas: false,
			expectedOutput:                s.CAN_COLLECT,
			testName:                      2,
		},
		{
			collectIncludeStopAreas:       []string{},
			collectExcludeStopAreas:       []string{"dummy"},
			collectUseDiscoveredStopAreas: false,
			expectedOutput:                s.CANNOT_COLLECT,
			testName:                      3,
		},
		{
			collectIncludeStopAreas:       []string{"dummy"},
			collectExcludeStopAreas:       []string{"dummy"},
			collectUseDiscoveredStopAreas: false,
			expectedOutput:                s.CAN_COLLECT,
			testName:                      4,
		},
		{
			collectIncludeStopAreas:       []string{},
			collectExcludeStopAreas:       []string{},
			collectUseDiscoveredStopAreas: true,
			discoveredStopArea:            "dummy",
			expectedOutput:                s.CAN_COLLECT,
			testName:                      5,
		},
		{
			collectIncludeStopAreas:       []string{},
			collectExcludeStopAreas:       []string{},
			collectUseDiscoveredStopAreas: true,
			discoveredStopArea:            "ANOTHER_DUMMY",
			expectedOutput:                s.CANNOT_COLLECT,
			testName:                      6,
		},
		{
			collectIncludeStopAreas:       []string{},
			collectExcludeStopAreas:       []string{},
			collectUseDiscoveredStopAreas: false,
			discoveredStopArea:            "dummy",
			expectedOutput:                s.COLLECT_UNKNOWN,
			testName:                      7,
		},
		{
			collectIncludeStopAreas:       []string{},
			collectExcludeStopAreas:       []string{},
			collectUseDiscoveredStopAreas: false,
			discoveredStopArea:            "ANOTHER_DUMMY",
			expectedOutput:                s.COLLECT_UNKNOWN,
			testName:                      8,
		},
		{
			collectIncludeStopAreas:       []string{"dummy"},
			collectExcludeStopAreas:       []string{},
			collectUseDiscoveredStopAreas: true,
			discoveredStopArea:            "ANOTHER_DUMMY",
			expectedOutput:                s.CAN_COLLECT,
			testName:                      9,
		},
		{
			collectIncludeStopAreas:       []string{},
			collectExcludeStopAreas:       []string{"dummy"},
			collectUseDiscoveredStopAreas: true,
			discoveredStopArea:            "dummy",
			expectedOutput:                s.CANNOT_COLLECT,
			testName:                      10,
		},
		{
			collectIncludeStopAreas:       []string{"other"},
			collectExcludeStopAreas:       []string{},
			collectUseDiscoveredStopAreas: false,
			expectedOutput:                s.CANNOT_COLLECT,
			testName:                      11,
		},
	}

	for _, tt := range TestCases {
		partner := NewPartner()

		settings := map[string]string{}

		if len(tt.collectIncludeStopAreas) != 0 {
			settings[s.COLLECT_INCLUDE_STOP_AREAS] = tt.collectIncludeStopAreas[0]
		}

		if len(tt.collectExcludeStopAreas) != 0 {
			settings[s.COLLECT_EXCLUDE_STOP_AREAS] = tt.collectExcludeStopAreas[0]
		}

		if tt.collectUseDiscoveredStopAreas {
			settings[s.COLLECT_USE_DISCOVERED_SA] = "true"
		}

		if tt.discoveredStopArea != "" {
			partner.discoveredStopAreas[tt.discoveredStopArea] = struct{}{}
		}

		partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)

		output := partner.CanCollectStop("dummy")
		assert.Equal(output, tt.expectedOutput)
	}
}

func Test_CanCollectLine(t *testing.T) {
	assert := assert.New(t)

	var TestCases = []struct {
		collectIncludeLines       []string
		collectExcludeLines       []string
		collectUseDiscoveredLines bool
		discoveredLine            string
		expectedOutput            bool
		testName                  int
	}{
		{
			collectIncludeLines:       []string{},
			collectExcludeLines:       []string{},
			collectUseDiscoveredLines: false,
			expectedOutput:            true,
			testName:                  1,
		},
		{
			collectIncludeLines:       []string{"dummy"},
			collectExcludeLines:       []string{},
			collectUseDiscoveredLines: false,
			expectedOutput:            true,
			testName:                  2,
		},
		{
			collectIncludeLines:       []string{},
			collectExcludeLines:       []string{"dummy"},
			collectUseDiscoveredLines: false,
			expectedOutput:            false,
			testName:                  3,
		},
		{
			collectIncludeLines:       []string{"dummy"},
			collectExcludeLines:       []string{"dummy"},
			collectUseDiscoveredLines: false,
			expectedOutput:            true,
			testName:                  4,
		},
		{
			collectIncludeLines:       []string{},
			collectExcludeLines:       []string{},
			collectUseDiscoveredLines: true,
			discoveredLine:            "dummy",
			expectedOutput:            true,
			testName:                  5,
		},
		{
			collectIncludeLines:       []string{},
			collectExcludeLines:       []string{},
			collectUseDiscoveredLines: true,
			discoveredLine:            "ANOTHER_DUMMY",
			expectedOutput:            false,
			testName:                  6,
		},
		{
			collectIncludeLines:       []string{},
			collectExcludeLines:       []string{},
			collectUseDiscoveredLines: false,
			discoveredLine:            "dummy",
			expectedOutput:            true,
			testName:                  7,
		},
		{
			collectIncludeLines:       []string{},
			collectExcludeLines:       []string{},
			collectUseDiscoveredLines: false,
			discoveredLine:            "ANOTHER_DUMMY",
			expectedOutput:            true,
			testName:                  8,
		},
		{
			collectIncludeLines:       []string{"dummy"},
			collectExcludeLines:       []string{},
			collectUseDiscoveredLines: true,
			discoveredLine:            "ANOTHER_DUMMY",
			expectedOutput:            true,
			testName:                  9,
		},
		{
			collectIncludeLines:       []string{},
			collectExcludeLines:       []string{"dummy"},
			collectUseDiscoveredLines: true,
			discoveredLine:            "dummy",
			expectedOutput:            false,
			testName:                  10,
		},
		{
			collectIncludeLines:       []string{"other"},
			collectExcludeLines:       []string{},
			collectUseDiscoveredLines: false,
			expectedOutput:            false,
			testName:                  11,
		},
	}

	for _, tt := range TestCases {
		partner := NewPartner()
		settings := map[string]string{}

		if len(tt.collectIncludeLines) != 0 {
			settings[s.COLLECT_INCLUDE_LINES] = tt.collectIncludeLines[0]
		}

		if len(tt.collectExcludeLines) != 0 {
			settings[s.COLLECT_EXCLUDE_LINES] = tt.collectExcludeLines[0]
		}

		if tt.collectUseDiscoveredLines {
			settings[s.COLLECT_USE_DISCOVERED_LINES] = "true"
		}

		if tt.discoveredLine != "" {
			partner.discoveredLines[tt.discoveredLine] = struct{}{}
		}

		partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)

		output := partner.CanCollectLine("dummy")
		assert.Equal(output, tt.expectedOutput, strconv.Itoa(tt.testName))
	}
}

func Test_Partners_FindAllByCollectPriority(t *testing.T) {
	partners := createTestPartnerManager()
	partner1 := &Partner{
		slug: "First",
	}

	settings1 := map[string]string{s.COLLECT_PRIORITY: "2"}
	partner1.PartnerSettings = s.NewPartnerSettings(partner1.UUIDGenerator, settings1)

	partner2 := &Partner{
		slug: "Second",
	}
	settings2 := map[string]string{s.COLLECT_PRIORITY: "1"}
	partner2.PartnerSettings = s.NewPartnerSettings(partner2.UUIDGenerator, settings2)

	partners.Save(partner1)
	partners.Save(partner2)

	orderedPartners := partners.FindAllByCollectPriority()
	if orderedPartners[0].Slug() != "First" {
		t.Errorf("Partners should be ordered")
	}
}

func Test_Partner_Subcription(t *testing.T) {
	partner := NewPartner()

	sub := partner.Subscriptions()
	sub.New("kind")

	if len(partner.Subscriptions().FindAll()) != 1 {
		t.Errorf("Wrong number of subcriptions want : %v got: %v", 1, len(sub.FindAll()))
	}
}

func Test_NewPartnerManager(t *testing.T) {
	partners := createTestPartnerManager()

	if partners.guardian == nil {
		t.Errorf("New PartnerManager should have a PartnersGuardian")
	}
}

func Test_PartnerManager_New(t *testing.T) {
	partners := createTestPartnerManager()
	partner := partners.New("partner")

	if partner.Id() != "" {
		t.Errorf("New Partner identifier should be an empty string, got: %s", partner.Id())
	}
}

func Test_PartnerManager_Save(t *testing.T) {
	partners := createTestPartnerManager()
	partner := partners.New("partner")

	if success := partners.Save(partner); !success {
		t.Errorf("Save should return true")
	}

	if partner.Id() == "" {
		t.Errorf("New Partner identifier should not be an empty string")
	}
}

func Test_PartnerManager_Find_NotFound(t *testing.T) {
	partners := createTestPartnerManager()
	partner := partners.Find("6ba7b814-9dad-11d1-0-00c04fd430c8")
	if partner != nil {
		t.Errorf("Find should return false when Partner isn't found")
	}
}

func Test_PartnerManager_Find(t *testing.T) {
	partners := createTestPartnerManager()

	existingPartner := partners.New("partner")
	partners.Save(existingPartner)
	partnerId := existingPartner.Id()

	partner := partners.Find(partnerId)
	if partner == nil {
		t.Fatal("Find should return true when Partner is found")
	}
	if partner.Id() != partnerId {
		t.Errorf("Find should return a Partner with the given Id")
	}
}

func Test_PartnerManager_FindByCredentials(t *testing.T) {
	partners := createTestPartnerManager()

	existingPartner := partners.New("partner")
	settings := map[string]string{
		s.LOCAL_CREDENTIAL:  "cred",
		s.LOCAL_CREDENTIALS: "cred2,cred3",
	}
	existingPartner.PartnerSettings = s.NewPartnerSettings(existingPartner.UUIDGenerator, settings)
	partners.Save(existingPartner)

	partner, ok := partners.FindByCredential("cred")
	if !ok {
		t.Fatal("FindByCredential should return true when Partner is found")
	}
	if partner.Id() != existingPartner.Id() {
		t.Errorf("FindByCredential should return a Partner with the given local_credential")
	}

	partner, ok = partners.FindByCredential("cred2")
	if !ok {
		t.Fatal("FindByCredential should return true when Partner is found")
	}
	if partner.Id() != existingPartner.Id() {
		t.Errorf("FindByCredential should return a Partner with the given local_credential")
	}

	partner, ok = partners.FindByCredential("cred3")
	if !ok {
		t.Fatal("FindByCredential should return true when Partner is found")
	}
	if partner.Id() != existingPartner.Id() {
		t.Errorf("FindByCredential should return a Partner with the given local_credential")
	}
}

func Test_PartnerManager_FindBySlug(t *testing.T) {
	partners := createTestPartnerManager()

	existingPartner := partners.New("partner")
	partners.Save(existingPartner)

	partner, ok := partners.FindBySlug("partner")
	if !ok {
		t.Fatal("FindBySlug should return true when Partner is found")
	}
	if partner.Id() != existingPartner.Id() {
		t.Errorf("FindBySlug should return a Partner with the given slug")
	}
}

func Test_PartnerManager_FindAll(t *testing.T) {
	partners := createTestPartnerManager()

	for i := 0; i < 5; i++ {
		existingPartner := partners.New(PartnerSlug(strconv.Itoa(i)))
		partners.Save(existingPartner)
	}

	foundPartners := partners.FindAll()

	if len(foundPartners) != 5 {
		t.Errorf("FindAll should return all partners")
	}
}

func Test_PartnerManager_Delete(t *testing.T) {
	partners := createTestPartnerManager()

	existingPartner := partners.New("partner")
	partners.Save(existingPartner)

	partnerId := existingPartner.Id()

	partners.Delete(existingPartner)

	partner := partners.Find(partnerId)
	if partner != nil {
		t.Errorf("Deleted Partner should not be findable")
	}
}

func Test_MemoryPartners_Load(t *testing.T) {
	model.InitTestDb(t)
	defer model.CleanTestDb(t)

	referentials := NewMemoryReferentials()
	referential := referentials.New("referential")
	referentials.Save(referential)

	// Insert Data in the test db
	dbPartner := model.DatabasePartner{
		Id:             "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
		ReferentialId:  string(referential.Id()),
		Slug:           "ratp",
		Name:           "RATP",
		Settings:       "{}",
		ConnectorTypes: "[]",
	}
	err := model.Database.Insert(&dbPartner)
	if err != nil {
		t.Fatal(err)
	}

	// Fetch data from the db
	partners := NewPartnerManager(referential)
	err = partners.Load()
	if err != nil {
		t.Fatal(err)
	}

	partnerId := PartnerId(dbPartner.Id)
	partner := partners.Find(partnerId)
	if partner == nil {
		t.Errorf("Loaded Partners should be found")
	} else if partner.Id() != partnerId {
		t.Errorf("Wrong Id:\n got: %v\n expected: %v", partner.Id(), partnerId)
	}
}

// Tested in Referential
// func Test_MemoryPartners_SaveToDatabase(t *testing.T) {}

func Test_Partners_StartStop(t *testing.T) {
	referentials := NewMemoryReferentials()
	referential := referentials.New(ReferentialSlug("referential"))
	referentials.Save(referential)
	partner := referential.Partners().New("partner")

	partner.ConnectorTypes = []string{TEST_STARTABLE_CONNECTOR}
	partner.RefreshConnectors()
	partner.Save()

	connector, ok := partner.Connector(TEST_STARTABLE_CONNECTOR)
	if !ok {
		t.Fatalf("Connector should have a TestStartableConnector")
	}
	if connector.(*TestStartableConnector).started {
		t.Errorf("Connector should be stoped")
	}

	referential.Start()
	if !connector.(*TestStartableConnector).started {
		t.Errorf("Connector should be started")
	}

	referential.Stop()
	if connector.(*TestStartableConnector).started {
		t.Errorf("Connector should be stoped")
	}
}

func Test_Partner_DefaultIdentifierGenerator(t *testing.T) {
	partner := &Partner{
		slug: "partner",
	}
	partner.PartnerSettings = s.NewEmptyPartnerSettings(partner.UUIDGenerator)

	g := partner.MessageIdentifierGenerator()
	if expected := "%{uuid}"; g.FormatString() != expected {
		t.Errorf("partner message_identifier IdentifierGenerator should be %v, got: %v ", expected, g.FormatString())
	}
	g = partner.ResponseMessageIdentifierGenerator()
	if expected := "%{uuid}"; g.FormatString() != expected {
		t.Errorf("partner response_message_identifier IdentifierGenerator should be %v, got: %v ", expected, g.FormatString())
	}
	g = partner.DataFrameIdentifierGenerator()
	if expected := "%{id}"; g.FormatString() != expected {
		t.Errorf("partner data_frame_identifier IdentifierGenerator should be %v, got: %v ", expected, g.FormatString())
	}
	g = partner.ReferenceIdentifierGenerator()
	if expected := "%{type}:%{id}"; g.FormatString() != expected {
		t.Errorf("partner reference_identifier IdentifierGenerator should be %v, got: %v ", expected, g.FormatString())
	}
	g = partner.ReferenceStopAreaIdentifierGenerator()
	if expected := "%{id}"; g.FormatString() != expected {
		t.Errorf("partner reference_stop_area_identifier IdentifierGenerator should be %v, got: %v ", expected, g.FormatString())
	}
}

func Test_Partner_IdentifierGenerator(t *testing.T) {
	partner := &Partner{
		slug: "partner",
	}

	settings := map[string]string{
		"generators.message_identifier":             "mid",
		"generators.response_message_identifier":    "rmid",
		"generators.data_frame_identifier":          "dfid",
		"generators.reference_identifier":           "rid",
		"generators.reference_stop_area_identifier": "rsaid",
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)

	g := partner.MessageIdentifierGenerator()
	if expected := "mid"; g.FormatString() != expected {
		t.Errorf("partner message_identifier IdentifierGenerator should be %v, got: %v ", expected, g.FormatString())
	}
	g = partner.ResponseMessageIdentifierGenerator()
	if expected := "rmid"; g.FormatString() != expected {
		t.Errorf("partner response_message_identifier IdentifierGenerator should be %v, got: %v ", expected, g.FormatString())
	}
	g = partner.DataFrameIdentifierGenerator()
	if expected := "dfid"; g.FormatString() != expected {
		t.Errorf("partner data_frame_identifier IdentifierGenerator should be %v, got: %v ", expected, g.FormatString())
	}
	g = partner.ReferenceIdentifierGenerator()
	if expected := "rid"; g.FormatString() != expected {
		t.Errorf("partner reference_identifier IdentifierGenerator should be %v, got: %v ", expected, g.FormatString())
	}
	g = partner.ReferenceStopAreaIdentifierGenerator()
	if expected := "rsaid"; g.FormatString() != expected {
		t.Errorf("partner reference_stop_area_identifier IdentifierGenerator should be %v, got: %v ", expected, g.FormatString())
	}
}

func Test_Partner_SIRIClient(t *testing.T) {
	partner := &Partner{
		slug: "partner",
	}
	partner.PartnerSettings = s.NewEmptyPartnerSettings(partner.UUIDGenerator)
	partner.SIRIClient()
	if partner.httpClient == nil {
		t.Error("partner.SIRIClient() should set Partner httpClient")
	}

	settings := map[string]string{
		s.REMOTE_URL: "remote_url",
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)

	partner.SIRIClient()
	if partner.httpClient.Url != "remote_url" {
		t.Error("Partner should have created a new SoapClient when partner setting changes")
	}

	settings = map[string]string{
		s.SUBSCRIPTIONS_REMOTE_URL: "sub_remote_url",
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)

	partner.SIRIClient()
	if partner.httpClient.SubscriptionsUrl != "sub_remote_url" {
		t.Error("Partner should have created a new SoapClient when partner setting changes")
	}
}

func Test_Partner_RequestorRef(t *testing.T) {
	partner := &Partner{
		slug: "partner",
	}

	settings := map[string]string{
		s.REMOTE_CREDENTIAL: "ara",
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)
	if partner.RequestorRef() != "ara" {
		t.Errorf("Wrong Partner RequestorRef:\n got: %s\n want: \"ara\"", partner.RequestorRef())
	}
}
