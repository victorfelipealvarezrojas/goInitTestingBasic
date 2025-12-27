package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Has(t *testing.T) {
	form := NewForm(nil)

	has := form.Has("whatever")
	if has {
		t.Error("Form shows has field when it should not")
	}

	postedData := url.Values{}
	postedData.Add("a", "a value")
	form = NewForm(postedData)

	has = form.Has("a")
	if !has {
		t.Error("Form shows does not have field when it should")
	}
}

func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := NewForm(r.PostForm)

	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("Form shows valid when required fields are missing")
	}

	postedData := url.Values{}
	postedData.Add("a", "a value")
	postedData.Add("b", "b value")
	postedData.Add("c", "c value")

	r, _ = http.NewRequest("POST", "/whatever", nil)
	r.PostForm = postedData

	form = NewForm(r.PostForm)
	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("Form shows invalid when required fields are present")
	}
}

func TestForm_RequiredV2(t *testing.T) {
	form := NewForm(url.Values{})
	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("Form shows valid when required fields are missing")
	}

	postedData := url.Values{}
	postedData.Add("a", "a value")
	postedData.Add("b", "b value")
	postedData.Add("c", "c value")

	form = NewForm(postedData)
	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("Form shows invalid when required fields are present")
	}
}

func TestForm_Check(t *testing.T) {
	form := NewForm(url.Values{})
	isEmail := form.Check(false, "email", "This field must be an email address")
	if isEmail {
		t.Error("Form shows valid email when it should not")
	}

	isEmail = form.Check(true, "email", "This field must be an email address")
	if !isEmail {
		t.Error("Form shows invalid email when it should not")
	}
}
