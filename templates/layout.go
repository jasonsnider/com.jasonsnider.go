package templates

import "database/sql"

type MetaData struct {
	Title       string
	Description sql.NullString
	Keywords    string
}

const MetaDataTemplate = `
{{define "meta"}}
    <title>{{.Title}}</title>
    <meta name="description" content="{{.Description}}">
	<meta name="keywords" content="{{.Keywords}}">
{{end}}
`

const MainLayoutTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    {{template "meta" .}}
	<base href="/">
	<link rel="stylesheet" href="/fonts/fonts.css?{{.BustCssCache}}">
	<link rel="stylesheet" href="/fontawesome/css/all.min.css">
	<link rel="stylesheet" href="/dist/css/main.min.css?{{.BustCssCache}}">
	<link rel="stylesheet" href="/highlight/styles/atom-one-dark.css">
	<meta name="viewport" content="width=device-width, initial-scale=1">
</head>
<body>
	<header>
		<nav id="MainNav">
			<a class="logo" href="/"><i class="fa fa-code"></i>&nbsp;Jason Snider</a>
			<button id="ShowMainNav"><i class="fas fa-bars fa-2x"></i><span class="sr-only">Menu</span></button>
			<ul>
				<li><a id="ToAbout" href="/#About">About Me</a></li>
				<li><a href="/articles">Blog</a></li>
				<li><a href="/games">Games</a></li>
			<li><a href="/tools">Tools</a></li>
			<li><a href="/contact">Contact</a></li>
			</ul>
		</nav>
	</header>
	<main>
		{{template "content" .}}
		<footer>
			<div class="left">Built with<i class="fa fa-heart" arial-hidden="true"></i><span class="sr-only">love</span>by Jason in Chicago</div>
			<div class="right"><a href="terms">Terms</a><a href="privacy">Privacy</a></div>
		</footer>
	</main>
	<script>
		var loadDeferredStyles = function() {
		var addStylesNode = document.getElementById("deferred-styles");
		var replacement = document.createElement("div");
		replacement.innerHTML = addStylesNode.textContent;
		document.body.appendChild(replacement)
		addStylesNode.parentElement.removeChild(addStylesNode);
		};
		var raf = window.requestAnimationFrame || window.mozRequestAnimationFrame ||
			window.webkitRequestAnimationFrame || window.msRequestAnimationFrame;
		if (raf) raf(function() { window.setTimeout(loadDeferredStyles, 0); });
		else window.addEventListener('load', loadDeferredStyles);
	</script>
	<script src="/highlight/highlight.pack.js" async=""></script>
	<script src="/dist/js/article.min.js?{{.BustJsCache}}" async=""></script>
</body>
</html>
`

const HomePageTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    {{template "meta" .}}
	<base href="/">
	<link rel="stylesheet" href="/fonts/fonts.css?{{.BustCssCache}}">
	<link rel="stylesheet" href="/fontawesome/css/all.min.css">
	<link rel="stylesheet" href="/dist/css/home-page.min.css?{{.BustCssCache}}">
	<meta name="viewport" content="width=device-width, initial-scale=1">
</head>
<body>
	<header class="cover">
		<nav id="MainNav">
			<a class="logo" href="/"><i class="fa fa-code"></i>&nbsp;Jason Snider</a>
			<button id="ShowMainNav"><i class="fas fa-bars fa-2x"></i><span class="sr-only">Menu</span></button>
			<ul>
				<li><a id="ToAbout" href="/#About">About Me</a></li>
				<li><a href="/articles">Blog</a></li>
				<li><a href="/games">Games</a></li>
			<li><a href="/tools">Tools</a></li>
			<li><a href="/contact">Contact</a></li>
			</ul>
		</nav>
		<h1>
			&nbsp;Jason Snider
			<div class="tagline">&nbsp;Builder of Things, Doer of Stuff</div>
		</h1>
		<div class="more">
			<a id="ToAboutNext" href="#About"><i class="fal fa-chevron-down"></i></a>
		</div>
	</header>
	<main>
		{{template "content" .}}
		<footer>
			<div class="left">Built with<i class="fa fa-heart" arial-hidden="true"></i><span class="sr-only">love</span>by Jason in Chicago</div>
			<div class="right"><a href="terms">Terms</a><a href="privacy">Privacy</a></div>
		</footer>
	</main>
	<script>
		var loadDeferredStyles = function() {
		var addStylesNode = document.getElementById("deferred-styles");
		var replacement = document.createElement("div");
		replacement.innerHTML = addStylesNode.textContent;
		document.body.appendChild(replacement)
		addStylesNode.parentElement.removeChild(addStylesNode);
		};
		var raf = window.requestAnimationFrame || window.mozRequestAnimationFrame ||
			window.webkitRequestAnimationFrame || window.msRequestAnimationFrame;
		if (raf) raf(function() { window.setTimeout(loadDeferredStyles, 0); });
		else window.addEventListener('load', loadDeferredStyles);
	</script>
	<script src="/dist/js/home-page.min.js?{{.BustJsCache}}"></script>
</body>
</html>
`

const AdminLayoutTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <title>{{.Title}}</title>
	<base href="/">
	<link rel="stylesheet" href="/fonts/fonts.css?{{.BustCssCache}}">
	<link rel="stylesheet" href="/fontawesome/css/all.min.css">
	<link rel="stylesheet" href="/dist/css/app.min.css?{{.BustCssCache}}">
	<meta name="viewport" content="width=device-width, initial-scale=1">
</head>
<body>

	<div class="wrapper">

	<aside>
		<ul>
			<li><a href="/admin/dashboard"><i class="fas fa-home"></i></a></li>
			<li><a href="/admin/users"><i class="fas fa-user"></i></a></li>
			<li><a href="/admin/articles"><i class="fas fa-newspaper"></i></a></li>
		</ul>
	</aside>

	<main>
		{{template "content" .}}
	</main>

	</div>

	<footer>
		<div class="container">&copy; 2007-2024 jasonsnider.com</div>
	</footer>
	<script src="/dist/js/app.min.js?{{.BustJsCache}}"></script>
</body>
</html>
`
