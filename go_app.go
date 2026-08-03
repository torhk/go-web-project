package main

import (
	"html/template"
	"log"
	"os"
	"net/http"
	"regexp"
	"fmt"
	"strings"
	"context"
	"go_serv/db"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Page struct {
	Title string
	Body  []byte
}
type Files struct {
	Title string
	Type string
	Name  []string
}

type Servers struct {
	Pg_db *db.Queries
}

func (p *Page) save() error {
	filename := "data/" + p.Title + ".txt"
	return os.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := "data/" + title + ".txt"
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	p := &Page{Title: "Login", Body: nil}
	err := templates.ExecuteTemplate(w, "login.html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}


func (s *Servers) loginPostHandler(w http.ResponseWriter, r *http.Request) {
	//TODO: Check for existing login token
	user := r.PostFormValue("uname")
	pass := r.PostFormValue("psw")
	db_user, err := s.Pg_db.GetUserFromName(r.Context(), user)
	if err != nil {
		log.Println("User not found: ", user)
	}
	db_hash, err := s.Pg_db.GetHash(r.Context(), db_user.ID)
	if err != nil {
		log.Println("Password not found: ", user)
	}
	if db_hash != pass{
		log.Println("Wrong pass or user: ", user)
		w.WriteHeader(http.StatusUnauthorized)
		loginHandler(w, r)
		return
	}
	tokenString := createJWT(db_user.ID.String())
	cookie := &http.Cookie{
		Name:     "user_token",
		Value:    tokenString,
		Path:     "/",
		MaxAge:   3600, // Expires in 1 hour (in seconds)
		HttpOnly: true, // Prevents JavaScript access (security best practice)
		Secure:   false, // Ensures cookie is only sent over HTTPS
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/", http.StatusFound)
}


func homeHandler(w http.ResponseWriter, r *http.Request) {
	dir := "./data"

	// ReadDir returns a slice of DirEntry
	folder, err := os.ReadDir(dir)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("homeHandler")
	files := Files{
		Title: "Files",
		Type:   "txt",
		Name: []string{},
	}
	// Loop through entries
	for _, file := range folder {
		fmt.Println(file.Name())
		files.Name = append(files.Name, strings.TrimSuffix(file.Name(), ".txt"))
	}
//	http.Redirect(w, r, "/view/FrontPage", http.StatusFound)	
	//files := []string{"apple", "banana", "cherry"}
	err = templates.ExecuteTemplate(w, "home.html", files)
	if err != nil {
		fmt.Println("homeHandler: nil")
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

func main() {
	//Static file server
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	//DB setup
	// ctx := context.Background()
	// conn, err := pgx.Connect(ctx, "user=pqgotest dbname=pqgotest sslmode=verify-full")
	// if err != nil {
	// 	return err
	// }
	// defer conn.Close(ctx)
	// queries := t.New(conn)

	//2
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbUrl := os.Getenv("DB_URL")
	connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
							dbUser, dbPass, dbUrl, dbName)
	pool, err := pgxpool.New(context.Background(), connStr)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
        os.Exit(1)
    }
    defer pool.Close()

    user_db := &Servers{
		Pg_db: db.New(pool),
	}
	err = pool.Ping(context.Background())
	if err != nil {
		log.Fatalf("Database connection failed! Could not ping: %v", err)
	}
	
	fmt.Println("Ping successful! Database is online and reachable.")
	user_db.Pg_db.GetUserFromName(context.Background(),"thk")

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

	log.Println("Starting server...")
	log.Fatal(http.ListenAndServe(":8090", mux))
}

