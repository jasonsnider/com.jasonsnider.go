package web

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/jasonsnider/com.jasonsnider.go/templates"
)

func (app *App) Home(w http.ResponseWriter, r *http.Request) {

	articleTemplate := `
	{{define "content"}}
		<section id="About">
			<h2>Hello, I'm Jason Snider</h2>
			<p>
				<img class="avatar" alt="64x64" src="/img/b449cfaa4cab48910ad8c0d9fdb812d0-128.jpg" height="128px" width="128px">
				Jason Snider is a full stack web and hybrid mobile application developer, 
				indie game dev, dev bootcamp instructor, systems architect, Linux aficionado,
				open source advocate, impromptu DBA, and security enthusiast with a sincere 
				passion for building and advancing knowledge while resolving complex problems  
				and business challenges through technical innovation.
			</p>
		</section>
		<section id="Projects">
			<h2>Things and Stuff</h2>
			<p>A high level overview of the things I build and the stuff I do.</p>
			
			<h3>Back Office</h3>
			<p>Extensive experience in designing and developing event, geo, and data driven CRM, CMS, ERP, BI, IO, and other miscellaneous back office systems.</p>

			<h3>Frontend Development</h3>
			<p>Extensive experience with a number of front end technologies including but not limited to HTML5, CSS3, JavaScript, jQuery, Angular, Ionic, Rxjs, Gulp, webpack Less, Sass, Material Design, Bootstrap, Vanilla JS, Ajax, API consumption.</p>
			
			<h3>Backend Development</h3>
			<p>Extensive experience with a number of front end technologies including but not limited to Linux, Apache, MySQL, MongoDB, Postgre, PHP, Python, BASH, NodeJS, Express, CakePHP, API development.</p>
			
			<h3>Mobile and Indie Game Development</h3>
			<p>Experience building building progressive web and hybrid mobile applications that are inclusive of indie gamesand general data management. Experience publishing to Google Play and Apple App store.</p>
			
			<h3>Portal Development</h3>
			<p>Experience in designing and developing portal driven monolithic systems. These systems provide a primary personal branding and career coaching platform that provides portals for students, instructors, administration, and white labeling for third party organizations.</p>
			
			<h3>Social Media</h3>
			<p>Experience in designing and developing social media platforms as well as integrating social features into existing platforms for which social aspects are not the primary concern of the system.</p>
			
			<h3>Systems Architect</h3>
			<p>Experience designing and implementing the overall technical architecture for a number of cloud based systems.</p>
			
			<h3>Open Source</h3>
			<p>Has contributed to and released a number of open source projects.</p>
			
			<h3>Security</h3>
			<p>Security enthusiast and hobbyist. In additional to concentrating my graduate work on InfoSec, I have listenednearly every episode of Security Now and strive to design and build secure solutions.</p>
			
			<h3>Instructor, Mentor and Leader</h3>
			<p>Experience designing and teaching development bootcamps, mentoring junior developers, and leading development teams and efforts.</p>
			
			<br><br>
			<div class="text-center"><a class="btn" href="resume">Resume</a></div></section><section id="SocialMedia">
				<a href="https://www.linkedin.com/in/jdsnider"><i class="fab fa-linkedin"><span class="sr-only">LinkedIn</span></i></a>
				<a href="https://github.com/jasonsnider"><i class="fab fa-github"><span class="sr-only">GitHub</span></i></a>
				<a href="https://twitter.com/jason_snider"><i class="fab fa-twitter"><span class="sr-only">Twitter</span></i></a>
		</section>
	{{end}}
    `

	tmpl := template.Must(template.New("layout").Parse(templates.HomePageTemplate))
	tmpl = template.Must(tmpl.New("meta").Parse(templates.MetaDataTemplate))
	tmpl = template.Must(tmpl.New("article").Parse(articleTemplate))

	pageData := ArticlePageData{
		Title:        "Jason Snider",
		Description:  "Jason Snider",
		Keywords:     "Jason Snider",
		Body:         "Jason Snider",
		BustCssCache: app.BustCssCache,
		BustJsCache:  app.BustJsCache,
	}

	err := tmpl.ExecuteTemplate(w, "layout", pageData)
	if err != nil {
		http.Error(w, fmt.Sprintf("Template execution failed: %v", err), http.StatusInternalServerError)
	}
}
