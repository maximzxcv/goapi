package tests

import (
	"goapi/app/api/handlers"

	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	succeed = "\u2713"
	failed  = "\u2717"
)

// TestGetUsers ....
func TestGetUsers(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "/users", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.GetUsers)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	t.Logf("Checking GetUsers handler")
	{
		status := rr.Code
		switch {
		case status != http.StatusOK:
			t.Errorf("\t%s\thandler returned wrong status code: got %v want %v", failed, status, http.StatusOK)
		default:
			t.Logf("\t%s\thandler returned correct status code: got %v want %v", succeed, status, http.StatusOK)
		}

		// Check the response body is what we expect.
		expected := `[{"id":"3896de2f-3df2-4a75-ac61-be2f287696f7","name":"Name 1"},{"id":"5d45c3cf-d225-41fa-a27e-9ab57541cad9","name":"Name 2"},{"id":"","name":""},{"id":"","name":""},{"id":"","name":""}]`
		if rr.Body == nil {
			t.Errorf("\t%s\thandler returned unexpected body: \ngot  %v \nwant %v", failed, rr.Body.String(), expected)
		} else {
			t.Logf("\t%s\thandler returned body as expected: got %v want %v", succeed, rr.Body.String(), expected)
		}
	}
}
