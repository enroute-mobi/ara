package model

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_StopAreaGroup_Id(t *testing.T) {
	assert := assert.New(t)
	stopAreaGroup := StopAreaGroup{
		id: "6ba7b814-9dad-11d1-0-00c04fd430c8",
	}

	assert.Equal(StopAreaGroupId("6ba7b814-9dad-11d1-0-00c04fd430c8"), stopAreaGroup.Id())
}

func Test_StopAreaGroup_MarshalJSON(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	stopAreaGroup := StopAreaGroup{
		id:          "6ba7b814-9dad-11d1-0-00c04fd430c8",
		Name:        "StopAreaGroupName",
		StopAreaIds: []StopAreaId{StopAreaId("d9efc4a7-6164-4d5a-905d-ab35b1de9f87")},
	}

	expected := `{"Name":"StopAreaGroupName","StopAreaIds":["d9efc4a7-6164-4d5a-905d-ab35b1de9f87"],"Id":"6ba7b814-9dad-11d1-0-00c04fd430c8"}`

	jsonBytes, err := stopAreaGroup.MarshalJSON()
	require.NoError(err)
	assert.JSONEq(expected, string(jsonBytes))
}

func Test_StopAreaGroup_UnmarshalJSON(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	test := `{
		"Name": "StopAreaGroupName",
        "ShortName": "short_name",
        "StopAreaIds": ["d9efc4a7-6164-4d5a-905d-ab35b1de9f87"]
		}`

	stopAreaGroup := StopAreaGroup{}
	err := json.Unmarshal([]byte(test), &stopAreaGroup)
	require.NoError(err)

	assert.Equal("StopAreaGroupName", stopAreaGroup.Name)
	assert.Equal("short_name", stopAreaGroup.ShortName)
	assert.ElementsMatch([]StopAreaId{StopAreaId("d9efc4a7-6164-4d5a-905d-ab35b1de9f87")}, stopAreaGroup.StopAreaIds)
}

func Test_StopAreaGroup_Save(t *testing.T) {
	assert := assert.New(t)

	model := NewTestMemoryModel()
	stopAreaGroup := model.StopAreaGroups().New()
	assert.Equal(model, stopAreaGroup.model, "New stopAreaGroup model should be MemoryStopAreaGroup model")

	ok := stopAreaGroup.Save()
	assert.True(ok, "stopAreaGroup.Save() should succeed")

	_, ok = model.StopAreaGroups().Find(stopAreaGroup.Id())
	assert.True(ok, "New stopAreaGroup should be found in MemoryStopAreaGroup")
}

func Test_MemoryStopAreaGroups_New(t *testing.T) {
	assert := assert.New(t)

	stopAreaGroups := NewMemoryStopAreaGroups()

	stopAreaGroup := stopAreaGroups.New()
	assert.Empty(string(stopAreaGroup.Id()), "New stopAreaGroup identifier should be an empty string")
}

func Test_MemoryStopAreaGroups_Save(t *testing.T) {
	assert := assert.New(t)

	stopAreaGroups := NewMemoryStopAreaGroups()

	stopAreaGroup := stopAreaGroups.New()
	success := stopAreaGroups.Save(stopAreaGroup)
	assert.True(success)
	assert.NotEqual("", stopAreaGroup.Id(), "New stopAreaGroup identifier shouldn't be an empty string")
}

func Test_MemoryStopAreaGroups_Find_NotFound(t *testing.T) {
	assert := assert.New(t)

	stopAreaGroups := NewMemoryStopAreaGroups()
	_, ok := stopAreaGroups.Find("6ba7b814-9dad-11d1-0-00c04fd430c8")
	assert.False(ok, "Find should return false when StopAreaGroup isn't found")
}

func Test_MemoryStopAreaGroups_Find(t *testing.T) {
	assert := assert.New(t)
	stopAreaGroups := NewMemoryStopAreaGroups()

	existingStopAreaGroup := stopAreaGroups.New()
	stopAreaGroups.Save(existingStopAreaGroup)

	stopAreaGroupId := existingStopAreaGroup.Id()

	stopAreaGroup, ok := stopAreaGroups.Find(stopAreaGroupId)
	assert.True(ok, "Find should return true when stopAreaGroup is found")
	assert.Equal(stopAreaGroupId, stopAreaGroup.Id(), "Find should return a stopAreaGroup with the given Id")
}

func Test_MemoryStopAreaGroups_FindAll(t *testing.T) {
	assert := assert.New(t)

	stopAreaGroups := NewMemoryStopAreaGroups()
	for i := 0; i < 5; i++ {
		existingStopAreaGroup := stopAreaGroups.New()
		stopAreaGroups.Save(existingStopAreaGroup)
	}

	foundStopAreaGroups := stopAreaGroups.FindAll()
	assert.Len(foundStopAreaGroups, 5, "FindAll should return all stopAreaGroups")
}

func Test_MemoryStopAreaGroups_Delete(t *testing.T) {
	assert := assert.New(t)

	stopAreaGroups := NewMemoryStopAreaGroups()
	existingStopAreaGroup := stopAreaGroups.New()

	stopAreaGroups.Save(existingStopAreaGroup)

	stopAreaGroups.Delete(existingStopAreaGroup)

	_, ok := stopAreaGroups.Find(existingStopAreaGroup.Id())
	assert.False(ok, "Deleted stopAreaGroup should not be findable")
}
