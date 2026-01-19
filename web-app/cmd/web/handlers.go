package main

import (
	"html/template"
	"log"
	"net/http"
	"path"
	"time"
	"webapp/pkg/data"
)

var pathToTemplates = "./templates/"

type templateData struct {
	IP    string
	Data  map[string]any
	Error string
	Flash string
	User  data.User
}

func (app *application) render(response http.ResponseWriter, request *http.Request,
	name string, dataSubmit *templateData) error {

	parseTemplate, err := template.ParseFiles(path.Join(pathToTemplates, name), path.Join(pathToTemplates, "Base.layout.gohtml"))
	if err != nil {
		http.Error(response, "Internal Server Error", http.StatusBadRequest)
		return err
	}

	dataSubmit.IP = app.ipFromContext(request.Context())

	dataSubmit.Error = app.Session.PopString(request.Context(), "errors")
	dataSubmit.Flash = app.Session.PopString(request.Context(), "flash")
	log.Println("Flash message:", dataSubmit.Flash)
	log.Println("Error message:", dataSubmit.Error)

	err = parseTemplate.Execute(response, dataSubmit)
	if err != nil {
		http.Error(response, "Internal Server Error"+err.Error(), http.StatusBadRequest)
		return err
	}

	return nil
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

func (app *application) Profile(w http.ResponseWriter, r *http.Request) {
	_ = app.render(w, r, "profile.page.gohtml", &templateData{})
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
		// redirect logiun page with errors message
		app.Session.Put(request.Context(), "errors", "invalid login credentials")
		http.Redirect(response, request, "/", http.StatusSeeOther)
		return
	}

	email := request.Form.Get("email")
	password := request.Form.Get("password")

	user, err := app.DB.GetUserByEmail(email)
	if err != nil {
		app.Session.Put(request.Context(), "errors", "invalid login")
		http.Redirect(response, request, "/", http.StatusSeeOther)
		return
	}

	if !app.Authenticate(request, user, password) {
		app.Session.Put(request.Context(), "errors", "invalid login credentials")
		http.Redirect(response, request, "/", http.StatusSeeOther)
		return
	}

	// prevent fixation attacks
	_ = app.Session.RenewToken(request.Context())

	app.Session.Put(request.Context(), "flash", "Successfully logged in!")

	log.Println("User authenticated successfully")
	http.Redirect(response, request, "/user/profile", http.StatusSeeOther)

}

func (app *application) Authenticate(r *http.Request, user *data.User, password string) bool {
	if valid, err := user.PasswordMatches(password); err != nil || !valid {
		return false
	}
	app.Session.Put(r.Context(), "user", user)
	return true
}
