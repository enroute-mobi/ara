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

func Test_Paginate(t *testing.T) {
	assert := assert.New(t)

	slice := []*model.Situation{}
	FULL_LIST_LENGTH := 35
	for i := 0; i <= FULL_LIST_LENGTH-1; i++ {
		s := &model.Situation{}
		slice = append(slice, s)
	}

	var TestCases = []struct {
		page               string
		per_page           string
		error              bool
		errorMessage       string
		expectedLength     int
		expectedTotalCount int
		message            string
	}{
		{
			error:              false,
			expectedLength:     FULL_LIST_LENGTH,
			expectedTotalCount: FULL_LIST_LENGTH,
			message:            "When no pagination is given, it should return the full list",
		},
		{
			page:               "1",
			error:              false,
			expectedLength:     DEFAULT_PER_PAGE,
			expectedTotalCount: FULL_LIST_LENGTH,
			message:            "When page=1 and no per_page, should return page 1 with the DEFAULT_PER_PAGE size",
		},
		{
			page:               "1",
			per_page:           "20",
			error:              false,
			expectedLength:     20,
			expectedTotalCount: FULL_LIST_LENGTH,
			message:            "When page=1 and per_page=20, should return page 1 with 20 items",
		},
		{
			page:               "1",
			per_page:           "80",
			error:              false,
			expectedLength:     DEFAULT_PER_PAGE,
			expectedTotalCount: FULL_LIST_LENGTH,
			message:            "When page=1 and per_page=80, should return page 1 with the DEFAULT_PER_PAGE size",
		},
		{
			page:         "WRONG",
			error:        true,
			errorMessage: "invalid request: query parameter \"page\": WRONG",
			message:      "When page is \"WRONG\" should return an error message",
		},
		{
			page:         "1",
			per_page:     "WRONG",
			error:        true,
			errorMessage: "invalid request: query parameter \"per_page\": WRONG",
			message:      "When page=1 and per_page is \"WRONG\" should return an error message",
		},
	}

	for _, tt := range TestCases {
		params := url.Values{}
		if tt.page != "" {
			params.Set("page", tt.page)
		}

		if tt.per_page != "" {
			params.Set("per_page", tt.per_page)
		}

		paginatedResource, err := paginate(slice, params)
		if tt.error == false {
			assert.NoError(err)
			assert.Len(paginatedResource.Models, tt.expectedLength, tt.message)
			assert.Equal(paginatedResource.TotalCount, tt.expectedTotalCount, tt.message)
		}

		if tt.error == true {
			assert.Equal(tt.errorMessage, err.Error())
		}
	}
}

func Test_Paginate_With_empty_models(t *testing.T) {
	assert := assert.New(t)

	slice := []*model.StopArea{}
	params := url.Values{}

	paginatedResource, err := paginate(slice, params)
	assert.NoError(err)
	assert.Equal([]*model.StopArea{}, paginatedResource.Models)
	assert.Equal(1, paginatedResource.CurrentPage)
	assert.Equal(0, paginatedResource.PerPage)
	assert.Equal(1, paginatedResource.TotalPages)
	assert.Equal(0, paginatedResource.TotalCount)
}
