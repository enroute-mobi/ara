package core

import "testing"

func Test_LocalCredentialsIndex_simple(t *testing.T) {
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
