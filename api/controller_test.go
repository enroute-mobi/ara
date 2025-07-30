package api

import (
	"net/url"
	"testing"

	"bitbucket.org/enroute-mobi/ara/model"
	"github.com/stretchr/testify/assert"
)

func Test_Paginate(t *testing.T) {
	assert := assert.New(t)

	slice := []*model.Situation{}
	FULL_LIST_LENGTH := 35
	for i := 0; i <= FULL_LIST_LENGTH-1; i++ {
		s := &model.Situation{}
		slice = append(slice, s)
	}

	var TestCases = []struct {
		page           string
		per_page       string
		error          bool
		errorMessage   string
		expectedLength int
		message        string
	}{
		{
			error:          false,
			expectedLength: FULL_LIST_LENGTH,
			message:        "When no pagination is given, it should return the full list",
		},
		{
			page:           "1",
			error:          false,
			expectedLength: DEFAULT_PER_PAGE,
			message:        "When page=1 and no per_page, should return page 1 with the DEFAULT_PER_PAGE size",
		},
		{
			page:           "1",
			per_page:       "20",
			error:          false,
			expectedLength: 20,
			message:        "When page=1 and per_page=20, should return page 1 with 20 items",
		},
		{
			page:           "1",
			per_page:       "80",
			error:          false,
			expectedLength: DEFAULT_PER_PAGE,
			message:        "When page=1 and per_page=80, should return page 1 with the DEFAULT_PER_PAGE size",
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
		}

		if tt.error == true {
			assert.Equal(tt.errorMessage, err.Error())
		}
	}
}
