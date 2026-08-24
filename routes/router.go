package routes

import (
	"fmt"
	"regexp"
	"strings"
	"net/http"
	"html/template"
	"path/filepath"
)

var templates = make(map[string]*template.Template)

func init() {
	pages, err := filepath.Glob("tmpl/pages/*.html")
	if err != nil {
		panic(err)
	}
	baseLayout := "tmpl/base.html"

	//Loop and blend each page with the base layout
	for _, page := range pages {
		fullFileName := filepath.Base(page)
		fileName := strings.TrimSuffix(fullFileName, filepath.Ext(fullFileName))
		
		tmpl, err := template.ParseFiles(baseLayout, page)
		if err != nil {
			panic(fmt.Sprintf("failed to parse %s: %v", fileName, err))
		}
		templates[fileName] = tmpl
	}
}

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

func renderTemplate(w http.ResponseWriter, pageName string, data interface{}) {
	tmpl, exists := templates[pageName]
	if !exists {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}

	err := tmpl.ExecuteTemplate(w, "base.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}