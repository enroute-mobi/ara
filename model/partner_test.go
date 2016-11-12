package model

import "testing"

func Test_Partner_Id(t *testing.T) {
	partner := Partner{
		id: "6ba7b814-9dad-11d1-0-00c04fd430c8",
	}

	if expected := PartnerId("6ba7b814-9dad-11d1-0-00c04fd430c8"); partner.Id() != expected {
		t.Errorf("Partner.Id() returns wrong value, got: %s, required: %s", partner.Id(), expected)
	}
}

func Test_Partner_Name(t *testing.T) {
	partner := Partner{
		name: "partner",
	}

	if expected := "partner"; partner.Name() != expected {
		t.Errorf("Partner.Name() returns wrong value, got: %s, required: %s", partner.Name(), expected)
	}
}

func Test_Partner_OperationnalStatus(t *testing.T) {
	partner := Partner{
		name: "partner",
	}

	if expected := OPERATIONNAL_STATUS_UNKNOWN; partner.OperationnalStatus() != expected {
		t.Errorf("Partner.OperationnalStatus() returns wrong status, got: %s, required: %s", partner.OperationnalStatus(), expected)
	}
}

func Test_Partner_MarshalJSON(t *testing.T) {
	partner := Partner{
		id:   "6ba7b814-9dad-11d1-0-00c04fd430c8",
		name: "partner",
	}
	expected := `{"Id":"6ba7b814-9dad-11d1-0-00c04fd430c8","Name":"partner"}`
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
	partners := NewPartnerManager()
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

func Test_PartnerManager_New(t *testing.T) {
	partners := NewPartnerManager()
	partner := partners.New("partner")

	if partner.Name() != "partner" {
		t.Errorf("New should create a partner with given name name:\n got: %s\n want: %s", partner.Name(), "partner")
	}
	if partner.Id() != "" {
		t.Errorf("New Partner identifier should be an empty string, got: %s", partner.Id())
	}
}

func Test_PartnerManager_Save(t *testing.T) {
	partners := NewPartnerManager()
	partner := partners.New("partner")

	if success := partners.Save(partner); !success {
		t.Errorf("Save should return true")
	}

	if partner.Id() == "" {
		t.Errorf("New Partner identifier should not be an empty string")
	}
}

func Test_PartnerManager_Find_NotFound(t *testing.T) {
	partners := NewPartnerManager()
	partner := partners.Find("6ba7b814-9dad-11d1-0-00c04fd430c8")
	if partner != nil {
		t.Errorf("Find should return false when Partner isn't found")
	}
}

func Test_PartnerManager_Find(t *testing.T) {
	partners := NewPartnerManager()

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
	partners := NewPartnerManager()

	existingPartner := partners.New("partner")
	partners.Save(existingPartner)

	partnerId := existingPartner.Id()

	partners.Delete(existingPartner)

	partner := partners.Find(partnerId)
	if partner != nil {
		t.Errorf("Deleted Partner should not be findable")
	}
}
