//2057008, 2624395, 9111696

package api

import (
	frontend "go-ticketsystem-alt/pkg/frontend"
	"net/http"
	"strings"
	"time"
)

type Mail struct {
	InternalID int
	Address    string
	Subject    string
	Text       string
}

//Handler der bei mailIn.go angesprochen wird
//liest Argumente der Mail aus der URL
func HandlerCreateTicket(w http.ResponseWriter, r *http.Request) {
	var mail Mail
	url := r.URL.String()
	mail.Address = strings.Split(url, "~")[1]
	mail.Subject = strings.Split(url, "~")[2]
	mail.Text = strings.Split(url, "~")[3]
	ticketExist(mail)
}

//Prüft ob ein Ticket zu Betreff schon existiert
// - ja -> Kommentar wird erstellt
// - nein -> neues Ticket wird erstellt
func ticketExist(mail Mail) bool {
	frontend.UpdateTickets()
	var tickets = frontend.TicketsByTicketID
	for i := 0; i < len(tickets); i++ {
		if tickets[i].Subject == mail.Subject {
			if tickets[i].Status == "geschlossen" {
				ticketTemp := tickets[i]
				ticketTemp.Status = "offen"
				tickets[i] = ticketTemp
			}
			entry := frontend.Entry{Date: time.Now().Local().Format("2006-01-02"), Author: mail.Address, Text: mail.Text, Visible: true}
			entrys := append(frontend.TicketsByTicketID[i].Entry, entry)
			ticket := &frontend.Ticket{ID: frontend.TicketsByTicketID[i].ID, Subject: frontend.TicketsByTicketID[i].Subject, Status: frontend.TicketsByTicketID[i].Status, Assigned: frontend.TicketsByTicketID[i].Assigned, IDEditor: frontend.TicketsByTicketID[i].IDEditor, Entry: entrys}
			frontend.SaveExternal(ticket)
			frontend.UpdateTickets()
			return true
		}
	}
	id := len(frontend.TicketsByTicketID) + 1
	newEntry := frontend.Entry{Date: time.Now().Local().Format("2006-01-02"), Author: mail.Address, Text: mail.Text, Visible: true}
	var entry []frontend.Entry
	entry = append(entry, newEntry)
	ticket := &frontend.Ticket{ID: id, Subject: mail.Subject, Status: "offen", Assigned: false, IDEditor: 0, Entry: entry}
	frontend.SaveExternal(ticket)
	frontend.UpdateTickets()
	return false
}
