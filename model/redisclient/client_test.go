package redisclient

import (
	"context"
	"encoding/json"
	"testing"

	"bitbucket.org/enroute-mobi/ara/logger"
)

type codes map[string]code

func (identifiers codes) MarshalJSON() ([]byte, error) {
	aux := map[string]string{}

	for codeSpace, code := range identifiers {
		aux[codeSpace] = code.Value
	}

	return json.Marshal(aux)
}

type code struct {
	CodeSpace string
	Value     string
}

type line struct {
	Id         string
	Codes      codes
	ReferentId string
}

func (tm *line) ModelId() string {
	return tm.Id
}

func TestRedisClient(t *testing.T) {
	logger.Log.Printf("%+v", TestClient.(*client).c.FT_List(context.Background()).Val())

	tm := &line{
		Id:         "1234",
		Codes:      codes{"internal": {"internal", "test:toto"}},
		ReferentId: "toto",
	}

	b, _ := json.Marshal(tm)
	logger.Log.Printf("%+v", string(b))

	if err := TestClient.Set(Line, tm); err != nil {
		t.Fatal(err)
	}

	r, err := TestClient.Get(Line, "1234")
	if err != nil {
		t.Fatal(err)
	}
	logger.Log.Printf("%+v", r)
	if r == "" {
		t.Error("Get should return a value")
	}

	docs, err := TestClient.FindAll(Line)
	if err != nil {
		t.Fatal(err)
	}
	logger.Log.Printf("%+v", docs)

	docs, err = TestClient.FindByCode(Line, "internal", "test:toto")
	if err != nil {
		t.Fatal(err)
	}
	logger.Log.Printf("%+v", docs)
	if len(docs) == 0 {
		t.Error("FindByCode should return a value")
	}

	docs, err = TestClient.FindBy(Line, ByReferentID, "toto")
	if err != nil {
		t.Fatal(err)
	}
	logger.Log.Printf("%+v", docs)
	if len(docs) == 0 {
		t.Error("FindByCode should return a value")
	}

	// t.Fatal("pouet")
}
