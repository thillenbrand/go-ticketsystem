package frontend

import (
	"net/http/httptest"
	"os"
	"testing"
)

func init() {
	err := os.Chdir("../../")
	if err != nil {
		panic(err)
	}
}

func TestFillTicket(t *testing.T) {
	if len(TicketsAll) != 0 {
		TicketsAll = TicketsAll[:0]
	}
	FillTicket()
	if len(TicketsAll) == 0 {
		t.Error()
	}
}

func TestOpenTickets(t *testing.T) {
	var tickets = openTickets()
	if len(tickets) == 0 {
		t.Error("Ausgelesene Tickets sind leer")
	}
}

func TestTicketID(t *testing.T) {
	request := httptest.NewRequest("", "/ticket1?1", nil)
	if ticketID(request) != 1 {
		t.Error()
	}
}

func TestSaveExternal(t *testing.T) {

}

func TestSave(t *testing.T) {

}

func TestUpdateTickets(t *testing.T) {

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
