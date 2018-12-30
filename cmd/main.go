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

type Entry struct {
	Date   string `json:"Date"`
	Author string `json:"Author"`
	Text   string `json:"Text"`
}

type Ticket struct {
	ID       int     `json:"ID"`
	Subject  string  `json:"Subject"`
	Status   string  `json:"Status"`
	IDEditor int     `json:"IDEditor"`
	Entry    []Entry `json:"Entry"`
}

type Tickets struct {
	Tickets []Ticket
}

type Page struct {
	ID       int
	Subject  string
	Status   string
	IDEditor int
	Entry    []Entry
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	http.FileServer(http.Dir("./pkg/frontend")).ServeHTTP(w, r)
}

func main() {
	/*
		http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, r.URL.Path[1:])
		})
	*/

	http.HandleFunc("/secure/dashboard.html", dashboardHandler)
	http.HandleFunc("/secure/ticketDetail", ticketDetailHandler)
	http.HandleFunc("/secure/tickets.html", ticketsHandler)
	http.HandleFunc("/", auth.Wrapper(mainHandler))

	err := http.ListenAndServeTLS(":443", "Server.crt", "Server.key", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}

func openTickets() []Ticket {
	files, err := ioutil.ReadDir("./pkg/tickets/")
	var tickets []Ticket
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		i := 0
		var temporary Ticket
		jsonFile, errorJ := os.Open("./pkg/tickets/" + file.Name())
		if errorJ != nil {
			fmt.Println(errorJ)
		}
		value, _ := ioutil.ReadAll(jsonFile)
		err := json.Unmarshal(value, &temporary)
		if err != nil {
			fmt.Println(err)
		}
		tickets = append(tickets, temporary)
		err = jsonFile.Close()
		if err != nil {
			fmt.Println(err)
		}
		i++
	}
	return tickets
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	var tickets = openTickets()
	p := Tickets{tickets}
	t, _ := template.ParseFiles("./pkg/frontend/secure/dashboard.html")
	err := t.Execute(w, p)
	if err != nil {
		fmt.Println(err)
	}
}

func ticketDetailHandler(w http.ResponseWriter, r *http.Request) {
	var tickets = openTickets()
	p := Tickets{tickets}
	t, _ := template.ParseFiles("./pkg/frontend/secure/ticketDetail.html")
	err := t.Execute(w, p)
	if err != nil {
		fmt.Println(err)
	}
}

func ticketsHandler(w http.ResponseWriter, r *http.Request) {
	var tickets = openTickets()
	p := Tickets{tickets}
	t, _ := template.ParseFiles("./pkg/frontend/secure/tickets.html")
	err := t.Execute(w, p)
	if err != nil {
		fmt.Println(err)
	}
}
