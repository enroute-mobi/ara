package core

import (
	"reflect"
	"testing"

	"github.com/af83/edwig/model"
)

func createTestPartnerManager() *PartnerManager {
	referentials := NewMemoryReferentials()
	referential := referentials.New(ReferentialSlug("referential"))
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
	expected := `{"ConnectorTypes":[],"Id":"6ba7b814-9dad-11d1-0-00c04fd430c8","OperationnalStatus":"unknown","Settings":{},"Slug":"partner"}`
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
