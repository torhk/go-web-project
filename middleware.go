package main
import (
	"log"
	"net/http"
	"go_serv/db"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

func Chain(mw ...Middleware) Middleware {
	return func(finalHandler http.HandlerFunc) http.HandlerFunc {
		for i := len(mw) - 1; i >= 0; i-- {
			finalHandler = mw[i](finalHandler)
		}
		return finalHandler
	}
}

func Logger(next http.HandlerFunc) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		log.Printf("Started %s %s", r.Method, r.URL.Path)
		next(w,r)
		log.Printf("Finished")
	}
}


func Auth(Pg_db *db.Queries) Middleware{
	return func(next http.HandlerFunc) http.HandlerFunc{
		return func(w http.ResponseWriter, r *http.Request){
			log.Printf("Auth")
			cookie, err := r.Cookie("user_token")
			if err != nil {
				http.Error(w, "Unauthrorized: missing cookie",http.StatusUnauthorized)
				return
			}
			if !verifyJWT(cookie.Value){
				//TODO: Redirect to login page?
				http.Redirect(w, r, "/", http.StatusFound)
				return
			}
			next(w,r)
		}
	}
}

