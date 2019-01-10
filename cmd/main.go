//2057008, 2624395, 9111696

package main

import (
	auth "go-ticketsystem/pkg/authentication"
	hand "go-ticketsystem/pkg/frontend"
	"log"
	"net/http"
	"os"
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

	//Checkt ob Folders tickets, users und mailQueue existieren, falls nicht werden diese erstellt.
	if _, err := os.Stat("./pkg/tickets"); os.IsNotExist(err) {
		os.Mkdir("./pkg/tickets", 0700)
	}

	if _, err := os.Stat("./pkg/users"); os.IsNotExist(err) {
		os.Mkdir("./pkg/users", 0700)
	}

	if _, err := os.Stat("./pkg/mailQueue"); os.IsNotExist(err) {
		os.Mkdir("./pkg/mailQueue", 0700)
	}

	hand.OpenTickets()

	http.HandleFunc("/", mainHandler)

	http.HandleFunc("/secure/dashboard.html", auth.Wrapper(hand.HandlerDashboard))
	http.HandleFunc("/secure/ticketDetail.html", auth.Wrapper(hand.HandlerTicketDet))
	http.HandleFunc("/secure/tickets.html", auth.Wrapper(hand.HandlerTickets))
	http.HandleFunc("/secure/ticketsOpen.html", auth.Wrapper(hand.HandlerOpenTickets))
	http.HandleFunc("/secure/ticketsProcessing.html", auth.Wrapper(hand.HandlerProTickets))
	http.HandleFunc("/secure/ticketsClosed.html", auth.Wrapper(hand.HandlerClosedTickets))
	http.HandleFunc("/secure/entry.html", auth.Wrapper(hand.HandlerEntry))
	http.HandleFunc("/secure/saveP/", auth.Wrapper(hand.HandlerSaveProfile))
	http.HandleFunc("/secure/save/", auth.Wrapper(hand.HandlerSave))
	http.HandleFunc("/save/", hand.HandlerSave)
	http.HandleFunc("/secure/close/", auth.Wrapper(hand.HandlerClose))
	http.HandleFunc("/secure/release/", auth.Wrapper(hand.HandlerRelease))
	http.HandleFunc("/secure/take/", auth.Wrapper(hand.HandlerTake))
	http.HandleFunc("/secure/add/", auth.Wrapper(hand.HandlerAdd))
	http.HandleFunc("/secure/assign/", auth.Wrapper(hand.HandlerAssign))
	http.HandleFunc("/secure/profile.html", auth.Wrapper(hand.HandlerProfile))
	http.HandleFunc("/register/", auth.Wrapper(auth.HandlerRegister))

	err := http.ListenAndServeTLS(":443", "Server.crt", "Server.key", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
