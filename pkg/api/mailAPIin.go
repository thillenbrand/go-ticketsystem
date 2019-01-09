package api

import (
	"go-ticketsystem/pkg/frontend"
	"time"
)

type Mail struct {
	InternalID int
	Address    string
	Subject    string
	Text       string
}

func GetMail(addr string, subj string, text string) {
	mail := Mail{Address: addr, Subject: subj, Text: text}
	ticketExist(mail)
}

func ticketExist(mail Mail) bool {
	frontend.UpdateTickets()
	for i := 0; i < len(frontend.TicketsAll); i++ {
		if frontend.TicketsAll[i].Subject == mail.Subject {
			if frontend.TicketsAll[i].Status == "geschlossen" {
				frontend.TicketsAll[i].Status = "offen"
			}
			entry := frontend.Entry{Date: time.Now().Local().Format("2006-01-02"), Author: mail.Address, Text: mail.Text, Visible: true}
			entrys := append(frontend.TicketsAll[i].Entry, entry)
			ticket := &frontend.Ticket{ID: frontend.TicketsAll[i].ID, Subject: frontend.TicketsAll[i].Subject, Status: frontend.TicketsAll[i].Status, Assigned: frontend.TicketsAll[i].Assigned, IDEditor: frontend.TicketsAll[i].IDEditor, Entry: entrys}
			frontend.SaveExternal(ticket)
			frontend.UpdateTickets()
			return true
		}
	}
	id := frontend.TicketsAll[len(frontend.TicketsAll)-1].ID + 1
	newEntry := frontend.Entry{Date: time.Now().Local().Format("2006-01-02"), Author: mail.Address, Text: mail.Text, Visible: true}
	var entry []frontend.Entry
	entry = append(entry, newEntry)
	ticket := &frontend.Ticket{ID: id, Subject: mail.Subject, Status: "offen", Assigned: false, IDEditor: 0, Entry: entry}
	frontend.SaveExternal(ticket)
	frontend.UpdateTickets()
	return false
}
