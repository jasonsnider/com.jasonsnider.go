package web

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jasonsnider/com.jasonsnider.go/internal/db"
	"github.com/jasonsnider/com.jasonsnider.go/internal/types"
	"github.com/jasonsnider/com.jasonsnider.go/templates"
)

func (app *App) ListGames(w http.ResponseWriter, r *http.Request) {

	db := db.DB{DB: app.DB}

	meta, err := db.FetchMetaDataBySlug("games")
	if err != nil {
		http.Error(w, fmt.Sprintf("FetchArticlesBySlug failed: %v", err), http.StatusInternalServerError)
		return
	}

	articles, err := db.FetchArticlesByType("game")

	if err != nil {
		http.Error(w, fmt.Sprintf("FetchArticlesByType failed: %v", err), http.StatusInternalServerError)
		return
	}

	funcMap := template.FuncMap{
		"safeValue": types.SafeValue,
	}

	articlesTemplate := `
        {{define "content"}}
            <h1>Games</h1>
            <div>
                {{range .Articles}}
                    <h2><a href="/games/{{.Slug}}">{{.Title}}</a></h2>
                    <p>{{safeValue .Description}}</p>
                {{end}}
            </div>
        {{end}}
    `
	tmpl := template.Must(template.New("layout").Funcs(funcMap).Parse(templates.MainLayoutTemplate))
	tmpl = template.Must(tmpl.New("meta").Funcs(funcMap).Parse(templates.MetaDataTemplate))
	tmpl = template.Must(tmpl.New("content").Parse(articlesTemplate))

	pageData := ArticlesPageData{
		Title:        meta.Title,
		Description:  meta.Description,
		Keywords:     meta.Keywords,
		Articles:     articles,
		BustCssCache: app.BustCssCache,
		BustJsCache:  app.BustJsCache,
	}

	err = tmpl.ExecuteTemplate(w, "layout", pageData)
	if err != nil {
		http.Error(w, fmt.Sprintf("Template execution failed: %v", err), http.StatusInternalServerError)
	}
}

func (app *App) ViewGame(w http.ResponseWriter, r *http.Request) {
	db := db.DB{DB: app.DB}
	vars := mux.Vars(r)
	slug := vars["slug"]

	var article types.Article
	article, err := db.FetchArticleBySlug(slug)

	if err != nil {
		http.Error(w, fmt.Sprintf("FetchArticleBySlug failed: %v", err), http.StatusInternalServerError)
		return
	}

	funcMap := template.FuncMap{
		"mdToHTML":  mdToHTML,
		"safeValue": types.SafeValue,
	}

	articleTemplate := `
		{{define "content"}}
			<h1>{{.Title}}</h1>
			<div>
				{{mdToHTML .Body}}
			</div>
		{{end}}
    `

	tmpl := template.Must(template.New("layout").Funcs(funcMap).Parse(templates.MainLayoutTemplate))
	tmpl = template.Must(tmpl.New("meta").Parse(templates.MetaDataTemplate))
	tmpl = template.Must(tmpl.New("article").Parse(articleTemplate))

	pageData := ArticlePageData{
		Title:        article.Title,
		Description:  article.Description,
		Keywords:     article.Keywords,
		Body:         article.Body.String,
		BustCssCache: app.BustCssCache,
		BustJsCache:  app.BustJsCache,
	}

	err = tmpl.ExecuteTemplate(w, "layout", pageData)
	if err != nil {
		http.Error(w, fmt.Sprintf("Template execution failed: %v", err), http.StatusInternalServerError)
	}
}
