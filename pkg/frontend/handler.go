package frontend

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
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

func WrapperDashboard(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var tickets = openTickets()
		var yourTicket []Ticket
		for i := 0; i < len(tickets); i++ {
			if tickets[i].IDEditor == 1234 {
				yourTicket = append(yourTicket, tickets[i])
			}
		}
		p := Tickets{yourTicket}
		t, _ := template.ParseFiles("./pkg/frontend/secure/dashboard.html")
		err := t.Execute(w, p)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func WrapperTickets(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var tickets = openTickets()
		p := Tickets{tickets}
		t, _ := template.ParseFiles("./pkg/frontend/secure/tickets.html")
		err := t.Execute(w, p)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func WrapperOpenTickets(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var tickets = openTickets()
		var openTicket []Ticket
		for i := 0; i < len(tickets); i++ {
			if tickets[i].Status == "offen" {
				openTicket = append(openTicket, tickets[i])
			}
		}
		p := Tickets{openTicket}
		t, _ := template.ParseFiles("./pkg/frontend/secure/ticketsOpen.html")
		err := t.Execute(w, p)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func WrapperProTickets(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var tickets = openTickets()
		var proTicket []Ticket
		for i := 0; i < len(tickets); i++ {
			if tickets[i].Status == "in Bearbeitung" {
				proTicket = append(proTicket, tickets[i])
			}
		}
		p := Tickets{proTicket}
		t, _ := template.ParseFiles("./pkg/frontend/secure/ticketsProcessing.html")
		err := t.Execute(w, p)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func WrapperClosedTickets(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var tickets = openTickets()
		var closedTicket []Ticket
		for i := 0; i < len(tickets); i++ {
			if tickets[i].Status == "geschlossen" {
				closedTicket = append(closedTicket, tickets[i])
			}
		}
		p := Tickets{closedTicket}
		t, _ := template.ParseFiles("./pkg/frontend/secure/ticketsClosed.html")
		err := t.Execute(w, p)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func WrapperTicketDet(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var tickets = openTickets()
		q := r.URL.String()
		q = strings.Split(q, "?")[1]
		id, err := strconv.Atoi(q)
		if err != nil {
			fmt.Println(err)
		}
		var ticketDet Ticket
		for i := 0; i < len(tickets); i++ {
			if tickets[i].ID == id {
				ticketDet = tickets[i]
			}
		}
		p := Ticket{ID: ticketDet.ID, Subject: ticketDet.Subject, Status: ticketDet.Status, IDEditor: ticketDet.IDEditor, Entry: ticketDet.Entry}
		t, _ := template.ParseFiles("./pkg/frontend/secure/ticketDetail.html")
		err = t.Execute(w, p)
		if err != nil {
			fmt.Println(err)
		}
	}
}
