package model

import "testing"

func Test_Referential_Id(t *testing.T) {
	referential := Referential{
		id: "6ba7b814-9dad-11d1-0-00c04fd430c8",
	}

	if expected := ReferentialId("6ba7b814-9dad-11d1-0-00c04fd430c8"); referential.Id() != expected {
		t.Errorf("Referential.Slug() returns wrong value, got: %s, required: %s", referential.Id(), expected)
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

func Test_Referential_Model(t *testing.T) {
	model := NewMemoryModel()
	referential := Referential{
		model: model,
	}
	if referential.Model() != model {
		t.Errorf("Referential.Model() returns wrong value, got: %v, required: %v", referential.Model(), model)
	}
}

func Test_Referential_MarshalJSON(t *testing.T) {
	referential := Referential{
		id:   "6ba7b814-9dad-11d1-0-00c04fd430c8",
		slug: "referential",
	}
	expected := `{"Id":"6ba7b814-9dad-11d1-0-00c04fd430c8","Slug":"referential"}`
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
	_, ok = referentials.Find(referential.Id())
	if !ok {
		t.Errorf("New Referential should be found in Referentials manager")
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

	if success := referentials.Save(&referential); !success {
		t.Errorf("Save should return true")
	}

	if referential.Id() == "" {
		t.Errorf("New Referential identifier should not be an empty string")
	}
}

func Test_MemoryReferentials_Find_NotFound(t *testing.T) {
	referentials := NewMemoryReferentials()
	_, ok := referentials.Find("6ba7b814-9dad-11d1-0-00c04fd430c8")
	if ok {
		t.Errorf("Find should return false when Referential isn't found")
	}
}

func Test_MemoryReferentials_Find(t *testing.T) {
	referentials := NewMemoryReferentials()

	existingReferential := referentials.New(ReferentialSlug("referential"))
	referentials.Save(&existingReferential)
	referentialId := existingReferential.Id()

	referential, ok := referentials.Find(referentialId)
	if !ok {
		t.Errorf("Find should return true when Referential is found")
	}
	if referential.Id() != referentialId {
		t.Errorf("Find should return a Referential with the given Id")
	}
}

func Test_MemoryReferentials_FindBySlug(t *testing.T) {
	referentials := NewMemoryReferentials()

	referentialSlug := ReferentialSlug("referential")
	existingReferential := referentials.New(referentialSlug)
	referentials.Save(&existingReferential)

	referential, ok := referentials.FindBySlug(referentialSlug)
	if !ok {
		t.Errorf("FindBySlug should return true when Referential is found")
	}
	if referential.Slug() != referentialSlug {
		t.Errorf("FindBySlug should return a Referential with the given Slug")
	}
}

func Test_MemoryReferentials_Delete(t *testing.T) {
	referentials := NewMemoryReferentials()

	existingReferential := referentials.New(ReferentialSlug("referential"))
	referentials.Save(&existingReferential)

	referentialId := existingReferential.Id()

	referentials.Delete(&existingReferential)

	_, ok := referentials.Find(referentialId)
	if ok {
		t.Errorf("Deleted Referential should not be findable")
	}
}

func Test_MemoryReferentials_Load(t *testing.T) {
	initTestDb(t)
	defer cleanTestDb(t)

	// Insert Data in the test db
	var databaseReferential = struct {
		Referential_id string `db:"referential_id"`
		Slug           string `db:"slug"`
	}{
		Referential_id: "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
		Slug:           "ratp",
	}
	Database.AddTableWithName(databaseReferential, "referentials")
	err := Database.Insert(&databaseReferential)
	if err != nil {
		t.Fatal(err)
	}

	// Fetch data from the db
	referentials := NewMemoryReferentials()
	referentials.Load()

	referentialId := ReferentialId(databaseReferential.Referential_id)
	referential, ok := referentials.Find(referentialId)
	if !ok {
		t.Errorf("Loaded Referentials should be found")
	}

	if referential.Id() != referentialId {
		t.Errorf("Wrong Id:\n got: %v\n expected: %v", referential.Id(), referentialId)
	}
}
