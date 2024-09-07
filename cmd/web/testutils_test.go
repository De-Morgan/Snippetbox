package main

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Create a newTestApplication helper which returns an instance of our
// application struct containing mocked dependencies.
func newTestApplication(_ testing.TB) *application {

	return &application{
		errorLog: log.New(io.Discard, "", 0),
		infoLog:  log.New(io.Discard, "", 0),
	}
}

// Define a custom testServer type which anonymously embeds a
// httptest.Server instance.
type testServer struct {
	*httptest.Server
}

func newTestServer(tb testing.TB, h http.Handler) *testServer {
	tb.Helper()
	ts := httptest.NewTLSServer(h)
	return &testServer{ts}
}

// Implement a get method on our custom testServer type. This makes a GET
// request to a given url path on the test server, and returns the response
// status code, headers and body.

func (ts *testServer) get(tb testing.TB, urlPath string) (statusCode int, header http.Header, body []byte) {

	rs, err := ts.Client().Get(ts.URL + urlPath)
	if err != nil {
		tb.Error(err)
		return
	}
	defer rs.Body.Close()

	body, err = io.ReadAll(rs.Body)
	if err != nil {
		tb.Error(err)
	}

	return rs.StatusCode, rs.Header, body

}
