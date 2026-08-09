package routes

import (
	"net/http"
	"regexp"
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
	
	
	mux.HandleFunc("GET /wiki/{entry}/view", applyMiddleware(makeHandler(viewHandler)))
	mux.HandleFunc("GET /wiki/{entry}/edit", applyMiddleware(makeHandler(editHandler)))
	mux.HandleFunc("POST /wiki/{entry}/save", applyMiddleware(makeHandler(saveHandler)))

	return mux
}

var validPath = regexp.MustCompile("^[a-zA-Z0-9]+$")

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 
		title := r.PathValue("entry")
		
		if !validPath.MatchString(title){
			http.NotFound(w, r)
			return
		}
		fn(w, r, title)
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