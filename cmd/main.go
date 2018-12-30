//2057008, 2624395, 9111696

package main

import (
	auth "go-ticketsystem/pkg/authentication"
	hand "go-ticketsystem/pkg/frontend"
	"log"
	"net/http"
)

func mainHandler(w http.ResponseWriter, r *http.Request) {
	http.FileServer(http.Dir("./pkg/frontend")).ServeHTTP(w, r)
}

func main() {
	/*
		http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, r.URL.Path[1:])
		})
	*/

	http.HandleFunc("/secure/dashboard.html", hand.WrapperDashboard(mainHandler))
	http.HandleFunc("/secure/ticketDetail.html", hand.WrapperTicketDet(mainHandler))
	http.HandleFunc("/secure/tickets.html", hand.WrapperTickets(mainHandler))
	http.HandleFunc("/secure/ticketsOpen.html", hand.WrapperOpenTickets(mainHandler))
	http.HandleFunc("/secure/ticketsProcessing.html", hand.WrapperProTickets(mainHandler))
	http.HandleFunc("/secure/ticketsClosed.html", hand.WrapperClosedTickets(mainHandler))
	http.HandleFunc("/", auth.Wrapper(mainHandler))

	err := http.ListenAndServeTLS(":443", "Server.crt", "Server.key", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
