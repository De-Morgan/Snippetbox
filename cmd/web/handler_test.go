package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_ping(t *testing.T) {
	t.Parallel()
	rr := httptest.NewRecorder()

	r, err := http.NewRequest("GET", "/", nil)

	if err != nil {
		t.Fatal(err)
	}
	ping(rr, r)

	rs := rr.Result()
	if rs.StatusCode != http.StatusOK {
		t.Errorf("want %d, got %d", http.StatusOK, rs.StatusCode)
	}
	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(body) != "OK" {
		t.Errorf("want body to equal %q, got %q", "OK", string(body))
	}

}
func TestPing(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	code, _, body := ts.get(t, "/ping")
	if code != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, code)
	}
	if string(body) != "OK" {
		t.Errorf("want body to equal %q", "OK")
	}

}
