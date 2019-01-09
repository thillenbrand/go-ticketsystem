package api

import (
	"encoding/json"
	"fmt"
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

func saveAllMails(m MailQueue) error {
	filename := "./pkg/users/users.json"
	mails, err := json.Marshal(m)
	if err != nil {
		fmt.Println(err)
	}
	return ioutil.WriteFile(filename, mails, 0600)
}

func FeedMailQueue(adress string, subject string, message string) { //TODO: bei Kommentar eines bearbeiters aufrufen, Kommentar und Zieladresse übergeben
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
	}

}

func confirmMailSent(mailIDs []int) {

}
