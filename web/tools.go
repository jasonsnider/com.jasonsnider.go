package web

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jasonsnider/go.jasonsnider.com/templates"
)

type ToolPageData struct {
	Title        string
	Description  string
	Keywords     string
	Body         template.HTML
	BustCssCache string
	BustJsCache  string
}

func (app *App) ListTools(w http.ResponseWriter, r *http.Request) {
	articles, err := app.FetchArticlesByType("tool")

	if err != nil {
		http.Error(w, fmt.Sprintf("FetchArticlesByType failed: %v", err), http.StatusInternalServerError)
		return
	}

	articlesTemplate := `
        {{define "content"}}
            <h1>Tools</h1>
            <div>
                {{range .Articles}}
                    <h2><a href="/tools/{{.Slug}}">{{.Title}}</a></h2>
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

func (app *App) ViewTool(w http.ResponseWriter, r *http.Request) {
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

	// Define the tools
	tools := map[string]string{
		"hash": `
            <div>
                <label for="InputCount">Iterations</label><input id="InputCount" type="text" value="1">
                <label for="InputString">Enter a string to return a list of hash values</label>
                <textarea id="InputString" rows="5"></textarea>
            </div>
            <div>
                <small>This utility uses only front end JavaScript no data is sent to the server.</small>
            </div>
            <div id="Hashes"></div>
            <script src="dist/js/tools/hash.min.js"></script>`,
		"strlen": `
			<div>
				<label for="InputString">Enter a string to calculate it's length
					<span class="hide" id="Results">&nbsp;(<strong id="StringLength"></strong>)</span>
				</label>
				<textarea id="InputString" rows="5" spellcheck="false"></textarea></div>
			<div>
				<small>This utility uses only front end JavaScript no data is sent to the server.</small>
			</div>
			<script src="/dist/js/tools/strlen.js"></script>`,
	}

	// Select the tool based on the article slug
	selectedTool := tools[article.Slug]

	// Combine the article body with the selected tool
	body := mdToHTML(article.Body) + template.HTML(selectedTool)

	articleTemplate := `
	{{define "content"}}
		<h1>{{.Title}}</h1>
		<div>
			{{.Body}}
		</div>
	{{end}}
`

	tmpl := template.Must(template.New("layout").Funcs(funcMap).Parse(templates.MainLayoutTemplate))
	tmpl = template.Must(tmpl.New("meta").Parse(templates.MetaDataTemplate))
	tmpl = template.Must(tmpl.New("article").Parse(articleTemplate))

	pageData := ToolPageData{
		Title:        article.Title,
		Description:  article.Description,
		Keywords:     article.Keywords,
		Body:         body,
		BustCssCache: app.BustCssCache,
		BustJsCache:  app.BustJsCache,
	}

	err = tmpl.ExecuteTemplate(w, "layout", pageData)
	if err != nil {
		http.Error(w, fmt.Sprintf("Template execution failed: %v", err), http.StatusInternalServerError)
	}
}
