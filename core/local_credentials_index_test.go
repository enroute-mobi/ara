package core

import "testing"

func Test_LocalCredentialsIndex_Simple(t *testing.T) {
	index := NewLocalCredentialsIndex()

	index.Index("id", "credential, credential2")

	if _, ok := index.Find("credential"); !ok {
		t.Error("Can't find after index: ", index)
	}
	if _, ok := index.Find("credential2"); !ok {
		t.Error("Can't find after index: ", index)
	}
}

func Test_LocalCredentialsIndex_Multiple(t *testing.T) {
	index := NewLocalCredentialsIndex()

	index.Index("id", "credential")
	index.Index("id2", "credential2")

	id, ok := index.Find("credential")
	if !ok {
		t.Error("Can't find after index: ", index)
	}
	if expected := "id"; string(id) != expected {
		t.Errorf("Wrong Id returned, got: %v want: %v", id, expected)
	}

	id, ok = index.Find("credential2")
	if !ok {
		t.Error("Can't find after index: ", index)
	}
	if expected := "id2"; string(id) != expected {
		t.Errorf("Wrong Id returned, got: %v want: %v", id, expected)
	}
}

func Test_LocalCredentialsIndex_Change(t *testing.T) {
	index := NewLocalCredentialsIndex()

	index.Index("id", "credential, credential2")
	index.Index("id", "credential ,credential3")

	if _, ok := index.Find("credential"); !ok {
		t.Error("Can't find after index: ", index)
	}
	if _, ok := index.Find("credential3"); !ok {
		t.Error("Can't find after index: ", index)
	}
	if _, ok := index.Find("credential2"); ok {
		t.Error("Shouldn't find after index: ", index)
	}
}

func Test_LocalCredentialsIndex_Delete(t *testing.T) {
	index := NewLocalCredentialsIndex()

	index.Index("id", "credential, credential2")
	index.Delete("id")

	if _, ok := index.Find("credential"); ok {
		t.Error("Shouldn't find after index: ", index)
	}
	if _, ok := index.Find("credential2"); ok {
		t.Error("Shouldn't find after index: ", index)
	}
}

func Test_LocalCredentialsIndex_UniqCredentials(t *testing.T) {
	index := NewLocalCredentialsIndex()

	index.Index("id", "credential, credential2")
	index.Index("id2", "credential3, credential4")

	if !index.UniqCredentials("id", "credential5") || !index.UniqCredentials("id3", "credential5") {
		t.Error("New credential should be uniq: ", index)
	}
	if !index.UniqCredentials("id", "credential") {
		t.Error("Credential for the same id should be uniq: ", index)
	}
	if index.UniqCredentials("id", "credential3") {
		t.Error("Credential for another id should not be uniq: ", index)
	}
}

func Test_LocalCredentialsIndex_EmptyString(t *testing.T) {
	index := NewLocalCredentialsIndex()

	index.Index("id", "")

	_, ok := index.Find("")
	if !ok {
		t.Error("Can't find after index: ", index)
	}
	if !index.UniqCredentials(PartnerId("id"), "") {
		t.Error("Empty string should return true to uniq for the same id: ", index)
	}
	if index.UniqCredentials(PartnerId("id2"), "") {
		t.Error("Empty string should return false to uniq for another id: ", index)
	}

	index.Index("id", "credential")
	_, ok = index.Find("")
	if ok {
		t.Error("Empty string shouldn't be found after modification: ", index)
	}
	_, ok = index.Find("credential")
	if !ok {
		t.Error("Credential should be found after modification: ", index)
	}
	if !index.UniqCredentials(PartnerId("id"), "") || !index.UniqCredentials(PartnerId("id2"), "") {
		t.Error("Empty string should return true to uniq: ", index)
	}

	index.Index("id", "")
	_, ok = index.Find("")
	if !ok {
		t.Error("Empty string should be found after modification: ", index)
	}
	_, ok = index.Find("credential")
	if ok {
		t.Error("Credential shouldn't be found after modification: ", index)
	}
	if !index.UniqCredentials(PartnerId("id"), "") {
		t.Error("Empty string should return true to uniq for the same id: ", index)
	}
	if index.UniqCredentials(PartnerId("id2"), "") {
		t.Error("Empty string should return false to uniq for another id: ", index)
	}
}
