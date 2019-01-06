//2057008, 2624395, 9111696

package frontend

import (
	"encoding/json"
	"fmt"
	"go-ticketsystem/pkg/authentication"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

//Struktur eines Kommentars
type Entry struct {
	Date    string `json:"Date"`
	Author  string `json:"Author"`
	Text    string `json:"Text"`
	Visible bool
}

//Struktur eines Tickets
type Ticket struct {
	ID       int     `json:"ID"`
	Subject  string  `json:"Subject"`
	Status   string  `json:"Status"`
	Assigned bool    `json:"Assigned"`
	IDEditor int     `json:"IDEditor"`
	Entry    []Entry `json:"Entry"`
}

type Tickets struct {
	Tickets []Ticket
}

//Struktur für die Seite TicketDetail, bei der alle Daten eines Tickets und die Auswahl aller anderen Tickets und User angezeigt wird
type TicketsDet struct {
	ID       int     `json:"ID"`
	Subject  string  `json:"Subject"`
	Status   string  `json:"Status"`
	Assigned bool    `json:"Assigned"`
	IDEditor int     `json:"IDEditor"`
	Entry    []Entry `json:"Entry"`
	Tickets  []Ticket
	Users    []authentication.User
}

//Struktur des Users kommt von userAuth.go
type User = authentication.User

//Struktur der Profildaten eines angemeldeten Users
//Value wird für das Laden der Profilseiten benötigt, damit der Slider für Vacation richtig angezeigt wird
type Profile struct {
	ID       int
	Name     string
	Pass     string
	Vacation bool
	Value    string
}

//Globale Variable, die alle vorhandenen Tickets enthält
var TicketsAll []Ticket

//wird von main.go aufgerufen, um TicketsAll zu befüllen
func FillTicket() {
	TicketsAll = openTickets()
}

//Funktion sucht alle vorhandenen Tickets in ./pkg/tickets, öffnet diese und gibt sie als []Ticket zurück
func openTickets() []Ticket {
	//Pfad von main.go, um zu den gespeicherten Tickets zu gelangen
	files, err := ioutil.ReadDir("./pkg/tickets/")
	if err != nil {
		log.Fatal(err)
	}
	//in tickets werden alle Tickets gespeichert, Variable wird am Ende zurück gegeben
	var tickets []Ticket
	//File sind alle gefundenen Tickets
	for _, file := range files {
		i := 0
		//in temporary wird das jeweils geöffnete Ticket gespeichert und später an tickets angehängt
		var temporary Ticket
		jsonFile, errorJ := os.Open("./pkg/tickets/" + file.Name())
		if errorJ != nil {
			fmt.Println(errorJ)
		}
		//Die .json-Dateinen müssen geöffnet werden und Unmarshalling wird ausgeführt
		value, _ := ioutil.ReadAll(jsonFile)
		err := json.Unmarshal(value, &temporary)
		if err != nil {
			fmt.Println(err)
		}
		tickets = append(tickets, temporary)
		err = jsonFile.Close()
		if err != nil {
			fmt.Println(err)
		}
		i++
	}
	return tickets
}

//Funktion um ID des gewünschten Tickets aus der URL auszulesen
func ticketID(r *http.Request) int {
	q := r.URL.String()
	q = strings.Split(q, "?")[1]
	id, err := strconv.Atoi(q)
	if err != nil {
		fmt.Println(err)
	}
	return id
}

//Funktion, um von anderen Klassen Tickets zu speichern
func SaveExternal(t *Ticket) {
	err := t.save()
	if err != nil {
		fmt.Println(err)
	}
}

//Funktion um Ticketänderungen/-erstellungen zu speichern
func (t *Ticket) save() error {
	filename := "./pkg/tickets/ticket" + strconv.Itoa(t.ID) + ".json"
	ticket, err := json.Marshal(t)
	if err != nil {
		fmt.Println(err)
	}
	return ioutil.WriteFile(filename, ticket, 0600)
}

//Funktion um globale Variable TicketsAll zu aktualisieren
func UpdateTickets() {
	TicketsAll = openTickets()
}

//Änderungen im Profil werden im .json gespeichert
func saveProfile() error {
	users := authentication.OpenUsers()
	//eingeloggter User wird ausgewählt und Urlaubswert verändert
	for i := 0; i < len(users.User); i++ {
		if users.User[i].ID == authentication.LoggedUserID {
			users.User[i].Vacation = authentication.LoggedUserVac
		}
	}
	filename := "./pkg/users/users.json"
	user, err := json.Marshal(users)
	if err != nil {
		fmt.Println(err)
	}
	return ioutil.WriteFile(filename, user, 0600)
}

//auf dem Dashboard werden alle Tickets angezeigt, die den eingeloggten User als Bearbeiter haben
func HandlerDashboard(w http.ResponseWriter, r *http.Request) {
	var tickets = TicketsAll
	var yourTicket []Ticket
	//alle tickets werden nach der ID des Users und den Status "in Bearbeitung" durchsucht und gesammelt
	for i := 0; i < len(tickets); i++ {
		if tickets[i].IDEditor == authentication.LoggedUserID {
			if tickets[i].Status == "in Bearbeitung" {
				yourTicket = append(yourTicket, tickets[i])
			}
		}
	}
	//dashboard.html wird mit den gefundenen tickets geladen
	p := Tickets{yourTicket}
	t, _ := template.ParseFiles("./pkg/frontend/secure/dashboard.html")
	err := t.Execute(w, p)
	if err != nil {
		fmt.Println(err)
	}
}

func HandlerTickets(w http.ResponseWriter, r *http.Request) {
	var tickets = TicketsAll
	p := Tickets{tickets}
	t, _ := template.ParseFiles("./pkg/frontend/secure/tickets.html")
	err := t.Execute(w, p)
	if err != nil {
		fmt.Println(err)
	}
}

//alle Tickets mit Status "offen" werden angezeigt
func HandlerOpenTickets(w http.ResponseWriter, r *http.Request) {
	var tickets = TicketsAll
	var openTicket []Ticket
	for i := 0; i < len(tickets); i++ {
		if tickets[i].Status == "offen" {
			openTicket = append(openTicket, tickets[i])
		}
	}
	p := Tickets{openTicket}
	t, _ := template.ParseFiles("./pkg/frontend/secure/ticketsOpen.html")
	err := t.Execute(w, p)
	if err != nil {
		fmt.Println(err)
	}
}

//dieselben Tickets wie im Dashboard-Handler werden ausgewählt und auf der Seite angezeigt
func HandlerProTickets(w http.ResponseWriter, r *http.Request) {
	var tickets = TicketsAll
	var yourTicket []Ticket
	for i := 0; i < len(tickets); i++ {
		if tickets[i].IDEditor == authentication.LoggedUserID {
			if tickets[i].Status == "in Bearbeitung" {
				yourTicket = append(yourTicket, tickets[i])
			}
		}
	}
	p := Tickets{yourTicket}
	t, _ := template.ParseFiles("./pkg/frontend/secure/ticketsProcessing.html")
	err := t.Execute(w, p)
	if err != nil {
		fmt.Println(err)
	}
}

//alle geschlossenen Tickets werden angezeigt
//dabei sieht ein Bearbeiter auch die geschlossenen Tickets der anderen Bearbeiter
func HandlerClosedTickets(w http.ResponseWriter, r *http.Request) {
	var tickets = TicketsAll
	var closedTicket []Ticket
	for i := 0; i < len(tickets); i++ {
		if tickets[i].Status == "geschlossen" {
			closedTicket = append(closedTicket, tickets[i])
		}
	}
	p := Tickets{closedTicket}
	t, _ := template.ParseFiles("./pkg/frontend/secure/ticketsClosed.html")
	err := t.Execute(w, p)
	if err != nil {
		fmt.Println(err)
	}
}

//Daten für Detailansicht eines Tickets
func HandlerTicketDet(w http.ResponseWriter, r *http.Request) {
	//Verfügbare User für Ticketzuweisung
	var users = authentication.OpenUsers()
	var user []authentication.User
	//Angemeldeter User wird aus Dropdown entfernt
	for i := 0; i < len(users.User); i++ {
		if users.User[i].ID == authentication.LoggedUserID {
			user = append(users.User[:i], users.User[i+1:]...)
		}
	}
	//User im Urlaub werden aus Dropdown entfernt
	for i := 0; i < len(user); i++ {
		if user[i].Vacation == true {
			user = append(user[:i], user[i+1:]...)
		}
	}
	id := ticketID(r)
	//alle verfügbaren Tickets
	var tickets = TicketsAll
	//Details des ausgewählten Tickets
	var ticketDet Ticket
	//alle Tickets außer ausgewähltes für Zusammenführungs-Dropdown
	var ticketsOther []Ticket
	//ausgewähltes Ticket wird befüllt und alle anderen Tickets werden an ticketsOther angehängt
	for i := 0; i < len(tickets); i++ {
		if tickets[i].ID == id {
			ticketDet = tickets[i]
			ticketsOther = append(ticketsOther, tickets[i+1:]...)
			ticketsOther = append(ticketsOther, tickets[i+1:]...)
			break
		}
	}
	//Seite wird nach dem Struct TicketsDet geladen
	p := TicketsDet{ID: ticketDet.ID, Subject: ticketDet.Subject, Status: ticketDet.Status, Assigned: ticketDet.Assigned, IDEditor: ticketDet.IDEditor, Entry: ticketDet.Entry, Tickets: ticketsOther, Users: user}
	t, _ := template.ParseFiles("./pkg/frontend/secure/ticketDetail.html")

	err := t.Execute(w, p)
	if err != nil {
		fmt.Println(err)
	}
}

//Seite zum hinzufügen eines neuen Eintrages
//Die ID des ausgewählten Tickets wird benötigt, um wieder auf das richtige Ticket zurück kommen zu können
func HandlerEntry(w http.ResponseWriter, r *http.Request) {
	var tickets = TicketsAll
	//ID des Tickets
	id := ticketID(r)
	var ticketDet Ticket
	for i := 0; i < len(tickets); i++ {
		if tickets[i].ID == id {
			ticketDet = tickets[i]
			break
		}
	}
	p := Ticket{ID: ticketDet.ID, Subject: ticketDet.Subject, Status: ticketDet.Status, Assigned: ticketDet.Assigned, IDEditor: ticketDet.IDEditor, Entry: ticketDet.Entry}
	t, _ := template.ParseFiles("./pkg/frontend/secure/entry.html")
	err := t.Execute(w, p)
	if err != nil {
		fmt.Println(err)
	}
}

//Handler um Kommentare und neue Tickets zu speichern
func HandlerSave(w http.ResponseWriter, r *http.Request) {
	//Werte aus HTML-Feldern; Datum wird automatisch ermittelt
	subject := r.FormValue("inputSubject")
	date := time.Now().Local().Format("2006-01-02")
	author := r.FormValue("inputName")
	text := r.FormValue("inputText")
	var visible bool
	//visible = true -> Kommentar auch für Ersteller sichtbar, visible = false -> nur für Bearbeiter sichtbar
	if r.FormValue("visible") == "" {
		visible = true
	} else {
		visible = false
	}
	//Bei Kommentaren wird der Name des eingeloggten Users verwendet
	if author == "" {
		author = authentication.LoggedUserName
	}
	newEntry := Entry{date, author, text, visible}
	tickets := TicketsAll
	var ticketDet Ticket
	var id int
	//wenn subject leer ist, wird ein neuer Kommentar erstellt, ansonsten ein neues Ticket
	if subject == "" {
		id = ticketID(r)
		for i := 0; i < len(tickets); i++ {
			if tickets[i].ID == id {
				ticketDet = tickets[i]
				break
			}
		}
		ticketDet.Entry = append(ticketDet.Entry, newEntry)
		ticket := &Ticket{ID: ticketDet.ID, Subject: ticketDet.Subject, Status: ticketDet.Status, Assigned: ticketDet.Assigned, IDEditor: ticketDet.IDEditor, Entry: ticketDet.Entry}
		err := ticket.save()
		if err != nil {
			fmt.Println(err)
		}
		//nach Speichern der Änderungen wird TicketsAll aktualisiert
		UpdateTickets()
		http.Redirect(w, r, "/secure/ticketDetail.html?"+strconv.Itoa(id), http.StatusFound)
	} else {
		//neues Ticket wird erstellt
		ticketDet.ID = tickets[len(tickets)-1].ID + 1
		ticketDet.Subject = subject
		ticketDet.Status = "offen"
		ticketDet.Assigned = false
		ticketDet.IDEditor = 0
		ticketDet.Entry = append(ticketDet.Entry, newEntry)
		ticket := &Ticket{ID: ticketDet.ID, Subject: ticketDet.Subject, Status: ticketDet.Status, Assigned: ticketDet.Assigned, IDEditor: ticketDet.IDEditor, Entry: ticketDet.Entry}
		err := ticket.save()
		if err != nil {
			fmt.Println(err)
		}
		UpdateTickets()
		http.Redirect(w, r, "/index.html", http.StatusFound)
	}
}

//Handler zum Freigeben von Tickets
func HandlerRelease(w http.ResponseWriter, r *http.Request) {
	var tickets = TicketsAll
	var id = ticketID(r)
	var ticketDet Ticket
	//ausgewähltes Ticket wird auf Status "offen" und unassigned gesetzt; der Bearbeiter wird gelöscht
	for i := 0; i < len(tickets); i++ {
		if tickets[i].ID == id {
			tickets[i].Status = "offen"
			tickets[i].Assigned = false
			tickets[i].IDEditor = 0
			ticketDet = tickets[i]
			break
		}
	}
	ticket := &Ticket{ID: ticketDet.ID, Subject: ticketDet.Subject, Status: ticketDet.Status, Assigned: ticketDet.Assigned, IDEditor: ticketDet.IDEditor, Entry: ticketDet.Entry}
	err := ticket.save()
	if err != nil {
		fmt.Println(err)
	}
	UpdateTickets()

	http.Redirect(w, r, "/secure/ticketDetail.html?"+strconv.Itoa(id), http.StatusFound)
}

//Handler zum Ticket übernehmen
func HandlerTake(w http.ResponseWriter, r *http.Request) {
	var tickets = TicketsAll
	var id = ticketID(r)
	var ticketDet Ticket
	//ausgewähltes Ticket wird auf Status "in Bearbeitung" und assigned gesetzt; als Bearbeiter wird der eingeloggte User gesetzt
	for i := 0; i < len(tickets); i++ {
		if tickets[i].ID == id {
			tickets[i].Status = "in Bearbeitung"
			tickets[i].Assigned = true
			tickets[i].IDEditor = authentication.LoggedUserID
			ticketDet = tickets[i]
			break
		}
	}
	ticket := &Ticket{ID: ticketDet.ID, Subject: ticketDet.Subject, Status: ticketDet.Status, Assigned: ticketDet.Assigned, IDEditor: ticketDet.IDEditor, Entry: ticketDet.Entry}
	err := ticket.save()
	if err != nil {
		fmt.Println(err)
	}
	UpdateTickets()
	http.Redirect(w, r, "/secure/ticketDetail.html?"+strconv.Itoa(id), http.StatusFound)
}

//Handler um Tickets zuzuweisen
func HandlerAssign(w http.ResponseWriter, r *http.Request) {
	var tickets = TicketsAll
	var id = ticketID(r)
	//idUser ist die ID des im Dropdown-Feldes ausgewählten Users
	idUser, err := strconv.Atoi(r.FormValue("userAssign"))
	if err != nil {
		fmt.Println(err)
	}
	var ticketDet Ticket
	//ausgewähltes Ticket wird dem ausgewählten User zugeordnet
	for i := 0; i < len(tickets); i++ {
		if tickets[i].ID == id {
			tickets[i].IDEditor = idUser
			ticketDet = tickets[i]
			break
		}
	}

	ticket := &Ticket{ID: ticketDet.ID, Subject: ticketDet.Subject, Status: ticketDet.Status, Assigned: ticketDet.Assigned, IDEditor: ticketDet.IDEditor, Entry: ticketDet.Entry}
	err = ticket.save()
	if err != nil {
		fmt.Println(err)
	}
	UpdateTickets()
	http.Redirect(w, r, "/secure/ticketDetail.html?"+strconv.Itoa(id), http.StatusFound)

}

//Handler um Tickets zusammen zu führen
func HandlerAdd(w http.ResponseWriter, r *http.Request) {
	var tickets = TicketsAll
	var id = ticketID(r)
	//ID des Tickets, das an das aktuelle Ticket angehängt werden soll
	idToAdd, err := strconv.Atoi(r.FormValue("ticketToAdd"))
	if err != nil {
		fmt.Println(err)
	}
	//aktuelles Ticket
	var ticketDet Ticket
	for i := 0; i < len(tickets); i++ {
		if tickets[i].ID == id {
			ticketDet = tickets[i]
			break
		}
	}
	//Ticket das angehängt werden soll
	var ticketToAdd Ticket
	for i := 0; i < len(tickets); i++ {
		if tickets[i].ID == idToAdd {
			ticketToAdd = tickets[i]
			break
		}
	}
	//Tickets dürfen nur zusammengefügt werden, wenn sie den gleichen Bearbeiter haben
	if ticketToAdd.IDEditor == ticketDet.IDEditor {
		for i := 0; i < len(ticketToAdd.Entry); i++ {
			ticketDet.Entry = append(ticketDet.Entry, ticketToAdd.Entry[i])
		}
		//setzt Ticket, das angehängt wird auf Status "geschlossen" und erstellt einen Systemkommentar warum das Ticket geschlossen wurde
		ticketToAdd.Status = "geschlossen"
		entryClosed := Entry{Date: time.Now().Local().Format("2006-01-02"), Author: "System", Text: "Das Ticket wurde wegen Zusammenführung geschlossen. Die Einträge wurden in Ticket Nr. " + strconv.Itoa(ticketDet.ID) + " übertragen."}
		ticketToAdd.Entry = append(ticketToAdd.Entry, entryClosed)
		//Speichert aktuelles Ticket mit angehängten Einträgen
		ticket := &Ticket{ID: ticketDet.ID, Subject: ticketDet.Subject, Status: ticketDet.Status, Assigned: ticketDet.Assigned, IDEditor: ticketDet.IDEditor, Entry: ticketDet.Entry}
		err = ticket.save()
		if err != nil {
			fmt.Println(err)
		}
		//Speichert geschlossenes Ticket
		ticket = &Ticket{ID: ticketToAdd.ID, Subject: ticketToAdd.Subject, Status: ticketToAdd.Status, Assigned: ticketToAdd.Assigned, IDEditor: ticketToAdd.IDEditor, Entry: ticketToAdd.Entry}
		err = ticket.save()
		if err != nil {
			fmt.Println(err)
		}
		UpdateTickets()
		http.Redirect(w, r, "/secure/ticketDetail.html?"+strconv.Itoa(id), http.StatusFound)
	} else {
		alert := "<script>alert('Zusammenführung fehlgeschlagen. Die Tickets haben nicht denselben Bearbeiter.');window.location = '/secure/ticketDetail.html?" + strconv.Itoa(ticketDet.ID) + "';</script>"
		i, err := fmt.Fprintf(w, alert)
		if err != nil {
			fmt.Println(i)
			fmt.Println(err)
		}
	}
}

//Handler um Tickets zu schließen
func HandlerClose(w http.ResponseWriter, r *http.Request) {
	var tickets = TicketsAll
	var id = ticketID(r)
	//ausgewähltes Ticket wird auf Status "geschlossen" und unassigned gesetzt
	var ticketDet Ticket
	for i := 0; i < len(tickets); i++ {
		if tickets[i].ID == id {
			tickets[i].Status = "geschlossen"
			tickets[i].Assigned = false
			ticketDet = tickets[i]
			break
		}
	}
	ticket := &Ticket{ID: ticketDet.ID, Subject: ticketDet.Subject, Status: ticketDet.Status, Assigned: ticketDet.Assigned, IDEditor: ticketDet.IDEditor, Entry: ticketDet.Entry}
	err := ticket.save()
	if err != nil {
		fmt.Println(err)
	}
	UpdateTickets()
	http.Redirect(w, r, "/secure/ticketDetail.html?"+strconv.Itoa(id), http.StatusFound)
}

//Handler um Profildaten des eingeloggten Users zu laden
func HandlerProfile(w http.ResponseWriter, r *http.Request) {
	vac := authentication.LoggedUserVac
	var value string
	//wenn der User im Urlaub ist wird der Slider mit checked versehen
	if vac == true {
		value = "checked"
	}
	user := Profile{ID: authentication.LoggedUserID, Name: authentication.LoggedUserName, Pass: "", Vacation: authentication.LoggedUserVac, Value: value}
	t, _ := template.ParseFiles("./pkg/frontend/secure/profile.html")
	err := t.Execute(w, user)
	if err != nil {
		fmt.Println(err)
	}
}

//Handler um Änderungen des Urlaubsstatus zu speichern
func HandlerSaveProfile(w http.ResponseWriter, r *http.Request) {
	vac := r.FormValue("vac")
	fmt.Println(vac)
	if vac == "" {
		authentication.LoggedUserVac = false
	} else {
		authentication.LoggedUserVac = true
	}
	err := saveProfile()
	if err != nil {
		fmt.Println(err)
	}
	http.Redirect(w, r, "/secure/profile.html", http.StatusFound)
}
