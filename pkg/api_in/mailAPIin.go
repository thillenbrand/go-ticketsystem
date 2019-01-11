//2057008, 2624395, 9111696

package api_in

import (
	"go-ticketsystem/pkg/backend"
	"net/http"
	"strings"
	"time"
)

type Mail struct {
	Address string
	Subject string
	Text    string
}

//Prüft ob ein Ticket zu Betreff schon existiert
// - ja -> Kommentar wird erstellt
// - nein -> neues Ticket wird erstellt
func TicketExist(mail Mail) bool {
	backend.UpdateTickets()
	for i := 1; i <= len(backend.TicketsByTicketID); i++ {
		if backend.TicketsByTicketID[i].Subject == mail.Subject {
			if backend.TicketsByTicketID[i].Status == "geschlossen" {
				ticketTemp := backend.TicketsByTicketID[i]
				ticketTemp.Status = "offen"
				backend.TicketsByTicketID[i] = ticketTemp
			}
			entry := backend.Entry{Date: time.Now().Local().Format("2006-01-02"), Author: mail.Address, Text: mail.Text, Visible: true}
			entrys := append(backend.TicketsByTicketID[i].Entry, entry)
			ticket := &backend.Ticket{ID: backend.TicketsByTicketID[i].ID, Subject: backend.TicketsByTicketID[i].Subject, Status: backend.TicketsByTicketID[i].Status, Assigned: backend.TicketsByTicketID[i].Assigned, IDEditor: backend.TicketsByTicketID[i].IDEditor, Entry: entrys}
			backend.SaveExternal(ticket)
			backend.UpdateTickets()
			return true
		}
	}
	id := len(backend.TicketsByTicketID) + 1
	newEntry := backend.Entry{Date: time.Now().Local().Format("2006-01-02"), Author: mail.Address, Text: mail.Text, Visible: true}
	var entry []backend.Entry
	entry = append(entry, newEntry)
	ticket := &backend.Ticket{ID: id, Subject: mail.Subject, Status: "offen", Assigned: false, IDEditor: 0, Entry: entry}
	backend.SaveExternal(ticket)
	backend.UpdateTickets()
	return false
}

//Handler der bei mailIn.go angesprochen wird
//liest Argumente der Mail aus der URL
func HandlerCreateTicket(w http.ResponseWriter, r *http.Request) {
	var mail Mail
	url := r.URL.String()
	mail.Address = strings.Split(url, "~")[1]
	mail.Subject = strings.Split(url, "~")[2]
	mail.Text = strings.Split(url, "~")[3]
	TicketExist(mail)
}
