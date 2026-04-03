package redisclient

import "testing"

type testModel struct {
	id string
	S  string
	I  int
}

func (tm *testModel) ModelId() string {
	return tm.id
}

func TestRedisClient(t *testing.T) {
	tm := &testModel{
		id: "1234",
		S:  "test string",
		I:  1,
	}

	if err := TestClient.Set("testModel", tm); err != nil {
		t.Fatal(err)
	}

	r, err := TestClient.Get("testModel", "1234")
	if err != nil {
		t.Fatal(err)
	}

	if r == "" {
		t.Error("Get should return a value")
	}
}
