package frontend

import (
	"encoding/json"
	"fmt"
	"go-ticketsystem/pkg/authentication"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type Entry struct {
	Date    string `json:"Date"`
	Author  string `json:"Author"`
	Text    string `json:"Text"`
	Visible bool
}

type Ticket struct {
	ID       int     `json:"ID"`
	Subject  string  `json:"Subject"`
	Status   string  `json:"Status"`
	Assigned bool    `json:"Assigned"`
	IDEditor int     `json:"IDEditor"`
	Entry    []Entry `json:"Entry"`
}

type Tickets struct {
	Tickets []Ticket
}

type TicketsDet struct {
	ID       int     `json:"ID"`
	Subject  string  `json:"Subject"`
	Status   string  `json:"Status"`
	Assigned bool    `json:"Assigned"`
	IDEditor int     `json:"IDEditor"`
	Entry    []Entry `json:"Entry"`
	Tickets  []Ticket
	Users    []authentication.User
}

type User = authentication.User

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

func HandlerDashboard(w http.ResponseWriter, r *http.Request) {

	var tickets = openTickets()
	var yourTicket []Ticket
	for i := 0; i < len(tickets); i++ {
		if tickets[i].IDEditor == authentication.LoggedUserID {

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

func HandlerTickets(w http.ResponseWriter, r *http.Request) {
	var tickets = openTickets()
	p := Tickets{tickets}
	t, _ := template.ParseFiles("./pkg/frontend/secure/tickets.html")
	err := t.Execute(w, p)
	if err != nil {
		fmt.Println(err)
	}
}

func HandlerOpenTickets(w http.ResponseWriter, r *http.Request) {
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

func HandlerProTickets(w http.ResponseWriter, r *http.Request) {
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

func HandlerClosedTickets(w http.ResponseWriter, r *http.Request) {
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

func HandlerTicketDet(w http.ResponseWriter, r *http.Request) {
	var tickets = openTickets()
	var users = authentication.OpenUsers()
	var user []authentication.User
	for i := 0; i < len(users.User); i++ {
		if users.User[i].ID == authentication.LoggedUserID {
			user = append(users.User[:i], users.User[i+1:]...)
		}
	}
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
			tickets = append(tickets[:i], tickets[i+1:]...)
			break
		}
	}
	p := TicketsDet{ID: ticketDet.ID, Subject: ticketDet.Subject, Status: ticketDet.Status, Assigned: ticketDet.Assigned, IDEditor: ticketDet.IDEditor, Entry: ticketDet.Entry, Tickets: tickets, Users: user}
	t, _ := template.ParseFiles("./pkg/frontend/secure/ticketDetail.html")

	err = t.Execute(w, p)
	if err != nil {
		fmt.Println(err)
	}
}

func HandlerEntry(w http.ResponseWriter, r *http.Request) {
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
			break
		}
	}
	p := Ticket{ID: ticketDet.ID, Subject: ticketDet.Subject, Status: ticketDet.Status, Assigned: ticketDet.Assigned, IDEditor: ticketDet.IDEditor, Entry: ticketDet.Entry}
	t, _ := template.ParseFiles("./pkg/frontend/secure/entry.html")
	err = t.Execute(w, p)
	if err != nil {
		fmt.Println(err)
	}
}

func HandlerSave(w http.ResponseWriter, r *http.Request) {
	subject := r.FormValue("inputSubject")
	date := time.Now().Local().Format("2006-01-02")
	author := r.FormValue("inputName")
	text := r.FormValue("inputText")
	var visible bool
	//visible = true -> auch für Ersteller sichtbar, visible = false -> nur für Bearbeiter sichtbar
	if r.FormValue("visible") == "" {
		visible = true
	} else {
		visible = false
	}
	newEntry := Entry{date, author, text, visible}
	tickets := openTickets()
	var ticketDet Ticket
	var id int
	var err error
	if subject == "" {
		q := r.URL.String()
		q = strings.Split(q, "?")[1]
		id, err = strconv.Atoi(q)
		if err != nil {
			fmt.Println(err)
		}

		for i := 0; i < len(tickets); i++ {
			if tickets[i].ID == id {
				ticketDet = tickets[i]
				break
			}
		}
		ticketDet.Entry = append(ticketDet.Entry, newEntry)
		ticket := &Ticket{ID: ticketDet.ID, Subject: ticketDet.Subject, Status: ticketDet.Status, Assigned: ticketDet.Assigned, IDEditor: ticketDet.IDEditor, Entry: ticketDet.Entry}
		err = ticket.save()
		if err != nil {
			fmt.Println(err)
		}

		http.Redirect(w, r, "/secure/ticketDetail.html?"+strconv.Itoa(id), http.StatusFound)
	} else {
		ticketDet.ID = tickets[len(tickets)-1].ID + 1
		ticketDet.Subject = subject
		ticketDet.Status = "offen"
		ticketDet.Assigned = false
		ticketDet.IDEditor = 0
		ticketDet.Entry = append(ticketDet.Entry, newEntry)
		ticket := &Ticket{ID: ticketDet.ID, Subject: ticketDet.Subject, Status: ticketDet.Status, Assigned: ticketDet.Assigned, IDEditor: ticketDet.IDEditor, Entry: ticketDet.Entry}
		err = ticket.save()
		if err != nil {
			fmt.Println(err)
		}
		http.Redirect(w, r, "/index.html", http.StatusFound)
	}
}

func HandlerRelease(w http.ResponseWriter, r *http.Request) {
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
			tickets[i].Status = "offen"
			tickets[i].Assigned = false
			tickets[i].IDEditor = 0
			ticketDet = tickets[i]
			break
		}
	}
	ticket := &Ticket{ID: ticketDet.ID, Subject: ticketDet.Subject, Status: ticketDet.Status, Assigned: ticketDet.Assigned, IDEditor: ticketDet.IDEditor, Entry: ticketDet.Entry}
	err = ticket.save()
	if err != nil {
		fmt.Println(err)
	}
	http.Redirect(w, r, "/secure/ticketDetail.html?"+strconv.Itoa(id), http.StatusFound)
}

func HandlerTake(w http.ResponseWriter, r *http.Request) {
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
			tickets[i].Status = "in Bearbeitung"
			tickets[i].Assigned = true
			tickets[i].IDEditor = authentication.LoggedUserID
			ticketDet = tickets[i]
			break
		}
	}
	ticket := &Ticket{ID: ticketDet.ID, Subject: ticketDet.Subject, Status: ticketDet.Status, Assigned: ticketDet.Assigned, IDEditor: ticketDet.IDEditor, Entry: ticketDet.Entry}
	err = ticket.save()
	if err != nil {
		fmt.Println(err)
	}
	http.Redirect(w, r, "/secure/ticketDetail.html?"+strconv.Itoa(id), http.StatusFound)
}

func HandlerAssign(w http.ResponseWriter, r *http.Request) {
	var tickets = openTickets()
	q := r.URL.String()
	q = strings.Split(q, "?")[1]
	id, err := strconv.Atoi(q)
	if err != nil {
		fmt.Println(err)
	}
	idUser, err := strconv.Atoi(r.FormValue("userAssign"))
	if err != nil {
		fmt.Println(err)
	}
	var ticketDet Ticket
	for i := 0; i < len(tickets); i++ {
		if tickets[i].ID == id {
			tickets[i].IDEditor = idUser
			ticketDet = tickets[i]
			break
		}
	}

	ticket := &Ticket{ID: ticketDet.ID, Subject: ticketDet.Subject, Status: ticketDet.Status, Assigned: ticketDet.Assigned, IDEditor: ticketDet.IDEditor, Entry: ticketDet.Entry}
	err = ticket.save()
	if err != nil {
		fmt.Println(err)
	}

	http.Redirect(w, r, "/secure/ticketDetail.html?"+strconv.Itoa(id), http.StatusFound)

}

func HandlerAdd(w http.ResponseWriter, r *http.Request) {
	var tickets = openTickets()
	q := r.URL.String()
	q = strings.Split(q, "?")[1]
	id, err := strconv.Atoi(q)
	if err != nil {
		fmt.Println(err)
	}
	idToAdd, err := strconv.Atoi(r.FormValue("ticketToAdd"))
	if err != nil {
		fmt.Println(err)
	}
	var ticketDet Ticket
	for i := 0; i < len(tickets); i++ {
		if tickets[i].ID == id {
			ticketDet = tickets[i]
			break
		}
	}
	var ticketToAdd Ticket
	for i := 0; i < len(tickets); i++ {
		if tickets[i].ID == idToAdd {
			ticketToAdd = tickets[i]
			break
		}
	}
	if ticketToAdd.IDEditor == ticketDet.IDEditor {
		for i := 0; i < len(ticketToAdd.Entry); i++ {
			ticketDet.Entry = append(ticketDet.Entry, ticketToAdd.Entry[i])
		}
		ticketToAdd.Status = "geschlossen"
		entryClosed := Entry{Date: time.Now().Local().Format("2006-01-02"), Author: "System", Text: "Das Ticket wurde wegen Zusammenführung geschlossen. Die Einträge wurden in Ticket Nr. " + strconv.Itoa(ticketDet.ID) + " übertragen."}
		ticketToAdd.Entry = append(ticketToAdd.Entry, entryClosed)

		ticket := &Ticket{ID: ticketDet.ID, Subject: ticketDet.Subject, Status: ticketDet.Status, Assigned: ticketDet.Assigned, IDEditor: ticketDet.IDEditor, Entry: ticketDet.Entry}
		err = ticket.save()
		if err != nil {
			fmt.Println(err)
		}
		ticket = &Ticket{ID: ticketToAdd.ID, Subject: ticketToAdd.Subject, Status: ticketToAdd.Status, Assigned: ticketToAdd.Assigned, IDEditor: ticketToAdd.IDEditor, Entry: ticketToAdd.Entry}
		err = ticket.save()
		if err != nil {
			fmt.Println(err)
		}
		http.Redirect(w, r, "/secure/ticketDetail.html?"+strconv.Itoa(id), http.StatusFound)
	} else {
		alert := "<script>alert('Zusammenführung fehlgeschlagen. Die Tickets haben nicht denselben Bearbeiter.');window.location = '/secure/ticketDetail.html?" + strconv.Itoa(ticketDet.ID) + "';</script>"
		fmt.Fprintf(w, alert)
	}
}

func HandlerProfile(w http.ResponseWriter, r *http.Request) {
	vac := r.FormValue("vac")
	fmt.Println(vac)
	user := authentication.User{ID: authentication.LoggedUserID, Name: authentication.LoggedUserName, Pass: "", Vacation: authentication.LoggedUserVac}
	t, _ := template.ParseFiles("./pkg/frontend/secure/profile.html")
	err := t.Execute(w, user)
	if err != nil {
		fmt.Println(err)
	}
}

func HandlerSaveProfile(w http.ResponseWriter, r *http.Request) {
	vac := r.FormValue("vac")
	fmt.Println(vac)
	http.Redirect(w, r, "/secure/profile.html", http.StatusFound)
}

func (t *Ticket) save() error {
	filename := "./pkg/tickets/ticket" + strconv.Itoa(t.ID) + ".json"
	ticket, err := json.Marshal(t)
	if err != nil {
		fmt.Println(err)
	}
	return ioutil.WriteFile(filename, ticket, 0600)
}
