package authentication

import (
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt" //TODO: externe Lib internalisieren?
	"io/ioutil"
	"log"
	"net/http"
)

// struct, das einen User darstellt
type User struct {
	ID       int    `json:"ID"`
	Name     string `json:"Name"`
	Pass     string `json:"Pass"`
	Vacation bool   `json:"Vacation"`
	// TODO: brauchen wir eine Liste von Tickets, die dem User zugewiesen sind?
}

// Users ist lediglich eine Auflistung meherer User
type Users struct {
	User []User `json:"Users"`
}

//TODO: Globale Variablen, die angeben, wer sich gerade eingeloggt - ersetzen durch prüfungen in den Funktionen selber - eigene Funktion schreiben
var LoggedUserName string
var LoggedUserID int
var LoggedUserVac bool

func CheckLoggedUserName(r *http.Request) string {
	user, _, _ := r.BasicAuth()

	return user
}

func CheckLoggedUserID(r *http.Request) int {
	user, _, _ := r.BasicAuth()
	var users = OpenUsers().User
	for _, u := range users {
		if u.Name == user {
			return u.ID
			break
		}
	}
	return 0
}

func CheckLoggedUserVac(r *http.Request) bool {
	user, _, _ := r.BasicAuth()
	var users = OpenUsers().User
	for _, u := range users {
		if u.Name == user {
			return u.Vacation
			break
		}
	}
	return false
}

// diese Funktion öffnet das User.json und gibt eine Liste der User in FOrm eines Structs zurück
func OpenUsers() Users {
	file, err := ioutil.ReadFile("./pkg/users/users.json")
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
	var users Users

	err = json.Unmarshal(file, &users)
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}

	return users
}

func checkUserValid(name, pswd string) bool {
	var users = OpenUsers().User

	for _, u := range users {
		if u.Name == name && u.Pass == pswd { //TODO: Passwortprüfung mit hash-compare ersetzen
			//fmt.Println("user: ", name, "| password: ", pswd)
			//fmt.Println("---")
			//TODO: check mit Test ersetzen
			LoggedUserID = u.ID
			LoggedUserName = u.Name
			LoggedUserVac = u.Vacation
			//fmt.Println("user: ", LoggedUserName, "| ID: ", LoggedUserID, "| Vacation: ", LoggedUserVac)
			//fmt.Println("---")
			//TODO: check mit Test ersetzen
			return true
			break
		}
	}
	LoggedUserID = 0
	LoggedUserName = ""
	LoggedUserVac = false

	return false
}

func Wrapper(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if ok && checkUserValid(user, pass) {
			handler(w, r)
		} else {
			w.Header().Set("WWW-Authenticate",
				"Basic realm=\"Please enter your username and password\"")
			http.Error(w,
				http.StatusText(http.StatusUnauthorized),
				http.StatusUnauthorized)
		}
	}
}

//TODO: Funktionen einsetzen um User anzulegen und dann zu prüfen
// nutzt eine Funktion aus golang.org/x/crypto/bcrypt, um ein eingegebenes Passwort zu 'salten' und zu 'hashen'
func saltAndHash(pass string) string {
	passByte := []byte(pass)

	hash, err := bcrypt.GenerateFromPassword(passByte, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}

	return string(hash)
}

// nutzt eine Funktion aus golang.org/x/crypto/bcrypt, um ein 'hashed' und 'salted' Passwort mit einem Passwort in Klartext zu vergleichen
func checkPass(passHash string, passPlain string) bool {
	passHashByte := []byte(passHash)
	passPlainByte := []byte(passPlain)

	err := bcrypt.CompareHashAndPassword(passHashByte, passPlainByte)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}
