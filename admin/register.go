package admin

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/jasonsnider/com.jasonsnider.go/internal/db"
	"github.com/jasonsnider/com.jasonsnider.go/internal/types"
	"github.com/jasonsnider/com.jasonsnider.go/templates"
)

func (app *App) RegisterUser(w http.ResponseWriter, r *http.Request) {

	db := db.DB{DB: app.DB}
	user := types.RegisterUser{}
	validationErrors := make(map[string]string)

	if r.Method == "POST" {
		validate := validator.New()

		user = types.RegisterUser{
			FirstName:       r.FormValue("first_name"),
			LastName:        r.FormValue("last_name"),
			Email:           r.FormValue("email"),
			Password:        r.FormValue("password"),
			ConfirmPassword: r.FormValue("confirm_password"),
		}

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
			// Register the user
			err = db.RegisterUser(user)
			if err != nil {
				log.Fatalf("failed to register user: %v", err)
			}

			log.Println("User registered successfully")
		}
	}

	pageTemplate := `
	{{define "content"}}
		<h1>Register</h1>
		<form action="/admin/register" method="POST">
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
			<div class="{{if index .ValidationErrors "Password"}}error{{end}}">
				<label for="body">Password</label>
				<input type="password" id="password" name="password" value="{{.User.Password}}">
				<div>{{if index .ValidationErrors "Password"}}{{index .ValidationErrors "Password"}}{{end}}</div>
			</div>
			<div class="{{if index .ValidationErrors "ConfirmPassword"}}error{{end}}">
				<label for="confirm_password">Confirm Password</label>
				<input type="password" id="confirm_password" name="confirm_password" value="{{.User.ConfirmPassword}}">
				<div>{{if index .ValidationErrors "ConfirmPassword"}}{{index .ValidationErrors "ConfirmPassword"}}{{end}}</div>
			</div>
			<button type="submit">Register</button>
		</form>
	{{end}}
	`

	tmpl := template.Must(template.New("layout").Parse(templates.AdminLayoutTemplate))
	tmpl = template.Must(tmpl.New("meta").Parse(templates.MetaDataTemplate))
	tmpl = template.Must(tmpl.New("registration").Parse(pageTemplate))

	pageData := UserRegistrationTemplate{
		Title:            "Register your account",
		Description:      "Register your account",
		Keywords:         "resgistration",
		Body:             pageTemplate,
		ValidationErrors: validationErrors,
		User:             user,
		BustCssCache:     app.BustCssCache,
		BustJsCache:      app.BustJsCache,
	}

	err := tmpl.ExecuteTemplate(w, "layout", pageData)
	if err != nil {
		http.Error(w, fmt.Sprintf("Template execution failed: %v", err), http.StatusInternalServerError)
	}
}
