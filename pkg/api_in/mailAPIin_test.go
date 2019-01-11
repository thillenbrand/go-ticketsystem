//2057008, 2624395, 9111696

package api_in

import (
	"fmt"
	"go-ticketsystem/pkg/backend"
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
	backend.UpdateTickets()
	var entry []backend.Entry
	entry = append(entry, backend.Entry{Date: "2019-01-01", Author: "Test", Text: "Test", Visible: true})
	ticket := &backend.Ticket{ID: len(backend.TicketsByTicketID) + 1, Subject: "Test", Status: "geschlossen", Assigned: false, IDEditor: 0, Entry: entry}
	backend.SaveExternal(ticket)
	backend.UpdateTickets()
}

func TestTicketExist(t *testing.T) {

	mail := Mail{Address: "test@test.de", Subject: "KeinTest", Text: "Dies ist ein Test"}
	if TicketExist(mail) != false {
		t.Error()
	}
	backend.UpdateTickets()
	id := strconv.Itoa(len(backend.TicketsByTicketID))
	err := os.Remove("./pkg/tickets/ticket" + id + ".json")
	if err != nil {
		fmt.Println(err)
	}
	backend.UpdateTickets()
}

func TestTicketExist2(t *testing.T) {
	mail := Mail{Address: "test@test.de", Subject: "Test", Text: "Dies ist ein Test"}
	if TicketExist(mail) != true {
		t.Error()
	}
	backend.UpdateTickets()
	id := strconv.Itoa(len(backend.TicketsByTicketID))
	err := os.Remove("./pkg/tickets/ticket" + id + ".json")
	if err != nil {
		fmt.Println(err)
	}
	backend.UpdateTickets()
}

func TestHandlerCreateTicket(t *testing.T) {
	req, err := http.NewRequest("GET", "/createTicket/~asd~asd~asd", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HandlerCreateTicket)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusFound {
		t.Error()
	}
}
