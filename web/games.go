package web

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jasonsnider/go.jasonsnider.com/templates"
)

func (app *App) ListGames(w http.ResponseWriter, r *http.Request) {
	articles, err := app.FetchArticlesByType("game")

	if err != nil {
		http.Error(w, fmt.Sprintf("FetchArticlesByType failed: %v", err), http.StatusInternalServerError)
		return
	}

	articlesTemplate := `
        {{define "content"}}
            <h1>Articles</h1>
            <div>
                {{range .Articles}}
                    <h2><a href="/games/{{.Slug}}">{{.Title}}</a></h2>
                    <p>{{.Description}}</p>
                {{end}}
            </div>
        {{end}}
    `
	tmpl := template.Must(template.New("layout").Parse(templates.MainLayoutTemplate))
	tmpl = template.Must(tmpl.New("meta").Parse(templates.MetaDataTemplate))
	tmpl = template.Must(tmpl.New("content").Parse(articlesTemplate))

	pageData := ArticlesPageData{
		Title:        "Articles",
		Description:  "A list of articles",
		Keywords:     "articles, blog",
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
	vars := mux.Vars(r)
	slug := vars["slug"]

	var article Article
	article, err := app.FetchArticleBySlug(slug)

	if err != nil {
		http.Error(w, fmt.Sprintf("FetchArticleBySlug failed: %v", err), http.StatusInternalServerError)
		return
	}

	funcMap := template.FuncMap{
		"mdToHTML": mdToHTML,
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
		Body:         article.Body,
		BustCssCache: app.BustCssCache,
		BustJsCache:  app.BustJsCache,
	}

	err = tmpl.ExecuteTemplate(w, "layout", pageData)
	if err != nil {
		http.Error(w, fmt.Sprintf("Template execution failed: %v", err), http.StatusInternalServerError)
	}
}
