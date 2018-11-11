//2057008, 2624395, 9111696

package main

import (
	auth "go-ticketsystem/pkg/authentication"
	"log"
	"net/http"
)

func main() {

	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:])
	})

	http.Handle("/", testWrapper(auth.LoginHandler)) //todo: testWrapper läuft, implementation fehlt - vllt besser über authenticator?
	http.HandleFunc("/login", auth.LoginHandler)

	err := http.ListenAndServeTLS(":443", "Server.crt", "Server.key", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}

func testWrapper(handler http.HandlerFunc) http.Handler {

	return http.FileServer(http.Dir("./pkg/frontend"))

}
