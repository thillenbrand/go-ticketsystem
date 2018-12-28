//2057008, 2624395, 9111696

package main

import (
	"go-ticketsystem/pkg/authentication"
	"log"
	"net/http"
)

func mainHandler(w http.ResponseWriter, r *http.Request) {
	http.FileServer(http.Dir("./pkg/frontend")).ServeHTTP(w, r)
}

func main() {

	http.HandleFunc("/", authentication.Wrapper(mainHandler))

	err := http.ListenAndServeTLS(":443", "Server.crt", "Server.key", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
