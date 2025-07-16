package model

import (
	"encoding/json"

	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Facility_Id(t *testing.T) {
	assert := assert.New(t)

	facility := Facility{id: "6ba7b814-9dad-11d1-0-00c04fd430c8"}

	assert.Equal(FacilityId("6ba7b814-9dad-11d1-0-00c04fd430c8"), facility.Id())
}

func Test_Facility_MarshalJSON(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	facility := Facility{
		Origin: "partner1",
		id:     "6ba7b814-9dad-11d1-0-00c04fd430c8",
		Status: "available",
	}

	facility.codes = make(Codes)
	code := NewCode("codeSpace", "value")
	facility.SetCode(code)

	expected := `{"Id":"6ba7b814-9dad-11d1-0-00c04fd430c8","Codes":{"codeSpace":"value"},"Status":"available","Origin":"partner1"}`

	jsonBytes, err := facility.MarshalJSON()

	require.NoError(err)
	assert.JSONEq(expected, string(jsonBytes))
}

func Test_Facility_UnmarshalJSON(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	text := `{
    "Codes": { "reflex": "FR:77491:ZDE:34004:STIF", "hastus": "sqypis" },
    "Status": "available"
  }`

	facility := Facility{}
	err := json.Unmarshal([]byte(text), &facility)
	require.NoError(err)

	expectedCodes := []Code{
		NewCode("reflex", "FR:77491:ZDE:34004:STIF"),
		NewCode("hastus", "sqypis"),
	}

	for _, expectedCode := range expectedCodes {
		code, found := facility.Code(expectedCode.CodeSpace())
		assert.True(found)
		assert.Equal(expectedCode, code)
	}

	assert.Equal(FacilityStatusAvailable, facility.Status)
}

func Test_Facility_Save(t *testing.T) {
	assert := assert.New(t)

	model := NewTestMemoryModel()
	facility := model.Facilities().New()
	code := NewCode("codeSpace", "value")
	facility.SetCode(code)

	assert.Equal(model, facility.model)

	ok := facility.Save()
	assert.True(ok)

	_, ok = model.Facilities().Find(facility.Id())
	assert.True(ok)

	_, ok = model.Facilities().FindByCode(code)
	assert.True(ok)
}

func Test_Facility_Code(t *testing.T) {
	assert := assert.New(t)

	facility := Facility{id: "6ba7b814-9dad-11d1-0-00c04fd430c8"}
	facility.codes = make(Codes)
	code := NewCode("codeSpace", "value")
	facility.SetCode(code)

	foundCode, ok := facility.Code("codeSpace")
	assert.True(ok)
	assert.Equal("value", foundCode.Value())

	_, ok = facility.Code("wrongkind")
	assert.False(ok)

	assert.Len(facility.Codes(), 1)
}

func Test_MemoryFacilities_New(t *testing.T) {
	assert := assert.New(t)

	facilities := NewMemoryFacilities()
	facility := facilities.New()

	assert.Equal(FacilityId(""), facility.Id())
}

func Test_MemoryFacilities_Save(t *testing.T) {
	assert := assert.New(t)

	facilities := NewMemoryFacilities()
	facility := facilities.New()

	ok := facilities.Save(facility)
	assert.True(ok)
	assert.NotZero(facility.Id())
}

func Test_MemoryFacilities_Find_NotFound(t *testing.T) {
	assert := assert.New(t)

	facilities := NewMemoryFacilities()
	_, ok := facilities.Find("6ba7b814-9dad-11d1-0-00c04fd430c8")
	assert.False(ok, "Find should return false when Facility isn't found")
}

func Test_MemoryFacilities_Find(t *testing.T) {
	assert := assert.New(t)

	facilities := NewMemoryFacilities()

	existingFacility := facilities.New()
	facilities.Save(existingFacility)

	facilityId := existingFacility.Id()

	facility, ok := facilities.Find(facilityId)
	assert.True(ok)
	assert.Equal(facilityId, facility.Id())
}

func Test_MemoryFacilities_FindAll(t *testing.T) {
	assert := assert.New(t)

	facilities := NewMemoryFacilities()

	for i := 0; i < 5; i++ {
		existingFacility := facilities.New()
		facilities.Save(existingFacility)
	}

	assert.Len(facilities.FindAll(), 5)
}

func Test_MemoryFacilities_Delete(t *testing.T) {
	assert := assert.New(t)

	facilities := NewMemoryFacilities()
	existingFacility := facilities.New()
	code := NewCode("codeSpace", "value")
	existingFacility.SetCode(code)
	ok := facilities.Save(existingFacility)
	assert.True(ok)

	ok = facilities.Delete(existingFacility)
	assert.True(ok)

	_, ok = facilities.Find(existingFacility.Id())
	assert.False(ok, "Deleted Facility should not be findable")

	_, ok = facilities.FindByCode(code)
	assert.False(ok, "Deleted Facility should not be findable by code")
}
