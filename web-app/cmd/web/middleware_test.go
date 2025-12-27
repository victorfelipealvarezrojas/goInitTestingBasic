package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
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

	// creat un handler para revisar el contexto
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
		handlerToTest := app.AddIPToContext(nextHandler)

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
	ctx := context.Background()

	ctx = context.WithValue(ctx, contextUserKey, "whatever")

	ip := app.ipFromContext(ctx)

	if ip != "whatever" {
		t.Errorf("expected 'whatever' but got %s", ip)
	}

}
