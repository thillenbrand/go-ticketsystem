//2057008, 2624395, 9111696

package authentication

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
)

// struct, das einen User darstellt
type User struct {
	ID       int    `json:"ID"`
	Name     string `json:"Name"`
	Pass     string `json:"Pass"`
	Vacation bool   `json:"Vacation"`
}

// Users ist lediglich eine Auflistung meherer User
type Users struct {
	User []User `json:"Users"`
}

// gibt den Namen des momentan eingeloggten Users aus
func CheckLoggedUserName(r *http.Request) string {
	user, _, _ := r.BasicAuth()

	return user
}

// gibt sie ID des momentan eingeloggten Users aus
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

// gibt den Uralubsstatus des momentan eingeloggten Users aus
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

// diese Funktion öffnet das User.json und gibt eine Liste der User in Form eines Structs zurück
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

// überprüft, ob die übergebene Logindaten mit einem User aus der Liste gespeicherter User übereinstimmen
func checkUserValid(name, pass string) bool {
	var users = OpenUsers().User

	for _, u := range users {
		if u.Name == name && checkPass(u.Pass, pass) {
			return true
			break
		}
	}

	return false
}

// Wrapper-Funktion, die zusätzlich zum übergebenen Handler die Athentifizierung durchführt
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

// nutzt eine Funktion aus golang.org/x/crypto/bcrypt, um ein eingegebenes Passwort zu 'salten' und zu 'hashen'
func saltAndHash(pass string) string {
	passByte := []byte(pass)

	hash, err := bcrypt.GenerateFromPassword(passByte, bcrypt.DefaultCost)
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

// speichert alle User im übergebenen Users-Struct in die users.json
func saveAllUsers(u Users) error {
	filename := "./pkg/users/users.json"
	users, err := json.Marshal(u)
	if err != nil {
		fmt.Println(err)
	}
	return ioutil.WriteFile(filename, users, 0600)
}

// trägt einen neuen User in das bestehende User-Dokument ein
func writeUser(users Users, newUser User) error {
	users.User = append(users.User, newUser)
	err := saveAllUsers(users)
	if err != nil {
		fmt.Println(err)
		return errors.New("userAuth: writing the userdata into the userfile failed")
	}
	return nil
}

// führt das Erstellen des neuen Users als struct durch un ruft anschließend die Funktion zum Eintragen des Users auf
func registerUser(username string, pass string) error {
	users := OpenUsers()
	oneUser := users.User
	var counter float64 = 0

	for _, u := range oneUser {
		counter = math.Max(float64(u.ID), counter)
	}

	var newUser User
	newUser.ID = int(counter) + 1
	newUser.Name = username
	newUser.Pass = saltAndHash(pass)
	newUser.Vacation = false

	for _, u := range oneUser {
		if u.Name == username {
			return errors.New("userAuth: username is already taken")
		}
	}

	writeUser(users, newUser)

	return nil
}

// handler, der sich um die Registierung neuer User kümmert
func HandlerRegister(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("inputName")
	pass := r.FormValue("inputPassword")

	err := registerUser(username, pass)
	if err != nil {
		log.Println(err)
		http.Error(w,
			"Der eingegebene Benutzername ist bereits vergeben!",
			http.StatusUnauthorized)
	}

	http.Redirect(w, r, "/index.html", http.StatusFound)
}
