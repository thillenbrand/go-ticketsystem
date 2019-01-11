//2057008, 2624395, 9111696

package api_out

import (
	"bytes"
	"errors"
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

type errReader int

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("test error")
}

func TestHandlerConfirmSend(t *testing.T) {
	mails := GetMailsFromQueue()

	var jsonIDs = []byte(`{"IDs": [{"ID":0}]}`)
	req, err := http.NewRequest("POST", "/confirmSend", bytes.NewBuffer(jsonIDs))
	if err != nil {
		t.Error()
	}

	request := httptest.NewRecorder()
	handler := http.HandlerFunc(HandlerConfirmSend)
	handler.ServeHTTP(request, req)

	status := request.Code
	if status != http.StatusOK {
		t.Error()
	}

	req2, err2 := http.NewRequest("PUT", "/confirmSend", errReader(0))
	if err2 != nil {
		t.Error()
	}

	request2 := httptest.NewRecorder()
	handler2 := http.HandlerFunc(HandlerConfirmSend)
	handler2.ServeHTTP(request2, req2)

	status2 := request2.Code
	if status2 != http.StatusBadRequest {
		t.Error()
	}
	SaveAllMails(mails)
}

func TestGetQueueFile(t *testing.T) {
	mails := GetMailsFromQueue()
	_, err := http.NewRequest("POST", "/getMailQueue", nil)
	if err != nil {
		t.Error()
	}
	SaveAllMails(mails)
}

func TestSendMailQueue(t *testing.T) {
	err := sendMailQueue()
	if err != nil {
		t.Error()
	}
}

func TestGetMailsFromQueue(t *testing.T) {
	var testMailQueue MailQueue
	if reflect.TypeOf(GetMailsFromQueue()) != reflect.TypeOf(testMailQueue) {
		t.Error()
	}
	if len(testMailQueue.Mail) != 0 {
		t.Error()
	}
}

func TestSaveAllMails(t *testing.T) {
	mails := GetMailsFromQueue()
	err := SaveAllMails(mails)
	if err != nil {
		t.Error()
	}
}

func TestFeedMailQueue(t *testing.T) {
	testAdress := "Dummy Adress"
	testSubject := "Dummy Subject"
	testMessage := "Dummy Message"
	mailQueue := GetMailsFromQueue()

	err := FeedMailQueue(testAdress, testSubject, testMessage)
	if err != nil {
		t.Error()
	}

	SaveAllMails(mailQueue)
}

func TestNotContains(t *testing.T) {
	testSlice := make([]int, 5)
	testIntTrue := 20
	testIntFalse := 0

	if notContains(testIntTrue, testSlice) != true {
		t.Error()
	}

	if notContains(testIntFalse, testSlice) != false {
		t.Error()
	}
}

func TestConfirmMailSent(t *testing.T) {
	testSlice := make([]int, 1)
	mailQueue := GetMailsFromQueue()

	err := ConfirmMailSent(testSlice)
	if err != nil {
		t.Error()
	}

	SaveAllMails(mailQueue)
}
