//2057008, 2624395, 9111696

package api

import (
	"fmt"
	"go-ticketsystem/pkg/frontend"
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

func TestTicketExist(t *testing.T) {
	frontend.UpdateTickets()
	var entry []frontend.Entry
	entry = append(entry, frontend.Entry{Date: "2019-01-01", Author: "Test", Text: "Test", Visible: true})
	ticket := &frontend.Ticket{ID: len(frontend.TicketsByTicketID) + 1, Subject: "Test", Status: "offen", Assigned: false, IDEditor: 0, Entry: entry}
	frontend.SaveExternal(ticket)
	frontend.UpdateTickets()
	mail := Mail{Address: "test@test.de", Subject: "KeinTest", Text: "Dies ist ein Test"}
	if ticketExist(mail) != false {
		t.Error()
	}

	frontend.UpdateTickets()
	id := strconv.Itoa(len(frontend.TicketsByTicketID))
	err := os.Remove("./pkg/tickets/ticket" + id + ".json")
	if err != nil {
		fmt.Println(err)
	}
	frontend.UpdateTickets()
}

func TestTicketExist2(t *testing.T) {
	frontend.UpdateTickets()
	var entry []frontend.Entry
	entry = append(entry, frontend.Entry{Date: "2019-01-01", Author: "Test", Text: "Test", Visible: true})
	ticket := &frontend.Ticket{ID: len(frontend.TicketsByTicketID) + 1, Subject: "Test", Status: "offen", Assigned: false, IDEditor: 0, Entry: entry}
	frontend.SaveExternal(ticket)
	frontend.UpdateTickets()
	mail := Mail{Address: "test@test.de", Subject: "Test", Text: "Dies ist ein Test"}
	if ticketExist(mail) != true {
		t.Error()
	}
	frontend.UpdateTickets()
	id := strconv.Itoa(len(frontend.TicketsByTicketID))
	err := os.Remove("./pkg/tickets/ticket" + id + ".json")
	if err != nil {
		fmt.Println(err)
	}
	frontend.UpdateTickets()
}
