package frontend

import (
	"fmt"
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

func TestSaveProfile(t *testing.T) {

}

func TestHandlerDashboard(t *testing.T) {

}

func TestHandlerTickets(t *testing.T) {

}

func TestHandlerOpenTickets(t *testing.T) {

}

func TestHandlerProTickets(t *testing.T) {

}

func TestHandlerClosedTickets(t *testing.T) {

}

func TestHandlerTicketDet(t *testing.T) {

}

func TestHandlerEntry(t *testing.T) {

}

func TestHandlerSave(t *testing.T) {

}

func TestHandlerRelease(t *testing.T) {

}

func TestHandlerTake(t *testing.T) {

}

func TestHandlerAssign(t *testing.T) {

}

func TestHandlerAdd(t *testing.T) {

}

func TestHandlerClose(t *testing.T) {

}

func TestHandlerProfile(t *testing.T) {

}

func TestHandlerSaveProfile(t *testing.T) {

}
