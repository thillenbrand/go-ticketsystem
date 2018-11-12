package authentication

import (
	"fmt"
	"net/http"
)

// todo: auth-Funktionen einführen
func LoginHandler(w http.ResponseWriter, r *http.Request) {

	//todo: nicht über URL.Query sondern über FormValue("Name") abfragen
	q := r.URL.Query()
	username := q.Get("username")
	pass := q.Get("pass")

	if username == "admin" && pass == "PraiseTheOmnissiah" { //todo: mit Daten aus DB/JSON abgleichen
		responseString := "<html><body>Hello " + username + "</body></html>"
		w.Write([]byte(responseString))
		//set logged in
	} else {
		responseString := "<html><body>Wrong username or password!</body></html>"
		w.Write([]byte(responseString))
	}

	fmt.Println("test")

	/*
			name := q.Get("name")
		if name == "" {
			name = "World"
		}
		responseString := "<html><body>Hello " + name + "</body></html>"
	*/
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

func checkUser(user, pass string) bool {
	fmt.Println(user, ",", pass)
	return true
}
