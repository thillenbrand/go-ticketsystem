//2057008, 2624395, 9111696

package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"go-ticketsystem/pkg/api"
	auth "go-ticketsystem/pkg/authentication"
	hand "go-ticketsystem/pkg/frontend"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var mail = api.Mail{}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	http.FileServer(http.Dir("./pkg/frontend")).ServeHTTP(w, r)
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
func sendmail(mail api.Mail) {
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

// gibt die Mails in der Warteschlange aus
func getMailQueue() {
	mailQueue := api.GetMailsFromQueue()
	oneMail := mailQueue.Mail

	for _, m := range oneMail {
		fmt.Println("ID der Mail: " + string(m.InternalID))
		fmt.Println("Betreff der Mail: " + m.Subject)
		fmt.Println("Zieladresse der Mail: " + m.Subject)
		fmt.Println(m.Text)
		fmt.Println("")
	}
}

func main() {
	//Checkt ob Folders tickets, users und mailQueue existieren, falls nicht werden diese erstellt.
	if _, err := os.Stat("./pkg/tickets"); os.IsNotExist(err) {
		os.Mkdir("./pkg/tickets", 0700)
	}

	if _, err := os.Stat("./pkg/users"); os.IsNotExist(err) {
		os.Mkdir("./pkg/users", 0700)
	}

	if _, err := os.Stat("./pkg/mailQueue"); os.IsNotExist(err) {
		os.Mkdir("./pkg/mailQueue", 0700)
	}

	//Flags zum Einstellen des Ports und Mode

	//default ist 0 für den Server
	var mode = flag.Int("mode", 0, "set 0 for server, 1 for mailIn, 2 for mailOut")

	var port = flag.Int("port", 443, "sets port for https server")

	flag.Parse()

	if *mode == 0 {

		hand.UpdateTickets()

		http.HandleFunc("/", mainHandler)

		http.HandleFunc("/createTicket/", auth.Wrapper(api.HandlerCreateTicket))

		http.HandleFunc("/secure/dashboard.html", auth.Wrapper(hand.HandlerDashboard))
		http.HandleFunc("/secure/ticketDetail.html", auth.Wrapper(hand.HandlerTicketDet))
		http.HandleFunc("/secure/ticketsOpen.html", auth.Wrapper(hand.HandlerOpenTickets))
		http.HandleFunc("/secure/ticketsProcessing.html", auth.Wrapper(hand.HandlerProTickets))
		http.HandleFunc("/secure/ticketsClosed.html", auth.Wrapper(hand.HandlerClosedTickets))
		http.HandleFunc("/secure/entry.html", auth.Wrapper(hand.HandlerEntry))
		http.HandleFunc("/secure/saveP/", auth.Wrapper(hand.HandlerSaveProfile))
		http.HandleFunc("/secure/save/", auth.Wrapper(hand.HandlerSave))
		http.HandleFunc("/save/", hand.HandlerSave)
		http.HandleFunc("/secure/close/", auth.Wrapper(hand.HandlerClose))
		http.HandleFunc("/secure/release/", auth.Wrapper(hand.HandlerRelease))
		http.HandleFunc("/secure/take/", auth.Wrapper(hand.HandlerTake))
		http.HandleFunc("/secure/add/", auth.Wrapper(hand.HandlerAdd))
		http.HandleFunc("/secure/assign/", auth.Wrapper(hand.HandlerAssign))
		http.HandleFunc("/secure/profile.html", auth.Wrapper(hand.HandlerProfile))
		http.HandleFunc("/register/", auth.HandlerRegister)

		err := http.ListenAndServeTLS(":"+strconv.Itoa(*port)+"", "Server.crt", "Server.key", nil)
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}

	} else if *mode == 1 {

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

	} else if *mode == 2 {

		fmt.Println("Dieses Kommandozeilentool kann die E-Mails anzeigen, die momentan in der Warteschlage stehen.")
		fmt.Println("Mögliche Befehle:")
		fmt.Println("show - alle Mails in der Warteschlage anzeigen")
		fmt.Println("send - Versenden durch den Mailservers nachahmen (entfernt die betreffenden Mails aus der Warteschlange!")
		fmt.Println("exit - Programm verlassen")

		for {

			inputReader := bufio.NewReader(os.Stdin)
			input, _ := inputReader.ReadString('\n')

			if strings.TrimRight(input, "\n") == "show" {
				fmt.Println("Die E-Mails in der Warteschlage werden nun angezeigt:")
				getMailQueue() //TODO: mit gefüllter json testen
				fmt.Println("")
				fmt.Println("Alle Mails in der Warteschlage wurden angezeigt, mögliche Befehle:")
				fmt.Println("show - alle Mails in der Warteschlage anzeigen")
				fmt.Println("send - Versenden durch den Mailservers nachahmen (entfernt die betreffenden Mails aus der Warteschlange!")
				fmt.Println("exit - Programm verlassen")
			} else if strings.TrimRight(input, "\n") == "send" {
				fmt.Println("Bitte die ID der zu versendenden Mails angeben (mit Kommatas getrennt, ohne Leerzeichen, Eingabe mit der Eingabetaste abschließen):")
				idReader := bufio.NewReader(os.Stdin)
				inputId, _ := idReader.ReadString('\n')
				ids := strings.Split(inputId, ",")
				ids[len(ids)-1] = strings.Split(ids[len(ids)-1], "\n")[0]
				var mailIds []int

				for _, id := range ids {
					number, _ := strconv.Atoi(id)
					if number != 0 { //TODO: 0 einfügen in json als default Mail
						mailIds = append(mailIds, number)
					}
				}
				err := api.ConfirmMailSent(mailIds)
				if err != nil {
					fmt.Println(err)
				}

				fmt.Println("Mail(s) wurde(n) versendet, mögliche Befehle:")
				fmt.Println("show - alle Mails in der Warteschlage anzeigen")
				fmt.Println("send - Versenden durch den Mailservers nachahmen (entfernt die betreffenden Mails aus der Warteschlange!")
				fmt.Println("exit - Programm verlassen")
			} else if strings.TrimRight(input, "\n") == "exit" {
				fmt.Println("Programm wird geschlossen.")
			} else {
				fmt.Println("Ungültige Eingabe, bitte erneut versuchen.")
			}
		}
	} else {
		fmt.Println("Invalid flag for mode!")
	}
}
