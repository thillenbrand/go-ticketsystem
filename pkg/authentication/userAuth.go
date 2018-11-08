package authentication

import (
	"fmt"
	"net/http"
)

// todo: auth-Funktionen einführen
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	name := q.Get("name")
	if name == "" {
		name = "World"
	}
	responseString := "<html><body>Hello " + name + "</body></html>"
	w.Write([]byte(responseString))

	fmt.Println("test")
}

// todo: Wrapper
/*
func Wrapper (handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fileServer := http.StripPrefix("/tickets", http.FileServer(http.Dir("./frontend")))
		fileServer.ServeHTTP(w, r)
	}
}
*/
