package admin

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/jasonsnider/com.jasonsnider.go/templates"
)

func (app *App) Dashboard(w http.ResponseWriter, r *http.Request) {

	user := RegisterUser{}
	validationErrors := make(map[string]string)

	pageTemplate := `
	{{define "content"}}
		<h1>Dashboard</h1>
	{{end}}
	`

	tmpl := template.Must(template.New("layout").Parse(templates.MainLayoutTemplate))
	tmpl = template.Must(tmpl.New("meta").Parse(templates.MetaDataTemplate))
	tmpl = template.Must(tmpl.New("dashboard").Parse(pageTemplate))

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
