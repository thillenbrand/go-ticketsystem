package main

import (
	"bufio"
	"fmt"
	"go-ticketsystem/pkg/api"
	"os"
	"strings"
)

var mail = api.Mail{}

func main() {
	fmt.Println("Do you want to send a mail to Ticketsystem? (y/n)")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	if strings.TrimRight(input, "\n") == "y" {
		entermail()
	} else {
		os.Exit(0)
	}
	fmt.Println(mail)
	Sendmail(mail)
}

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

func Sendmail(mail api.Mail) {
	api.GetMail(mail.Address, mail.Subject, mail.Text)
}
