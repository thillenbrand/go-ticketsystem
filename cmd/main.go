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

	http.Handle("/", http.FileServer(http.Dir("./pkg/frontend")))
	http.HandleFunc("/login", auth.LoginHandler)
	//http.HandleFunc("/ticket", auth.Wrapper(auth.LoginHandler)) todo: Wrapper läuft noch nicht

	err := http.ListenAndServeTLS(":443", "Server.crt", "Server.key", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
