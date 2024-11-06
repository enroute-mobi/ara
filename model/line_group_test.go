package model

import (
	"encoding/json"
	"testing"
	"time"

	"bitbucket.org/enroute-mobi/ara/clock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_LineGroup_Id(t *testing.T) {
	assert := assert.New(t)
	lineGroup := LineGroup{
		id: "6ba7b814-9dad-11d1-0-00c04fd430c8",
	}

	assert.Equal(LineGroupId("6ba7b814-9dad-11d1-0-00c04fd430c8"), lineGroup.Id())
}

func Test_LineGroup_MarshalJSON(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	lineGroup := LineGroup{
		id:      "6ba7b814-9dad-11d1-0-00c04fd430c8",
		Name:    "LineGroupName",
		LineIds: []LineId{LineId("d9efc4a7-6164-4d5a-905d-ab35b1de9f87")},
	}

	expected := `{"Name":"LineGroupName","LineIds":["d9efc4a7-6164-4d5a-905d-ab35b1de9f87"],"Id":"6ba7b814-9dad-11d1-0-00c04fd430c8"}`

	jsonBytes, err := lineGroup.MarshalJSON()
	require.NoError(err)
	assert.JSONEq(expected, string(jsonBytes))
}

func Test_LineGroup_UnmarshalJSON(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	test := `{
		"Name": "LineGroupName",
        "ShortName": "short_name",
        "LineIds": ["d9efc4a7-6164-4d5a-905d-ab35b1de9f87"]
		}`

	lineGroup := LineGroup{}
	err := json.Unmarshal([]byte(test), &lineGroup)
	require.NoError(err)

	assert.Equal("LineGroupName", lineGroup.Name)
	assert.Equal("short_name", lineGroup.ShortName)
	assert.ElementsMatch([]LineId{LineId("d9efc4a7-6164-4d5a-905d-ab35b1de9f87")}, lineGroup.LineIds)
}

func Test_LineGroup_Save(t *testing.T) {
	assert := assert.New(t)

	model := NewTestMemoryModel()
	lineGroup := model.LineGroups().New()
	assert.Equal(model, lineGroup.model, "New lineGroup model should be MemoryLineGroup model")

	ok := lineGroup.Save()
	assert.True(ok, "lineGroup.Save() should succeed")

	_, ok = model.LineGroups().Find(lineGroup.Id())
	assert.True(ok, "New lineGroup should be found in MemoryLineGroup")
}

func Test_MemoryLineGroups_New(t *testing.T) {
	assert := assert.New(t)

	lineGroups := NewMemoryLineGroups()

	lineGroup := lineGroups.New()
	assert.Empty(string(lineGroup.Id()), "New lineGroup identifier should be an empty string")
}

func Test_MemoryLineGroups_Save(t *testing.T) {
	assert := assert.New(t)

	lineGroups := NewMemoryLineGroups()

	lineGroup := lineGroups.New()
	success := lineGroups.Save(lineGroup)
	assert.True(success)
	assert.NotEqual("", lineGroup.Id(), "New lineGroup identifier shouldn't be an empty string")
}

func Test_MemoryLineGroups_Find_NotFound(t *testing.T) {
	assert := assert.New(t)

	lineGroups := NewMemoryLineGroups()
	_, ok := lineGroups.Find("6ba7b814-9dad-11d1-0-00c04fd430c8")
	assert.False(ok, "Find should return false when LineGroup isn't found")
}

func Test_MemoryLineGroups_Find(t *testing.T) {
	assert := assert.New(t)
	lineGroups := NewMemoryLineGroups()

	existingLineGroup := lineGroups.New()
	lineGroups.Save(existingLineGroup)

	lineGroupId := existingLineGroup.Id()

	lineGroup, ok := lineGroups.Find(lineGroupId)
	assert.True(ok, "Find should return true when lineGroup is found")
	assert.Equal(lineGroupId, lineGroup.Id(), "Find should return a lineGroup with the given Id")
}

func Test_MemoryLineGroups_FindAll(t *testing.T) {
	assert := assert.New(t)

	lineGroups := NewMemoryLineGroups()
	for i := 0; i < 5; i++ {
		existingLineGroup := lineGroups.New()
		lineGroups.Save(existingLineGroup)
	}

	foundLineGroups := lineGroups.FindAll()
	assert.Len(foundLineGroups, 5, "FindAll should return all lineGroups")
}

func Test_MemoryLineGroups_Delete(t *testing.T) {
	assert := assert.New(t)

	lineGroups := NewMemoryLineGroups()
	existingLineGroup := lineGroups.New()

	lineGroups.Save(existingLineGroup)

	lineGroups.Delete(existingLineGroup)

	_, ok := lineGroups.Find(existingLineGroup.Id())
	assert.False(ok, "Deleted lineGroup should not be findable")
}

func Test_MemoryLineGroups_Load(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	InitTestDb(t)
	defer CleanTestDb(t)

	clock.SetDefaultClock(clock.NewFakeClock())
	defer clock.SetDefaultClock(clock.NewRealClock())

	// Insert Data in the test db
	databaseLineGroup := DatabaseLineGroup{
		Id:              "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
		Name:            "lineGroup1",
		ModelName:       "2017-01-01",
		ReferentialSlug: "referential",
		ShortName:       "short_name",
		LineIds:         `["d0eebc99-9c0b","e0eebc99-9c0b"]`,
	}

	Database.AddTableWithName(databaseLineGroup, "line_groups")
	err := Database.Insert(&databaseLineGroup)
	require.NoError(err)

	// Fetch data from the db
	model := NewTestMemoryModel()
	model.date = Date{
		Year:  2017,
		Month: time.January,
		Day:   1,
	}
	lineGroups := model.LineGroups().(*MemoryLineGroups)
	err = lineGroups.Load("referential")
	require.NoError(err)

	lineGroup, ok := lineGroups.Find(LineGroupId(databaseLineGroup.Id))
	require.True(ok)

	assert.Equal(LineGroupId("a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"), lineGroup.Id())
	assert.Equal("lineGroup1", lineGroup.Name)
	assert.Equal("short_name", lineGroup.ShortName)
	assert.Len(lineGroup.LineIds, 2)
	assert.Equal(LineId("d0eebc99-9c0b"), lineGroup.LineIds[0])
	assert.Equal(LineId("e0eebc99-9c0b"), lineGroup.LineIds[1])
}
