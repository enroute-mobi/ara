package core

import (
	"encoding/json"
	"reflect"
	"strconv"
	"testing"
	"time"

	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/uuid"
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
	partner.SetSettingsDefinition(map[string]string{
		"local_credential":     "loc",
		"remote_objectid_kind": "_internal",
	})
	partner.ConnectorTypes = []string{"push-collector"}
	partners.Save(partner)

	// No Connectors
	ps, err := partner.CheckStatus()
	if err == nil {
		t.Fatalf("should have an error when partner doesn't have any connectors")
	}

	// Push collector but old collect
	partner.RefreshConnectors()

	ps, err = partner.CheckStatus()
	if err != nil {
		t.Fatalf("should not have an error during checkstatus: %v", err)
	}
	if expected := OPERATIONNAL_STATUS_DOWN; ps.OperationnalStatus != expected {
		t.Errorf("partner.PartnerStatus.OperationnalStatus returns wrong status, got: %s, required: %s", partner.PartnerStatus.OperationnalStatus, expected)
	}

	// Push collector but recent collect
	partner.lastPush = clock.DefaultClock().Now()

	ps, err = partner.CheckStatus()
	if err != nil {
		t.Fatalf("should not have an error during checkstatus: %v", err)
	}
	if expected := OPERATIONNAL_STATUS_UP; ps.OperationnalStatus != expected {
		t.Errorf("partner.PartnerStatus.OperationnalStatus returns wrong status, got: %s, required: %s", partner.PartnerStatus.OperationnalStatus, expected)
	}
}

func Test_Partner_SubcriptionCancel(t *testing.T) {
	partners := createTestPartnerManager()
	partner := partners.New("slug")
	partner.SetSettingsDefinition(map[string]string{
		"remote_url":           "une url",
		"remote_objectid_kind": "_internal",
	})
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

	subscription := partner.Subscriptions().FindOrCreateByKind("StopMonitoringCollect")
	subscription.CreateAddNewResource(ref)
	subscription.Save()

	partner.CancelSubscriptions()
	if len(partner.Subscriptions().FindAll()) != 0 {
		t.Errorf("Subscriptions should not be found \n")
	}

	// stopVisit = referential.Model().StopVisits().FindByStopAreaId(stopArea.Id())[0]
	// if stopVisit.IsCollected() != false {
	// 	t.Errorf("stopVisit should be false but got %v\n", stopVisit.IsCollected())
	// }
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
	partner.PartnerSettings = NewPartnerSettings(partner)
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

func Test_Partner_CanCollectTrue(t *testing.T) {
	partner := &Partner{}
	partner.PartnerSettings = NewPartnerSettings(partner)
	stopAreaObjectId := "NINOXE:StopPoint:SP:24:LOC"
	partner.SetSetting(COLLECT_INCLUDE_STOP_AREAS, "NINOXE:StopPoint:SP:24:LOC")
	partner.setCollectSettings()
	if !partner.CanCollect(stopAreaObjectId, map[string]struct{}{}) {
		t.Errorf("Partner can collect should return true")
	}

	partner.SetSetting(COLLECT_USE_DISCOVERED_SA, "true")
	partner.setCollectSettings()
	if !partner.CanCollect(stopAreaObjectId, map[string]struct{}{}) {
		t.Errorf("Partner can collect should return true")
	}
}

func Test_Partner_CanCollectTrueLine(t *testing.T) {
	partner := &Partner{}
	partner.PartnerSettings = NewPartnerSettings(partner)
	stopAreaObjectId := "NINOXE:StopPoint:SP:24:LOC"
	lines := map[string]struct{}{"NINOXE:Line:SP:24:": struct{}{}}
	partner.SetSetting(COLLECT_INCLUDE_LINES, "NINOXE:Line:SP:24:")
	partner.setCollectSettings()
	if !partner.CanCollect(stopAreaObjectId, lines) {
		t.Errorf("Partner can collect should return true")
	}

	partner.SetSetting(COLLECT_USE_DISCOVERED_SA, "true")
	partner.setCollectSettings()
	if !partner.CanCollect(stopAreaObjectId, lines) {
		t.Errorf("Partner can collect should return true")
	}
}

func Test_Partner_CanCollectTrue_EmptySettings(t *testing.T) {
	partner := &Partner{}
	partner.PartnerSettings = NewPartnerSettings(partner)
	stopAreaObjectId := "NINOXE:StopPoint:SP:24:LOC"
	partner.setCollectSettings()
	if !partner.CanCollect(stopAreaObjectId, map[string]struct{}{}) {
		t.Errorf("Partner can collect should return true")
	}

	partner.SetSetting(COLLECT_USE_DISCOVERED_SA, "true")
	partner.setCollectSettings()
	if partner.CanCollect(stopAreaObjectId, map[string]struct{}{}) {
		t.Errorf("Partner can collect should return false")
	}
}

func Test_Partner_CanCollectFalse(t *testing.T) {
	partner := &Partner{}
	partner.PartnerSettings = NewPartnerSettings(partner)
	stopAreaObjectId := "BAD_VALUE"
	partner.SetSetting(COLLECT_INCLUDE_STOP_AREAS, "NINOXE:StopPoint:SP:24:LOC")
	partner.setCollectSettings()
	if partner.CanCollect(stopAreaObjectId, map[string]struct{}{}) {
		t.Errorf("Partner can collect should return flase")
	}
}

func Test_Partner_CanCollectFalseLine(t *testing.T) {
	partner := &Partner{}
	partner.PartnerSettings = NewPartnerSettings(partner)
	stopAreaObjectId := "BAD_VALUE"
	partner.SetSetting(COLLECT_INCLUDE_LINES, "NINOXE:Line:SP:24:")
	partner.setCollectSettings()
	if partner.CanCollect(stopAreaObjectId, map[string]struct{}{}) {
		t.Errorf("Partner can collect should return flase")
	}
}

func Test_Partner_CanCollectTrueExcluded(t *testing.T) {
	partner := &Partner{}
	partner.PartnerSettings = NewPartnerSettings(partner)
	stopAreaObjectId := "NINOXE:StopPoint:SP:24:LOC"
	partner.SetSetting(COLLECT_INCLUDE_STOP_AREAS, "NINOXE:StopPoint:SP:24:LOC")
	partner.SetSetting(COLLECT_EXCLUDE_STOP_AREAS, "NINOXE:StopPoint:SP:25:LOC")
	partner.setCollectSettings()
	if !partner.CanCollect(stopAreaObjectId, map[string]struct{}{}) {
		t.Errorf("Partner can collect should return true")
	}

	partner.SetSetting(COLLECT_USE_DISCOVERED_SA, "true")
	partner.setCollectSettings()
	if !partner.CanCollect(stopAreaObjectId, map[string]struct{}{}) {
		t.Errorf("Partner can collect should return true")
	}
}

func Test_Partner_CanCollectFalseExcluded(t *testing.T) {
	partner := &Partner{}
	partner.PartnerSettings = NewPartnerSettings(partner)
	stopAreaObjectId := "NINOXE:StopPoint:SP:24:LOC"
	partner.SetSetting(COLLECT_INCLUDE_STOP_AREAS, "NINOXE:StopPoint:SP:24:LOC")
	partner.SetSetting(COLLECT_EXCLUDE_STOP_AREAS, "NINOXE:StopPoint:SP:24:LOC")
	partner.setCollectSettings()
	if partner.CanCollect(stopAreaObjectId, map[string]struct{}{}) {
		t.Errorf("Partner can collect should return false")
	}

	partner.SetSetting(COLLECT_USE_DISCOVERED_SA, "true")
	partner.setCollectSettings()
	if partner.CanCollect(stopAreaObjectId, map[string]struct{}{}) {
		t.Errorf("Partner can collect should return false")
	}
}

func Test_Partner_CanCollectFalseSPD(t *testing.T) {
	partner := &Partner{}
	partner.PartnerSettings = NewPartnerSettings(partner)
	partner.discoveredStopAreas = make(map[string]struct{})
	stopAreaObjectId := "NINOXE:StopPoint:SP:24:LOC"
	partner.SetSetting(COLLECT_USE_DISCOVERED_SA, "true")
	partner.setCollectSettings()
	if partner.CanCollect(stopAreaObjectId, map[string]struct{}{}) {
		t.Errorf("Partner can collect should return false")
	}
}

func Test_Partner_CanCollectTrueSPD(t *testing.T) {
	partner := &Partner{}
	partner.PartnerSettings = NewPartnerSettings(partner)
	partner.discoveredStopAreas = make(map[string]struct{})
	stopAreaObjectId := "NINOXE:StopPoint:SP:24:LOC"
	partner.SetSetting(COLLECT_USE_DISCOVERED_SA, "true")
	partner.discoveredStopAreas["NINOXE:StopPoint:SP:24:LOC"] = struct{}{}
	partner.setCollectSettings()
	if !partner.CanCollect(stopAreaObjectId, map[string]struct{}{}) {
		t.Errorf("Partner can collect should return true")
	}
}

func Test_Partner_CanCollectTrueSPDButExcluded(t *testing.T) {
	partner := &Partner{}
	partner.PartnerSettings = NewPartnerSettings(partner)
	partner.discoveredStopAreas = make(map[string]struct{})
	stopAreaObjectId := "NINOXE:StopPoint:SP:24:LOC"
	partner.SetSetting(COLLECT_USE_DISCOVERED_SA, "true")
	partner.discoveredStopAreas["NINOXE:StopPoint:SP:24:LOC"] = struct{}{}
	partner.SetSetting(COLLECT_EXCLUDE_STOP_AREAS, "NINOXE:StopPoint:SP:24:LOC")
	partner.setCollectSettings()
	if partner.CanCollect(stopAreaObjectId, map[string]struct{}{}) {
		t.Errorf("Partner can collect should return false")
	}
}

func Test_Partners_FindAllByCollectPriority(t *testing.T) {
	partners := createTestPartnerManager()
	partner1 := &Partner{
		slug: "First",
	}
	partner1.PartnerSettings = NewPartnerSettings(partner1)
	partner2 := &Partner{
		slug: "Second",
	}
	partner2.PartnerSettings = NewPartnerSettings(partner2)

	partner1.SetSetting("collect.priority", "2")

	partner2.SetSetting("collect.priority", "1")

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

func Test_APIPartner_SetFactories(t *testing.T) {
	partner := &APIPartner{
		ConnectorTypes: []string{"unexistant-factory", "test-check-status-client"},
		factories:      make(map[string]ConnectorFactory),
	}
	partner.setFactories()

	if len(partner.factories) != 1 {
		t.Errorf("Factories should have been successfully created by setFactories")
	}
}

func Test_APIPartner_Validate(t *testing.T) {
	partners := createTestPartnerManager()
	// Check empty Slug
	apiPartner := &APIPartner{
		manager: partners,
	}
	valid := apiPartner.Validate()

	if valid {
		t.Errorf("Validate should return false")
	}
	if len(apiPartner.Errors) != 1 {
		t.Errorf("apiPartner Errors should not be empty")
	}
	if len(apiPartner.Errors.Get("Slug")) != 1 || apiPartner.Errors.Get("Slug")[0] != ERROR_BLANK {
		t.Errorf("apiPartner should have Error for Slug, got %v", apiPartner.Errors)
	}

	// Check wrong format slug
	apiPartner.Slug = "Wrong_format"
	valid = apiPartner.Validate()

	if valid {
		t.Errorf("Validate should return false")
	}
	if len(apiPartner.Errors) != 1 {
		t.Errorf("apiPartner Errors should not be empty")
	}
	if len(apiPartner.Errors.Get("Slug")) != 1 || apiPartner.Errors.Get("Slug")[0] != ERROR_SLUG_FORMAT {
		t.Errorf("apiPartner should have Error for Slug, got %v", apiPartner.Errors)
	}

	jsonBytes, _ := json.Marshal(apiPartner)
	expected := `{"Slug":"Wrong_format","Errors":{"Slug":["Invalid format: only lowercase alphanumeric characters and _"]}}`
	if string(jsonBytes) != expected {
		t.Fatalf("Invalid JSON, expected %v, got %v", expected, string(jsonBytes))
	}

	// Check Already Used Slug and local_credential
	partner := partners.New("slug")
	partner.SetSetting("local_credential", "cred")
	partners.Save(partner)

	apiPartner = &APIPartner{
		Slug:     "slug",
		Settings: map[string]string{"local_credential": "cred"},
		manager:  partners,
	}
	valid = apiPartner.Validate()

	if valid {
		t.Errorf("Validate should return false")
	}
	if len(apiPartner.Errors) != 2 {
		t.Errorf("apiPartner Errors should not be empty")
	}
	if len(apiPartner.Errors.Get("Slug")) != 1 || apiPartner.Errors.Get("Slug")[0] != ERROR_UNIQUE {
		t.Errorf("apiPartner should have Error for Slug, got %v", apiPartner.Errors)
	}
	if len(apiPartner.Errors.GetSettingError(LOCAL_CREDENTIAL)) != 1 || apiPartner.Errors.GetSettingError(LOCAL_CREDENTIAL)[0] != ERROR_UNIQUE {
		t.Errorf("apiPartner should have Error for local_credential, got %v", apiPartner.Errors)
	}

	jsonBytes, _ = json.Marshal(apiPartner)
	expected = `{"Slug":"slug","Settings":{"local_credential":"cred"},"Errors":{"Settings":{"local_credential":["Is already in use"]},"Slug":["Is already in use"]}}`
	if string(jsonBytes) != expected {
		t.Fatalf("Invalid JSON, expected %v, got %v", expected, string(jsonBytes))
	}

	// Check ok
	apiPartner = &APIPartner{
		Slug:     "slug_2",
		Settings: map[string]string{"local_credential": "cred2"},
		manager:  partners,
	}
	valid = apiPartner.Validate()

	if !valid {
		t.Errorf("Validate should return true")
	}
	if len(apiPartner.Errors) != 0 {
		t.Errorf("apiPartner Errors should be empty, got %v", apiPartner.Errors)
	}

	// Check settings errors
	apiPartner = &APIPartner{
		Slug:           "",
		Settings:       map[string]string{},
		ConnectorTypes: []string{SIRI_STOP_POINTS_DISCOVERY_REQUEST_BROADCASTER},
		manager:        partners,
		factories:      make(map[string]ConnectorFactory),
	}
	valid = apiPartner.Validate()

	if valid {
		t.Errorf("Validate should return false")
	}
	if len(apiPartner.Errors) != 2 {
		t.Errorf("apiPartner Errors should not be empty, got %v", apiPartner.Errors)
	}
	if len(apiPartner.Errors.getSettings()) != 2 {
		t.Errorf("apiPartner Setting Errors should have 2 errors, got %v", apiPartner.Errors.getSettings())
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

func Test_PartnerManager_FindByCredential(t *testing.T) {
	partners := createTestPartnerManager()

	existingPartner := partners.New("partner")
	existingPartner.SetSetting("local_credential", "cred")
	partners.Save(existingPartner)

	partner, ok := partners.FindBySetting(LOCAL_CREDENTIAL, "cred")
	if !ok {
		t.Fatal("FindBySetting should return true when Partner is found")
	}
	if partner.Id() != existingPartner.Id() {
		t.Errorf("FindBySetting should return a Partner with the given local_credential")
	}
}

func Test_PartnerManager_FindByCredentials(t *testing.T) {
	partners := createTestPartnerManager()

	existingPartner := partners.New("partner")
	existingPartner.SetSetting(LOCAL_CREDENTIAL, "cred")
	existingPartner.SetSetting(LOCAL_CREDENTIALS, "cred2,cred3")
	partners.Save(existingPartner)

	partner, ok := partners.FindByCredential("cred")
	if !ok {
		t.Fatal("FindBySetting should return true when Partner is found")
	}
	if partner.Id() != existingPartner.Id() {
		t.Errorf("FindBySetting should return a Partner with the given local_credential")
	}

	partner, ok = partners.FindByCredential("cred2")
	if !ok {
		t.Fatal("FindBySetting should return true when Partner is found")
	}
	if partner.Id() != existingPartner.Id() {
		t.Errorf("FindBySetting should return a Partner with the given local_credential")
	}

	partner, ok = partners.FindByCredential("cred3")
	if !ok {
		t.Fatal("FindBySetting should return true when Partner is found")
	}
	if partner.Id() != existingPartner.Id() {
		t.Errorf("FindBySetting should return a Partner with the given local_credential")
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

func Test_Partner_IdentifierGenerator(t *testing.T) {
	partner := &Partner{
		slug: "partner",
	}
	partner.PartnerSettings = NewPartnerSettings(partner)

	midGenerator := partner.IdentifierGenerator("message_identifier")
	if expected := "%{uuid}"; midGenerator.formatString != expected {
		t.Errorf("partner message_identifier IdentifierGenerator should be %v, got: %v ", expected, midGenerator.formatString)
	}
	midGenerator = partner.IdentifierGenerator("response_message_identifier")
	if expected := "%{uuid}"; midGenerator.formatString != expected {
		t.Errorf("partner response_message_identifier IdentifierGenerator should be %v, got: %v ", expected, midGenerator.formatString)
	}
	midGenerator = partner.IdentifierGenerator("data_frame_identifier")
	if expected := "%{id}"; midGenerator.formatString != expected {
		t.Errorf("partner data_frame_identifier IdentifierGenerator should be %v, got: %v ", expected, midGenerator.formatString)
	}
	midGenerator = partner.IdentifierGenerator("reference_identifier")
	if expected := "%{type}:%{id}"; midGenerator.formatString != expected {
		t.Errorf("partner reference_identifier IdentifierGenerator should be %v, got: %v ", expected, midGenerator.formatString)
	}
	midGenerator = partner.IdentifierGenerator("reference_stop_area_identifier")
	if expected := "%{id}"; midGenerator.formatString != expected {
		t.Errorf("partner reference_stop_area_identifier IdentifierGenerator should be %v, got: %v ", expected, midGenerator.formatString)
	}

	partner.SetSettingsDefinition(map[string]string{
		"generators.message_identifier":             "mid",
		"generators.response_message_identifier":    "rmid",
		"generators.data_frame_identifier":          "dfid",
		"generators.reference_identifier":           "rid",
		"generators.reference_stop_area_identifier": "rsaid",
	})

	midGenerator = partner.IdentifierGenerator("message_identifier")
	if expected := "mid"; midGenerator.formatString != expected {
		t.Errorf("partner message_identifier IdentifierGenerator should be %v, got: %v ", expected, midGenerator.formatString)
	}
	midGenerator = partner.IdentifierGenerator("response_message_identifier")
	if expected := "rmid"; midGenerator.formatString != expected {
		t.Errorf("partner response_message_identifier IdentifierGenerator should be %v, got: %v ", expected, midGenerator.formatString)
	}
	midGenerator = partner.IdentifierGenerator("data_frame_identifier")
	if expected := "dfid"; midGenerator.formatString != expected {
		t.Errorf("partner data_frame_identifier IdentifierGenerator should be %v, got: %v ", expected, midGenerator.formatString)
	}
	midGenerator = partner.IdentifierGenerator("reference_identifier")
	if expected := "rid"; midGenerator.formatString != expected {
		t.Errorf("partner reference_identifier IdentifierGenerator should be %v, got: %v ", expected, midGenerator.formatString)
	}
	midGenerator = partner.IdentifierGenerator("reference_stop_area_identifier")
	if expected := "rsaid"; midGenerator.formatString != expected {
		t.Errorf("partner reference_stop_area_identifier IdentifierGenerator should be %v, got: %v ", expected, midGenerator.formatString)
	}
}
