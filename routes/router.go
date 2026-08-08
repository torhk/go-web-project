package routes

import (
	"net/http"
	"regexp"
	"fmt"
	"html/template"
)

func AddRoutes(user_db *Servers) *http.ServeMux {
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	//Middleware pipelines
	applyMiddleware := Chain(Logger, Auth(user_db.Pg_db))
	openPageMiddleware := Chain(Logger)

	//Handlers
	mux.HandleFunc("GET /{$}", openPageMiddleware(homeHandler))
	mux.HandleFunc("GET /login", openPageMiddleware(loginHandler))
	mux.HandleFunc("POST /login", openPageMiddleware(user_db.loginPostHandler))
	
	
	mux.HandleFunc("GET /view/", applyMiddleware(makeHandler(viewHandler)))
	mux.HandleFunc("GET /edit/", applyMiddleware(makeHandler(editHandler)))
	mux.HandleFunc("POST /save/", applyMiddleware(makeHandler(saveHandler)))

	return mux
}

var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fmt.Println("makeHandler")
		fn(w, r, m[2])
	}
}

//var templates = template.Must(template.ParseFiles("edit.html", "view.html"))
var templates = template.Must(template.ParseGlob("tmpl/*.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}