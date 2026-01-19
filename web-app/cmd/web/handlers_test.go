package main

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func Test_application_handlers(t *testing.T) {
	var theTest = []struct {
		name       string
		url        string
		statusCode int
	}{
		{"home", "/", http.StatusOK},
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

	pathToTemplates = "./../../templates/"

}

func getCtx(req *http.Request) context.Context {
	ctx := context.WithValue(req.Context(), contextUserKey, "inknown")
	return ctx
}

func Test_app_login(t *testing.T) {
	var Tests = []struct {
		name                     string
		postedData               url.Values
		expectedStatusCode       int
		expectedLocationRedirect string
	}{
		{
			name: "valid credemtial test",
			postedData: url.Values{
				"email":    {"admin@example.com"},
				"password": {"secret"},
			},
			expectedStatusCode:       http.StatusSeeOther,
			expectedLocationRedirect: "/user/profile",
		},
		{
			name: "invalid credemtial test",
			postedData: url.Values{
				"email":    {"invalid@example.com"},
				"password": {"invalid"},
			},
			expectedStatusCode:       http.StatusSeeOther,
			expectedLocationRedirect: "/",
		},
	}

	for _, e := range Tests {
		req, _ := http.NewRequest("POST", "/login", strings.NewReader(e.postedData.Encode()))
		req = addContextAndSessionToRequest(req, app)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder() // httptest.NewRecorder() crea un ResponseRecorder que actúa como un http.ResponseWriter falso para capturar la respuesta en tests.

		handler := http.HandlerFunc(app.Login) // crea un manejador HTTP a partir de la función app.Login
		handler.ServeHTTP(rr, req)             // sirve la solicitud HTTP utilizando el manejador y el request simulado

		if rr.Code != e.expectedStatusCode {
			t.Errorf("for %s expected status code %d but got %d", e.name, e.expectedStatusCode, rr.Code)
		}

		actualLocation, err := rr.Result().Location()

		if err == nil {
			if actualLocation.String() != e.expectedLocationRedirect {
				t.Errorf("for %s expected redirect to %s but got %s", e.name, e.expectedLocationRedirect, actualLocation.String())
			}
		} else {
			t.Errorf("%s: no location header set", e.name)
		}

	}

}

func addContextAndSessionToRequest(req *http.Request, app application) *http.Request {
	req = req.WithContext(getCtx(req))
	ctx, _ := app.Session.Load(req.Context(), req.Header.Get("X-Session"))
	return req.WithContext(ctx)
}
