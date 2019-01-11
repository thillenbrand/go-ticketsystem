//2057008, 2624395, 9111696

package main

import (
	"bufio"
	"fmt"
	"go-ticketsystem/pkg/api"
	"os"
	"strconv"
	"strings"
)

// Kommandozeilentool mit einem begrenzten Befehlssatz - ermöglicht das Einsehen der Mails in der Warteschlage und das Simunlieren des Versendens
func main() {

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
			fmt.Println("Ungültige Eingabe, bitte ernet versuchen.")
		}
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
