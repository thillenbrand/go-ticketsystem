//2057008, 2624395, 9111696

package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"go-ticketsystem/pkg/api_in"
	"log"
	"net/http"
	"os"
	"strings"
)

var mail = api_in.Mail{}

//Erstellen eines Tickets oder Kommentares mit der Kommandozeile
func main() {
	fmt.Println("Do you want to send a mail to Ticketsystem? (y/n)")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	if strings.TrimRight(input, "\n") == "y" {
		entermail()
	} else {
		os.Exit(0)
	}
	sendmail(mail)
	fmt.Println("Message successfully created.")
}

//Auslesen der Benutzereingaben
func entermail() {
	fmt.Println("Please enter your e-mail-address:")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	mail.Address = strings.TrimRight(input, "\n")

	fmt.Println("Please enter the subject of email:")
	input, _ = reader.ReadString('\n')
	mail.Subject = strings.TrimRight(input, "\n")

	fmt.Println("Please enter the text of email:")
	input, _ = reader.ReadString('\n')
	mail.Text = strings.TrimRight(input, "\n")
}

//Request an Webserver, um Ticket zu erstellen
func sendmail(mail api_in.Mail) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	client := &http.Client{}
	username := "admin"
	passwd := "supersecret"
	req, err := http.NewRequest("PUT", "https://localhost/createTicket/~"+mail.Address+"~"+mail.Subject+"~"+mail.Text, nil)
	req.SetBasicAuth(username, passwd)
	_, err = client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
}
