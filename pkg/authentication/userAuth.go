package authentication

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type User struct {
	ID       int    `json:"ID"`
	Name     string `json:"Name"`
	Pass     string `json:"Pass"`
	Vacation bool   `json:"Vacation"`
	// TODO: brauchen wir eine Liste von Tickets, die dem User zugewiesen sind?
}

type Users struct {
	User []User `json:"Users"`
}

func openUsers() Users {
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
	var users = openUsers().User

	for _, u := range users {
		if u.Name == name && u.Pass == pswd {
			//fmt.Println("user: ", name, ", password: ", pswd)
			//TODO: check mit Test ersetzen
			return true
			break
		}
	}
	return false
}

func Wrapper(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pswd, ok := r.BasicAuth()
		if ok && checkUserValid(user, pswd) {
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
