package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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

func FeedMailQueue(adress string, subject string, message string) { //TODO: bei Kommentar eines bearbeiters aufrufen, Kommentar und Zieladresse übergeben

	//TODO: Mails in einer JSON abspeichern
}

func confirmMailSent(mailIDs []int) {

}
