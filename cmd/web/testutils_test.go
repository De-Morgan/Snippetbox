package main

import (
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/golangcollege/sessions"
	"github.com/morgan/snippetbox/pkg/models/mock"
)

// Create a newTestApplication helper which returns an instance of our
// application struct containing mocked dependencies.
func newTestApplication(tb testing.TB) *application {
	// Create an instance of the template cache.
	templateCache, err := newTemplateCache("./../../ui/html/")
	if err != nil {
		tb.Fatal(err)
	}
	session := sessions.New([]byte("3dSm5MnygFHh7XidAtbskXrjbwfoJcbJ"))
	session.Lifetime = 12 * time.Hour
	session.Secure = true

	return &application{
		errorLog:      log.New(io.Discard, "", 0),
		infoLog:       log.New(io.Discard, "", 0),
		templateCache: templateCache,
		session:       session,
		snippets:      &mock.SnippetModel{},
		users:         &mock.UserModel{},
	}
}

// Define a custom testServer type which anonymously embeds a
// httptest.Server instance.
type testServer struct {
	*httptest.Server
}

func newTestServer(tb testing.TB, h http.Handler) *testServer {
	tb.Helper()
	jar, err := cookiejar.New(nil)
	if err != nil {
		tb.Fatal(err)
	}
	ts := httptest.NewTLSServer(h)
	ts.Client().Jar = jar

	// Disable redirect-following for the client. Essentially this function
	// is called after a 3xx response is received by the client, and returning
	// the http.ErrUseLastResponse error forces it to immediately return the
	// received response
	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
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

	return ts.extractResponseInfo(tb, rs)

}

func (ts *testServer) post(tb testing.TB, urlPath string, form url.Values) (statusCode int, header http.Header, body []byte) {
	rs, err := ts.Client().PostForm(ts.URL+urlPath, form)
	if err != nil {
		tb.Fatal(err)
	}
	defer rs.Body.Close()
	return ts.extractResponseInfo(tb, rs)

}

func (*testServer) extractResponseInfo(tb testing.TB, rs *http.Response) (int, http.Header, []byte) {
	tb.Helper()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		tb.Fatal(err)
	}
	return rs.StatusCode, rs.Header, body
}
