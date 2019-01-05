package frontend

import (
	"os"
	"testing"
)

func TestOpenTickets(t *testing.T) {
	err := os.Chdir("../../")
	if err != nil {
		panic(err)
	}
	var tickets = openTickets()
	if len(tickets) == 0 {
		t.Error()
	}
}
