//2057008, 2624395, 9111696

package main

import (
	"encoding/json"
	"fmt"
	auth "go-ticketsystem/pkg/authentication"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Commentary struct {
	Name string `json:"Name"`
	Text string `json:"Text"`
}

type Ticket struct {
	ID           int          `json:"ID"`
	SDescription string       `json:"SDescription"`
	Description  string       `json:"Description"`
	UName        string       `json:"UName"`
	Email        string       `json:"Email"`
	Commentary   []Commentary `json:"Commentary"`
}

type Tickets struct {
	Tickets []Ticket `json:"ticket"`
}

type Page struct {
	ID           int
	SDescription string
	Description  string
	UName        string
	Email        string
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	http.FileServer(http.Dir("./pkg/frontend")).ServeHTTP(w, r)
}

func main() {

	openTickets()

	http.HandleFunc("/", auth.Wrapper(mainHandler))
	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:])
	})

	http.HandleFunc("/dashboard.html", dashboardHandler)
	http.HandleFunc("/ticketDetail.html", ticketDetailHandler)
	http.HandleFunc("/tickets.html", ticketsHandler)

	err := http.ListenAndServeTLS(":443", "Server.crt", "Server.key", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}

func testWrapper(handler http.HandlerFunc) http.Handler {

	return http.FileServer(http.Dir("./pkg/frontend"))

}

func openTickets() {
	files, err := ioutil.ReadDir("./pkg/tickets/")
	var tickets []Ticket
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		fmt.Println(file.Name())
		jsonFile, errorJ := os.Open("./pkg/tickets/" + file.Name())
		if errorJ != nil {
			fmt.Println(errorJ)
		}
		fmt.Println("Successfully Opened " + file.Name())
		defer jsonFile.Close()
		value, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(value, &tickets)
		fmt.Println(tickets)

	}

	fmt.Println(tickets)
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	jsonFile, errorJ := os.Open("./pkg/tickets/ticket1.json")
	if errorJ != nil {
		fmt.Println(errorJ)
	}
	fmt.Println("Successfully Opened ticket1.json")
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var tickets Tickets
	json.Unmarshal(byteValue, &tickets)
	p := Page{ID: tickets.Tickets[0].ID, SDescription: tickets.Tickets[0].SDescription,
		Description: tickets.Tickets[0].Description, UName: tickets.Tickets[0].UName, Email: tickets.Tickets[0].Email}
	t, _ := template.ParseFiles("./pkg/frontend/dashboard.html")
	t.Execute(w, p)
}

func ticketDetailHandler(w http.ResponseWriter, r *http.Request) {
	jsonFile, errorJ := os.Open("./pkg/tickets/ticket1.json")
	if errorJ != nil {
		fmt.Println(errorJ)
	}
	fmt.Println("Successfully Opened ticket1.json")
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var tickets Tickets
	json.Unmarshal(byteValue, &tickets)
	p := Page{ID: tickets.Tickets[0].ID, SDescription: tickets.Tickets[0].SDescription,
		Description: tickets.Tickets[0].Description, UName: tickets.Tickets[0].UName, Email: tickets.Tickets[0].Email}
	t, _ := template.ParseFiles("./pkg/frontend/ticketDetail.html")
	t.Execute(w, p)
}

func ticketsHandler(w http.ResponseWriter, r *http.Request) {
	jsonFile, errorJ := os.Open("./pkg/tickets/ticket1.json")
	if errorJ != nil {
		fmt.Println(errorJ)
	}
	fmt.Println("Successfully Opened ticket1.json")
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var tickets Tickets
	json.Unmarshal(byteValue, &tickets)
	p := Page{ID: tickets.Tickets[0].ID, SDescription: tickets.Tickets[0].SDescription,
		Description: tickets.Tickets[0].Description, UName: tickets.Tickets[0].UName, Email: tickets.Tickets[0].Email}
	t, _ := template.ParseFiles("./pkg/frontend/tickets.html")
	t.Execute(w, p)
}
