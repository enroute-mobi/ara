package core

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/af83/edwig/model"
)

func createTestPartnerManager() *PartnerManager {
	referentials := NewMemoryReferentials()
	referential := referentials.New(ReferentialSlug("referential"))
	referentials.Save(referential)
	return NewPartnerManager(referential)
}

func Test_Partner_Id(t *testing.T) {
	partner := Partner{
		id: "6ba7b814-9dad-11d1-0-00c04fd430c8",
	}

	if expected := PartnerId("6ba7b814-9dad-11d1-0-00c04fd430c8"); partner.Id() != expected {
		t.Errorf("Partner.Id() returns wrong value, got: %s, required: %s", partner.Id(), expected)
	}
}

func Test_Partner_Slug(t *testing.T) {
	partner := Partner{
		slug: "partner",
	}

	if expected := PartnerSlug("partner"); partner.Slug() != expected {
		t.Errorf("Partner.Slug() returns wrong value, got: %s, required: %s", partner.Id(), expected)
	}
}

func Test_Partner_OperationnalStatus(t *testing.T) {
	partner := NewPartner()

	if expected := OPERATIONNAL_STATUS_UNKNOWN; partner.OperationnalStatus() != expected {
		t.Errorf("Partner.OperationnalStatus() returns wrong status, got: %s, required: %s", partner.OperationnalStatus(), expected)
	}
}

func Test_Partner_MarshalJSON(t *testing.T) {
	partner := Partner{
		id:                 "6ba7b814-9dad-11d1-0-00c04fd430c8",
		operationnalStatus: OPERATIONNAL_STATUS_UNKNOWN,
		slug:               "partner",
		Settings:           make(map[string]string),
		ConnectorTypes:     []string{},
	}
	expected := `{"Id":"6ba7b814-9dad-11d1-0-00c04fd430c8","Slug":"partner","OperationnalStatus":"unknown","ConnectorTypes":[],"Settings":{}}`
	jsonBytes, err := partner.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}

	jsonString := string(jsonBytes)
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
	partner := Partner{connectors: make(map[string]Connector)}
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
	partner.Settings = make(map[string]string)
	stopAreaObjectId := model.NewObjectID("internal", "NINOXE:StopPoint:SP:24:LOC")

	partner.Settings["collect.include_stop_areas"] = "NINOXE:StopPoint:SP:24:LOC"
	if partner.CanCollect(stopAreaObjectId) != true {
		t.Errorf("Partner can collect should return true")
	}
}

func Test_Partner_CanCollectTrue2(t *testing.T) {
	partner := &Partner{}
	partner.Settings = make(map[string]string)
	stopAreaObjectId := model.NewObjectID("internal", "NINOXE:StopPoint:SP:24:LOC")

	if partner.CanCollect(stopAreaObjectId) != true {
		t.Errorf("Partner can collect should return true")
	}
}

func Test_Partner_CanCollectFalse(t *testing.T) {
	partner := &Partner{}
	partner.Settings = make(map[string]string)
	stopAreaObjectId := model.NewObjectID("internal", "BAD_VALUE")

	partner.Settings["collect.include_stop_areas"] = "NINOXE:StopPoint:SP:24:LOC"
	if partner.CanCollect(stopAreaObjectId) != false {
		t.Errorf("Partner can collect should return flase")
	}
}

func Test_Partners_FindAllByCollectPriority(t *testing.T) {
	partners := createTestPartnerManager()
	partner1 := Partner{}
	partner2 := Partner{}

	partner1.Settings = make(map[string]string)
	partner2.Settings = make(map[string]string)

	partner1.Settings["collect.priority"] = "2"
	partner1.SetSlug("First")

	partner2.Settings["collect.priority"] = "1"
	partner2.SetSlug("Second")

	partners.Save(&partner1)
	partners.Save(&partner2)

	orderedPartners := partners.FindAllByCollectPriority()
	if orderedPartners[0].Slug() != "First" {
		t.Errorf("Partners should be ordered")
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
	if len(apiPartner.Errors["Slug"]) != 1 || apiPartner.Errors["Slug"][0] != ERROR_BLANK {
		t.Errorf("apiPartner should have Error for Slug, got %v", apiPartner.Errors)
	}

	// Check Already Used Slug and local_credential
	partner := partners.New("slug")
	partner.Settings["local_credential"] = "cred"
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
	if len(apiPartner.Errors["Slug"]) != 1 || apiPartner.Errors["Slug"][0] != ERROR_UNIQUE {
		t.Errorf("apiPartner should have Error for Slug, got %v", apiPartner.Errors)
	}
	if len(apiPartner.Errors["Settings[\"local_credential\"]"]) != 1 || apiPartner.Errors["Settings[\"local_credential\"]"][0] != ERROR_UNIQUE {
		t.Errorf("apiPartner should have Error for local_credential, got %v", apiPartner.Errors)
	}

	// Check ok
	apiPartner = &APIPartner{
		Slug:     "slug2",
		Settings: map[string]string{"local_credential": "cred2"},
		manager:  partners,
	}
	valid = apiPartner.Validate()

	if !valid {
		t.Errorf("Validate should return true")
	}
	if len(apiPartner.Errors) != 0 {
		t.Errorf("apiPartner Errors should be empty")
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
	existingPartner.Settings["local_credential"] = "cred"
	partners.Save(existingPartner)

	partner, ok := partners.FindByLocalCredential("cred")
	if !ok {
		t.Fatal("FindByLocalCredential should return true when Partner is found")
	}
	if partner.Id() != existingPartner.Id() {
		t.Errorf("FindByLocalCredential should return a Partner with the given local_credential")
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
	var databasePartner = struct {
		Id             string `db:"id"`
		ReferentialId  string `db:"referential_id"`
		Slug           string `db:"slug"`
		Settings       string `db:"settings"`
		ConnectorTypes string `db:"connector_types"`
	}{
		Id:             "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
		ReferentialId:  string(referential.Id()),
		Slug:           "ratp",
		Settings:       "{}",
		ConnectorTypes: "[]",
	}

	model.Database.AddTableWithName(databasePartner, "partners")
	err := model.Database.Insert(&databasePartner)
	if err != nil {
		t.Fatal(err)
	}

	// Fetch data from the db
	partners := NewPartnerManager(referential)
	err = partners.Load()
	if err != nil {
		t.Fatal(err)
	}

	partnerId := PartnerId(databasePartner.Id)
	partner := partners.Find(partnerId)
	if partner == nil {
		t.Errorf("Loaded Partners should be found")
	} else if partner.Id() != partnerId {
		t.Errorf("Wrong Id:\n got: %v\n expected: %v", partner.Id(), partnerId)
	}
}
