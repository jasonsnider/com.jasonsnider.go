package admin

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gomarkdown/markdown"
	"github.com/gorilla/mux"
	"github.com/jasonsnider/com.jasonsnider.go/internal/db"
	"github.com/jasonsnider/com.jasonsnider.go/internal/types"
	"github.com/jasonsnider/com.jasonsnider.go/templates"
)

type ArticlesPageData struct {
	Title        string
	Description  string
	Keywords     string
	Articles     []types.Article
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

type ArticleUpdateTemplate struct {
	Title            string
	Description      string
	Keywords         string
	Body             string
	ValidationErrors map[string]string
	Article          types.Article
	BustCssCache     string
	BustJsCache      string
}

func mdToHTML(md string) template.HTML {
	return template.HTML(markdown.ToHTML([]byte(md), nil, nil))
}

func (app *App) CreateArticle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Create Article")
}

func (app *App) ListArticles(w http.ResponseWriter, r *http.Request) {

	db := db.DB{DB: app.DB}
	articles, err := db.FetchArticlesByType("post")

	if err != nil {
		http.Error(w, fmt.Sprintf("FetchArticlesByType failed: %v", err), http.StatusInternalServerError)
		return
	}

	articlesTemplate := `
        {{define "content"}}
            <h1>Articles</h1>
			<div class="menu">
				<a href="/admin/articles/create">Create</a>
			</div>
			{{range .Articles}}
				<div class="row rotate">
					<div class="col"><a href="/admin/articles/{{.ID}}">{{.Title}}</a></div>
					<div class="col-end">
						<a href="/admin/articles/{{.ID}}"><i class="fas fa-eye"></i></a>
						<a href="/admin/articles/{{.ID}}/edit"><i class="fas fa-edit"></i></a>
						<a href="/admin/articles/{{.ID}}/delete"><i class="fas fa-trash"></i></a>
					</div>
				</div>
			{{end}}
        {{end}}
    `
	tmpl := template.Must(template.New("layout").Parse(templates.AdminLayoutTemplate))
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
	db := db.DB{DB: app.DB}
	vars := mux.Vars(r)
	id := vars["id"]

	var article types.Article
	article, err := db.FetchArticleByID(id)

	if err != nil {
		http.Error(w, fmt.Sprintf("FetchArticleByID failed: %v", err), http.StatusInternalServerError)
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

	tmpl := template.Must(template.New("layout").Funcs(funcMap).Parse(templates.AdminLayoutTemplate))
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

func (app *App) UpdateArticle(w http.ResponseWriter, r *http.Request) {
	db := db.DB{DB: app.DB}
	vars := mux.Vars(r)
	id := vars["id"]

	var article types.Article
	article, err := db.FetchArticleByID(id)

	if err != nil {
		http.Error(w, fmt.Sprintf("FetchArticleById failed: %v", err), http.StatusInternalServerError)
		return
	}

	validationErrors := make(map[string]string)

	if r.Method == "POST" {
		validate := validator.New()

		article.ID = r.FormValue("id")
		article.Title = r.FormValue("title")
		article.Description = r.FormValue("description")
		article.Body = r.FormValue("body")
		article.Keywords = r.FormValue("keywords")

		err := validate.Struct(article)

		if err != nil {
			for _, err := range err.(validator.ValidationErrors) {
				fieldName := err.Field()
				tag := err.Tag()

				var errorMessage string
				switch tag {
				case "required":
					errorMessage = fmt.Sprintf("%s is required", fieldName)
				default:
					errorMessage = fmt.Sprintf("%s is invalid", fieldName)
				}

				validationErrors[fieldName] = errorMessage
			}
		} else {

			query := `
				UPDATE articles
				SET title = $1, description = $2, body = $3, keywords = $4
				WHERE id = $5
			`

			tx, err := app.DB.Begin(context.Background())
			if err != nil {
				log.Fatalf("begin transaction failed: %v", err)
			}

			_, err = tx.Exec(context.Background(), query, article.Title, article.Description, article.Body, article.Keywords, article.ID)
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
		<form action="/admin/articles/{{.Article.ID}}/edit" method="POST">
			<input type="hidden" name="id" value="{{.Article.ID}}">
			<div class="{{if index .ValidationErrors "Title"}}error{{end}}">
				<label for="title">Article</label>
				<input type="text" id="Title" name="title" value="{{.Article.Title}}">
				<div>{{if index .ValidationErrors "Title"}}{{index .ValidationErrors "Title"}}{{end}}</div>
			</div>
			<div class="{{if index .ValidationErrors "Slug"}}error{{end}}">
				<label for="slug">Slug</label>
				<input type="text" id="Slug" name="slug" value="{{.Article.Slug}}">
				<div>{{if index .ValidationErrors "Slug"}}{{index .ValidationErrors "Slug"}}{{end}}</div>
			</div>
			<div>
				<label for="body">Article</label>
				<textarea id="Body" name="body" rows="40">{{.Article.Body}}</textarea>
			</div>
			<div>
				<label for="description">Description</label>
				<textarea id="Description" name="description">{{.Article.Description}}</textarea>
			</div>
			<div>
				<label for="keywords">Keywords</label>
				<textarea id="Keywords" name="keywords">{{.Article.Keywords}}</textarea>
			</div>
			<button type="submit">Submit</button>
		</form>
	{{end}}
	`

	tmpl := template.Must(template.New("layout").Parse(templates.AdminLayoutTemplate))
	tmpl = template.Must(tmpl.New("meta").Parse(templates.MetaDataTemplate))
	tmpl = template.Must(tmpl.New("article").Parse(pageTemplate))

	pageData := ArticleUpdateTemplate{
		Title:            "Update Article",
		Description:      "Register your account",
		Keywords:         "resgistration",
		Body:             pageTemplate,
		ValidationErrors: validationErrors,
		Article:          article,
		BustCssCache:     app.BustCssCache,
		BustJsCache:      app.BustJsCache,
	}

	err = tmpl.ExecuteTemplate(w, "layout", pageData)
	if err != nil {
		http.Error(w, fmt.Sprintf("Template execution failed: %v", err), http.StatusInternalServerError)
	}
}

func (app *App) DeleteArticle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Delete Article")
}
