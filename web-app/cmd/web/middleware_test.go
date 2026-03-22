package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"webapp/pkg/data"
)

func Test_application_addIPToContext(t *testing.T) {

	tests := []struct {
		headerName  string
		headerValue string
		addr        string
		emptyAddr   bool
	}{
		{"", "", "", false},
		{"", "", "", true},
		{"X-Forwarded-For", "192.3.2.1", "", false}, // encabezado de proxy en caso de que venga de un proxy
		{"", "", "hellow:world", false},
	}

	// crea un handler para revisar el contexto
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// obtener la IP del contexto
		val := r.Context().Value(contextUserKey)
		if val == nil {
			t.Error(contextUserKey, "not found in context")
		}

		ip, ok := val.(string)
		if !ok {
			t.Error("value in context not a string")
		}
		t.Log(ip)
	})

	for _, e := range tests {
		// crear handle to test
		handlerToTest := app.AddIPToContext(nextHandler)

		// crear solicitud
		req := httptest.NewRequest("GET", "http://testing", nil)
		if e.emptyAddr {
			req.RemoteAddr = ""
		}

		if len(e.headerName) > 0 {
			req.Header.Add(e.headerName, e.headerValue)
		}

		if len(e.addr) > 0 {
			req.RemoteAddr = e.addr
		}

		handlerToTest.ServeHTTP(httptest.NewRecorder(), req)

	}
}

func Test_application_ipFromContext(t *testing.T) {
	ctx := context.Background() //  — context.Background() te da un contexto vacío, raíz, sin valores ni deadlines.

	ctx = context.WithValue(ctx, contextUserKey, "whatever")

	ip := app.ipFromContext(ctx)

	if ip != "whatever" {
		t.Errorf("expected 'whatever' but got %s", ip)
	}
}

func Test_application_auth(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})

	var test = []struct {
		name   string
		isAuth bool
	}{
		{"authenticated user", true},
		{"non authenticated user", false},
	}

	for _, e := range test {
		handlerToTest := app.auth(nextHandler)

		req := httptest.NewRequest("GET", "http://testing", nil)
		req = addContextAndSessionToRequest(req, app)
		if e.isAuth {
			app.Session.Put(req.Context(), "user", data.User{ID: 1})
		}
		rr := httptest.NewRecorder()
		handlerToTest.ServeHTTP(rr, req)

		if e.isAuth {
			if rr.Code != http.StatusOK {
				t.Errorf("for %s expected status code 200 but got %d", e.name, rr.Code)
			}
		} else {
			if rr.Code != http.StatusSeeOther {
				t.Errorf("for %s expected status code 303 but got %d", e.name, rr.Code)
			}
		}

	}
}
