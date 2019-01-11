//2057008, 2624395, 9111696

package api

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"math"
)

type MailQueue struct {
	Mail []Mail `json:"Mail"`
}

// gibt die Mails in der Warteschlage als struct zurück, alternativ könnte der Mailserver auch einfach die mailQueue.json anfordern
func GetMailsFromQueue() MailQueue {
	file, err := ioutil.ReadFile("./pkg/mailQueue/mailQueue.json")
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
	var mailQueue MailQueue

	err = json.Unmarshal(file, &mailQueue)
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
