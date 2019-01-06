package frontend

import (
	"os"
	"testing"
)

func TestFillTicket(t *testing.T) {

}

func TestOpenTickets(t *testing.T) {
	err := os.Chdir("../../")
	if err != nil {
		panic(err)
	}
	var tickets = openTickets()
	if len(tickets) == 0 {
		t.Error("Ausgelesene Tickets sind leer")
	}
}

func TestTicketID(t *testing.T) {

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
