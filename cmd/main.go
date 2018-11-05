//2057008, 2624395, 9111696

package main

import (
	"log"
	"net/http"
)

func main() {

	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:])
	})

	http.Handle("/", http.FileServer(http.Dir("./pkg/frontend")))

	err := http.ListenAndServeTLS(":443", "server.crt", "server.key", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
