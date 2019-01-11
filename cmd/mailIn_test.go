//2057008, 2624395, 9111696

package cmd

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
