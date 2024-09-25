package web

import (
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/jasonsnider/com.jasonsnider.go/internal/types"
	"github.com/jasonsnider/com.jasonsnider.go/templates"
	"github.com/mailgun/mailgun-go"
)

type Contact struct {
	Subject string `json:"subject"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Body    string `json:"body"`
}

func SendSimpleMessage(contact Contact) (string, error) {

	mailgunDomain := os.Getenv("MAILGUN_DOMAIN")
	mailgunApiKey := os.Getenv("MAILGUN_API_KEY")
	supportEmail := os.Getenv("SUPPORT_EMAIL")

	fmt.Println("DOMAIN:", mailgunDomain)
	fmt.Println("API KEY:", mailgunApiKey)
	fmt.Println("SUPPORT EMAIL:", supportEmail)

	mg := mailgun.NewMailgun(mailgunDomain, mailgunApiKey)
	m := mg.NewMessage(
		contact.Email,
		contact.Subject,
		contact.Body,
		supportEmail,
	)
	_, id, err := mg.Send(m)
	return id, err
}

func (app *App) Contact(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		contact := Contact{
			Subject: r.FormValue("_subject"),
			Name:    r.FormValue("name"),
			Email:   r.FormValue("_replyto"),
			Body:    r.FormValue("body"),
		}
		SendSimpleMessage(contact)
	}

	contactTemplate := `
	{{define "content"}}
		<h1>Contact</h1>
		<form action="/contact" method="POST">
			<input type="hidden" name="_next" value="https://jasonsnider.com/thanks">
			<div>
				<label for="subject">Subject</label>
				<select id="subject" name="_subject">
					<option value="CONTACT: jasonsnider.com">General Contact</option>
					<option value="SUPPORT: jasonsnider.com">Support</option>
				</select>
			</div>
			<div>
				<label for="subject">Name</label><input type="text" id="name" name="name">
			</div>
			<div>
				<label for="subject">Email</label>
				<input type="email" id="email" name="_replyto">
			</div>
			<div>
				<label for="body">Body</label>
				<textarea id="body" name="body" rows="5" spellcheck="false"></textarea>
			</div>
			<div>
				<input type="submit">
			</div>
		</form>
	{{end}}
    `
	tmpl := template.Must(template.New("layout").Parse(templates.MainLayoutTemplate))
	tmpl = template.Must(tmpl.New("meta").Parse(templates.MetaDataTemplate))
	tmpl = template.Must(tmpl.New("article").Parse(contactTemplate))

	pageData := ArticlePageData{
		Title:        "Contact",
		Description:  types.TypeSqlNullString("Contact Jason Snider"),
		Keywords:     types.TypeSqlNullString("contact, email"),
		Body:         contactTemplate,
		BustCssCache: app.BustCssCache,
		BustJsCache:  app.BustJsCache,
	}

	err := tmpl.ExecuteTemplate(w, "layout", pageData)
	if err != nil {
		http.Error(w, fmt.Sprintf("Template execution failed: %v", err), http.StatusInternalServerError)
	}
}
