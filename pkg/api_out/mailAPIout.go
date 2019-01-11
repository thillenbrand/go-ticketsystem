//2057008, 2624395, 9111696

package api_out

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"math"
	"net/http"
)

type Mail struct {
	InternalID int
	Address    string
	Subject    string
	Text       string
}

type IDs struct {
	ID []int
}

// diese struct stellt die Warteschlage dar - sie enthält eine slice aus Mails
type MailQueue struct {
	Mail []Mail `json:"Mail"`
}

// der Handler nimmt REST-POST-Nachrichten (idempotent) an und löscht die dadurch als versendet bestätigten Mails aus der Queue
// jede Nachricht enthält eine JSON, in der die IDs der abgesendeten Mails angegeben werden
func HandlerConfirmSend(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Printf("mailAPIout: Error reading body: %v", err)
		http.Error(w, "mailAPIout: Can not read request body", http.StatusBadRequest)
		return
	}
	var IDs IDs
	json.Unmarshal(reqBody, &IDs)
	var ids []int
	for _, i := range IDs.ID {
		ids = append(ids, i)
	}

	ConfirmMailSent(ids)
}

// der Handler sendet die momentane Warteschlage an den Mailserver, unabhängig von der eingehenden Nachricht - es gibt keine Änderungen auf dem Server hier selbst (nullipotent)
func HandlerSendMail(w http.ResponseWriter, r *http.Request) {
	sendMailQueue()
}

// gibt die Warteschlage aus Mails als byte-slice aus
func getQueueFile() []byte {
	file, err := ioutil.ReadFile("./pkg/mailQueue/mailQueue.json")
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}

	return file
}

// sendet die moemtane Warteschlage an den Mailserver, nutzt https und Basic Authentication (admin:supersecret)
func sendMailQueue() error {
	var jsonStr = []byte(getQueueFile())

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	client := &http.Client{}
	username := "admin"
	passwd := "supersecret"
	req, err := http.NewRequest("POST", "https://example.com/mails/toSend", bytes.NewBuffer(jsonStr)) //hier wird die Adresse des Mailservers statt example.com eingetragen
	req.SetBasicAuth(username, passwd)
	_, err = client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

// gibt die Mails in der Warteschlage als struct zurück - die API selbst
func GetMailsFromQueue() MailQueue {

	var mailQueue MailQueue

	err := json.Unmarshal(getQueueFile(), &mailQueue)
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}

	return mailQueue
}

// speichert noch abzuschickende Mails in der Warteschlangen-Datei
func saveAllMails(m MailQueue) error {
	filename := "./pkg/mailQueue/mailQueue.json"
	mails, err := json.Marshal(m)
	if err != nil {
		fmt.Println(err)
	}
	return ioutil.WriteFile(filename, mails, 0600)
}

// wird genutzt um neue Mails auf die Warteschlage zu setzen
func FeedMailQueue(adress string, subject string, message string) error { //TODO: bei Kommentar eines bearbeiters aufrufen, Kommentar und Zieladresse übergeben
	mailQueue := GetMailsFromQueue()
	oneMail := mailQueue.Mail
	var counter float64 = 0

	for _, m := range oneMail {
		counter = math.Max(float64(m.InternalID), counter)
	}

	var newMail Mail
	newMail.InternalID = int(counter) + 1
	newMail.Address = adress
	newMail.Subject = subject
	newMail.Text = message

	mailQueue.Mail = append(mailQueue.Mail, newMail)
	err := saveAllMails(mailQueue)
	if err != nil {
		fmt.Println(err)
		return errors.New("mailAPIout: updating the queue failed")
	}

	return nil
}

// prüft, ob ein int-slice ein bestimmtes int nicht enthält
func notContains(data int, intSlice []int) bool {
	for _, x := range intSlice {
		if data == x {
			return false
		}
	}
	return true
}

// wird vom Mailserver angesprochen, löscht gesendete Mails aus der Warteschlange
func ConfirmMailSent(mailIDs []int) error {
	mailQueue := GetMailsFromQueue()
	oneMail := mailQueue.Mail

	var newMailQueue MailQueue

	for _, m := range oneMail {
		if notContains(m.InternalID, mailIDs) {
			newMailQueue.Mail = append(newMailQueue.Mail, m)
		}
	}

	err := saveAllMails(newMailQueue)
	if err != nil {
		fmt.Println(err)
		return errors.New("mailAPIout: updating the queue failed")
	}

	return nil
}
