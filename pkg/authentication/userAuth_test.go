//2057008, 2624395, 9111696

package authentication

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
)

func init() {
	err := os.Chdir("../../")
	if err != nil {
		panic(err)
	}
}

var testUser = "admin"
var testPass = "supersecret"

func TestCheckLoggedUserName(t *testing.T) {
	request := httptest.NewRequest("", "/", nil)
	request.SetBasicAuth(testUser, testPass)
	if CheckLoggedUserName(request) != testUser {
		t.Error()
	}
}

func TestCheckLoggedUserID(t *testing.T) {
	request := httptest.NewRequest("", "/", nil)
	request.SetBasicAuth(testUser, testPass)
	if CheckLoggedUserID(request) != 1 {
		t.Error()
	}
}

func TestCheckLoggedUserVac(t *testing.T) {
	request := httptest.NewRequest("", "/", nil)
	request.SetBasicAuth(testUser, testPass)
	if CheckLoggedUserVac(request) != false {
		t.Error()
	}
}

func TestOpenUsers(t *testing.T) {
	var users = OpenUsers()
	if len(users.User) == 0 {
		t.Error()
	}
	if reflect.TypeOf(OpenUsers()) != reflect.TypeOf(users) {
		t.Error()
	}
}

func TestCheckUserValid(t *testing.T) {
	if checkUserValid(testUser, testPass) != true {
		t.Error()
	}
	if checkUserValid(testUser, "random") == true {
		t.Error()
	}
}

func TestWrapper(t *testing.T) {
	var testHandler http.HandlerFunc
	if reflect.TypeOf(Wrapper(HandlerRegister)) != reflect.TypeOf(testHandler) {
		t.Error()
	}
}

func TestSaltAndHash(t *testing.T) {
	var testString = "testPass"
	err := bcrypt.CompareHashAndPassword([]byte(saltAndHash(testString)), []byte(testString))
	if err != nil {
		t.Error()
	}
}

func TestCheckPass(t *testing.T) {
	testHash, _ := bcrypt.GenerateFromPassword([]byte(testPass), bcrypt.DefaultCost)
	if checkPass(string(testHash), testPass) != true {
		t.Error()
	}

	if checkPass(string(testHash), testUser) == true {
		t.Error()
	}
}

func TestSaveAllUsers(t *testing.T) {
	users := OpenUsers()
	err := saveAllUsers(users)
	if err != nil {
		t.Error()
	}
}

func TestWriteUser(t *testing.T) {
	users := OpenUsers()
	var newUser User
	err := writeUser(users, newUser)
	if err != nil {
		t.Error()
	}

	saveAllUsers(users)
}

func TestRegisterUser(t *testing.T) {
	users := OpenUsers()
	err := registerUser("testAdmin", testPass)
	if err != nil {
		t.Error()
	}

	saveAllUsers(users)
}

func TestHandlerRegister(t *testing.T) {
	users := OpenUsers()
	req, err := http.NewRequest("GET", "/register.html", nil)

	if err != nil {
		t.Fatal(err)
	}

	request := httptest.NewRecorder()
	handler := http.HandlerFunc(HandlerRegister)
	handler.ServeHTTP(request, req)

	status := request.Code
	fmt.Println(status)
	if status != http.StatusFound {
		t.Error()
	}
	saveAllUsers(users)
}
