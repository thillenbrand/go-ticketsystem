package backend

import (
	"fmt"
)

var ticketIDCounter int = 0
var tickets = []ticket{}

type entry struct {
	date            string
	entryCreatorID  string
	text            string
	visibleForUsers bool
}

type ticket struct {
	ticketID int
	subject  string
	status   string
	editorID string
	entries  []entry
	visible  bool
}

func ticketCreate(subject string, status string, editorID string) (newTicket ticket) {

	if subject == "" {
		fmt.Println("Subject empty!")
		return
	} else if status == "" {
		fmt.Println("Status empty!")
		return
	} else if status != "offen" && status != "in Bearbeitung" && status != "geschlossen" {
		fmt.Println("Invalid status!")
		return
	} else {
		newTicket.ticketID = ticketIDCounter
		ticketIDCounter++
		newTicket.subject = subject

		newTicket.status = status

		newTicket.editorID = editorID

		newTicket.entries = make([]entry, 0)

		newTicket.visible = true

		tickets = append(tickets, newTicket)
		return newTicket
	}
}

func ticketAssign(ticketID int, newEditorID string) {
	tickets[ticketID].editorID = newEditorID
	tickets[ticketID].status = "in Bearbeitung"
	tickets[ticketID].visible = false
}

func ticketRelease(ticketID int) {
	tickets[ticketID].editorID = ""
	tickets[ticketID].status = "offen"
	tickets[ticketID].visible = true
}

func entryCreate(ticketID int, date string, entryCreatorID string, text string, visible bool) (newEntry entry) {
	if date == "" {
		fmt.Println("Invalid date!")
	} else if text == "" {
		fmt.Println("Invalid text input!")
	} else {
		newEntry.date = date

		newEntry.entryCreatorID = entryCreatorID

		newEntry.text = text

		newEntry.visibleForUsers = visible

		tickets[ticketID].entries = append(tickets[ticketID].entries, newEntry)
	}
	return
}

func ticketConcat(ticketID int, concatTicketID int) {
	if tickets[ticketID].editorID == tickets[concatTicketID].editorID {
		for i := 0; i < len(tickets[concatTicketID].entries); i++ {
			tickets[ticketID].entries = append(tickets[ticketID].entries, tickets[concatTicketID].entries[i])
		}
	} else {
		fmt.Println("Ticket cant be concated, different Editors!")
	}
}
