package core

import "testing"

func TestErrorsAdd(t *testing.T) {
	e := NewErrors()

	e.Add("a", "m")
	if len(e.Get("a")) != 1 {
		t.Fatalf("Wrong number of errors, want 1 got %v", len(e.Get("a")))
	}
	e.Add("a", "m")
	if len(e.Get("a")) != 1 {
		t.Fatalf("Wrong number of errors, want 1 got %v", len(e.Get("a")))
	}
	e.Add("a", "m2")
	if len(e.Get("a")) != 2 {
		t.Fatalf("Wrong number of errors, want 2 got %v", len(e.Get("a")))
	}
}

func TestErrorsAddSettingError(t *testing.T) {
	e := NewErrors()

	e.AddSettingError("a", "m")
	if len(e.GetSettingError("a")) != 1 {
		t.Fatalf("Wrong number of errors, want 1 got %v", len(e.GetSettingError("a")))
	}
	e.AddSettingError("a", "m")
	if len(e.GetSettingError("a")) != 1 {
		t.Fatalf("Wrong number of errors, want 1 got %v", len(e.GetSettingError("a")))
	}
	e.AddSettingError("a", "m2")
	if len(e.GetSettingError("a")) != 2 {
		t.Fatalf("Wrong number of errors, want 2 got %v", len(e.GetSettingError("a")))
	}
}
