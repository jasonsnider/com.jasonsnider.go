package admin

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/jasonsnider/com.jasonsnider.go/internal/db"
	"github.com/jasonsnider/com.jasonsnider.go/internal/types"
	"github.com/jasonsnider/com.jasonsnider.go/pkg/passwords"
	"github.com/jasonsnider/com.jasonsnider.go/templates"
)

type AuthTemplate struct {
	Title            string
	Description      string
	Keywords         string
	Body             string
	ValidationErrors map[string]string
	Auth             types.Auth
	BustCssCache     string
	BustJsCache      string
}

func (app *App) Authenticate(w http.ResponseWriter, r *http.Request) {

	db := db.DB{DB: app.DB}
	auth := types.Auth{}
	validationErrors := make(map[string]string)

	if r.Method == "POST" {
		validate := validator.New()

		auth = types.Auth{
			Email:    r.FormValue("email"),
			Password: r.FormValue("password"),
		}

		err := validate.Struct(auth)

		if err != nil {
			for _, err := range err.(validator.ValidationErrors) {
				fieldName := err.Field()
				tag := err.Tag()

				var errorMessage string
				switch tag {
				case "required":
					errorMessage = fmt.Sprintf("%s is required", fieldName)
				case "email":
					errorMessage = fmt.Sprintf("%s must be a valid email address", fieldName)
				default:
					errorMessage = fmt.Sprintf("%s is invalid", fieldName)
				}

				validationErrors[fieldName] = errorMessage
			}
		} else {

			user, err := db.FetchAuth(auth.Email)

			if err == nil {

				confirm := passwords.CheckPasswordHash(auth.Password, user.Hash)
				if confirm {

					sessionExpiryStr := os.Getenv("SESSION_EXPIRY")
					sessionExpiry, err := strconv.Atoi(sessionExpiryStr)
					if err != nil {
						log.Printf("Invalid session expiry value: %v", err)
						http.Error(w, "Internal Server Error", http.StatusInternalServerError)
						return
					}

					session, _ := app.SessionStore.Get(r, "com-jasonsnider-go")
					session.Values["authenticated"] = true
					session.Values["user_email"] = user.Email
					session.Options.MaxAge = sessionExpiry
					err = session.Save(r, w)

					if err != nil {
						log.Printf("Failed to save session: %v", err)
					} else {
						log.Printf("Session saved for user: %s", user.Email)
					}

					// Read back the session data
					session, _ = app.SessionStore.Get(r, "com-jasonsnider-go")
					authenticated := session.Values["authenticated"]
					userEmail := session.Values["user_email"]
					log.Printf("Session data - Authenticated: %v, User Email: %s", authenticated, userEmail)

					http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)

					return
				} else {
					http.Error(w, fmt.Sprintf("Hash compare failed: %v", err), http.StatusInternalServerError)
					return
				}
			}

		}
	}

	pageTemplate := `
	{{define "content"}}
		<h1>Login</h1>
		<form action="/admin/login" method="POST">
			<div class="{{if index .ValidationErrors "Email"}}error{{end}}">
				<label for="subject">Email</label>
				<input type="email" id="email" name="email" value="{{.Auth.Email}}">
				<div>{{if index .ValidationErrors "Email"}}{{index .ValidationErrors "Email"}}{{end}}</div>
			</div>
			<div class="{{if index .ValidationErrors "Password"}}error{{end}}">
				<label for="body">Password</label>
				<input type="password" id="password" name="password" value="{{.Auth.Password}}">
				<div>{{if index .ValidationErrors "Password"}}{{index .ValidationErrors "Password"}}{{end}}</div>
			</div>
			<button type="submit">Login</button>
		</form>
	{{end}}
	`

	funcMap := template.FuncMap{
		"safeValue": types.SafeValue,
	}

	tmpl := template.Must(template.New("layout").Parse(templates.MainLayoutTemplate))
	tmpl = template.Must(tmpl.New("meta").Funcs(funcMap).Parse(templates.MetaDataTemplate))
	tmpl = template.Must(tmpl.New("registration").Parse(pageTemplate))

	pageData := AuthTemplate{
		Title:            "Login",
		Description:      "Login",
		Keywords:         "login",
		Body:             pageTemplate,
		ValidationErrors: validationErrors,
		Auth:             auth,
		BustCssCache:     app.BustCssCache,
		BustJsCache:      app.BustJsCache,
	}

	err := tmpl.ExecuteTemplate(w, "layout", pageData)
	if err != nil {
		http.Error(w, fmt.Sprintf("Template execution failed: %v", err), http.StatusInternalServerError)
	}
}

func (app *App) Logout(w http.ResponseWriter, r *http.Request) {

	session, err := app.SessionStore.Get(r, "com-jasonsnider-go")
	if err != nil {
		log.Printf("Failed to get session: %v", err)
		http.Error(w, "Failed to get session", http.StatusInternalServerError)
		return
	}

	// This tells Redis to delete the session
	session.Options.MaxAge = -1

	err = session.Save(r, w)
	if err != nil {
		log.Printf("Failed to delete session during logout: %v", err)
		http.Error(w, "Failed to log out", http.StatusInternalServerError)
		return
	}

	log.Println("User logged out successfully and session deleted")
	http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
}
