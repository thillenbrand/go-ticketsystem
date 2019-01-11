//2057008, 2624395, 9111696

package api

import (
	"reflect"
	"testing"
)

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
	err := saveAllMails(mails)
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

	saveAllMails(mailQueue)
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

	saveAllMails(mailQueue)
}
