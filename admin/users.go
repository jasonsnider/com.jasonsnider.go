package admin

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/jasonsnider/com.jasonsnider.go/internal/db"
	"github.com/jasonsnider/com.jasonsnider.go/internal/types"
	"github.com/jasonsnider/com.jasonsnider.go/templates"
)

type UserUpdateTemplate struct {
	Title            string
	Description      string
	Keywords         string
	Body             string
	ValidationErrors map[string]string
	User             types.User
	BustCssCache     string
	BustJsCache      string
}

func (app *App) CreateUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Create User")
}

func (app *App) ListUsers(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Read Users")
}

func (app *App) ViewUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Read User")
}

func (app *App) UpdateUser(w http.ResponseWriter, r *http.Request) {

	db := db.DB{DB: app.DB}
	vars := mux.Vars(r)
	id := vars["id"]

	var user types.User
	user, err := db.FetchUserById(id)

	if err != nil {
		http.Error(w, fmt.Sprintf("FetchUserById failed: %v", err), http.StatusInternalServerError)
		return
	}

	validationErrors := make(map[string]string)

	fmt.Println("GET")
	if r.Method == "POST" {
		fmt.Println("POST")
		validate := validator.New()

		user.ID = r.FormValue("id")
		user.FirstName = r.FormValue("first_name")
		user.LastName = r.FormValue("last_name")
		user.Email = r.FormValue("email")

		err := validate.Struct(user)

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
				case "min":
					errorMessage = fmt.Sprintf("%s must be at least %s characters long", fieldName, err.Param())
				case "eqfield":
					errorMessage = fmt.Sprintf("%s must match %s", fieldName, err.Param())
				default:
					errorMessage = fmt.Sprintf("%s is invalid", fieldName)
				}

				validationErrors[fieldName] = errorMessage
			}
		} else {

			query := `
				UPDATE users
				SET first_name = $1, last_name = $2, email = $3, username='bob', admin=true
				WHERE id = $4
			`

			tx, err := app.DB.Begin(context.Background())
			if err != nil {
				log.Fatalf("begin transaction failed: %v", err)
			}

			_, err = tx.Exec(context.Background(), query, user.FirstName, user.LastName, user.Email, user.ID)
			if err != nil {
				tx.Rollback(context.Background())
				log.Fatalf("update failed: %v", err)
			}

			err = tx.Commit(context.Background())
			if err != nil {
				log.Fatalf("commit transaction failed: %v", err)
			}

		}
	}

	pageTemplate := `
	{{define "content"}}
		<h1>Edit</h1>
		<form action="/admin/users/edit/{{.User.ID}}" method="POST">
			<input type="hidden" name="id" value="{{.User.ID}}">
			<div class="{{if index .ValidationErrors "FirstName"}}error{{end}}">
				<label for="first_name">First Name</label>
				<input type="text" id="FirstName" name="first_name" value="{{.User.FirstName}}">
				<div>{{if index .ValidationErrors "FirstName"}}{{index .ValidationErrors "FirstName"}}{{end}}</div>
			</div>
			<div class="{{if index .ValidationErrors "LastName"}}error{{end}}">
				<label for="last_name">Last Name</label>
				<input type="text" id="LastName" name="last_name" value="{{.User.LastName}}">
				<div>{{if index .ValidationErrors "LastName"}}{{index .ValidationErrors "LastName"}}{{end}}</div>
			</div>
			<div class="{{if index .ValidationErrors "Email"}}error{{end}}">
				<label for="subject">Email</label>
				<input type="email" id="email" name="email" value="{{.User.Email}}">
				<div>{{if index .ValidationErrors "Email"}}{{index .ValidationErrors "Email"}}{{end}}</div>
			</div>
			<button type="submit">Submit</button>
		</form>
	{{end}}
	`

	tmpl := template.Must(template.New("layout").Parse(templates.MainLayoutTemplate))
	tmpl = template.Must(tmpl.New("meta").Parse(templates.MetaDataTemplate))
	tmpl = template.Must(tmpl.New("registration").Parse(pageTemplate))

	pageData := UserUpdateTemplate{
		Title:            "Register your account",
		Description:      "Register your account",
		Keywords:         "resgistration",
		Body:             pageTemplate,
		ValidationErrors: validationErrors,
		User:             user,
		BustCssCache:     app.BustCssCache,
		BustJsCache:      app.BustJsCache,
	}

	err = tmpl.ExecuteTemplate(w, "layout", pageData)
	if err != nil {
		http.Error(w, fmt.Sprintf("Template execution failed: %v", err), http.StatusInternalServerError)
	}
}

func (app *App) DeleteUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Delete User")
}
