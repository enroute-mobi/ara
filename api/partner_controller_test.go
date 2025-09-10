package api

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"bitbucket.org/enroute-mobi/ara/core"
	"bitbucket.org/enroute-mobi/ara/core/partners"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type RequestTestData struct {
	Server *Server
	Id     string
	Method string
	Action string
	Body   []byte
}

func checkPartnerResponseStatus(responseRecorder *httptest.ResponseRecorder, t *testing.T) {
	require := require.New(t)

	require.Equal(http.StatusOK, responseRecorder.Code)
	require.Equal("application/json", responseRecorder.Header().Get("Content-Type"))
}

func preparePartnerRequest(method string, sendIdentifier bool, body []byte, t *testing.T) (partner *core.Partner, responseRecorder *httptest.ResponseRecorder, referential *core.Referential) {
	// Create a referential
	referentials := core.NewMemoryReferentials()
	server := &Server{}
	server.SetReferentials(referentials)
	referential = referentials.New("default")
	referential.Tokens = []string{"testToken"}
	referential.Save()

	// Set the fake UUID generator
	uuid.SetDefaultUUIDGenerator(uuid.NewFakeUUIDGenerator())
	// Save a new Partner
	referential.Partners().SetUUIDGenerator(uuid.NewFakeUUIDGenerator())
	partner = referential.Partners().New("first_partner")
	referential.Partners().Save(partner)

	address := []byte("/default/partners")
	if sendIdentifier {
		address = append(address, fmt.Sprintf("/%s", partner.Id())...)
	}
	request, err := http.NewRequest(method, string(address), bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("Authorization", "Token token=testToken")

	// Create a ResponseRecorder
	responseRecorder = httptest.NewRecorder()

	request.SetPathValue("referential_slug", string(referential.Slug()))
	request.SetPathValue("model", "partners")
	switch method {
	case "GET":
		if sendIdentifier == false && body == nil {
			server.handleReferentialModelIndex(responseRecorder, request)
		} else {
			request.SetPathValue("id", string(partner.Id()))
			server.handleReferentialModelShow(responseRecorder, request)
		}
	case "POST":
		server.handleReferentialModelCreate(responseRecorder, request)
	case "PUT":
		request.SetPathValue("id", string(partner.Id()))
		server.handleReferentialModelUpdate(responseRecorder, request)
	case "DELETE":
		request.SetPathValue("id", string(partner.Id()))
		server.handleReferentialModelDelete(responseRecorder, request)
	default:
		t.Fatalf("Unknown method: %s", method)
	}

	return
}

func Test_PartnerController_SubscriptionsCreate(t *testing.T) {
	assert := assert.New(t)

	// Create a referential
	referentials := core.NewMemoryReferentials()
	server := &Server{}
	server.SetReferentials(referentials)
	referential := referentials.New("default")
	referential.Tokens = []string{"testToken"}
	referential.Save()

	// Set the fake UUID generator
	uuid.SetDefaultUUIDGenerator(uuid.NewFakeUUIDGenerator())
	// Save a new Partner
	referential.Partners().SetUUIDGenerator(uuid.NewFakeUUIDGenerator())
	partner := referential.Partners().New("first_partner")
	referential.Partners().Save(partner)

	address := []byte("/default/partners")
	address = append(address, fmt.Sprintf("/%s", partner.Slug())...)

	body := []byte(`{ "Kind": "kind"}`)

	request, err := http.NewRequest("GET", string(address), bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("Authorization", "Token token=testToken")

	// Create a ResponseRecorder
	responseRecorder := httptest.NewRecorder()

	// sending request
	request.SetPathValue("referential_slug", string(referential.Slug()))
	request.SetPathValue("id", "first_partner")

	server.handlePartnerSubscriptionsCreate(responseRecorder, request)

	assert.Equal(http.StatusOK, responseRecorder.Code)
	assert.Len(partner.Subscriptions().FindAll(), 1)

	subscription := partner.Subscriptions().FindAll()[0]
	assert.Equal("kind", subscription.Kind())
}

func Test_PartnerController_SubscriptionsIndex(t *testing.T) {
	assert := assert.New(t)

	// Create a referential
	referentials := core.NewMemoryReferentials()
	server := &Server{}
	server.SetReferentials(referentials)
	referential := referentials.New("default")
	referential.Tokens = []string{"testToken"}
	referential.Save()

	// Set the fake UUID generator
	uuid.SetDefaultUUIDGenerator(uuid.NewFakeUUIDGenerator())
	// Save a new Partner
	referential.Partners().SetUUIDGenerator(uuid.NewFakeUUIDGenerator())
	partner := referential.Partners().New("first_partner")
	referential.Partners().Save(partner)

	// Create Subscription
	subscripions := partner.Subscriptions()
	subscription := subscripions.New("kind")
	subscription.Save()

	address := []byte("/default/partners")
	address = append(address, fmt.Sprintf("/%s", partner.Slug())...)

	request, err := http.NewRequest("GET", string(address), nil)
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("Authorization", "Token token=testToken")

	// Create a ResponseRecorder
	responseRecorder := httptest.NewRecorder()

	// sending request
	request.SetPathValue("referential_slug", string(referential.Slug()))
	request.SetPathValue("id", "first_partner")

	server.handlePartnerSubscriptionsIndex(responseRecorder, request)

	assert.Equal(http.StatusOK, responseRecorder.Code)
	expected := `[{"Kind":"kind","SubscriptionRef":"6ba7b814-9dad-11d1-0-00c04fd430c8"}]`
	assert.JSONEq(expected, responseRecorder.Body.String())
}

func Test_PartnerController_SubscriptionsDelete(t *testing.T) {
	assert := assert.New(t)

	// Create a referential
	referentials := core.NewMemoryReferentials()
	server := &Server{}
	server.SetReferentials(referentials)
	referential := referentials.New("default")
	referential.Tokens = []string{"testToken"}
	referential.Save()

	// Set the fake UUID generator
	uuid.SetDefaultUUIDGenerator(uuid.NewFakeUUIDGenerator())
	// Save a new Partner
	referential.Partners().SetUUIDGenerator(uuid.NewFakeUUIDGenerator())
	partner := referential.Partners().New("first_partner")
	referential.Partners().Save(partner)

	// Create Subscription
	subscripions := partner.Subscriptions()
	subscription := subscripions.New("kind")
	subscription.Save()

	// Test subscriptions
	assert.Len(partner.Subscriptions().FindAll(), 1)

	address := []byte("/default/partners")
	address = append(address, fmt.Sprintf("/%s", partner.Slug())...)

	request, err := http.NewRequest("DELETE", string(address), nil)
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("Authorization", "Token token=testToken")

	// Create a ResponseRecorder
	responseRecorder := httptest.NewRecorder()

	// sending request
	request.SetPathValue("referential_slug", string(referential.Slug()))
	request.SetPathValue("id", "first_partner")
	request.SetPathValue("subscription_id", string(subscription.Id()))

	server.handlePartnerSubscriptionsDelete(responseRecorder, request)

	assert.Equal(http.StatusOK, responseRecorder.Code)
	assert.Empty(partner.Subscriptions().FindAll())
}

func Test_PartnerController_Delete(t *testing.T) {
	assert := assert.New(t)

	// Send request
	partner, responseRecorder, referential := preparePartnerRequest("DELETE", true, nil, t)

	// Test response
	checkPartnerResponseStatus(responseRecorder, t)

	//Test Results
	deletedPartner := referential.Partners().Find(partner.Id())
	assert.Nil(deletedPartner, "Partner shouldn't be found after DELETE request")

	expected, err := partner.MarshalJSON()
	assert.NoError(err)
	assert.JSONEq(string(expected), responseRecorder.Body.String())
}

func Test_PartnerController_Update(t *testing.T) {
	assert := assert.New(t)

	// Prepare and send request
	body := []byte(`{ "Slug": "another_test", "Name": "test" }`)
	partner, responseRecorder, referential := preparePartnerRequest("PUT", true, body, t)

	// Check response
	checkPartnerResponseStatus(responseRecorder, t)

	// Test Results
	updatedPartner := referential.Partners().Find(partner.Id())
	assert.NotNil(updatedPartner, "Partner should be found after PUT request")
	assert.Equal(partners.Slug("another_test"), updatedPartner.Slug())
	assert.Equal("test", updatedPartner.Name)

	expected, err := updatedPartner.MarshalJSON()
	assert.NoError(err)
	assert.JSONEq(string(expected), responseRecorder.Body.String())

	assert.Empty(partner.Subscriptions().FindAll(), "All subscription should be deleted after a partner Edit")
}

func Test_PartnerController_UpdateConnectorTypes(t *testing.T) {
	assert := assert.New(t)

	// Prepare and send request
	body := []byte(`{ "ConnectorTypes": ["test-check-status-client"] }`)
	partner, responseRecorder, referential := preparePartnerRequest("PUT", true, body, t)

	// Check response
	checkPartnerResponseStatus(responseRecorder, t)

	// Test Results
	updatedPartner := referential.Partners().Find(partner.Id())
	assert.NotNil(updatedPartner, "Partner should be found after PUT request")
	assert.Equal(partners.Slug("first_partner"), updatedPartner.Slug())
	assert.Len(updatedPartner.ConnectorTypes, 1, "ConnectorTypes should have been updated by POST request")
}

func Test_PartnerController_Show(t *testing.T) {
	assert := assert.New(t)

	// Send request
	partner, responseRecorder, _ := preparePartnerRequest("GET", true, nil, t)

	// Check response
	checkPartnerResponseStatus(responseRecorder, t)

	//Test Results
	expected, err := partner.MarshalJSON()
	assert.NoError(err)
	assert.JSONEq(string(expected), responseRecorder.Body.String())
}

func Test_PartnerController_Create(t *testing.T) {
	assert := assert.New(t)

	// Prepare and send request
	body := []byte(`{ "Slug": "test" }`)
	_, responseRecorder, referential := preparePartnerRequest("POST", false, body, t)

	// Check response
	checkPartnerResponseStatus(responseRecorder, t)

	// Test Results
	assert.Len(referential.Partners().FindAll(), 2)

	partner, _ := referential.Partners().FindBySlug(partners.Slug("test"))
	assert.NotNil(partner)
	assert.Equal("test", string(partner.Slug()))

	expected, err := partner.MarshalJSON()
	assert.NoError(err)
	assert.JSONEq(string(expected), responseRecorder.Body.String())
}

func Test_PartnerController_Create_Invalid(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	// Prepare and send request
	body := []byte(`{ "Slug": "invalid_slug", "ConnectorTypes": ["test-validation-connector"] }`)
	_, responseRecorder, _ := preparePartnerRequest("POST", false, body, t)

	// Test response
	require.Equal(http.StatusBadRequest, responseRecorder.Code)
	require.Equal("application/json", responseRecorder.Header().Get("Content-Type")) //

	// Test Results
	expected := `{"Slug":"invalid_slug","ConnectorTypes":["test-validation-connector"],"Errors":{"slug":["Invalid format"]}}`
	assert.JSONEq(expected, responseRecorder.Body.String())
}

func Test_PartnerController_Index(t *testing.T) {
	assert := assert.New(t)

	// Send request
	_, responseRecorder, _ := preparePartnerRequest("GET", false, nil, t)

	// Test response
	checkPartnerResponseStatus(responseRecorder, t)

	//Test Results
	expected := `[{"Id":"6ba7b814-9dad-11d1-0-00c04fd430c8","Slug":"first_partner","PartnerStatus":{"OperationnalStatus":"unknown","RetryCount":0,"ServiceStartedAt":"0001-01-01T00:00:00Z"},"ConnectorTypes":[],"Settings":{}}]`
	assert.JSONEq(expected, responseRecorder.Body.String())
}

func Test_PartnerController_Save(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	model.InitTestDb(t)
	defer model.CleanTestDb(t)

	// Create a referential
	referentials := core.NewMemoryReferentials()
	referentials.SetUUIDGenerator(uuid.NewRealUUIDGenerator())
	server := &Server{}
	server.SetReferentials(referentials)
	referential := referentials.New("default")
	referential.Tokens = []string{"testToken"}
	referential.Save()
	status, refErr := referentials.SaveToDatabase()
	require.NoError(refErr)
	require.Equal(200, status, "Cannot save referentials to Database")

	// Initialize the partners manager
	referential.Partners().SetUUIDGenerator(uuid.NewRealUUIDGenerator())
	// Save a new partner
	partner := referential.Partners().New("First Partner")
	referential.Partners().Save(partner)

	// Create a request
	request, err := http.NewRequest("POST", "/default/partners/save", nil)
	require.NoError(err)

	request.Header.Set("Authorization", "Token token=testToken")

	// Create a ResponseRecorder
	responseRecorder := httptest.NewRecorder()

	request.SetPathValue("referential_slug", string(referential.Slug()))
	server.handleReferentialPartnerSave(responseRecorder, request)

	// Test response
	checkPartnerResponseStatus(responseRecorder, t)

	//Test Results
	selectPartners := []model.SelectPartner{}
	sqlQuery := "select * from partners"
	_, err = model.Database.Select(&selectPartners, sqlQuery)
	assert.Nil(err)
	assert.Len(selectPartners, 1)
	assert.Equal(selectPartners[0].Id, string(partner.Id()))
	assert.Equal(selectPartners[0].ReferentialId, string(referential.Id()))
}

func Test_PartnerController_FindPartner(t *testing.T) {
	assert := assert.New(t)

	// Create a referential
	referentials := core.NewMemoryReferentials()
	referential := referentials.New("default")
	referential.Save()

	// Save a new partner
	partner := referential.Partners().New("First Partner")
	referential.Partners().Save(partner)

	controller := &PartnerController{referential: referential}

	foundPartner := controller.findPartner("First Partner")
	assert.NotNil(foundPartner, "Can't find Partner by Slug")

	foundPartner = controller.findPartner(string(partner.Id()))
	assert.NotNil(foundPartner, "Can't find Partner by Id")
}
