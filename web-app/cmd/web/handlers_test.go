package main

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_application_handlers(t *testing.T) {
	var theTest = []struct {
		name       string
		url        string
		statusCode int
	}{
		{"home", "/", http.StatusBadRequest},
		{"404", "/not-found", http.StatusNotFound},
	}

	routes := app.routes()

	// create a test server
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, e := range theTest {
		resp, err := ts.Client().Get(ts.URL + e.url)
		if err != nil {
			t.Log(err)
			t.Fatal(err)
		}

		if resp.StatusCode != e.statusCode {
			t.Errorf("for %s expected %d but got %d", e.name, e.statusCode, resp.StatusCode)
		}
	}
}

func Test_application_home(t *testing.T) {

	req, _ := http.NewRequest("GET", "/", nil)
	req = addContextAndSessionToRequest(req, app)
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(app.Home)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("returned wrong status ciode; expected 200")
	}

	body, _ := io.ReadAll(rr.Body)
	if !strings.Contains(string(body), `<small>From session`) {
		t.Errorf("did not find expected string in body")
	}

}

func Test_application_homev2(t *testing.T) {
	var tests = []struct {
		name         string
		putInSession string
		expectedHtml string
	}{
		{"first visit", "", `<small>From session:`},
		{"second visit", "session_new_value_1", `<small>From session: session_new_value_1`},
		{"third visit", "session_new_value_2", `<small>From session: session_new_value_2`},
		{"fourth visit", "session_new_value_3", `<small>From session: session_new_value_3`},
	}

	for _, e := range tests {
		req, _ := http.NewRequest("GET", "/", nil)
		req = addContextAndSessionToRequest(req, app)

		_ = app.Session.Destroy(req.Context())

		if e.putInSession != "" {
			app.Session.Put(req.Context(), "test", e.putInSession)
		}

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(app.Home)
		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("returned wrong status ciode; expected 200")
		}

		body, _ := io.ReadAll(rr.Body)
		if !strings.Contains(string(body), e.expectedHtml) {
			t.Errorf("for %s did not find expected string in body", e.name)
		}
	}
}

func Test_application_render_bad_template(t *testing.T) {
	pathToTemplates = "./testdata/"
	req, _ := http.NewRequest("GET", "/", nil)
	req = addContextAndSessionToRequest(req, app)

	rr := httptest.NewRecorder()

	err := app.render(rr, req, "bad.page.gohtml", &templateData{})
	if err == nil {
		t.Error("expected error but did not get one")
	}

}

func getCtx(req *http.Request) context.Context {
	ctx := context.WithValue(req.Context(), contextUserKey, "inknown")
	return ctx
}

func addContextAndSessionToRequest(req *http.Request, app application) *http.Request {
	req = req.WithContext(getCtx(req))
	ctx, _ := app.Session.Load(req.Context(), req.Header.Get("X-Session"))
	return req.WithContext(ctx)
}
