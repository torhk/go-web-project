package routes

import (
	"log"
	"time"
	"net/http"
	"go_serv/db"
	"go_serv/auth"
)

type Servers struct {
	Pg_db *db.Queries
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	p := &Page{Title: "Login", Body: nil}
	renderTemplate(w, "login", p)
	// TODO: Redirect if alredy authenticated.
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:     "user_token", // Must match the original cookie name
		Value:    "",              // Clear the value
		Path:     "/",             // Must match the original path
		MaxAge:   -1,              // Tells the browser to delete the cookie immediately
		Expires:  time.Unix(0, 0), // A date in the past (Jan 1, 1970) for older browsers
		HttpOnly: true,            // Security best practice
		Secure:   true,            // Secure flag (if using HTTPS)
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
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
	tokenString := auth.CreateJWT(db_user.ID.String())
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