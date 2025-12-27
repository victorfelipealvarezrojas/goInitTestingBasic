package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path"
	"time"
)

var pathToTemplates = "./templates/"

type templateData struct {
	IP   string
	Data map[string]any
}

func (app *application) Home(response http.ResponseWriter, request *http.Request) {
	var td = make(map[string]any)

	if !app.Session.Exists(request.Context(), "test") {
		app.Session.Put(request.Context(), "test", time.Now().UTC().String())
	}

	if app.Session.Exists(request.Context(), "test") {
		msg := app.Session.GetString(request.Context(), "test")
		td["test"] = msg
	}

	_ = app.render(response, request, "Home.page.gohtml", &templateData{Data: td})
}

func (app *application) Login(response http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		log.Println("Error parsing form:", err)
		http.Error(response, "bad request", http.StatusBadRequest)
		return
	}

	// validate form inputs
	form := NewForm(request.PostForm)
	form.Required("email", "password")

	if !form.Valid() {
		log.Println("Form is not valid")
		http.Error(response, "invalid form submission", http.StatusBadRequest)
		return
	}

	email := request.Form.Get("email")
	password := request.Form.Get("password")

	log.Println(email, password)

	fmt.Fprint(response, email)
}

func (app *application) render(response http.ResponseWriter, request *http.Request,
	name string, data *templateData) error {

	parseTemplate, err := template.ParseFiles(path.Join(pathToTemplates, name), path.Join(pathToTemplates, "Base.layout.gohtml"))
	if err != nil {
		http.Error(response, "Internal Server Error", http.StatusBadRequest)
		return err
	}

	data.IP = app.ipFromContext(request.Context())

	err = parseTemplate.Execute(response, data)
	if err != nil {
		http.Error(response, "Internal Server Error"+err.Error(), http.StatusBadRequest)
		return err
	}

	return nil
}
