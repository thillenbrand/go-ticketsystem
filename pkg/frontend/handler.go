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

type TBTID struct {
	Tickets map[int]Ticket
}

//Struktur für die Seite TicketDetail, bei der alle Daten eines Tickets und die Auswahl aller anderen Tickets und User angezeigt wird
type TicketsDet struct {
	ID       int     `json:"ID"`
	Subject  string  `json:"Subject"`
	Status   string  `json:"Status"`
	Assigned bool    `json:"Assigned"`
	IDEditor int     `json:"IDEditor"`
	Entry    []Entry `json:"Entry"`
	Tickets  map[int]Ticket
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

//Globale Maps, die alle Tickets nach ID, EditorID und Status enthalten
var TicketsByTicketID map[int]Ticket

var TicketsByEditorID map[int][]Ticket

var TicketsByStatus map[string][]Ticket

//Funktion sucht alle vorhandenen Tickets in ./pkg/tickets, öffnet diese und speichert sie in verschiedenen Maps
func OpenTickets() {
	files, err := ioutil.ReadDir("./pkg/tickets/")
	if err != nil {
		log.Fatal(err)
	}
	TicketsByTicketID = make(map[int]Ticket)
	TicketsByEditorID = make(map[int][]Ticket)
	TicketsByStatus = make(map[string][]Ticket)
	for _, file := range files {
		var temporary Ticket
		jsonFile, errorJ := os.Open("./pkg/tickets/" + file.Name())
		if errorJ != nil {
			fmt.Println(errorJ)
		}
		value, _ := ioutil.ReadAll(jsonFile)
		err := json.Unmarshal(value, &temporary)
		if err != nil {
			fmt.Println(err)
		}
		//füllen Map nach Ticket ID
		TicketsByTicketID[temporary.ID] = temporary

		//füllen Map nach Editor ID
		var ticketsEditorID []Ticket = TicketsByEditorID[temporary.IDEditor]
		ticketsEditorID = append(ticketsEditorID, temporary)
		TicketsByEditorID[temporary.IDEditor] = ticketsEditorID

		//füllen Map nach Status
		var ticketsStatus []Ticket = TicketsByStatus[temporary.Status]
		ticketsStatus = append(ticketsStatus, temporary)
		TicketsByStatus[temporary.Status] = ticketsStatus

		err = jsonFile.Close()
		if err != nil {
			fmt.Println(err)
		}
	}
}

//auf dem Dashboard werden alle Tickets angezeigt, die den eingeloggten User als Bearbeiter haben
func HandlerDashboard(w http.ResponseWriter, r *http.Request) {
	OpenTickets()
	var tickets = TicketsByEditorID[authentication.CheckLoggedUserID(r)]
	var ticketsDashboard []Ticket

	for _, t := range tickets {
		if t.Status == "in Bearbeitung" {
			ticketsDashboard = append(ticketsDashboard, t)
		}
	}

	p := Tickets{ticketsDashboard}
	t, _ := template.ParseFiles("./pkg/frontend/secure/dashboard.html")
	err := t.Execute(w, p)
	if err != nil {
		fmt.Println(err)
	}
}

//alle Tickets mit Status "offen" werden angezeigt
func HandlerOpenTickets(w http.ResponseWriter, r *http.Request) {
	OpenTickets()
	var tickets = TicketsByStatus["offen"]

	p := Tickets{tickets}
	t, _ := template.ParseFiles("./pkg/frontend/secure/ticketsOpen.html")
	err := t.Execute(w, p)
	if err != nil {
		fmt.Println(err)
	}
}

//dieselben Tickets wie im Dashboard-Handler werden ausgewählt und auf der Seite angezeigt
func HandlerProTickets(w http.ResponseWriter, r *http.Request) {
	OpenTickets()
	var tickets = TicketsByEditorID[authentication.CheckLoggedUserID(r)]
	var ticketsProgress []Ticket

	for _, t := range tickets {
		if t.Status == "in Bearbeitung" {
			ticketsProgress = append(ticketsProgress, t)
		}
	}

	p := Tickets{ticketsProgress}
	t, _ := template.ParseFiles("./pkg/frontend/secure/ticketsProcessing.html")
	err := t.Execute(w, p)
	if err != nil {
		fmt.Println(err)
	}
}

//alle geschlossenen Tickets werden angezeigt
//dabei sieht ein Bearbeiter auch die geschlossenen Tickets der anderen Bearbeiter
func HandlerClosedTickets(w http.ResponseWriter, r *http.Request) {
	OpenTickets()
	var tickets = TicketsByStatus["geschlossen"]

	p := Tickets{tickets}
	t, _ := template.ParseFiles("./pkg/frontend/secure/ticketsClosed.html")
	err := t.Execute(w, p)
	if err != nil {
		fmt.Println(err)
	}
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

//Daten für Detailansicht eines Tickets
func HandlerTicketDet(w http.ResponseWriter, r *http.Request) {

	//Verfügbare User für Ticketzuweisung
	var users = authentication.OpenUsers()
	var user []authentication.User

	//Angemeldeter User wird aus Dropdown entfernt
	for i := 0; i < len(users.User); i++ {
		if users.User[i].ID == authentication.CheckLoggedUserID(r) {
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
	var tickets = TicketsByTicketID

	//Details des ausgewählten Tickets
	var ticketDet Ticket = TicketsByTicketID[id]

	//Seite wird nach dem Struct TicketsDet geladen
	p := TicketsDet{ID: ticketDet.ID, Subject: ticketDet.Subject, Status: ticketDet.Status, Assigned: ticketDet.Assigned, IDEditor: ticketDet.IDEditor, Entry: ticketDet.Entry, Tickets: tickets, Users: user}
	t, _ := template.ParseFiles("./pkg/frontend/secure/ticketDetail.html")

	err := t.Execute(w, p)
	if err != nil {
		fmt.Println(err)
	}
}

//Seite zum Hinzufügen eines neuen Eintrages
//Die ID des ausgewählten Tickets wird benötigt, um wieder auf das richtige Ticket zurück kommen zu können
func HandlerEntry(w http.ResponseWriter, r *http.Request) {

	id := ticketID(r)
	var ticketDet Ticket = TicketsByTicketID[id]

	p := Ticket{ID: ticketDet.ID, Subject: ticketDet.Subject, Status: ticketDet.Status, Assigned: ticketDet.Assigned, IDEditor: ticketDet.IDEditor, Entry: ticketDet.Entry}
	t, _ := template.ParseFiles("./pkg/frontend/secure/entry.html")

	w.WriteHeader(http.StatusOK)
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
		author = authentication.CheckLoggedUserName(r)
	}
	newEntry := Entry{date, author, text, visible}
	tickets := TicketsByTicketID

	var ticketDet Ticket
	var id int

	//wenn subject leer ist, wird ein neuer Kommentar erstellt, ansonsten ein neues Ticket
	if subject == "" {

		id = ticketID(r)

		ticketDet = tickets[id]
		ticketDet.Entry = append(ticketDet.Entry, newEntry)
		TicketsByTicketID[id] = ticketDet

		ticket := &Ticket{ID: ticketDet.ID, Subject: ticketDet.Subject, Status: ticketDet.Status, Assigned: ticketDet.Assigned, IDEditor: ticketDet.IDEditor, Entry: ticketDet.Entry}

		err := ticket.save()
		if err != nil {
			fmt.Println(err)
		}
		w.WriteHeader(http.StatusOK)
		http.Redirect(w, r, "/secure/ticketDetail.html?"+strconv.Itoa(id), http.StatusFound)
	} else {

		//neues Ticket wird erstellt
		ticketDet.ID = len(tickets) + 1
		ticketDet.Subject = subject
		ticketDet.Status = "offen"
		ticketDet.Assigned = false
		ticketDet.IDEditor = 0
		ticketDet.Entry = append(ticketDet.Entry, newEntry)

		TicketsByTicketID[ticketDet.ID] = ticketDet

		ticket := &Ticket{ID: ticketDet.ID, Subject: ticketDet.Subject, Status: ticketDet.Status, Assigned: ticketDet.Assigned, IDEditor: ticketDet.IDEditor, Entry: ticketDet.Entry}

		err := ticket.save()
		if err != nil {
			fmt.Println(err)
		}
		w.WriteHeader(http.StatusOK)
		http.Redirect(w, r, "/index.html", http.StatusFound)
	}
}

//Handler zum Freigeben von Tickets
func HandlerRelease(w http.ResponseWriter, r *http.Request) {
	var id = ticketID(r)
	var ticketTemp = TicketsByTicketID[id]

	ticketTemp.Status = "offen"
	ticketTemp.Assigned = false
	ticketTemp.IDEditor = 0

	TicketsByTicketID[id] = ticketTemp

	ticket := &Ticket{ID: ticketTemp.ID, Subject: ticketTemp.Subject, Status: ticketTemp.Status, Assigned: ticketTemp.Assigned, IDEditor: ticketTemp.IDEditor, Entry: ticketTemp.Entry}
	err := ticket.save()
	if err != nil {
		fmt.Println(err)
	}
	http.Redirect(w, r, "/secure/ticketDetail.html?"+strconv.Itoa(id), http.StatusFound)
	w.WriteHeader(http.StatusOK)
}

//Handler zum Ticket übernehmen
func HandlerTake(w http.ResponseWriter, r *http.Request) {
	var id = ticketID(r)
	var ticketTemp = TicketsByTicketID[id]

	ticketTemp.Status = "in Bearbeitung"
	ticketTemp.Assigned = true
	ticketTemp.IDEditor = authentication.CheckLoggedUserID(r)

	TicketsByTicketID[id] = ticketTemp

	ticket := &Ticket{ID: ticketTemp.ID, Subject: ticketTemp.Subject, Status: ticketTemp.Status, Assigned: ticketTemp.Assigned, IDEditor: ticketTemp.IDEditor, Entry: ticketTemp.Entry}
	err := ticket.save()
	if err != nil {
		fmt.Println(err)
	}
	http.Redirect(w, r, "/secure/ticketDetail.html?"+strconv.Itoa(id), http.StatusFound)
	w.WriteHeader(http.StatusOK)
}

//Handler um Tickets zuzuweisen
func HandlerAssign(w http.ResponseWriter, r *http.Request) {
	var id = ticketID(r)
	//idUser ist die ID des im Dropdown-Feldes ausgewählten Users
	idUser, err := strconv.Atoi(r.FormValue("userAssign"))
	if err != nil {
		fmt.Println(err)
	}

	var ticketTemp = TicketsByTicketID[id]

	ticketTemp.Status = "in Bearbeitung"
	ticketTemp.Assigned = true
	ticketTemp.IDEditor = idUser

	TicketsByTicketID[id] = ticketTemp

	ticket := &Ticket{ID: ticketTemp.ID, Subject: ticketTemp.Subject, Status: ticketTemp.Status, Assigned: ticketTemp.Assigned, IDEditor: ticketTemp.IDEditor, Entry: ticketTemp.Entry}
	err = ticket.save()
	if err != nil {
		fmt.Println(err)
	}
	http.Redirect(w, r, "/secure/ticketDetail.html?"+strconv.Itoa(id), http.StatusFound)
	w.WriteHeader(http.StatusOK)

}

//Handler um Tickets zusammen zu führen
func HandlerAdd(w http.ResponseWriter, r *http.Request) {

	var tickets = TicketsByTicketID

	var id = ticketID(r)

	//ID des Tickets, das an das aktuelle Ticket angehängt werden soll
	idToAdd, err := strconv.Atoi(r.FormValue("ticketToAdd"))
	if err != nil {
		fmt.Println(err)
	}

	//aktuelles Ticket
	var ticketDet Ticket = tickets[id]

	//Ticket das angehängt werden soll
	var ticketToAdd Ticket = tickets[idToAdd]

	//Tickets dürfen nur zusammengefügt werden, wenn sie den gleichen Bearbeiter haben
	if ticketToAdd.IDEditor == ticketDet.IDEditor {
		for i := 0; i < len(ticketToAdd.Entry); i++ {
			ticketDet.Entry = append(ticketDet.Entry, ticketToAdd.Entry[i])
		}
		TicketsByTicketID[id] = ticketDet
		//setzt Ticket, das angehängt wird auf Status "geschlossen" und erstellt einen Systemkommentar warum das Ticket geschlossen wurde
		ticketToAdd.Status = "geschlossen"

		entryClosed := Entry{Date: time.Now().Local().Format("2006-01-02"), Author: "System", Text: "Das Ticket wurde wegen Zusammenführung geschlossen. Die Einträge wurden in Ticket Nr. " + strconv.Itoa(ticketDet.ID) + " übertragen."}

		ticketToAdd.Entry = append(ticketToAdd.Entry, entryClosed)
		TicketsByTicketID[idToAdd] = ticketToAdd
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
		http.Redirect(w, r, "/secure/ticketDetail.html?"+strconv.Itoa(id), http.StatusFound)
		w.WriteHeader(http.StatusOK)
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
	var tickets = TicketsByTicketID
	var id = ticketID(r)

	//ausgewähltes Ticket wird auf Status "geschlossen" und unassigned gesetzt
	var ticketDet Ticket = tickets[id]
	ticketDet.Status = "geschlossen"
	ticketDet.Assigned = false

	TicketsByTicketID[id] = ticketDet

	ticket := &Ticket{ID: ticketDet.ID, Subject: ticketDet.Subject, Status: ticketDet.Status, Assigned: ticketDet.Assigned, IDEditor: ticketDet.IDEditor, Entry: ticketDet.Entry}
	err := ticket.save()
	if err != nil {
		fmt.Println(err)
	}
	http.Redirect(w, r, "/secure/ticketDetail.html?"+strconv.Itoa(id), http.StatusFound)
	w.WriteHeader(http.StatusOK)
}

//Handler um Profildaten des eingeloggten Users zu laden
func HandlerProfile(w http.ResponseWriter, r *http.Request) {
	vac := authentication.CheckLoggedUserVac(r)
	var value string
	//wenn der User im Urlaub ist wird der Slider mit checked versehen
	if vac == true {
		value = "checked"
	}
	user := Profile{ID: authentication.CheckLoggedUserID(r), Name: authentication.CheckLoggedUserName(r), Pass: "", Vacation: authentication.CheckLoggedUserVac(r), Value: value}
	t, _ := template.ParseFiles("./pkg/frontend/secure/profile.html")
	err := t.Execute(w, user)
	if err != nil {
		fmt.Println(err)
	}
}

// Handler um Änderungen des Urlaubsstatus zu speichern
func HandlerSaveProfile(w http.ResponseWriter, r *http.Request) {
	vac := r.FormValue("vac")
	var vacBool bool
	if vac == "" {
		vacBool = false
	} else {
		vacBool = true
	}
	err := setVacation(r, vacBool)
	if err != nil {
		fmt.Println(err)
	}
	http.Redirect(w, r, "/secure/profile.html", http.StatusFound)
	w.WriteHeader(http.StatusOK)
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

//Funktion um alle Maps zu aktualisieren
func UpdateTickets() {
	OpenTickets()
}

// geänderter Urlaubsstatus wird in unsers.json gespeichert
func setVacation(r *http.Request, vac bool) error {
	users := authentication.OpenUsers()
	//eingeloggter User wird ausgewählt und Urlaubswert verändert
	for i := 0; i < len(users.User); i++ {
		if users.User[i].ID == authentication.CheckLoggedUserID(r) {
			users.User[i].Vacation = vac
		}
	}
	filename := "./pkg/users/users.json"
	user, err := json.Marshal(users)
	if err != nil {
		fmt.Println(err)
	}
	return ioutil.WriteFile(filename, user, 0600)
}
