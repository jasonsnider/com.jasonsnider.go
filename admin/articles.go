package admin

import (
	"context"
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gomarkdown/markdown"
	"github.com/gorilla/mux"
	"github.com/jasonsnider/com.jasonsnider.go/internal/db"
	"github.com/jasonsnider/com.jasonsnider.go/internal/types"
	"github.com/jasonsnider/com.jasonsnider.go/pkg/inflection"
	"github.com/jasonsnider/com.jasonsnider.go/templates"
)

type ArticlesPageData struct {
	Title        string
	Description  sql.NullString
	Keywords     sql.NullString
	Articles     []types.Article
	BustCssCache string
	BustJsCache  string
}

type ArticlePageData struct {
	ID           string
	Title        string
	Description  sql.NullString
	Keywords     sql.NullString
	Body         string
	BustCssCache string
	BustJsCache  string
}

type ArticleUpdateTemplate struct {
	Title            string
	Description      sql.NullString
	Keywords         sql.NullString
	Body             string
	ValidationErrors map[string]string
	Article          types.Article
	BustCssCache     string
	BustJsCache      string
}

func mdToHTML(md string) template.HTML {
	return template.HTML(markdown.ToHTML([]byte(md), nil, nil))
}

func parseTime(timeStr string) *time.Time {
	if timeStr == "" {
		return nil
	}
	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		log.Printf("failed to parse time: %v", err)
		return nil
	}
	return &t
}

func (app *App) CreateArticle(w http.ResponseWriter, r *http.Request) {
	db := db.DB{DB: app.DB}
	article := types.Article{}
	validationErrors := make(map[string]string)

	if r.Method == "POST" {
		validate := validator.New()
		validate.RegisterValidation("uniqueEmail", db.UniqueEmail)

		article = types.Article{
			Title: r.FormValue("title"),
		}

		err := validate.Struct(article)

		if err != nil {
			for _, err := range err.(validator.ValidationErrors) {
				fieldName := err.Field()
				fieldNameHuman := inflection.Humanize(fieldName)
				tag := err.Tag()

				var errorMessage string
				switch tag {
				case "required":
					errorMessage = fmt.Sprintf("%s is required", fieldNameHuman)
				default:
					errorMessage = fmt.Sprintf("%s is invalid", fieldNameHuman)
				}

				validationErrors[fieldName] = errorMessage
			}
		} else {
			articleID, err := db.CreateArticle(article)
			if err != nil {
				log.Fatalf("failed to create article: %v", err)
			}

			log.Println("Article created successfully")
			http.Redirect(w, r, "/admin/articles/"+articleID+"/edit", http.StatusSeeOther)
		}
	}

	pageTemplate := `
	{{define "content"}}
		<header class="row">
			<h1 class="col">Create an Article</h1>
			<div class="col-end">
				<a class="btn" href="/admin/articles">Articles</a>
			</div>
		</header>

		<form action="/admin/articles/create" method="POST" novalidate>
			<div class="{{if index .ValidationErrors "Title"}}error{{end}}">
				<label for="title">Article</label>
				<input type="text" id="Title" name="title" value="{{.Article.Title}}">
				<div>{{if index .ValidationErrors "Title"}}{{index .ValidationErrors "Title"}}{{end}}</div>
			</div>
			<button type="submit">Submit</button>
		</form>
	{{end}}
	`

	tmpl := template.Must(template.New("layout").Parse(templates.AdminLayoutTemplate))
	tmpl = template.Must(tmpl.New("create_user").Parse(pageTemplate))

	pageData := ArticleUpdateTemplate{
		Title:            "Create a user",
		Body:             pageTemplate,
		ValidationErrors: validationErrors,
		Article:          article,
		BustCssCache:     app.BustCssCache,
		BustJsCache:      app.BustJsCache,
	}

	err := tmpl.ExecuteTemplate(w, "layout", pageData)
	if err != nil {
		http.Error(w, fmt.Sprintf("Template execution failed: %v", err), http.StatusInternalServerError)
	}
}

func (app *App) ListArticles(w http.ResponseWriter, r *http.Request) {

	db := db.DB{DB: app.DB}
	articles, err := db.FetchArticles()

	if err != nil {
		http.Error(w, fmt.Sprintf("FetchArticlesByType failed: %v", err), http.StatusInternalServerError)
		return
	}

	articlesTemplate := `
        {{define "content"}}
            
			<header class="row">
				<h1 class="col">Articles</h1>
				<div class="col-end">
					<a class="btn" href="/admin/articles/create">Create</a>
				</div>
			</header>

			{{range .Articles}}
				<div class="row rotate">
					<div class="col-4"><a href="/admin/articles/{{.ID}}">{{.Title}}</a></div>
					<div class="col">{{safeValue .Type}}</div>
					<div class="col">{{safeValue .Format}}</div>
					<div class="col-end">
						<a href="/admin/articles/{{.ID}}"><i class="fas fa-eye"></i></a>
						<a href="/admin/articles/{{.ID}}/edit"><i class="fas fa-edit"></i></a>
						<a href="/admin/articles/{{.ID}}/delete"><i class="fas fa-trash"></i></a>
					</div>
				</div>
			{{end}}
        {{end}}
    `
	funcMap := template.FuncMap{
		"safeValue": types.SafeValue,
	}
	tmpl := template.Must(template.New("layout").Funcs(funcMap).Parse(templates.AdminLayoutTemplate))
	tmpl = template.Must(tmpl.New("content").Parse(articlesTemplate))

	pageData := ArticlesPageData{
		Title:        "Articles",
		Description:  types.TypeSqlNullString("A list of articles"),
		Keywords:     types.TypeSqlNullString("articles, blog"),
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
			<header class="row">
				<h1 class="col">{{.Title}}</h1>
				<div class="col-end">
					<a class="btn" href="/admin/articles/{{.ID}}/edit">Edit</a>
					<a class="btn" href="/admin/articles/{{.ID}}/delete">Delete</a>
				</div>
			</header>
			<div>
				{{mdToHTML .Body}}
			</div>
		{{end}}
    `

	tmpl := template.Must(template.New("layout").Funcs(funcMap).Parse(templates.AdminLayoutTemplate))
	tmpl = template.Must(tmpl.New("article").Parse(articleTemplate))

	pageData := ArticlePageData{
		ID:           article.ID,
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

		publishedTime, _ := types.ParseSqlNullTime(r.FormValue("published"))

		article.ID = r.FormValue("id")
		article.Title = r.FormValue("title")
		article.Slug = r.FormValue("slug")
		article.Description = types.TypeSqlNullString(r.FormValue("description"))
		article.Body = types.TypeSqlNullString(r.FormValue("body"))
		article.Keywords = types.TypeSqlNullString(r.FormValue("keywords"))
		article.Type = types.TypeSqlNullString(r.FormValue("type"))
		article.Format = types.TypeSqlNullString(r.FormValue("format"))
		article.Published = publishedTime
		// if published := parseTime(r.FormValue("published")); published != nil {
		// 	article.Published = *published
		// }

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
				SET title = $1, description = $2, keywords = $3, body = $4, type = $5, format = $6, published = $7
				WHERE id = $8
			`

			tx, err := app.DB.Begin(context.Background())
			if err != nil {
				log.Fatalf("begin transaction failed: %v", err)
			}

			_, err = tx.Exec(context.Background(), query, article.Title, article.Description, article.Keywords, article.Body, article.Type, article.Format, article.Published, article.ID)
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
		<header class="row">
			<h1 class="col">{{.Title}}</h1>
			<div class="col-end">
				<a class="btn" href="/admin/articles/{{.Article.ID}}">View</a>
				<a class="btn" href="/admin/articles/{{.Article.ID}}/delete">Delete</a>
			</div>
		</header>
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
				<textarea id="Body" name="body" rows="40">{{safeValue .Article.Body}}</textarea>
			</div>
			<div>
				<label for="description">Description</label>
				<textarea id="Description" name="description">{{safeValue .Article.Description}}</textarea>
			</div>
			<div>
				<label for="keywords">Keywords</label>
				<textarea id="Keywords" name="keywords">{{safeValue .Article.Keywords}}</textarea>
			</div>
			<div>
				<label for="type">Type</label>
				<input type="text" id="Type" name="type" value="{{safeValue .Article.Type}}">
			</div>
			<div>
				<label for="format">Format</label>
				<input type="text" id="Format" name="format" value="{{safeValue .Article.Format}}">
			</div>
			<div>
				<label for="published">Published</label>
				<input type="text" id="Published" name="published" value="{{safeValue .Article.Published}}">
			</div>
			<button type="submit">Submit</button>
		</form>
	{{end}}
	`
	funcMap := template.FuncMap{
		"safeValue": types.SafeValue,
	}

	tmpl := template.Must(template.New("layout").Funcs(funcMap).Parse(templates.AdminLayoutTemplate))
	tmpl = template.Must(tmpl.New("article").Parse(pageTemplate))

	pageData := ArticleUpdateTemplate{
		Title:            "Update Article",
		Description:      types.TypeSqlNullString("Register your account"),
		Keywords:         types.TypeSqlNullString("resgistration"),
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
	db := db.DB{DB: app.DB}
	vars := mux.Vars(r)
	id := vars["id"]

	err := db.DeleteArticle(id)

	if err != nil {
		http.Error(w, fmt.Sprintf("DeleteArticleByID failed: %v", err), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/articles", http.StatusSeeOther)
}
