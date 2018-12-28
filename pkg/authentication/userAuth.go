package authentication

import (
	"fmt"
	"net/http"
)

func checkUserValid(user, pswd string) bool {
	if user == "Käse" && pswd == "KäseKäse" {
		return true
	}
	return false

	fmt.Println("user: ", user, ", password: ", pswd)
	return true
}

func Wrapper(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pswd, ok := r.BasicAuth()
		if ok && checkUserValid(user, pswd) {
			handler(w, r)
		} else {
			w.Header().Set("WWW-Authenticate",
				"Basic realm=\"Please enter your username and password\"")
			http.Error(w,
				http.StatusText(http.StatusUnauthorized),
				http.StatusUnauthorized)
		}
	}
}
