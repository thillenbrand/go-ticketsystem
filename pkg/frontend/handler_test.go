package frontend

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
)

func init() {
	err := os.Chdir("../../")
	if err != nil {
		panic(err)
	}
}

func TestOpenTickets(t *testing.T) {
	UpdateTickets()
	var entry []Entry
	entry = append(entry, Entry{Date: "2019-01-01", Author: "Test", Text: "Test", Visible: true})
	ticket := &Ticket{ID: len(TicketsAll) + 1, Subject: "Test", Status: "offen", Assigned: false, IDEditor: 0, Entry: entry}
	SaveExternal(ticket)
	var tickets = openTickets()
	if len(tickets) == 0 {
		t.Error()
	}
	err := os.Remove("./pkg/tickets/ticket" + strconv.Itoa(ticket.ID) + ".json")
	if err != nil {
		fmt.Println(err)
	}
}

func TestTicketID(t *testing.T) {
	request := httptest.NewRequest("", "/ticket1?1", nil)
	if ticketID(request) != 1 {
		t.Error()
	}
}

func TestSaveExternal(t *testing.T) {
	UpdateTickets()
	start := len(TicketsAll)
	var entry []Entry
	entry = append(entry, Entry{Date: "2019-01-01", Author: "Test", Text: "Test", Visible: true})
	ticket := &Ticket{ID: len(TicketsAll) + 1, Subject: "Test", Status: "offen", Assigned: false, IDEditor: 0, Entry: entry}
	SaveExternal(ticket)
	UpdateTickets()
	if len(TicketsAll) == start {
		t.Error()
	}
	err := os.Remove("./pkg/tickets/ticket" + strconv.Itoa(ticket.ID) + ".json")
	if err != nil {
		fmt.Println(err)
	}
	UpdateTickets()
}

func TestSave(t *testing.T) {
	UpdateTickets()
	start := len(TicketsAll)
	var entry []Entry
	entry = append(entry, Entry{Date: "2019-01-01", Author: "Test", Text: "Test", Visible: true})
	ticket := &Ticket{ID: len(TicketsAll) + 1, Subject: "Test", Status: "offen", Assigned: false, IDEditor: 0, Entry: entry}
	err := ticket.save()
	if err != nil {
		fmt.Println(err)
	}
	UpdateTickets()
	if len(TicketsAll) == start {
		t.Error()
	}
	err = os.Remove("./pkg/tickets/ticket" + strconv.Itoa(ticket.ID) + ".json")
	if err != nil {
		fmt.Println(err)
	}
	UpdateTickets()
}

func TestUpdateTickets(t *testing.T) {
	start := len(TicketsAll)
	var entry []Entry
	entry = append(entry, Entry{Date: "2019-01-01", Author: "Test", Text: "Test", Visible: true})
	ticket := &Ticket{ID: len(TicketsAll) + 1, Subject: "Test", Status: "offen", Assigned: false, IDEditor: 0, Entry: entry}
	SaveExternal(ticket)
	UpdateTickets()
	if len(TicketsAll) == start {
		t.Error()
	}
	err := os.Remove("./pkg/tickets/ticket" + strconv.Itoa(ticket.ID) + ".json")
	if err != nil {
		fmt.Println(err)
	}
}

func TestHandlerDashboard(t *testing.T) {
	req, err := http.NewRequest("GET", "/secure/dashboard.html", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HandlerDashboard)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Error()
	}
}

func TestHandlerTickets(t *testing.T) {
	req, err := http.NewRequest("GET", "/secure/tickets.html", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HandlerTickets)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Error()
	}
}

func TestHandlerOpenTickets(t *testing.T) {
	req, err := http.NewRequest("GET", "/secure/ticketsOpen.html", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HandlerOpenTickets)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Error()
	}
}

func TestHandlerProTickets(t *testing.T) {
	req, err := http.NewRequest("GET", "/secure/ticketsProcessing.html", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HandlerProTickets)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Error()
	}
}

func TestHandlerClosedTickets(t *testing.T) {
	req, err := http.NewRequest("GET", "/secure/ticketsClosed.html", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HandlerClosedTickets)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Error()
	}
}

func TestHandlerTicketDet(t *testing.T) {
	req, err := http.NewRequest("GET", "/secure/ticketDetail.html?1", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HandlerTicketDet)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Error()
	}
}

func TestHandlerEntry(t *testing.T) {
	req, err := http.NewRequest("GET", "/secure/entry.html?1", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HandlerEntry)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Error()
	}
}

func TestHandlerSave(t *testing.T) {
	req, err := http.NewRequest("GET", "/secure/save/ticket?1", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HandlerSave)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Error()
	}
}

func TestHandlerRelease(t *testing.T) {
	req, err := http.NewRequest("GET", "/secure/release/ticket?1", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HandlerRelease)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Error()
	}
}

func TestHandlerTake(t *testing.T) {
	req, err := http.NewRequest("GET", "/secure/take/ticket?1", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HandlerTake)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Error()
	}
}

func TestHandlerAssign(t *testing.T) {
	req, err := http.NewRequest("GET", "/secure/assign/ticket?1", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HandlerAssign)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Error()
	}
}

func TestHandlerAdd(t *testing.T) {
	req, err := http.NewRequest("GET", "/secure/add/ticket?1", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HandlerAdd)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Error()
	}
}

func TestHandlerClose(t *testing.T) {
	req, err := http.NewRequest("GET", "/secure/close/ticket?1", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HandlerClose)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Error()
	}
}

func TestHandlerProfile(t *testing.T) {
	req, err := http.NewRequest("GET", "/secure/dashboard.html", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HandlerProfile)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Error()
	}
}

func TestHandlerSaveProfile(t *testing.T) {
	req, err := http.NewRequest("GET", "/secure/profile.html", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HandlerSaveProfile)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Error()
	}
}
