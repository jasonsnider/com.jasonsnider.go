package admin

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/jasonsnider/com.jasonsnider.go/templates"
)

func (app *App) Dashboard(w http.ResponseWriter, r *http.Request) {

	pageTemplate := `
	{{define "content"}}
		<h1>Dashboard</h1>
		<div>
			<a href="/admin/articles">Articles</a>&nbsp;|&nbsp; 
			<a href="/admin/users">Users</a>
		</div>
	{{end}}
	`

	tmpl := template.Must(template.New("layout").Parse(templates.AdminLayoutTemplate))
	tmpl = template.Must(tmpl.New("dashboard").Parse(pageTemplate))

	pageData := ArticlePageData{
		Title:        "Dashboard",
		Body:         "",
		BustCssCache: app.BustCssCache,
		BustJsCache:  app.BustJsCache,
	}

	err := tmpl.ExecuteTemplate(w, "layout", pageData)
	if err != nil {
		http.Error(w, fmt.Sprintf("Template execution failed: %v", err), http.StatusInternalServerError)
	}
}
