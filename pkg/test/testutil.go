package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/day0ops/randomise-route-keys/pkg/config"
	"github.com/day0ops/randomise-route-keys/pkg/server"
)

func SetupTestRouter() *gin.Engine {
	router := gin.Default()
	gin.SetMode(gin.TestMode)
	server.RegisterRoutes(router)
	return router
}

type WrappedTestClient struct {
	Router *gin.Engine
}

func (tc *WrappedTestClient) PerformRequest(method, path string, body interface{}, headers map[string]string) *httptest.ResponseRecorder {
	var reqBody *bytes.Reader
	if body != nil {
		jsonBody, _ := json.Marshal(body)
		reqBody = bytes.NewReader(jsonBody)
	} else {
		reqBody = bytes.NewReader([]byte{})
	}

	req, _ := http.NewRequest(method, path, reqBody)

	// Set default headers
	req.Header.Set("Content-Type", "application/json")

	// Set additional headers if provided
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	w := httptest.NewRecorder()
	tc.Router.ServeHTTP(w, req)
	return w
}

func AssertResponseStatus(t *testing.T, response *httptest.ResponseRecorder, expectedStatus int) {
	assert.Equal(t, expectedStatus, response.Code)
}

func AssertResponse(t *testing.T, response *httptest.ResponseRecorder, expectedStatus int, expectedBody string) {
	assert.Equal(t, expectedStatus, response.Code)
	if expectedBody != "" {
		assert.JSONEq(t, expectedBody, response.Body.String())
	}
}

func AssertJsonKeyPresentInResponse(t *testing.T, response *httptest.ResponseRecorder, expectedStatus int, expectedKey string) {
	assert.Equal(t, expectedStatus, response.Code)
	if expectedKey != "" {
		assert.NotNil(t, response.Body)
		var resp map[string]string
		err := json.Unmarshal([]byte(response.Body.String()), &resp)
		_, exists := resp[expectedKey]
		assert.Nil(t, err)
		assert.True(t, exists)
	}
}

func ReadTestData(t *testing.T) {
	t.Setenv(config.RouteListFilePathEnvVar, "testdata/route-list.json")
	_ = server.ReadRouteListFile()
}

func NewTestClient() *WrappedTestClient {
	router := SetupTestRouter()
	client := &WrappedTestClient{Router: router}

	return client
}

func WrappedTestCase(test func(t *testing.T, tc *WrappedTestClient)) func(*testing.T) {
	return func(t *testing.T) {
		context := NewTestClient()
		context.beforeEach()
		defer context.afterEach()
		test(t, context)
	}
}

func (tc *WrappedTestClient) beforeEach() {
}

func (tc *WrappedTestClient) afterEach() {
	server.ClearCachedRouteList()
}
