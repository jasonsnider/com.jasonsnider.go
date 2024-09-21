package web

import (
	"context"
	"fmt"
	"html/template"
	"net/http"

	"github.com/gomarkdown/markdown"
	"github.com/gorilla/mux"
	"github.com/jasonsnider/go.jasonsnider.com/templates"
)

type Article struct {
	ID          string `json:"id"`
	Slug        string `json:"slug"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Keywords    string `json:"keywords"`
	Body        string `json:"body"`
}

type ArticlesPageData struct {
	Title        string
	Description  string
	Keywords     string
	Articles     []Article
	BustCssCache string
	BustJsCache  string
}

type ArticlePageData struct {
	Title        string
	Description  string
	Keywords     string
	Body         string
	BustCssCache string
	BustJsCache  string
}

func mdToHTML(md string) template.HTML {
	return template.HTML(markdown.ToHTML([]byte(md), nil, nil))
}

func (app *App) FetchArticlesByType(articleType string) ([]Article, error) {
	rows, err := app.DB.Query(context.Background(), "SELECT id, slug, title, description, keywords, body FROM articles WHERE type=$1", articleType)
	if err != nil {
		return nil, fmt.Errorf("query failed: %v", err)
	}
	defer rows.Close()

	var articles []Article
	for rows.Next() {
		var article Article
		err := rows.Scan(&article.ID, &article.Slug, &article.Title, &article.Description, &article.Keywords, &article.Body)
		if err != nil {
			return nil, fmt.Errorf("row scan failed: %v", err)
		}
		articles = append(articles, article)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("rows iteration failed: %v", rows.Err())
	}

	return articles, nil
}

func (app *App) FetchArticleBySlug(slug string) (Article, error) {
	var article Article
	err := app.DB.QueryRow(context.Background(), "SELECT id, slug, title, description, keywords, body FROM articles WHERE slug=$1", slug).Scan(&article.ID, &article.Slug, &article.Title, &article.Description, &article.Keywords, &article.Body)
	if err != nil {
		return article, fmt.Errorf("query failed: %v", err)
	}

	return article, nil
}

func (app *App) ListArticles(w http.ResponseWriter, r *http.Request) {

	articles, err := app.FetchArticlesByType("post")

	if err != nil {
		http.Error(w, fmt.Sprintf("FetchArticlesByType failed: %v", err), http.StatusInternalServerError)
		return
	}

	articlesTemplate := `
        {{define "content"}}
            <h1>Articles</h1>
            <div>
                {{range .Articles}}
                    <h2><a href="/articles/{{.Slug}}">{{.Title}}</a></h2>
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

func (app *App) ViewArticle(w http.ResponseWriter, r *http.Request) {
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
