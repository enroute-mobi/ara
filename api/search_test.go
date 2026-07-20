package api

import (
	"net/url"
	"testing"

	"bitbucket.org/enroute-mobi/ara/model"
	"github.com/stretchr/testify/assert"
)

func Test_SearchByCode_Errors(t *testing.T) {
	assert := assert.New(t)

	s := &model.StopArea{}
	slice := []*model.StopArea{s}

	params := url.Values{}
	values, err := searchByCode(slice, params)
	assert.NoError(err)
	assert.ElementsMatch(values, slice, "should return the full list if there is no params \"code\"")

	params.Set("code", "fake")
	_, err = searchByCode(slice, params)
	assert.Error(err)
	assert.Equal("invalid request: query parameter \"code\" : fake", err.Error())

	params.Set("code", ":")
	_, err = searchByCode(slice, params)
	assert.Error(err)
	assert.Equal("code space or value should not be empty", err.Error())

	params.Set("code", "external:xx")
	_, err = searchByCode(slice, params)
	assert.Error(err)
	assert.Equal("length of search value must be at least 3 characters, got: xx", err.Error())

	params.Set("code", "external:*$#$&(&#@*@")
	_, err = searchByCode(slice, params)
	assert.Error(err)
	assert.Equal("cannot create search pattern: error parsing regexp: missing argument to repetition operator: `*`", err.Error())
}

func Test_SearchByName_Errors(t *testing.T) {
	assert := assert.New(t)

	s := &model.StopArea{}
	slice := []*model.StopArea{s}

	params := url.Values{}
	values, err := searchByName(slice, params)
	assert.NoError(err)
	assert.ElementsMatch(values, slice, "should return the full list if there is no params \"code\"")

	params.Set("name", "(*+")
	_, err = searchByName(slice, params)
	assert.Error(err)
	assert.Equal("cannot create search pattern: error parsing regexp: missing argument to repetition operator: `*`", err.Error())

	params.Set("name", "xx")
	_, err = searchByName(slice, params)
	assert.Error(err)
	assert.Equal("length of search name must be at least 3 characters, got: xx", err.Error())
}

func Test_SearchByLineIds_Errors(t *testing.T) {
	assert := assert.New(t)

	v := &model.Vehicle{}
	slice := []*model.Vehicle{v}

	params := url.Values{}
	values, err := searchByLineIds(slice, params)
	assert.NoError(err)
	assert.ElementsMatch(values, slice, "should return the full list if there is no params \"code\"")

	params.Set("line_ids[]", "32132321")
	_, err = searchByLineIds(slice, params)
	assert.Error(err)
	assert.Equal("line id is not a valid UUID: 32132321", err.Error())
}

func Test_SearchByText_Errors(t *testing.T) {
	assert := assert.New(t)

	s := &model.Situation{}
	slice := []*model.Situation{s}

	params := url.Values{}
	values, err := searchByText(slice, params)
	assert.NoError(err)
	assert.ElementsMatch(values, slice, "should return the full list if there is no params \"text\"")

	params.Set("text", "xx")
	_, err = searchByText(slice, params)
	assert.Error(err)
	assert.Equal("length of search text must be at least 3 characters, got: xx", err.Error())

	params.Set("text", "(*+")
	_, err = searchByText(slice, params)
	assert.Error(err)
	assert.Equal("cannot create search pattern: error parsing regexp: missing argument to repetition operator: `*`", err.Error())
}

func Test_SearchByAffectedLineIds_Errors(t *testing.T) {
	assert := assert.New(t)

	s := &model.Situation{}
	slice := []*model.Situation{s}

	params := url.Values{}
	values, err := searchByAffectedLineIds(slice, params)
	assert.NoError(err)
	assert.ElementsMatch(values, slice, "should return the full list if there is no params \"line_ids[]\"")

	params.Set("line_ids[]", "32132321")
	_, err = searchByAffectedLineIds(slice, params)
	assert.Error(err)
	assert.Equal("line id is not a valid UUID: 32132321", err.Error())
}

func Test_SearchByAffectedStopAreaIds_Errors(t *testing.T) {
	assert := assert.New(t)

	s := &model.Situation{}
	slice := []*model.Situation{s}

	params := url.Values{}
	values, err := searchByAffectedStopAreaIds(slice, params)
	assert.NoError(err)
	assert.ElementsMatch(values, slice, "should return the full list if there is no params \"stop_area_ids[]\"")

	params.Set("stop_area_ids[]", "32132321")
	_, err = searchByAffectedStopAreaIds(slice, params)
	assert.Error(err)
	assert.Equal("stop area id is not a valid UUID: 32132321", err.Error())
}
