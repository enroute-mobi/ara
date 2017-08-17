package model

import "testing"

func Test_StopAreaLineIds_Add(t *testing.T) {
	ids := StopAreaLineIds{}
	ids.Add(LineId("1234"))
	if len(ids) != 1 {
		t.Errorf("Ids shoud have len 1, got: %v", len(ids))
	}
	// Second time to test when already exists
	ids.Add(LineId("1234"))
	if len(ids) != 1 {
		t.Errorf("Ids shoud have len 1, got: %v", len(ids))
	}
}
