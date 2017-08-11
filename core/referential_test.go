package core

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/af83/edwig/model"
)

func Test_Referential_Id(t *testing.T) {
	referential := Referential{
		id: "6ba7b814-9dad-11d1-0-00c04fd430c8",
	}

	if expected := ReferentialId("6ba7b814-9dad-11d1-0-00c04fd430c8"); referential.Id() != expected {
		t.Errorf("Referential.Id() returns wrong value, got: %s, required: %s", referential.Id(), expected)
	}
}

func Test_Referential_Slug(t *testing.T) {
	referential := Referential{
		slug: "referential",
	}

	if expected := ReferentialSlug("referential"); referential.Slug() != expected {
		t.Errorf("Referential.Slug() returns wrong value, got: %s, required: %s", referential.Slug(), expected)
	}
}

func Test_Referential_StartedAt(t *testing.T) {
	model.SetDefaultClock(model.NewFakeClock())
	referentials := NewMemoryReferentials()
	referential := referentials.New("slug")
	referential.Start()
	referential.Stop()

	if expected := time.Date(1984, time.April, 4, 0, 0, 0, 0, time.UTC); referential.StartedAt() != expected {
		t.Errorf("Referential.StartedAt() returns wrong value, got: %s, required: %s", referential.StartedAt(), expected)
	}
}

func Test_Referential_Model(t *testing.T) {
	model := model.NewMemoryModel()
	referential := Referential{
		model: model,
	}
	if referential.Model() != model {
		t.Errorf("Referential.Model() returns wrong value, got: %v, required: %v", referential.Model(), model)
	}
}

func Test_Referential_Partners(t *testing.T) {
	partners := createTestPartnerManager()
	referential := Referential{
		partners: partners,
	}
	if referential.Partners() != partners {
		t.Errorf("Referential.Partners() returns wrong value, got: %v, required: %v", referential.Partners(), partners)
	}
}

func Test_Referential_MarshalJSON(t *testing.T) {

	referential := &Referential{
		id:       "6ba7b814-9dad-11d1-0-00c04fd430c8",
		slug:     "referential",
		Settings: map[string]string{"key": "value"},
	}
	referential.partners = NewPartnerManager(referential)
	expected := `{"Id":"6ba7b814-9dad-11d1-0-00c04fd430c8","Slug":"referential","Settings":{"key":"value"}}`
	jsonBytes, err := referential.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}

	jsonString := string(jsonBytes)
	if jsonString != expected {
		t.Errorf("Referential.MarshalJSON() returns wrong json:\n got: %s\n want: %s", jsonString, expected)
	}
}

func Test_Referential_Save(t *testing.T) {
	referentials := NewMemoryReferentials()
	referential := referentials.New(ReferentialSlug("referential"))

	if referential.manager != referentials {
		t.Errorf("New referential manager should be referentials")
	}

	ok := referential.Save()
	if !ok {
		t.Errorf("referential.Save() should succeed")
	}
	referential = referentials.Find(referential.Id())
	if referential == nil {
		t.Errorf("New Referential should be found in Referentials manager")
	}
}

func Test_Referential_setNextReloadAt(t *testing.T) {
	var conditions = []struct {
		setting        string
		clockHour      int
		clockMinute    int
		expectedDay    int
		expectedHour   int
		expectedMinute int
	}{
		{"04:00", 15, 0, 2, 4, 0},
		{"04:00", 3, 0, 1, 4, 0},
		{"04:00", 4, 0, 2, 4, 0},
		{"abc", 15, 0, 2, 4, 0},
	}

	for _, condition := range conditions {
		referential := Referential{Settings: map[string]string{"model.reloadAt": condition.setting}}

		fakeClock := model.NewFakeClockAt(time.Date(2017, time.January, 1, condition.clockHour, condition.clockMinute, 0, 0, time.UTC))
		referential.SetClock(fakeClock)

		referential.setNextReloadAt()

		expected := time.Date(2017, time.January, condition.expectedDay, condition.expectedHour, condition.expectedMinute, 0, 0, time.UTC)
		if !referential.NextReloadAt().Equal(expected) {
			t.Errorf("Wrong NextReloadAt:\n expected: %v\n got: %v", expected, referential.NextReloadAt())
		}
	}
}

func Test_APIReferential_Validate(t *testing.T) {
	referentials := NewMemoryReferentials()
	// Check empty Slug
	apiReferential := &APIReferential{
		manager: referentials,
	}
	valid := apiReferential.Validate()

	if valid {
		t.Errorf("Validate should return false")
	}
	if len(apiReferential.Errors) != 1 {
		t.Errorf("apiReferential Errors should not be empty")
	}
	if len(apiReferential.Errors["Slug"]) != 1 || apiReferential.Errors["Slug"][0] != ERROR_BLANK {
		t.Errorf("apiReferential should have Error for Slug, got %v", apiReferential.Errors)
	}

	// Check Already Used Slug
	referential := referentials.New("slug")
	referentials.Save(referential)
	apiReferential = &APIReferential{
		Slug:    "slug",
		manager: referentials,
	}
	valid = apiReferential.Validate()

	if valid {
		t.Errorf("Validate should return false")
	}
	if len(apiReferential.Errors) != 1 {
		t.Errorf("apiReferential Errors should not be empty")
	}
	if len(apiReferential.Errors["Slug"]) != 1 || apiReferential.Errors["Slug"][0] != ERROR_UNIQUE {
		t.Errorf("apiReferential should have Error for Slug, got %v", apiReferential.Errors)
	}

	// Check ok
	apiReferential = &APIReferential{
		Slug:    "slug2",
		manager: referentials,
	}
	valid = apiReferential.Validate()

	if !valid {
		t.Errorf("Validate should return true")
	}
	if len(apiReferential.Errors) != 0 {
		t.Errorf("apiReferential Errors should be empty")
	}
}

func Test_MemoryReferentials_New(t *testing.T) {
	referentials := NewMemoryReferentials()
	referential := referentials.New(ReferentialSlug("referential"))

	if referential.Slug() != "referential" {
		t.Errorf("New should create a referential with given slug slug:\n got: %s\n want: %s", referential.Slug(), "referential")
	}
	if referential.Id() != "" {
		t.Errorf("New Referential identifier should be an empty string, got: %s", referential.Id())
	}

	if referential.NextReloadAt().IsZero() {
		t.Errorf("New Referential should have a defined NextReloadAt time, got: %s", referential.NextReloadAt())
	}
}

func Test_MemoryReferentials_Save(t *testing.T) {
	referentials := NewMemoryReferentials()
	referential := referentials.New(ReferentialSlug("referential"))

	if success := referentials.Save(referential); !success {
		t.Errorf("Save should return true")
	}

	if referential.Id() == "" {
		t.Errorf("New Referential identifier should not be an empty string")
	}
}

func Test_MemoryReferentials_Find_NotFound(t *testing.T) {
	referentials := NewMemoryReferentials()
	referential := referentials.Find("6ba7b814-9dad-11d1-0-00c04fd430c8")
	if referential != nil {
		t.Errorf("Find should return nil when Referential isn't found")
	}
}

func Test_MemoryReferentials_Find(t *testing.T) {
	referentials := NewMemoryReferentials()

	existingReferential := referentials.New(ReferentialSlug("referential"))
	referentials.Save(existingReferential)
	referentialId := existingReferential.Id()

	referential := referentials.Find(referentialId)
	if referential == nil {
		t.Errorf("Find should return a Referential")
	}
	if referential.Id() != referentialId {
		t.Errorf("Find should return a Referential with the given Id")
	}
}

func Test_MemoryReferentials_FindBySlug(t *testing.T) {
	referentials := NewMemoryReferentials()

	referentialSlug := ReferentialSlug("referential")
	existingReferential := referentials.New(referentialSlug)
	referentials.Save(existingReferential)

	referential := referentials.FindBySlug(referentialSlug)
	if referential == nil {
		t.Errorf("FindBySlug should return a Referential")
	}
	if referential.Slug() != referentialSlug {
		t.Errorf("FindBySlug should return a Referential with the given Slug")
	}
}

func Test_MemoryReferentials_Delete(t *testing.T) {
	referentials := NewMemoryReferentials()

	existingReferential := referentials.New(ReferentialSlug("referential"))
	referentials.Save(existingReferential)

	referentialId := existingReferential.Id()

	referentials.Delete(existingReferential)

	referential := referentials.Find(referentialId)
	if referential != nil {
		t.Errorf("Deleted Referential should not be findable")
	}
}

func Test_MemoryReferentials_Load(t *testing.T) {
	model.InitTestDb(t)
	defer model.CleanTestDb(t)

	// Insert Data in the test db
	dbRef := model.DatabaseReferential{
		ReferentialId: "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
		Slug:          "ratp",
		Settings:      "{ \"test.key\": \"test-value\" }",
		Tokens:        "[\"apiToken\"]",
	}
	err := model.Database.Insert(&dbRef)
	if err != nil {
		t.Fatal(err)
	}

	// Fetch data from the db
	referentials := NewMemoryReferentials()
	err = referentials.Load()
	if err != nil {
		t.Fatal(err)
	}

	referentialId := ReferentialId(dbRef.ReferentialId)
	referential := referentials.Find(referentialId)
	if referential == nil {
		t.Errorf("Loaded Referentials should be found")
	}

	if referential.Id() != referentialId {
		t.Errorf("Wrong Id:\n got: %v\n expected: %v", referential.Id(), referentialId)
	}
	if expected := map[string]string{"test.key": "test-value"}; !reflect.DeepEqual(referential.Settings, expected) {
		t.Errorf("Wrong Settings:\n got: %#v\n expected: %#v", referential.Settings, expected)
	}
	if expected := "ratp"; referential.Slug() != ReferentialSlug(expected) {
		t.Errorf("Wrong Slug:\n got: %v\n expected: %v", referential.Slug(), expected)
	}
	if expected := "apiToken"; len(referential.Tokens) != 1 || referential.Tokens[0] != expected {
		t.Errorf("Wrong Tokens:\n got: %v\n expected: %v", referential.Tokens, expected)
	}
}

func Test_MemoryReferentials_SaveToDatabase(t *testing.T) {
	model.InitTestDb(t)
	defer model.CleanTestDb(t)

	// Insert Referential in the test db
	referentials := NewMemoryReferentials()
	ref := referentials.New("slug")
	ref.Save()

	errors := referentials.SaveToDatabase()
	if len(errors) != 0 {
		t.Fatalf("Error while saving Referentials: %v", errors)
	}

	// Insert to times to check uniqueness constraints
	ref2 := referentials.New("slug2")
	ref2.Settings = map[string]string{"setting": "value"}
	ref2.Tokens = []string{"token"}
	ref2.Save()

	errors = referentials.SaveToDatabase()
	if len(errors) != 0 {
		t.Fatalf("Error while saving Referentials: %v", errors)
	}

	// Check Referentials
	referentials2 := NewMemoryReferentials()
	err := referentials2.Load()
	if err != nil {
		t.Fatal(err)
	}

	if referential := referentials2.Find(ReferentialId(ref.id)); referential == nil {
		t.Errorf("Loaded Referentials should be found")
	}
	referential := referentials2.Find(ReferentialId(ref2.id))
	if referential == nil {
		t.Fatalf("Loaded Referentials should be found")
	}
	if referential.slug != "slug2" {
		t.Errorf("Wrong Referential Slug, got: %v want: slug2", referential.slug)
	}
	if len(referential.Settings) != 1 || referential.Setting("setting") != "value" {
		t.Errorf("Wrong Referential Settings, got: %v want {\"setting\":\"value\"}", referential.Settings)
	}
	if len(referential.Tokens) != 1 || referential.Tokens[0] != "token" {
		t.Errorf("Wrong Referential tokens, got: %v want: [token]", referential.Tokens)
	}
}

func Test_MemoryReferentials_SaveToDatabase_CleanPartners(t *testing.T) {
	model.InitTestDb(t)
	defer model.CleanTestDb(t)

	// Insert Partner in the test db
	dbPartner := model.DatabasePartner{
		Id:             "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
		ReferentialId:  "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
		Slug:           "ratp",
		Settings:       "{}",
		ConnectorTypes: "[]",
	}
	err := model.Database.Insert(&dbPartner)
	if err != nil {
		t.Fatal(err)
	}

	// Insert Referential in the test db
	referentials := NewMemoryReferentials()
	ref := referentials.New("slug")
	ref.Save()

	errors := referentials.SaveToDatabase()
	if len(errors) != 0 {
		t.Fatalf("Error while saving Referentials: %v", errors)
	}

	// Check Partner
	selectPartners := []model.SelectPartner{}
	sqlQuery := fmt.Sprintf("select * from partners")
	_, err = model.Database.Select(&selectPartners, sqlQuery)
	if err != nil {
		t.Fatalf("Error while fetching partners: %v", err)
	}

	if len(selectPartners) != 0 {
		t.Errorf("Partner should not be found")
	}
}

func Test_MemoryReferentials_SaveToDatabase_PartnerWithoutReferential(t *testing.T) {
	model.InitTestDb(t)
	defer model.CleanTestDb(t)

	referentials := NewMemoryReferentials()
	ref := referentials.New("slug")
	ref.Save()

	partner := ref.partners.New("slug")
	partner.Save()

	errors := ref.partners.SaveToDatabase()
	if len(errors) != 1 {
		t.Fatalf("Partner save should return an error, got: %v", errors)
	}
}

func Test_MemoryReferentials_SaveToDatabase_SavePartner(t *testing.T) {
	model.InitTestDb(t)
	defer model.CleanTestDb(t)

	// Insert Referential in the test db
	referentials := NewMemoryReferentials()
	ref := referentials.New("slug")
	ref.Save()

	errors := referentials.SaveToDatabase()
	if len(errors) != 0 {
		t.Fatalf("Error while saving Referentials: %v", errors)
	}

	// Insert Partner in the test db
	partners := ref.partners
	partner := partners.New("slug")
	partner.Save()

	errors = partners.SaveToDatabase()
	if len(errors) != 0 {
		t.Fatalf("Error while saving Partners: %v", errors)
	}

	// Save data in the DB 2 times to check uniqueness constraints
	partner2 := partners.New("slug2")
	partner2.Settings = map[string]string{"setting": "value"}
	partner2.ConnectorTypes = []string{"connector"}
	partner2.Save()

	errors = partners.SaveToDatabase()
	if len(errors) != 0 {
		t.Fatalf("Error while saving Partners: %v", errors)
	}

	// Check Partners
	partners2 := NewPartnerManager(ref)
	err := partners2.Load()
	if err != nil {
		t.Fatal(err)
	}

	if p := partners2.Find(PartnerId(partner.id)); p == nil {
		t.Errorf("Loaded Partners should be found")
	}
	testPartner := partners2.Find(PartnerId(partner2.id))
	if testPartner == nil {
		t.Fatalf("Loaded Partners should be found")
	}
	if testPartner.slug != "slug2" {
		t.Errorf("Wrong Partner Slug, got: %v want: slug2", testPartner.slug)
	}
	if len(testPartner.Settings) != 1 || testPartner.Setting("setting") != "value" {
		t.Errorf("Wrong Partner Settings, got: %v want {\"setting\":\"value\"}", testPartner.Settings)
	}
	if len(testPartner.ConnectorTypes) != 1 || testPartner.ConnectorTypes[0] != "connector" {
		t.Errorf("Wrong Partner ConnectorTypes, got: %v want [connector]", testPartner.ConnectorTypes)
	}
}
