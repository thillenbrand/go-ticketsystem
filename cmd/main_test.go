package main

import "testing"

func TestEnterMail(t *testing.T) {
	entermail()
	if mail.Address == "" {
		t.Error()
	}
	if mail.Subject == "" {
		t.Error()
	}
	if mail.Text == "" {
		t.Error()
	}
}
