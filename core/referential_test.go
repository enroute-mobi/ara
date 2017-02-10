package core

import (
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
	referential := &Referential{
		slug: "referential",
	}
	referential.partners = NewPartnerManager(referential)
	referential.modelGuardian = NewModelGuardian(referential)
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
	expected := `{"Id":"6ba7b814-9dad-11d1-0-00c04fd430c8","Partners":[],"Settings":{"key":"value"},"Slug":"referential"}`
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
	var databaseReferential = struct {
		Referential_id string `db:"referential_id"`
		Slug           string `db:"slug"`
		Settings       string `db:"settings"`
	}{
		Referential_id: "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
		Slug:           "ratp",
		Settings:       "{}",
	}

	model.Database.AddTableWithName(databaseReferential, "referentials")
	err := model.Database.Insert(&databaseReferential)
	if err != nil {
		t.Fatal(err)
	}

	// Fetch data from the db
	referentials := NewMemoryReferentials()
	err = referentials.Load()
	if err != nil {
		t.Fatal(err)
	}

	referentialId := ReferentialId(databaseReferential.Referential_id)
	referential := referentials.Find(referentialId)
	if referential == nil {
		t.Errorf("Loaded Referentials should be found")
	}

	if referential.Id() != referentialId {
		t.Errorf("Wrong Id:\n got: %v\n expected: %v", referential.Id(), referentialId)
	}
}
