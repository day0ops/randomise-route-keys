package test

import (
	"net/http"
	"testing"
)

func TestHealthzHandler(t *testing.T) {
	// Set up the Test Client using the reusable function
	client := NewTestClient()
	// Perform a request to the route
	response := client.PerformRequest("GET", "/healthz", nil, nil)
	// Assert the response
	AssertResponse(t, response, http.StatusOK, "")
}

func TestNonExistentHandler(t *testing.T) {
	// Set up the Test Client using the reusable function
	client := NewTestClient()
	// Perform a request to the route
	response := client.PerformRequest("GET", "/non-existence", nil, nil)
	// Assert the response
	AssertResponseStatus(t, response, http.StatusNotFound)
}

func TestRootHandler(t *testing.T) {
	t.Run("handle an empty file", WrappedTestCase(func(t *testing.T, tc *WrappedTestClient) {
		// Perform a request to the route
		response := tc.PerformRequest("GET", "/", nil, nil)
		// Assert the response
		AssertResponseStatus(t, response, http.StatusInternalServerError)
	}))

	t.Run("handle valid file", WrappedTestCase(func(wt *testing.T, tc *WrappedTestClient) {
		// Setup test data
		ReadTestData(wt)
		// Perform a request to the route
		response := tc.PerformRequest("GET", "/", nil, nil)
		// Assert the response
		AssertJsonKeyPresentInResponse(wt, response, http.StatusOK, "decision")
	}))
}
